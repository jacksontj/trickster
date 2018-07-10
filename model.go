/**
* Copyright 2018 Comcast Cable Communications Management, LLC
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/prometheus/common/model"
)

type PrometheusEnvelope struct {
	Status string         `json:"status"`
	Data   PrometheusData `json:"data"`
}

type PrometheusData struct {
	ResultType string      `json:"resultType"`
	Result     model.Value `json:"result"`
}

func (d *PrometheusData) UnmarshalJSON(data []byte) error {
	type Raw struct {
		ResultType model.ValueType  `json:"resultType"`
		Result     *json.RawMessage `json:"result"`
	}
	rawData := &Raw{}
	if err := json.Unmarshal(data, rawData); err != nil {
		return err
	}
	d.ResultType = rawData.ResultType.String()
	switch rawData.ResultType {
	case model.ValNone:
		return nil
	case model.ValScalar:
		tmp := &model.Scalar{}
		if err := json.Unmarshal(*rawData.Result, tmp); err != nil {
			return err
		}
		d.Result = tmp
	case model.ValVector:
		tmp := model.Vector{}
		if err := json.Unmarshal(*rawData.Result, &tmp); err != nil {
			return err
		}
		d.Result = tmp
	case model.ValMatrix:
		tmp := model.Matrix{}
		if err := json.Unmarshal(*rawData.Result, &tmp); err != nil {
			return err
		}
		d.Result = tmp
	case model.ValString:
		tmp := &model.String{}
		if err := json.Unmarshal(*rawData.Result, tmp); err != nil {
			return err
		}
		d.Result = tmp
	default:
		return fmt.Errorf("Unknown result type: %v", d.ResultType)
	}

	return nil
}

// ClientRequestContext contains the objects needed to fulfull a client request
type ClientRequestContext struct {
	Request           *http.Request
	Writer            http.ResponseWriter
	CacheKey          string
	CacheLookupResult string
	// TODO: do we need this?
	Matrix             PrometheusEnvelope
	Origin             PrometheusOriginConfig
	RequestParams      url.Values
	RequestExtents     MatrixExtents
	OriginUpperExtents MatrixExtents
	OriginLowerExtents MatrixExtents
	StepParam          string
	StepMS             int64
	Time               int64
	WaitGroup          sync.WaitGroup
}

// MatrixExtents describes the start and end epoch times (in ms) for a given range of data
type MatrixExtents struct {
	Start int64
	End   int64
}

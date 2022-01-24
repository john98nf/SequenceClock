// Copyright © 2021 Giannis Fakinos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

/*
	This file is not included in the build process
	and will be used by the deployer.
	Please, do not edit!!
*/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/apache/openwhisk-client-go/whisk"
)

type controller (func(map[string]interface{}) map[string]interface{})

var (
	client         *whisk.Client
	controllerType = map[string]controller{
		"greedy": greedyControl,
		"dummy":  dummyControl,
	}
)

/*
	Main serveless function.
*/
func Main(obj map[string]interface{}) map[string]interface{} {
	wskConfig := &whisk.Config{
		Host:      os.Getenv("__OW_API_HOST"),
		Namespace: os.Getenv("__OW_NAMESPACE"),
		AuthToken: os.Getenv("__OW_API_KEY"),
		Insecure:  true,
	}
	client, _ = whisk.NewClient(http.DefaultClient, wskConfig)
	return controllerType[ALGORITHM_TYPE](obj)
}

/*
	No actual control over function invocation,
	as it just invokes each function.
	Used for benchmarking purposes and referrence point.
*/
func dummyControl(obj map[string]interface{}) map[string]interface{} {
	var (
		id         string
		status     string
		latency    int64
		fullRes    map[string]interface{}
		err        error
		errMetrics error
	)
	aRes := obj
	for i, f := range functionList {
		fullRes, _, err = client.Actions.Invoke(f, aRes, true, false)
		if err != nil {
			panic(err)
		}
		id, status, latency, aRes, errMetrics = extractMetrics(fullRes)
		if errMetrics != nil {
			panic(errMetrics.Error())
		}
		fmt.Println(i, id, latency, strings.Replace(status, " ", "", -1))
	}
	return aRes
}

/*
	Greedy control.
	Each function invokation execution time is been monitored and its actual elapsed time is been compared
	with a referrence point (pre-configured). Their difference is called slack.
	A positive slack means that the next function invokation may run with fewer resources (slow down).
	A negative slack means that the next fuction must run with more resources (speed up).
*/
func greedyControl(obj map[string]interface{}) map[string]interface{} {
	var (
		tStart     time.Time
		tEnd       time.Time
		elapsed    time.Duration
		id         string
		status     string
		latency    int64
		fullRes    map[string]interface{}
		r          *Request = NewRequest("", &Metrics{})
		errMetrics error
	)

	watcherClient := NewWatcherClient(KUBE_MAIN_IP)
	aRes := obj
	for i, f := range functionList {
		tStart = time.Now()
		r.Function = functionList[i]
		reset, err := watcherClient.RequestResources(r)
		if err != nil {
			panic(err)
		}

		fullRes, _, err = client.Actions.Invoke(f, aRes, true, false)

		if err := watcherClient.ResetResources(reset); err != nil {
			panic(err)
		}
		if err != nil {
			panic(err)
		}
		id, status, latency, aRes, errMetrics = extractMetrics(fullRes)
		if errMetrics != nil {
			panic(errMetrics)
		}
		fmt.Println(i, id, latency, strings.Replace(status, " ", "", -1), r.Metrics.Slack)
		if status != "success" {
			panic(fmt.Errorf("invocation terminated with status: %s", status))
		}

		tEnd = time.Now()
		elapsed = tEnd.Sub(tStart)
		r.Metrics.PreviousSlack = r.Metrics.Slack
		r.Metrics.Slack += profiledExecutionTimes[i] - int64(elapsed)
		r.Metrics.SumOfSlack += r.Metrics.Slack
	}
	return aRes
}

/*
	Helper function for extracting information
	from OpenWhisk API output.
*/
func extractMetrics(r map[string]interface{}) (string, string, int64, map[string]interface{}, error) {
	end, err1 := r["end"].(json.Number).Int64()
	start, err2 := r["start"].(json.Number).Int64()
	id, ok3 := r["activationId"].(string)
	status, ok4 := r["response"].(map[string]interface{})["status"].(string)
	res, ok5 := r["response"].(map[string]interface{})["result"].(map[string]interface{})

	if err1 != nil || err2 != nil || !ok3 || !ok4 || !ok5 {
		return id, status, (end - start), res, fmt.Errorf("problem with type assertion (end, start, activationId, status, result): (%v, %v, %v, %v, %v)", err1.Error(), err2.Error(), ok3, ok4, ok5)
	}
	return id, status, (end - start), res, nil

}

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
	"log"
	"net/http"
	"os"
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
	aRes := obj
	for i, f := range functionList {
		log.Printf("Invoking function-%v %v\n", i, f)
		aRes, _, _ = client.Actions.Invoke(f, aRes, true, true)
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
		tStart  time.Time
		tEnd    time.Time
		elapsed time.Duration
		r       *Request = NewRequest("", &Metrics{})
	)

	watcherClient := NewWatcherClient(KUBE_MAIN_IP)
	aRes := obj
	for i, f := range functionList {
		tStart = time.Now()
		log.Printf("Slack: %v\n", r.Metrics.Slack)
		r.Function = functionList[i]
		r.Metrics.ProfiledExecutionTime = profiledExecutionTimes[i]
		reset, err := watcherClient.RequestResources(r)
		if err != nil {
			log.Printf("Error from watcher client: %v\n", err.Error())
		}

		log.Printf("Invoking function-%v %v\n", i, f)
		aRes, _, _ = client.Actions.Invoke(f, aRes, true, true)

		log.Println("Sending reset request")
		if err := watcherClient.ResetResources(reset); err != nil {
			log.Printf("Error from watcher client: %v\n", err.Error())
		}
		tEnd = time.Now()
		elapsed = tEnd.Sub(tStart)
		log.Printf("Elapsed time: %v\n", elapsed)
		r.Metrics.PreviousSlack = r.Metrics.Slack
		r.Metrics.Slack += profiledExecutionTimes[i] - int64(elapsed)
		r.Metrics.SumOfSlack += r.Metrics.Slack
	}
	log.Printf("Last slack %v\n", r.Metrics.Slack)
	return aRes
}

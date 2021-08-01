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

package main

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	SpeedUpRequest int = iota + 1
	SlowDownRequest
	ResetRequest
)

type RequestType int

type watcherClientInterface interface {
	SpeedUpRequest(function string) (bool, error)
	SlowDownRequest(function string) (bool, error)
	ResetRequest(function string) (bool, error)
	ExecuteRequest(function string, requestType RequestType) (bool, error)
}

type WatcherClient struct {
	endpoint string `json:endpoint,omitempty`
}

func NewWatcherClient(host string) *WatcherClient {
	return &WatcherClient{
		endpoint: "http://" + host + ":32042/api/function/requestResources",
	}
}

func (client *WatcherClient) SpeedUpRequest(function string) (bool, error) {
	return client.ExecuteRequest(function, SpeedUpRequest)
}

func (client *WatcherClient) SlowDownRequest(function string) (bool, error) {
	return client.ExecuteRequest(function, SlowDownRequest)
}

func (client *WatcherClient) ResetRequest(function string) (bool, error) {
	return client.ExecuteRequest(function, ResetRequest)
}

func (client *WatcherClient) ExecuteRequest(f string, rT int) (bool, error) {
	params := url.Values{}
	params.Add("Function", f)
	params.Add("Type", fmt.Sprint(rT))
	resp, err := http.PostForm(client.endpoint, params)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true, nil
	} else {
		return false, nil
	}
}

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/iris-contrib/schema"
)

type watcherClientInterface interface {
	RequestResources(r *Request) (*ResetRequest, error)
	ResetResources(r *ResetRequest) error
}

type WatcherClient struct {
	endpoint string
}

func NewWatcherClient(host string) *WatcherClient {
	return &WatcherClient{
		endpoint: "http://" + host + ":32042/api/function",
	}
}

func (client *WatcherClient) RequestResources(r *Request) (*ResetRequest, error) {
	body, err := postHTTPRequest(client.endpoint+"/requestResources", *r)
	if err != nil {
		return nil, err
	} else {
		res := NewResetRequest(0, r.Function)
		err := json.Unmarshal(body, res)
		return res, err
	}
}

func (client *WatcherClient) ResetResources(r *ResetRequest) error {
	_, err := postHTTPRequest(client.endpoint+"/resetResources", *r)
	return err
}

func postHTTPRequest(endpoint string, data interface{}) ([]byte, error) {
	var encoder = schema.NewEncoder()
	params := url.Values{}
	if err := encoder.Encode(data, params); err != nil {
		return nil, err
	}
	resp, err := http.PostForm(endpoint, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		return body, nil
	} else {
		return nil, fmt.Errorf(string(body))
	}
}

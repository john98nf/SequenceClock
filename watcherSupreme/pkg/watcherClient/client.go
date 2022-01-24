// Copyright © 2021 Giannis Fakinos

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

package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	wrq "github.com/john98nf/SequenceClock/watcher/pkg/request"

	"github.com/iris-contrib/schema"
)

type WatcherInterface interface {
	SendRequest(r *wrq.Request) (bool, error)
	SendResetRequest(r wrq.ResetRequest) (bool, error)
}

type WatcherClient struct {
	Node    string
	BaseURL string
}

func NewWatcherClient(node string) *WatcherClient {
	return &WatcherClient{
		Node:    node,
		BaseURL: "http://" + node + ":8080/api/function",
	}
}

func (w *WatcherClient) SendRequest(r *wrq.Request) (bool, error) {
	return w.executeRequest(w.BaseURL+"/requestResources", *r)
}

func (w *WatcherClient) SendResetRequest(r *wrq.ResetRequest) (bool, error) {
	return w.executeRequest(w.BaseURL+"/resetRequest", *r)
}

func (w *WatcherClient) executeRequest(endpoint string, msg interface{}) (bool, error) {
	var encoder = schema.NewEncoder()
	params := url.Values{}
	if err := encoder.Encode(msg, params); err != nil {
		return false, err
	}
	resp, err := http.PostForm(endpoint, params)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true, nil
	} else if resp.StatusCode == 404 {
		return false, nil
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf(string(body))
	}
}

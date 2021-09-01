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
	"net/http"
	"net/url"

	wrq "github.com/john98nf/SequenceClock/watcher/pkg/request"

	"github.com/iris-contrib/schema"
)

type WatcherInterface interface {
	ExecuteRequest(r *wrq.Request) (bool, error)
}

type WatcherClient struct {
	Node        string
	RequestsURL string
}

func NewWatcherClient(node string) *WatcherClient {
	return &WatcherClient{
		Node:        node,
		RequestsURL: "http://" + node + ":8080/api/function/requestResources",
	}
}

func (w *WatcherClient) ExecuteRequest(msg interface{}) (bool, error) {
	var encoder = schema.NewEncoder()
	params := url.Values{}
	if err := encoder.Encode(msg, params); err != nil {
		return false, err
	}
	resp, err := http.PostForm(w.RequestsURL, params)
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

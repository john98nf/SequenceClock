package client

import (
	wrq "john98nf/SequenceClock/watcher/pkg/request"
	"net/http"
	"net/url"

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
		RequestsURL: "http://" + node + ":8080/api/function/requestResoures",
	}
}

func (w *WatcherClient) ExecuteRequest(r *wrq.Request) (bool, error) {
	var encoder = schema.NewEncoder()
	params := url.Values{}
	if err := encoder.Encode(r, params); err != nil {
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

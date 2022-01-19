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

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	wrq "github.com/john98nf/SequenceClock/watcher/pkg/request"
	wrc "github.com/john98nf/SequenceClock/watcherSupreme/pkg/watcherClient"

	"github.com/gin-gonic/gin"
)

var (
	clients         []*wrc.WatcherClient
	counterID       uint64
	mutex           = sync.RWMutex{}
	requestCatalog  = map[uint64]*wrc.WatcherClient{}
	functionCatalog = map[string]*wrc.WatcherClient{}
)

func main() {
	router := gin.New()

	apiWatcher := router.Group("/api")
	{
		// GET Request http://localhost:8080/api/check
		apiWatcher.GET("/check", check)
		// POST Request http://localhost:8080/api/function/requestResources
		apiWatcher.POST("/function/requestResources", requestHandler)
		// POST Request http://localhost:8080/api/function/resetResources
		apiWatcher.POST("/function/resetResources", resetHandler)
		// GET Request http://localhost:8080/api/catalogs
		apiWatcher.GET("/catalogs", getCatalogs)
	}

	clients = connectWatchers()
	router.Run(":8080")
}

/*
	Liveness & Readiness probe for watcher.
*/
func check(c *gin.Context) {
	c.String(http.StatusOK, "Hello from watcher supreme!")
}

/*
	SpeedUp/Normal/SlowDown Request for certain function.
	Watcher supreme wil inform watchers, which will
	search their docker runtime for the actual docker container.
*/
func requestHandler(c *gin.Context) {
	var req wrq.Request
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	mutex.Lock()
	reset := wrq.NewResetRequest(counterID, req.Function)
	req.ID = reset.ID
	counterID++
	mutex.Unlock()
	go requestResourceAllocationFromWatchers(req)

	c.JSON(http.StatusOK, *reset)
}

/*
	Reset Request
	Watcher Supreme informs only involved watcher nodes,
	in order for them to deactivate the chosen request.
*/
func resetHandler(c *gin.Context) {
	var rs wrq.ResetRequest
	if err := c.ShouldBind(&rs); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	go resetRequestToWatchers(rs)

	c.String(http.StatusOK, "ok")
}

/*
	Iterate through all nodes.
	Request for resource allocation in every runtime
	that the certain docker container is found.
*/
func requestResourceAllocationFromWatchers(req wrq.Request) {
	mutex.RLock()
	c, ok := functionCatalog[req.Function]
	mutex.RUnlock()
	if ok {
		found, err := c.SendRequest(&req)
		if err != nil && err != fmt.Errorf("{\"message\":\"Not Found\"}") {
			log.Println(err)
		} else if found {
			mutex.Lock()
			requestCatalog[req.ID] = c
			mutex.Unlock()
			return
		}
	}
	for _, c := range clients {
		if res, err := c.SendRequest(&req); err != nil {
			log.Println(err)
			continue
		} else if res {
			mutex.Lock()
			requestCatalog[req.ID] = c
			functionCatalog[req.Function] = c
			mutex.Unlock()
			return
		}
	}
}

/*
	Reach only the neccessary cluster node and
	send a reset request for a specific serverless function.
*/
func resetRequestToWatchers(rs wrq.ResetRequest) {
	var (
		ok     bool
		client *wrc.WatcherClient
	)
	for {
		mutex.RLock()
		client, ok = requestCatalog[rs.ID]
		mutex.RUnlock()
		if ok {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if res, err := client.SendResetRequest(&rs); err != nil {
		log.Println("Problem with watcher:", err.Error())
	} else if !res {
		log.Println("Problem with watcher:", client.Node)
	}
	mutex.Lock()
	delete(requestCatalog, rs.ID)
	mutex.Unlock()
}

/*
	Get information for catalog data structures.
	Development oriented api call.
*/
func getCatalogs(c *gin.Context) {
	req := make(map[uint64]interface{})
	fc := make(gin.H)
	mutex.RLock()
	for k, c := range requestCatalog {
		req[k] = c.Node
	}
	for k, c := range functionCatalog {
		fc[k] = c.Node
	}
	mutex.RUnlock()
	c.JSON(http.StatusOK, gin.H{"requests": req, "functions": fc})
}

/*
	Initiation function. Creates watcher clients
	for cluster nodes.
*/
func connectWatchers() []*wrc.WatcherClient {
	// TODO: Call kubernetes api for automatic
	// node discovery.
	var nodes []string = []string{
		"192.168.1.244",
		"192.168.1.245",
		"192.168.1.246",
	}
	res := make([]*wrc.WatcherClient, len(nodes))
	for i, n := range nodes {
		res[i] = wrc.NewWatcherClient(n)
	}
	return res
}

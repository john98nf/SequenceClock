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
	"log"
	"net/http"

	wrq "github.com/john98nf/SequenceClock/watcher/pkg/request"
	wrc "github.com/john98nf/SequenceClock/watcherSupreme/pkg/watcherClient"

	"github.com/gin-gonic/gin"
)

var (
	clients   []*wrc.WatcherClient
	placement map[uint64][]*wrc.WatcherClient
	counterID uint64
)

func main() {
	router := gin.Default()

	apiWatcher := router.Group("/api")
	{
		// GET Request http://localhost:8080/api/check
		apiWatcher.GET("/check", check)
		// POST Request http://localhost:8080/api/function/requestResources
		apiWatcher.POST("/function/requestResources", requestHandler)
		// POST Request http://localhost:8080/api/function/resetResources
		apiWatcher.POST("/function/resetResources", resetHandler)
	}

	placement = make(map[uint64][]*wrc.WatcherClient)
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
	reset := wrq.NewResetRequest(counterID, req.Function)
	req.ID = reset.ID
	counterID++
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
	// TO DO: in case that similar request has already been received
	// send immediate requests to watchers.
	if _, ok := placement[req.ID]; !ok {
		placement[req.ID] = make([]*wrc.WatcherClient, len(clients))
	}
	for found := false; !found; {
		for _, c := range clients {
			res, err := c.ExecuteRequest(&req)
			if err != nil {
				log.Println(err)
				continue
			}
			if res {
				placement[req.ID] = append(placement[req.ID], c)
			}
			found = found || res
		}
	}
}

/*
	Reach only the neccessary cluster nodes and
	send them a reset request for a specific serverless function.
*/
func resetRequestToWatchers(rs wrq.ResetRequest) {
	subClients, ok := placement[rs.ID]
	if !ok {
		log.Printf("No #%v recent request was found\n", rs.ID)
		return
	}

	for _, c := range subClients {
		if res, err := c.ExecuteRequest(rs); err != nil {
			log.Println("Problem with watcher:", err.Error())
		} else if !res {
			log.Println("Problem with watcher:", c.Node)
		}
	}
	delete(placement, rs.ID)
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

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
	"net/http"
	"os"

	"github.com/john98nf/SequenceClock/watcher/internal/conflicts"
	wrq "github.com/john98nf/SequenceClock/watcher/pkg/request"

	"github.com/gin-gonic/gin"
)

var (
	hostIP           string = os.Getenv("HOST_IP")
	conflictResolver conflicts.ConflictResolver
)

func main() {
	router := gin.Default()

	apiWatcher := router.Group("/api")
	{
		// GET Request http://localhost:8080/api/check
		apiWatcher.GET("/check", check)
		// GET Request http://localhost:8080/api/function/{name}
		apiWatcher.GET("/function/:name", getContainer)
		// POST Request http://localhost:8080/api/function/requestResources
		apiWatcher.POST("/function/requestResources", requestHandler)
		// POST ResetRequest http://localhost:8080/api/function/resetRequest
		apiWatcher.POST("/function/resetRequest", resetHandler)
	}
	conflictResolver = conflicts.NewConflictResolver()
	router.Run(":8080")
}

/*
	Liveness & Readiness probe for watcher.
*/
func check(c *gin.Context) {
	c.String(http.StatusOK, "Hello from watcher inside node %v!", hostIP)
}

/*
	Delete previous request with specified id.
	Docker container returns to initial state of
	CPU_QUOTAS = -1.
*/
func resetHandler(c *gin.Context) {
	var rs wrq.ResetRequest
	if err := c.ShouldBind(&rs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := conflictResolver.RemoveRequest(rs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

/*
	Provides or Removes resources from openwhisk function
	docker container resources.
*/
func requestHandler(c *gin.Context) {
	var req wrq.Request
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	container, err := searchFunctionContainer(req.Function, "user-action")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if container == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No related docker container was found"})
		return
	}
	// TO DO: Introduce not manual cpu-quota change
	var reqCPUQuota int64
	switch req.Type {
	case wrq.SpeedUpRequest:
		reqCPUQuota = int64(150000)
	case wrq.SlowDownRequest:
		reqCPUQuota = int64(50000)
	case wrq.ResetRequest:
		reqCPUQuota = int64(-1)
	}
	if err := UpdateContainerCPUQuota(container.ID, reqCPUQuota); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Docker container updated",
	})
}

/*
	Get information for function related
	container. Development oriented api call.
*/
func getContainer(c *gin.Context) {
	fName := c.Param("Function")
	podType := c.DefaultQuery("type", "0")
	pdt := map[string]string{
		"0": "user-action",
		"1": "POD",
	}
	if _, ok := pdt[podType]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not supported pod type."})
		return
	}
	cnt, err := conflictResolver.SearchRegistry(fName, pdt[podType])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if cnt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No such container"})
		return
	}
	c.JSON(http.StatusOK, cnt)
}

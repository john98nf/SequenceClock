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
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

const (
	REG_EXP1 string = `(^\/k8s_user-action_wskowdev-invoker-00-[0-9][0-9]-guest-%v_openwhisk_)`
	REG_EXP2 string = `(^\/k8s_POD_wskowdev-invoker-00-[0-9][0-9]-guest-%v_openwhisk_)`
)

var (
	hostIP string = os.Getenv("HOST_IP")
	cli    *client.Client
)

func main() {
	router := gin.Default()

	apiWatcher := router.Group("/api")
	{
		// GET Request http://localhost:8080/api/check
		apiWatcher.GET("/check", check)
		// GET Request http://localhost:8080/api/function/:name
		apiWatcher.GET("/function/:name", getContainer)
		// PATCH Request http://localhost:8080/api/function/:name/speedUp
		apiWatcher.PATCH("/function/:name/speedUp", speedUp)
		// PATCH Request http://localhost:8080/api/function/:name/slowDown
		apiWatcher.PATCH("/function/:name/slowDown", slowDown)
	}

	cliLocal, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	cli = cliLocal

	router.Run(":8080")
}

/*
	Liveness & Readiness probe for watcher.
*/
func check(c *gin.Context) {
	c.String(http.StatusOK, "Hello from watcher inside node %v!", hostIP)
}

/*
	Provides more reqources to specified
	function pod.
*/
func speedUp(c *gin.Context) {
	functionName := c.Param("name")
	c.String(http.StatusOK, "Watcher %v: SpeedUp request for function %v.", hostIP, functionName)
}

/*
	Removes resources from sepcified
	function pod.
*/
func slowDown(c *gin.Context) {
	functionName := c.Param("name")
	c.String(http.StatusOK, "Watcher %v: SpeedUp request for function %v.", hostIP, functionName)
}

/*
	Get information for function related
	container.
*/
func getContainer(c *gin.Context) {
	f := c.Param("name")
	podType := c.DefaultQuery("type", "running")
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data := make(gin.H)
	var exp string
	if podType == "running" {
		exp = fmt.Sprintf(REG_EXP1, f)
	} else if podType == "paused" {
		exp = fmt.Sprintf(REG_EXP2, f)
	} else {
		c.String(http.StatusBadRequest, "Not supported pod type.")
		return
	}

	for _, cnt := range containers {
		if matched, err := regexp.MatchString(exp, cnt.Names[0]); matched {
			data["ID"] = cnt.ID[:10]
			data["Names"] = cnt.Names
			data["Image"] = cnt.Image
			data["Command"] = cnt.Command
			data["Created"] = cnt.Created
			data["Ports"] = cnt.Ports
			data["SizeRw"] = cnt.SizeRw
			data["SizeRootFs"] = cnt.SizeRootFs
			data["Labels"] = cnt.Labels
			data["State"] = cnt.State
			data["Status"] = cnt.Status
			data["NetworkMode"] = cnt.HostConfig.NetworkMode
			// Ignored NetworkSettings
			data["Mount"] = cnt.Mounts
		} else if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.JSON(200, data)

}

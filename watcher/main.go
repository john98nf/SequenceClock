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

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
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
		// PATCH Request http://localhost:8080/api/function/:name/speedUp
		apiWatcher.PATCH("/function/:name/speedUp", speedUp)
		// PATCH Request http://localhost:8080/api/function/:name/slowDown?name=x
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

// func test(c *gin.Context) {

// 	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	data := make(gin.H)
// 	for _, c := range containers {
// 		data[c.ID[:10]] = c.Names
// 	}
// 	c.JSON(200, data)

// }

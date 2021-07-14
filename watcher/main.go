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

	"github.com/gin-gonic/gin"
)

var hostIP string = os.Getenv("HOST_IP")

func main() {
	router := gin.Default()

	apiWatcher := router.Group("/api")
	{
		// GET Request http://localhost:8080/api/check
		apiWatcher.GET("/check", check)
		// PATCH Request http://localhost:8080/api/function/speedUp?name=x
		apiWatcher.PATCH("/function/speedUp", speedUp)
		// PATCH Request http://localhost:8080/api/function/slowDown?name=x
		apiWatcher.PATCH("/function/slowDown", slowDown)
	}

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
	functionName := c.Query("name")
	c.String(http.StatusOK, "Watcher %v: SpeedUp request for function %v.", hostIP, functionName)
}

/*
	Removes resources from sepcified
	function pod.
*/
func slowDown(c *gin.Context) {
	functionName := c.Query("name")
	c.String(http.StatusOK, "Watcher %v: SpeedUp request for function %v.", hostIP, functionName)
}

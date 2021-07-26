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
	REG_EXP string = `^wskowdev-invoker-00-[0-9][0-9]-guest-%v`
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
		// GET Request http://localhost:8080/api/function
		apiWatcher.GET("/function/:name", getContainer)
		// POST Request http://localhost:8080/api/function/speedUp
		apiWatcher.POST("/function/speedUp", speedUp)
		// POST Request http://localhost:8080/api/function/slowDown
		apiWatcher.POST("/function/slowDown", slowDown)
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
	fName := c.PostForm("name")
	if fName == "" {
		c.String(http.StatusBadRequest, "No function name was provided")
		return
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	exp := fmt.Sprintf(REG_EXP, fName)
	for _, cnt := range containers {
		if l, ok := cnt.Labels["io.kubernetes.container.name"]; ok && l == "user-action" {
			if matched, err := regexp.MatchString(exp, cnt.Labels["io.kubernetes.pod.name"]); matched {
				// TO DO: actual resource allocation for docker container
				c.JSON(http.StatusOK, gin.H{
					"function": fName,
					"message":  "Docker container found",
					"node":     hostIP,
				})
				return
			} else if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
	}
	c.JSON(http.StatusNotFound, gin.H{})
}

/*
	Removes resources from sepcified
	function pod.
*/
func slowDown(c *gin.Context) {
	fName := c.PostForm("name")
	if fName == "" {
		c.String(http.StatusBadRequest, "No function name was provided")
		return
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	exp := fmt.Sprintf(REG_EXP, fName)
	for _, cnt := range containers {
		if l, ok := cnt.Labels["io.kubernetes.container.name"]; ok && l == "user-action" {
			if matched, err := regexp.MatchString(exp, cnt.Labels["io.kubernetes.pod.name"]); matched {
				// TO DO: actual resource allocation for docker container
				c.JSON(http.StatusOK, gin.H{
					"function": fName,
					"message":  "Docker container found",
					"node":     hostIP,
				})
				return
			} else if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
	}
	c.JSON(http.StatusNotFound, gin.H{})
}

/*
	Get information for function related
	container. Test oriented api call.
*/
func getContainer(c *gin.Context) {
	fName := c.Param("name")
	podType := c.DefaultQuery("type", "user-action")
	if podType != "user-action" && podType != "POD" {
		c.String(http.StatusBadRequest, "Not supported pod type.")
		return
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	exp := fmt.Sprintf(REG_EXP, fName)
	for _, cnt := range containers {
		if l, ok := cnt.Labels["io.kubernetes.container.name"]; ok && l == podType {
			if matched, err := regexp.MatchString(exp, cnt.Labels["io.kubernetes.pod.name"]); matched {
				c.JSON(http.StatusOK, gin.H{
					"ID":          cnt.ID[:10],
					"Image":       cnt.Image,
					"Command":     cnt.Command,
					"Created":     cnt.Created,
					"Porst":       cnt.Ports,
					"SizeRw":      cnt.SizeRw,
					"SizeRootfS":  cnt.SizeRootFs,
					"Labels":      cnt.Labels,
					"State":       cnt.State,
					"Statue":      cnt.Status,
					"NetworkMode": cnt.HostConfig.NetworkMode,
					"Mount":       cnt.Mounts,
				})
				return
			} else if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
	}
	c.JSON(http.StatusNotFound, gin.H{})
}

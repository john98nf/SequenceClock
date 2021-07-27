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
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

const (
	REG_EXP                      string = `^wskowdev-invoker-00-[0-9][0-9]-guest-%v`
	BLKIO_WEIGHT_DEFAULT         uint16 = 0
	CPUSET_CPUS_DEFAULT          string = ""
	CPUSET_MEMS_DEFAULT          string = ""
	CPU_SHARES_DEFAULT           int64  = 0 // Openwhisk default 2
	CPU_PERIOD_DEFAULT           int64  = 0 // Openwhisk default 100000
	CPU_QUOTA_DEFAULT            int64  = 0
	CPU_REALTIME_PERIOD_DEFAULT  int64  = 0
	CPU_REALTIME_RUNTIME_DEFAULT int64  = 0
	MEMORY_DEFAULT               int64  = 0 // Openwhisk default 268435456
	MEMORY_RESERVATION_DEFAULT   int64  = 0
	MEMORY_SWAP_DEFAULT          int64  = 0 // Openwhisk default -1
	KERNEL_MEMORY_DEFAULT        int64  = 0
	NANO_CPUS_DEFAULT            int64  = 0
	RESTART_POLICY_DEFAULT       string = "" // Openwhisk default {"Name": "no",MaximumRetryCount: 0}
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
	// slack := c.PostForm("slack")
	// expectedElapsedTime := c.PostForm("elapsedTime")
	if fName == "" {
		c.String(http.StatusBadRequest, "No function name was provided")
		return
	}
	container, err := searchFunctionContainer(fName, "user-action")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} else if container == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	// TO DO: Introduce not manual cpu-quota change
	if err := updateContainerCPUQuota(container.ID, 150000); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"function": fName,
		"message":  "Docker container updated",
		"node":     hostIP,
	})
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
	container, err := searchFunctionContainer(fName, "user-action")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} else if container == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	// TO DO: Introduce not manual cpu-quota change
	if err := updateContainerCPUQuota(container.ID, 50000); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"function": fName,
		"message":  "Docker container updated",
		"node":     hostIP,
	})
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
	cnt, err := searchFunctionContainer(fName, podType)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} else if cnt == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ID":          cnt.ID[:10],
		"Image":       cnt.Image,
		"Command":     cnt.Command,
		"Created":     cnt.Created,
		"Ports":       cnt.Ports,
		"SizeRw":      cnt.SizeRw,
		"SizeRootfS":  cnt.SizeRootFs,
		"Labels":      cnt.Labels,
		"State":       cnt.State,
		"Status":      cnt.Status,
		"NetworkMode": cnt.HostConfig.NetworkMode,
		"Mount":       cnt.Mounts,
	})
}

/*
	Helper function for searching docker runtime
	for an openwhisk action container.
*/
func searchFunctionContainer(name, podType string) (*types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	exp := fmt.Sprintf(REG_EXP, name)
	for _, cnt := range containers {
		if l, ok := cnt.Labels["io.kubernetes.container.name"]; ok && l == podType {
			if matched, err := regexp.MatchString(exp, cnt.Labels["io.kubernetes.pod.name"]); matched {
				return &cnt, nil
			} else if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func updateContainerCPUQuota(containerID string, cpuQuota int64) error {
	updateConfig := containertypes.UpdateConfig{
		Resources: containertypes.Resources{
			BlkioWeight:        BLKIO_WEIGHT_DEFAULT,
			CpusetCpus:         CPUSET_CPUS_DEFAULT,
			CpusetMems:         CPUSET_MEMS_DEFAULT,
			CPUShares:          CPU_SHARES_DEFAULT,
			Memory:             MEMORY_DEFAULT,
			MemoryReservation:  MEMORY_RESERVATION_DEFAULT,
			MemorySwap:         MEMORY_SWAP_DEFAULT,
			KernelMemory:       KERNEL_MEMORY_DEFAULT,
			CPUPeriod:          CPU_PERIOD_DEFAULT,
			CPUQuota:           cpuQuota,
			CPURealtimePeriod:  CPU_REALTIME_PERIOD_DEFAULT,
			CPURealtimeRuntime: CPU_REALTIME_RUNTIME_DEFAULT,
			NanoCPUs:           NANO_CPUS_DEFAULT,
		},
		RestartPolicy: containertypes.RestartPolicy{
			Name:              "no",
			MaximumRetryCount: 0,
		},
	}
	if _, err := cli.ContainerUpdate(context.Background(), containerID, updateConfig); err != nil {
		return err
	}
	return nil
}

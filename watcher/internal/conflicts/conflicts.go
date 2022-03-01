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

package conflicts

import (
	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"sync"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	wfs "github.com/john98nf/SequenceClock/watcher/internal/state"
	wrq "github.com/john98nf/SequenceClock/watcher/pkg/request"
)

const (
	REG_EXP                      string = `^wskowdev-invoker-[0-9]*-[0-9]*-guest-%v`
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
	CPU_PERIOD_OPENWHISK_DEFAULT int64  = 100000
	CPU_QUOTAS_LOWER_BOUND       int64  = 1000
	Kp                           int64  = 10
	Ki                           int64  = 1
	Kd                           int64  = 0
)

var found bool

type ConflictResolverInterface interface {
	SearchDockerRuntime(function, podType string) (*types.Container, error)
	RegistryContains(function string) bool
	InsertToRegistry()
	UpdateRegistry(id uint64, function string, quotas int64) error
	RemoveFromRegistry(rs wrq.ResetRequest) error
	ReconfigureRegistry() error
	ExportRegistry() map[string]interface{}
}

type ConflictResolver struct {
	mutex          sync.RWMutex
	Registry       map[string]*wfs.FunctionState
	DockerClient   *client.Client
	Cores          int64
	lambdaPrevious float64
}

func NewConflictResolver(cores int64) ConflictResolver {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	registry := make(map[string]*wfs.FunctionState)
	return ConflictResolver{
		mutex:        sync.RWMutex{},
		Registry:     registry,
		DockerClient: cli,
		Cores:        cores,
	}
}

/*
	Searches Registry data stucture
	for existing function state.
*/
func (cr *ConflictResolver) RegistryContains(function string) bool {
	cr.mutex.RLock()
	_, ok := cr.Registry[function]
	cr.mutex.RUnlock()
	return ok
}

/*
	Places new request into registry and
	update container resources.
*/
func (cr *ConflictResolver) UpdateRegistry(req *wrq.Request) (bool, error) {
	cr.mutex.Lock()
	state, ok := cr.Registry[req.Function]
	if !ok {
		// TO DO: Solve Openwhisk autoscaling problem
		container, err := cr.SearchDockerRuntime(req.Function, "user-action")
		if err != nil {
			cr.mutex.Unlock()
			return false, err
		} else if container == nil {
			cr.mutex.Unlock()
			return false, nil
		}
		state = wfs.NewFunctionState(container.ID)
		cr.Registry[req.Function] = state
	}

	quotas := retainCPUThreshold(computePIDControllerOutput(req)+CPU_PERIOD_OPENWHISK_DEFAULT, cr.Cores)

	if quotas > state.DesiredQuotas {
		if state.DesiredQuotas != 0 {
			state.Requests.Active[state.Requests.Current] = state.DesiredQuotas
		}
		quotas_old := state.DesiredQuotas
		state.Requests.Current, state.DesiredQuotas = req.ID, quotas
		state.Quotas = quotas
		if err := cr.updateContainerCPUQuota(state.Container, quotas); err != nil {
			log.Println(err.Error())
		}
		cr.ReconfigureRegistry(quotas, quotas_old)
	} else {
		state.Requests.Active[req.ID] = quotas
	}
	cr.mutex.Unlock()
	return true, nil
}

/*
	Removes a request from registry
	and resets function state.
*/
func (cr *ConflictResolver) RemoveFromRegistry(rs wrq.ResetRequest) error {
	cr.mutex.Lock()
	state, ok := cr.Registry[rs.Function]
	if !ok {
		cr.mutex.Unlock()
		return fmt.Errorf("request for '%v' function not found", rs.Function)
	}
	if state.Requests.Current == rs.ID {
		if len(state.Requests.Active) == 0 {
			// TO DO: Solve Openwhisk autoscaling problem
			if err := cr.updateContainerCPUQuota(state.Container, -1); err != nil {
				log.Println(err.Error())
			}
			delete(cr.Registry, rs.Function)
			if len(cr.Registry) != 0 {
				cr.ReconfigureRegistry(0, state.DesiredQuotas)
			} else {
				cr.lambdaPrevious = 0
			}
		} else {
			quotas_old := state.DesiredQuotas
			state.Requests.Current, state.DesiredQuotas = nextRequest(state)
			state.Quotas = state.DesiredQuotas
			delete(state.Requests.Active, state.Requests.Current)
			if err := cr.updateContainerCPUQuota(state.Container, state.Quotas); err != nil {
				log.Println(err.Error())
			}
			cr.ReconfigureRegistry(state.DesiredQuotas, quotas_old)
		}
	} else {
		if _, ok := state.Requests.Active[rs.ID]; !ok {
			cr.mutex.Unlock()
			return fmt.Errorf("request %v not found", rs.ID)
		}
		delete(state.Requests.Active, rs.ID)
	}
	cr.mutex.Unlock()
	return nil
}

/*
	Finds appropriate function (maximum DesiredCPUQuotas)
	request to enable.
*/
func nextRequest(state *wfs.FunctionState) (uint64, int64) {
	var (
		maxID uint64
		maxq  int64
	)
	for id, quotas := range state.Requests.Active {
		if quotas > maxq || quotas == maxq && id < maxID {
			maxID = id
			maxq = quotas
		}
	}
	return maxID, maxq
}

/*
	Recompute CPU quotas for each container,
	using the formula λ*DesiredCPUQuotas
	where λ = CPU_Cores / sum of DesiredCPUQuotas.
*/
func (cr *ConflictResolver) ReconfigureRegistry(quotas_new int64, quotas_old int64) {
	var (
		sum    int64
		lambda float64
	)
	if cr.lambdaPrevious == 0 {
		for _, s := range cr.Registry {
			sum += s.DesiredQuotas
		}
		lambda = float64(cr.Cores*CPU_PERIOD_OPENWHISK_DEFAULT) / float64(sum)
	} else {
		nt := float64(cr.Cores * CPU_PERIOD_OPENWHISK_DEFAULT)
		lambda = nt / (nt/cr.lambdaPrevious + float64(quotas_new-quotas_old))
	}

	if !((cr.lambdaPrevious >= 1) && (lambda >= 1)) {
		for f, s := range cr.Registry {
			if lambda < 1.0 {
				s.Quotas = lowerBound(int64(lambda * float64(s.DesiredQuotas)))
			} else {
				s.Quotas = s.DesiredQuotas
			}
			if err := cr.updateContainerCPUQuota(cr.Registry[f].Container, s.Quotas); err != nil {
				log.Println(err.Error())
			}
		}
	}
	cr.lambdaPrevious = lambda
}

/*
	Helper method for searching docker runtime
	for an openwhisk action container.
*/
func (cr *ConflictResolver) SearchDockerRuntime(function, podType string) (*types.Container, error) {
	containers, err := cr.DockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	exp := fmt.Sprintf(REG_EXP, function)
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

/*
	Export Registry information.
	Method used for debugging purposes.
*/
func (cr *ConflictResolver) ExportRegistry() map[string]interface{} {
	res := make(map[string]interface{})
	cr.mutex.RLock()
	for k, s := range cr.Registry {
		res[k] = *s
	}
	cr.mutex.RUnlock()
	return res
}

/*
	Helper method for updating CPU quotas
	of specified docker container.
*/
func (cr *ConflictResolver) updateContainerCPUQuota(containerID string, cpuQuota int64) error {
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
	if _, err := cr.DockerClient.ContainerUpdate(context.Background(), containerID, updateConfig); err != nil {
		return err
	}
	return nil
}

/*
	PID controller function.
	Input: slack in nanoseconds
	Output: Δcpu_quotas in miliseconds
*/
func computePIDControllerOutput(req *wrq.Request) int64 {
	return -1 * mseconds(Kp*req.Metrics.Slack+
		Ki*req.Metrics.SumOfSlack+
		Kd*req.Metrics.PreviousSlack)
}

/*
	Apply upper and lower bound to PID output.
*/
func retainCPUThreshold(quotas, cores int64) int64 {
	threshold := CPU_PERIOD_OPENWHISK_DEFAULT * cores
	if quotas > threshold {
		return threshold
	} else if quotas <= 5*CPU_QUOTAS_LOWER_BOUND {
		return 5 * CPU_QUOTAS_LOWER_BOUND
	} else {
		return quotas
	}
}

/*
	Convert nanoseconds into milliseconds.
*/
func mseconds(x int64) int64 {
	return int64(math.Round(float64(x) * 0.000001))
}

/*
	Apply lower bound to lambda correction.
	Used by ReconfigureRegistry
*/
func lowerBound(q int64) int64 {
	if q <= CPU_QUOTAS_LOWER_BOUND {
		return CPU_QUOTAS_LOWER_BOUND
	} else {
		return q
	}
}

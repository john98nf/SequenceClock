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
	"regexp"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	wfs "github.com/john98nf/SequenceClock/watcher/internal/state"
	wrq "github.com/john98nf/SequenceClock/watcher/pkg/request"
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
	CORES                        int    = 4
)

type ConflictResolverInterface interface {
	SearchRegistry(function, podType string) (*types.Container, error)
	UpdateRegistry(id uint64, function string, quotas int64) error
	RemoveFromRegistry(rs wrq.ResetRequest) error
	ReconfigureRegistry() error
}

type ConflictResolver struct {
	Registry     map[string]*wfs.FunctionState
	DockerClient *client.Client
}

func NewConflictResolver() ConflictResolver {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	registry := make(map[string]*wfs.FunctionState)
	return ConflictResolver{
		Registry:     registry,
		DockerClient: cli,
	}
}

/*
	Places new request into registry and
	update container resources.
*/
func (cr *ConflictResolver) UpdateRegistry(id uint64, function string, quotas int64) error {
	state := cr.Registry[function]
	if quotas > state.Quotas {
		state.Requests.Active[state.Requests.Current] = state.DesiredQuotas
		state.Requests.Current, state.DesiredQuotas = id, quotas
		cr.ReconfigureRegistry()
	} else {
		state.Requests.Active[id] = quotas
	}
	return nil
}

/*
	Removes a request from registry
	and resets function state.
*/
func (cr *ConflictResolver) RemoveFromRegistry(rs wrq.ResetRequest) error {
	state, ok := cr.Registry[rs.Function]
	if !ok {
		return fmt.Errorf("request for '%v' function not found", rs.Function)
	}
	if state.Requests.Current == rs.ID {
		if len(state.Requests.Active) == 0 {
			// TO DO: Solve Openwhisk autoscaling problem
			if err := cr.updateContainerCPUQuota(state.Container, -1); err != nil {
				return err
			}
			delete(cr.Registry, rs.Function)
		} else {
			state.Requests.Current, state.DesiredQuotas = nextRequest(state)
		}
		cr.ReconfigureRegistry()
	} else {
		if _, ok := state.Requests.Active[rs.ID]; !ok {
			return fmt.Errorf("request %v not found", rs.ID)
		}
		delete(state.Requests.Active, rs.ID)
	}
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
func (cr *ConflictResolver) ReconfigureRegistry() error {
	var sum int64
	for _, s := range cr.Registry {
		sum += s.DesiredQuotas
	}
	lamda := int64(float64(CORES) / float64(sum) * 100)
	for f, s := range cr.Registry {
		s.Quotas = lamda * s.DesiredQuotas
		if err := cr.updateContainerCPUQuota(cr.Registry[f].Container, s.Quotas); err != nil {
			log.Println(err.Error())
		}
	}
	return nil
}

/*
	Helper method for searching docker runtime
	for an openwhisk action container.
*/
func (cr *ConflictResolver) SearchRegistry(function, podType string) (*types.Container, error) {
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

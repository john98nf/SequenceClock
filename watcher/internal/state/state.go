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

package state

/*
	Struct for tracking function container/s
	state from inside watcher.
*/
type FunctionState struct {
	Container     string // TO DO: Solve Openwhisk autoscaling problem
	Quotas        int64
	DesiredQuotas int64
	Requests      RequestsInfo
}

type RequestsInfo struct {
	Current uint64
	Active  map[uint64]int64
}

func NewFunctionState(container string) *FunctionState {
	return &FunctionState{
		Container: container,
		Requests: RequestsInfo{
			Current: 0,
			Active:  map[uint64]int64{},
		},
	}
}

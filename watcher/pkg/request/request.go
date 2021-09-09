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

/*
	Initial Request struct made by sequence controller
	and passed to watcher supreme.
*/
type Request struct {
	ID       uint64   `form:"id" schema:"-"`
	Function string   `form:"function" binding:"required" schema:"function"`
	Metrics  *Metrics `form:"metrics" binding:"required" schema:"metrics"`
}

/*
	ResetRequest carries the id given to sequence controller
	by the watcher supreme.
*/
type ResetRequest struct {
	ID       uint64 `form:"id" binding:"exists" schema:"id"`
	Function string `form:"function" binding:"required" schema:"function"`
}

/*
	Struct send as part of sequence controller Requests
	to watchers.
*/
type Metrics struct {
	Slack                 int64 `form:"slack" binding:"exists" schema:"slack"`                 // Used by P controller
	SumOfSlack            int64 `form:"sumOfSlack" binding:"exists" schema:"sumOfSlack"`       // Used by I controller
	PreviousSlack         int64 `form:"previousSlack" binding:"exists" schema:"previousSlack"` // Used by D controller
	ProfiledExecutionTime int64 `form:"profiledExecutionTime" binding:"exists" schema:"profiledExecutionTime"`
}

/*
	Returns a new Reset Request.
*/
func NewResetRequest(id uint64, function string) *ResetRequest {
	return &ResetRequest{
		ID:       id,
		Function: function,
	}
}

/*
	Returns a new Request.
*/
func NewRequest(f string, m *Metrics) *Request {
	return &Request{
		Function: f,
		Metrics:  m,
	}
}

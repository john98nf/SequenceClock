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

module github.com/john98nf/SequenceClock/watcher

go 1.15

replace github.com/john98nf/SequenceClock/watcher/internal/conflicts => ./internal/conflicts

replace github.com/john98nf/SequenceClock/watcher/pkg/request => ./pkg/request

replace github.com/john98nf/SequenceClock/watcher/internal/state => ./internal/state

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/john98nf/SequenceClock/watcher/internal/conflicts v0.0.0-00010101000000-000000000000
	github.com/john98nf/SequenceClock/watcher/internal/state v0.0.0-00010101000000-000000000000
	github.com/john98nf/SequenceClock/watcher/pkg/request v0.0.0-20210820205221-369ee2bc9c4d
	github.com/morikuni/aec v1.0.0 // indirect
)

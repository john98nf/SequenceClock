// Copyright © 2021 Giannis Fakinos
//
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

package controllerTemplates

import (
	"fmt"
	"strings"
)

func ImportsTemplate() []byte {
	return []byte(`package main

import (
	"fmt"
	"log"
)

`)
}

func FunctionListTemplate(list []string) []byte {
	return []byte(fmt.Sprintf("var functionList = [...]string{%v}\n\n", strings.Join(list, ",")))
}

func MainTemplate() []byte {
	return []byte(`func Main(obj map[string]interface{}) map[string]interface{} {
	name, ok := obj["name"].(string)
	if !ok {
	name = "world"
	}
	
	msg := make(map[string]interface{})
	msg["message"] = "Hello, " + name + "!"
	log.Printf("name=%s\n", name)
	
	return msg
}
`)
}

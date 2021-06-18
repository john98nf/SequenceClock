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

/*
	This file is not include in build process
	and will be copied by 'sc create' command.
	Please, do not edit!!
*/

package main

import (
	"log"
	"net/http"

	"github.com/apache/openwhisk-client-go/whisk"
)

var functionList = [...]string{}

func Main(obj map[string]interface{}) map[string]interface{} {
	wskConfig := &whisk.Config{
		Host:      "",
		Namespace: "",
		AuthToken: "",
		Insecure:  true,
	}
	client, _ := whisk.NewClient(http.DefaultClient, wskConfig)

	log.Printf("Invoking function-%v %v\n", 0, functionList[0])
	aRes, _, _ := client.Actions.Invoke(functionList[0], obj, true, true)
	for i, f := range functionList[1:] {
		log.Printf("Invoking function-%v %v\n", i+1, f)
		aRes, _, _ = client.Actions.Invoke(f, aRes, true, true)
	}

	return aRes
}

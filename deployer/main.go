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
	"fmt"

	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	deployerAPI := app.Party("/api")
	{
		deployerAPI.Use(iris.Compression)

		// GET: http://localhost:8080/api/check
		deployerAPI.Get("/check", check)
		// POST: http://localhost:8080/api/create?name=x
		deployerAPI.Post("/create", create)
		// DELETE: http://localhost:8080/api/delete?name=x
		deployerAPI.Delete("/delete", delete)
	}

	app.Listen(":42000")
}

func check(ctx iris.Context) {
	ctx.Text("SC-Deployer is fully functional!")
}

func create(ctx iris.Context) {
	sequence := ctx.URLParam("name")
	ctx.WriteString(fmt.Sprintf("Create request for %v", sequence))
}

func delete(ctx iris.Context) {
	sequence := ctx.URLParam("name")
	ctx.WriteString(fmt.Sprintf("Delete request for %v", sequence))
}

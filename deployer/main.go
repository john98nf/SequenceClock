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
	"log"
	"net/http"
	"os"

	tpl "github.com/john98nf/SequenceClock/deployer/internal/templateHandler"
	"github.com/john98nf/SequenceClock/deployer/pkg/sequence"

	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	deployerAPI := router.Group("/api")
	{
		// GET: http://localhost:8080/api/check
		deployerAPI.GET("/check", check)
		// POST: http://localhost:8080/api/create?name=x
		deployerAPI.POST("/create", create)
		// DELETE: http://localhost:8080/api/delete?name=x
		deployerAPI.DELETE("/delete", delete)
	}

	router.Run(":42000")
}

/*
	Liveness & Readiness call.
*/
func check(c *gin.Context) {
	c.String(http.StatusOK, "SC-Deployer is fully functional!")
}

/*
	API call for creating a new sequence
	and deploying it to cluster.
*/
func create(c *gin.Context) {
	var seq sequence.Sequence
	if err := c.ShouldBind(&seq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := seq.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wskConfig := &whisk.Config{
		Host:      os.Getenv("API_HOST"),
		Namespace: os.Getenv("NAMESPACE"),
		AuthToken: os.Getenv("OPENWHISK_AUTH_TOKEN"),
		Insecure:  true,
	}
	client, errCl := whisk.NewClient(http.DefaultClient, wskConfig)
	if errCl != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errCl.Error()})
		return
	}

	template := tpl.NewTemplate(&seq, client)
	if err := template.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := template.Deploy(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := template.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("sequence '%v' created.", seq.Name)})
}

/*
	API call for deleting a sequence
*/
func delete(c *gin.Context) {
	sequence := c.Query("name")
	if sequence == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no sequence name provided"})
		return
	}

	wskConfig := &whisk.Config{
		Host:      os.Getenv("API_HOST"),
		Namespace: os.Getenv("NAMESPACE"),
		AuthToken: os.Getenv("OPENWHISK_AUTH_TOKEN"),
		Insecure:  true,
	}
	client, errCl := whisk.NewClient(http.DefaultClient, wskConfig)
	if errCl != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errCl.Error()})
		return
	}
	actions, _, errL := client.Actions.List("", nil)
	if errL != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errL.Error()})
		return
	}

	if func(a []whisk.Action, s string) bool {
		for _, a := range a {
			if s == a.Name {
				return true
			}
		}
		return false
	}(actions, sequence) {
		if _, err := client.Actions.Delete(sequence); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("sequence '%v' deleted", sequence)})
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprint("no '%v' sequence detected", sequence)})
	}
}

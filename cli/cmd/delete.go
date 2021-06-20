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

package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/spf13/cobra"
)

func init() {
	delete.Flags().StringVarP(&deletionSequence, "name", "n", "", "Sequence name (required)")
	delete.MarkFlagRequired("name")
}

var (
	deletionSequence string
	delete           = &cobra.Command{
		Use:   "delete",
		Short: "Delete an existing function sequence.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			wskConfig := &whisk.Config{
				Host:      apihost,
				Namespace: namespace,
				AuthToken: authToken,
				Insecure:  true,
			}
			client, _ := whisk.NewClient(http.DefaultClient, wskConfig)
			actions, _, _ := client.Actions.List("", nil)
			if contains(actions, deletionSequence) {
				_, err := client.Actions.Delete(deletionSequence)
				if err != nil {
					log.Fatal(err)
				} else {
					fmt.Printf("Sequence '%v' deleted.\n", deletionSequence)
				}
			} else {
				fmt.Printf("No '%v' sequence detected.\n", deletionSequence)
			}
		},
	}
)

func contains(actionSlice []whisk.Action, s string) bool {
	for _, a := range actionSlice {
		if s == a.Name {
			return true
		}
	}
	return false
}

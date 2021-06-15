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
	"SequenceClock/controllerTemplates"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	create.Flags().StringVarP(&sequenceName, "name", "n", "", "Sequence name (required)")
	create.MarkFlagRequired("name")
	create.Flags().StringSliceVarP(&sequenceList, "list", "l", []string{}, "List of functions in the sequence (required)")
	create.MarkFlagRequired("list")
}

var (
	sequenceName string
	sequenceList []string
	create       = &cobra.Command{
		Use:   "create",
		Short: "Create a new function sequence.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TO DO: implementation
			if len(sequenceList) <= 1 {
				return fmt.Errorf("please provide at least 2 functions for the creation.")
			}

			homedir, err1 := os.UserHomeDir()
			if err1 != nil {
				return err1
			}
			appDirectory := homedir + "/sequenceClock"
			_, err2 := os.Stat(appDirectory)
			if os.IsNotExist(err2) {
				err3 := os.Mkdir(appDirectory, 0755)
				if err3 != nil {
					return err3
				}
			}
			f, err4 := os.Create(fmt.Sprintf("%v/%v.go", appDirectory, sequenceName))
			if err4 != nil {
				return err4
			}
			defer f.Close()

			_, err := f.Write(controllerTemplates.ImportsTemplate())
			if err != nil {
				return err
			}
			_, err = f.Write(controllerTemplates.FunctionListTemplate(sequenceList))
			if err != nil {
				return err
			}
			_, err = f.Write(controllerTemplates.MainTemplate())
			if err != nil {
				return err
			}

			fmt.Printf("New sequence '%v' generated.\n", sequenceName)
			fmt.Println("Pipeline:", strings.Join(sequenceList, ", "))
			return nil
		},
	}
)

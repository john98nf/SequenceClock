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
	"os"
	"strings"
)

/*
	Add proper var definitions to
	main.go of template.
*/
func FunctionListTemplate(
	fileName string,
	functionList []string,
	apihost,
	namespace,
	authToken string) error {

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	varDefinitions := fmt.Sprintf(`var (
		apihost string = "%v"
		namespace string = "%v"
		authToken string = "%v"
		var functionList = [...]string{%v}
)`, apihost, namespace, authToken, strings.Join(functionList, "\",\""))

	if _, err = f.WriteString(varDefinitions); err != nil {
		return err
	}
	return nil
}

/*
	Copy given template to appFolder.
*/
func CreateTemplate(sequenceName, appFolder string) error {
	// To Do get code from github directly
	// Temporary solution get template from source code
	if errDir := os.Mkdir(appFolder+"/"+sequenceName, 0755); errDir != nil {
		if !os.IsExist(errDir) {
			return errDir
		}
	}
	return copyTemplate("./internal/controllerTemplates/wskTemplate", appFolder+"/"+sequenceName)
}

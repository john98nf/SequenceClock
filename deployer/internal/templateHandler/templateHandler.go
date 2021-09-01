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

package templateHandler

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/apache/openwhisk-client-go/whisk"

	sq "github.com/john98nf/SequenceClock/deployer/pkg/sequence"
)

const (
	OPENWHISK_CONTROLLER_TEMPLATE string = "/tmp/controller/openwhisk/"
	OPENFAAS_CONTROLLER_TEMPLATE  string = "/tmp/controller/openfaas/"
	GO_RUNTIME                    string = "go:1.15"
)

type TemplateInterface interface {
	Create() error
	Deploy() error
	Delete() error
}

type Template struct {
	Sequence *sq.Sequence
	Client   *whisk.Client
	Location string
}

/*
	Creates a new Template struct.
*/
func NewTemplate(sequence *sq.Sequence, client *whisk.Client) *Template {
	return &Template{
		Sequence: sequence,
		Client:   client,
		Location: "",
	}
}

/*
	Copies controller template and
	creates a zip folder <sequenceName>.zip.
*/
func (tpl *Template) Create() error {
	execPath, errP := execPath()
	if errP != nil {
		log.Println(errP)
		return fmt.Errorf("couldn't found executable path")
	}
	fziper := NewFileZiper(execPath, OPENWHISK_CONTROLLER_TEMPLATE)
	zipFile, errZ := fziper.zipTemplate(*tpl.Sequence)
	if errZ != nil {
		log.Println(errZ)
		return fmt.Errorf("couldn't create zip archive")
	}

	tpl.Location = zipFile
	return nil
}

/*
	Uses zip archive created from Create() method
	and deployes it to openwhisk client.
*/
func (tpl *Template) Deploy() error {
	newAction := whisk.Action{
		Name:        tpl.Sequence.Name,
		Namespace:   os.Getenv("NAMESPACE"),
		Annotations: whisk.KeyValueArr{whisk.KeyValue{Key: "provide-api-key", Value: "true"}},
	}
	newAction.Exec = new(whisk.Exec)
	newAction.Exec.Kind = GO_RUNTIME
	zipCnt, err := ioutil.ReadFile(tpl.Location)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("couldn't read zip file")
	}
	code := base64.StdEncoding.EncodeToString(zipCnt)
	newAction.Exec.Code = &code

	if _, _, errI := tpl.Client.Actions.Insert(&newAction, true); errI != nil {
		log.Println(errI)
		return fmt.Errorf("couldn't deploy new sequence")
	}
	return nil
}

/*
	Deletes template.
*/
func (tpl *Template) Delete() error {
	if tpl.Location == "" {
		return fmt.Errorf("deletion of non existing template")
	}
	if err := os.Remove(tpl.Location); err != nil {
		log.Println(err)
		return fmt.Errorf("couldn't delete zip template")
	}
	return nil
}

/*
	Mini function that finds execution path.
*/
func execPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

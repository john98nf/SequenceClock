package templateHandler

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/apache/openwhisk-client-go/whisk"

	sq "john98nf/SequenceClock/deployer/internal/sequence"
)

const (
	OPENWHISK_CONTROLLER_TEMPLATE string = "controller/openwhisk/"
	OPENFAAS_CONTROLLER_TEMPLATE  string = "controller/openfaas/"
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
	fziper := &fileZiper{
		dstFolder:  execPath,
		baseFolder: OPENWHISK_CONTROLLER_TEMPLATE,
	}
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

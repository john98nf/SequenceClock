package templateHandler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/apache/openwhisk-client-go/whisk"

	sq "john98nf/SequenceClock/deployer/internal/sequence"
)

const (
	OPENWHISK_CONTROLLER_TEMPLATE string = "controller/openwhisk/"
	OPENFAAS_CONTROLLER_TEMPLATE  string = "controller/openfaas/"
)

type TemplateInterface interface {
	CreateBase() error
}

type Template struct {
	Sequence *sq.Sequence
	Client   *whisk.Client
}

/*
	Creates a new Template struct.
*/
func NewTemplate(sequence *sq.Sequence, client *whisk.Client) *Template {
	return &Template{
		Sequence: sequence,
		Client:   client,
	}
}

/*
	Copies controller template and
	creates a zip folder <sequenceName>.zip.
*/
func (tpl *Template) CreateBase() error {
	execPath, errP := execPath()
	if errP != nil {
		return fmt.Errorf("couldn't found executable path")
	}
	fziper := &fileZiper{
		dstFolder:  execPath,
		baseFolder: OPENWHISK_CONTROLLER_TEMPLATE,
		name:       tpl.Sequence.Name,
	}
	zipFile, errZ := fziper.zipTemplate(*tpl.Sequence)
	if errZ != nil {
		return errZ
	} else {
		log.Println(zipFile)
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

package templateHandler

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/apache/openwhisk-client-go/whisk"

	sq "john98nf/SequenceClock/deployer/internal/sequence"
)

const (
	MAIN_HANDLER string = "main.go"
	CONSTANTS    string = `var (
	apihost string = "%v"
	namespace string = "%v"
	authToken string = "%v"
	functionList = [...]string{%v}
)
`
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
	baseFolder := execPath + "/controller/" + tpl.Sequence.Framework + "/"
	outFile, err := os.Create(fmt.Sprintf("%v/%v.zip", execPath, tpl.Sequence.Name))
	if err != nil {
		return fmt.Errorf("couldn't create zip archive for controller")
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	if errZ := addFiles(w, baseFolder, ""); errZ != nil {
		return fmt.Errorf("couldn't add files to archive")
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("couldn't close zip writer")
	}
	return nil
}

/*
	Recursive function that reads the contents of basePath
	and moves them to a zip folder, provided by w *zip.Writer.
*/
func addFiles(w *zip.Writer, basePath, baseInZip string) error {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, errR := ioutil.ReadFile(basePath + file.Name())
			if errR != nil {
				log.Println(errR)
				return errR
			}
			f, errF := w.Create(baseInZip + file.Name())
			if errF != nil {
				log.Println(errF)
				return errF
			}
			if file.Name() == MAIN_HANDLER {
				constDef := []byte(fmt.Sprintf(CONSTANTS, "0.0.0.0:31001", "_", "1234abcd", "\""+strings.Join([]string{"f1", "f2"}, "\",\"")+"\""))
				dat = append(dat, constDef...)
			}
			_, errW := f.Write(dat)
			if errW != nil {
				log.Println(errW)
				return errW
			}
		} else if file.IsDir() {
			newBase := basePath + file.Name() + "/"
			return addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
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

package templateHandler

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	sq "john98nf/SequenceClock/deployer/internal/sequence"
	"os"
	"strings"
)

const (
	MAIN_HANDLER  string = "main.go"
	FUNCTION_LIST string = "var functionList = [...]string{\"%v\"}"
)

type fileZiperInterface interface {
	zipTemplate(seq sq.Sequence) (string, error)
}

type fileZiper struct {
	dstFolder  string
	baseFolder string
}

func (obj *fileZiper) zipTemplate(seq sq.Sequence) (string, error) {
	zipFile := fmt.Sprintf("%v/%v.zip", obj.dstFolder, seq.Name)
	outFile, err := os.Create(zipFile)
	if err != nil {
		return "", fmt.Errorf("couldn't create zip archive")
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	if errZ := addFiles(w, obj.baseFolder, "", seq.Functions); errZ != nil {
		return "", fmt.Errorf("couldn't add files to archive")
	}

	if err = w.Close(); err != nil {
		return "", fmt.Errorf("couldn't close zip writer")
	}
	return zipFile, nil
}

/*
	Recursive function that reads the contents of basePath
	and moves them to a zip folder, provided by w *zip.Writer.
*/
func addFiles(w *zip.Writer, basePath, baseInZip string, functionList []string) error {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, errR := ioutil.ReadFile(basePath + file.Name())
			if errR != nil {
				return errR
			}
			f, errF := w.Create(baseInZip + file.Name())
			if errF != nil {
				return errF
			}
			if file.Name() == MAIN_HANDLER {
				varsDef := []byte(fmt.Sprintf(FUNCTION_LIST, strings.Join(functionList, "\",\"")))
				dat = append(dat, varsDef...)
			}
			_, errW := f.Write(dat)
			if errW != nil {
				return errW
			}
		} else if file.IsDir() {
			newBase := basePath + file.Name() + "/"
			if errRec := addFiles(w, newBase, baseInZip+file.Name()+"/", functionList); errRec != nil {
				return errRec
			}
		}
	}
	return nil
}

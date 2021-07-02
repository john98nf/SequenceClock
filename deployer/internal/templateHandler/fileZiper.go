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
	PACKAGE_DEFINITION string = `package main
`
	CONFIG_CONTROLLER_FILE string = "config.go"
	ZIP_ARCHIVE_PATH       string = "%v/%v.zip"
	FUNCTION_SLICE         string = "var functionList = [...]string{\"%v\"}"
)

type fileZiperInterface interface {
	zipTemplate(seq sq.Sequence) (string, error)
}

type fileZiper struct {
	dstFolder  string
	baseFolder string
}

func (obj *fileZiper) zipTemplate(seq sq.Sequence) (string, error) {
	zipFile := fmt.Sprintf(ZIP_ARCHIVE_PATH, obj.dstFolder, seq.Name)
	outFile, err := os.Create(zipFile)
	if err != nil {
		return "", fmt.Errorf("couldn't create zip archive")
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	if errZ := addFiles(w, obj.baseFolder, ""); errZ != nil {
		return "", fmt.Errorf("couldn't add files to archive")
	}
	if errC := addConfig(w, seq); errC != nil {
		return "", fmt.Errorf("couldn't add config file to archive")
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
func addFiles(w *zip.Writer, basePath, baseInZip string) error {
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
			_, errW := f.Write(dat)
			if errW != nil {
				return errW
			}
		} else if file.IsDir() {
			newBase := basePath + file.Name() + "/"
			if errRec := addFiles(w, newBase, baseInZip+file.Name()+"/"); errRec != nil {
				return errRec
			}
		}
	}
	return nil
}

/*
	Add config variables and constants
	to zip archive.
*/
func addConfig(w *zip.Writer, seq sq.Sequence) error {
	f, errF := w.Create(CONFIG_CONTROLLER_FILE)
	if errF != nil {
		return errF
	}

	dat := []byte(PACKAGE_DEFINITION + fmt.Sprintf(FUNCTION_SLICE, strings.Join(seq.Functions, "\",\"")))

	_, errW := f.Write(dat)
	if errW != nil {
		return errW
	}
	return nil
}

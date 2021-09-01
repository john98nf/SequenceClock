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
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	sq "github.com/john98nf/SequenceClock/deployer/pkg/sequence"
)

const (
	PACKAGE_DEFINITION string = `package main
	
	import "time"

`
	CONFIG_CONTROLLER_FILE string = "config.go"
	ZIP_ARCHIVE_PATH       string = "%v/%v.zip"
	CONSTANTS              string = `const (
		ALGORITHM_TYPE string = "%v"
		KUBE_MAIN_IP string = "%v"
)
`
	VARIABLES string = `var (
		functionList = [...]string{"%v"}
		profiledExecutionTimes = [...]time.Duration{%v}
)
`
)

type fileZiperInterface interface {
	zipTemplate(seq sq.Sequence) (string, error)
}

type fileZiper struct {
	dstFolder  string
	baseFolder string
}

func NewFileZiper(dstFolder, baseFolder string) *fileZiper {
	return &fileZiper{
		dstFolder:  dstFolder,
		baseFolder: baseFolder,
	}
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

	dat := []byte(PACKAGE_DEFINITION +
		fmt.Sprintf(CONSTANTS, seq.AlgorithmType, os.Getenv("HOST_IP")) +
		fmt.Sprintf(VARIABLES,
			strings.Join(seq.Functions, "\",\""),
			strings.Join(stringify(seq.ProfiledExecutionTimes), ",")))

	_, errW := f.Write(dat)
	if errW != nil {
		return errW
	}
	return nil
}

func stringify(l []time.Duration) []string {
	s := make([]string, len(l))
	for i, elem := range l {
		s[i] = fmt.Sprint(elem.Nanoseconds())
	}
	return s
}

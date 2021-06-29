package templateHandler

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	sq "john98nf/SequenceClock/deployer/internal/sequence"
	"log"
	"os"
	"strings"
)

const (
	MAIN_HANDLER string = "main.go"
	VARIABLES    string = `var (
	apihost string = "%v"
	namespace string = "%v"
	authToken string = "%v"
	functionList = [...]string{%v}
)
`
)

type fileZiper struct {
	dstFolder  string
	baseFolder string
	name       string
}

func (obj *fileZiper) zipTemplate(seq sq.Sequence) (string, error) {
	zipFile := fmt.Sprintf("%v/%v.zip", obj.dstFolder, obj.name)
	outFile, err := os.Create(zipFile)
	if err != nil {
		return "", fmt.Errorf("couldn't create zip archive")
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	if errZ := addFiles(w, obj.baseFolder, ""); errZ != nil {
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
				varsDef := []byte(fmt.Sprintf(VARIABLES, "0.0.0.0:31001", "_", "1234abcd", "\""+strings.Join([]string{"f1", "f2"}, "\",\"")+"\""))
				dat = append(dat, varsDef...)
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

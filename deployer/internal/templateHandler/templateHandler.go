package templateHandler

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type TemplateInterface interface {
	CreateBase() error
	// AddServerlessFunction() error
}

type Template struct {
	// Sequence *sq.Sequence
	// Client   *whisk.Client
}

/*
	Creates a new Template struct.
*/
func NewTemplate() *Template {
	return &Template{}
}

/*
	Copies controller template and
	creates a zip folder <sequenceName>.zip.
*/
func (tpl *Template) CreateBase(name string) error {
	execPath, errP := execPath()
	if errP != nil {
		return fmt.Errorf("couldn't found executable path")
	}
	baseFolder := execPath + "/controller"
	outFile, err := os.Create(fmt.Sprintf("%v/%v.zip", execPath, name))
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

/*
	Add proper var definitions to
	main.go of template.
*/
// func FunctionListTemplate(
// 	fileName string,
// 	functionList []string,
// 	apihost,
// 	namespace,
// 	authToken string) error {

// 	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	varDefinitions := fmt.Sprintf(`var (
// 		apihost string = "%v"
// 		namespace string = "%v"
// 		authToken string = "%v"
// 		var functionList = [...]string{%v}
// )`, apihost, namespace, authToken, strings.Join(functionList, "\",\""))

// 	if _, err = f.WriteString(varDefinitions); err != nil {
// 		return err
// 	}
// 	return nil
// }

/*
	Copy given template to appFolder.
*/
// func CreateTemplate(sequenceName, appFolder string) error {
// 	// To Do get code from github directly
// 	// Temporary solution get template from source code
// 	if errDir := os.Mkdir(appFolder+"/"+sequenceName, 0755); errDir != nil {
// 		if !os.IsExist(errDir) {
// 			return errDir
// 		}
// 	}
// 	return copyTemplate("./internal/controllerTemplates/wskTemplate", appFolder+"/"+sequenceName)
// }

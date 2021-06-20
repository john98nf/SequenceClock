package controllerTemplates

import (
	"io"
	"io/ioutil"
	"os"
)

/*
	Move recursively contents of src
	folder to dst folder, provided dst exists.
*/
func copyTemplate(src, dst string) error {
	files, errDir := ioutil.ReadDir(src)
	if errDir != nil {
		return errDir
	}
	for _, file := range files {
		if !file.IsDir() {
			errFile := copyFile(file.Name(), src, dst)
			if errFile != nil {
				return errFile
			}
		} else if file.IsDir() {
			newFolder := dst + "/" + file.Name()
			if errMakeDir := os.Mkdir(newFolder, 0755); errMakeDir != nil {
				if os.IsExist(errMakeDir) {
					continue
				} else {
					return errMakeDir
				}
			}
			if err := copyTemplate(src+"/"+file.Name(), newFolder); err != nil {
				return err
			}
		}
	}
	return nil
}

/*
	Copy file from src to dst folder.
	Both src and dst must represent
	actual paths.
*/
func copyFile(file, src, dst string) error {
	source, errO := os.Open(src + "/" + file)
	if errO != nil {
		return errO
	}
	defer source.Close()

	destination, errC := os.Create(dst + "/" + file)
	if errC != nil {
		return errC
	}
	defer destination.Close()
	_, errCopy := io.Copy(destination, source)
	return errCopy
}

// func ZipWriter() {
// 	baseFolder := "./test/"

// 	// Get a Buffer to Write To
// 	outFile, err := os.Create(`test.zip`)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer outFile.Close()

// 	// Create a new zip archive.
// 	w := zip.NewWriter(outFile)

// 	// Add some files to the archive.
// 	addFiles(w, baseFolder, "")

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Make sure to check the error on Close.
// 	err = w.Close()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func addFiles(w *zip.Writer, basePath, baseInZip string) {
// 	fmt.Println("Call of add files for", basePath, baseInZip)
// 	// Open the Directory
// 	files, err := ioutil.ReadDir(basePath)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	for _, file := range files {
// 		fmt.Println(basePath + file.Name())
// 		if !file.IsDir() {
// 			dat, err := ioutil.ReadFile(basePath + file.Name())
// 			if err != nil {
// 				fmt.Println(err)
// 			}

// 			// Add some files to the archive.
// 			f, err := w.Create(baseInZip + file.Name())
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			_, err = f.Write(dat)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		} else if file.IsDir() {

// 			// Recurse
// 			newBase := basePath + file.Name() + "/"
// 			fmt.Println("Recursing and Adding SubDir: " + file.Name())
// 			fmt.Println("Recursing and Adding SubDir: " + newBase)

// 			addFiles(w, newBase, baseInZip+file.Name()+"/")
// 		}
// 	}
// }

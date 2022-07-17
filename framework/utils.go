package framework

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

func GetFilesAndDirRecursively(path string) (htmlFiles []string) {
	dirFiles, err := ioutil.ReadDir(path)

	if err != nil {
		log.Fatalln("error reading directory: " + err.Error())
	}

	for _, dirFile := range dirFiles {
		fullPath := filepath.Join(path, dirFile.Name())

		if dirFile.IsDir() {
			htmlFiles = append(htmlFiles, GetFilesAndDirRecursively(fullPath)...)
		} else if dirFile.Mode().IsRegular() {
			// do something with the file
			ext := filepath.Ext(dirFile.Name())
			if ext == ".gohtml" {
				htmlFiles = append(htmlFiles, fullPath)
			}
		}
	}
	return
}

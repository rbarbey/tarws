package main

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"

	"tarws/cmd"
)

func main() {
	cmd.Execute()

	//err := filepath.Walk("/Users/robert/Development/golang/src/", tar)
	// tarWriter := tar.NewWriter(ioutil.Discard)
	// defer tarWriter.Close()
	//
	// iterate("/Users/robert/Development/golang/src/", tarWriter)
}

func iterate(path string, tarWriter *tar.Writer) {
	dir, err := os.Open(path)
	handle(err)
	defer dir.Close()

	// read all file entries in dir in one slice
	fileInfos, err := dir.Readdir(0)
	handle(err)

	for _, fileInfo := range fileInfos {
		currentPath := filepath.Join(path, fileInfo.Name())
		if fileInfo.IsDir() {
			iterate(currentPath, tarWriter)
		} else {
			log.Printf("Adding %s\n", currentPath)
			write(currentPath, fileInfo, tarWriter)
		}
	}
}

func write(path string, fileInfo os.FileInfo, tarWriter *tar.Writer) {
	err := tarWriter.WriteHeader(&tar.Header{
		Name:    path,
		Size:    fileInfo.Size(),
		Mode:    int64(fileInfo.Mode()),
		ModTime: fileInfo.ModTime(),
	})
	handle(err)

	file, err := os.Open(path)
	handle(err)
	defer file.Close()

	_, err = io.Copy(tarWriter, file)
	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

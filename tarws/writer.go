package tarws

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"
)

type RecursiveTarWriter struct {
	tarWriter *tar.Writer
}

func NewWriter(w io.Writer) *RecursiveTarWriter {
	return &RecursiveTarWriter{tarWriter: tar.NewWriter(w)}
}

func (tw *RecursiveTarWriter) Close() error {
	return tw.tarWriter.Close()
}

func (tw *RecursiveTarWriter) WriteRecursively(path string) error {
	return tw.recurseInto(path)
}

func (tw *RecursiveTarWriter) recurseInto(path string) (err error) {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	// read all file entries in dir in one slice
	fileInfos, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		currentPath := filepath.Join(path, fileInfo.Name())
		if fileInfo.IsDir() {
			err = tw.recurseInto(currentPath)
			if err != nil {
				return err
			}
		} else {
			log.Printf("Adding %s\n", currentPath)
			err = tw.write(currentPath, fileInfo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (tw *RecursiveTarWriter) write(path string, fileInfo os.FileInfo) (err error) {
	err = tw.tarWriter.WriteHeader(&tar.Header{
		Name:    path,
		Size:    fileInfo.Size(),
		Mode:    int64(fileInfo.Mode()),
		ModTime: fileInfo.ModTime(),
	})
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(tw.tarWriter, file)
	return err
}

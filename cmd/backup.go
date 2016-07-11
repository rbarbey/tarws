package cmd

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup data",
		Run:   backup,
	}
)

func backup(cmd *cobra.Command, args []string) {
	fmt.Printf("Backupd command %+v\n", args)

	tarWriter := tar.NewWriter(ioutil.Discard)
	defer tarWriter.Close()

	iterate(args[0], tarWriter)
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

func init() {
	backupCmd.Flags().StringVarP(&region, "region", "r", "", "Region in which the target S3 bucket is located")
	backupCmd.Flags().StringVarP(&bucket, "bucket", "b", "", "S3 bucket to which the resulting tar should be uploaded")
	backupCmd.Flags().StringVarP(&key, "key", "k", "", "name of the uploaded file in the target S3 bucket")

	TarwsCmd.AddCommand(backupCmd)
}

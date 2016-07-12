package cmd

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup data",
		Run:   backup,
	}
	region, bucket, key string
)

func backup(cmd *cobra.Command, args []string) {
	// create a pipe so that s3uploader can read from tar's writer
	reader, writer := io.Pipe()

	go func() {
		tarWriter := tar.NewWriter(writer)
		defer tarWriter.Close()

		iterate(args[0], tarWriter)
		writer.Close()
	}()

	session := session.New(&aws.Config{
		Region: aws.String(region),
	})
	session.Handlers.Send.PushFront(func(r *request.Request) {
		log.Printf("Request: %s/%s, Payload: %s\n", r.ClientInfo.ServiceName, r.Operation, r.Params)
	})
	uploader := s3manager.NewUploader(session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		log.Fatalf("Error uploading to s3: %+v", err)
	}

	log.Println("Successfully uploaded to", result.Location)
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

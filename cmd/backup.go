package cmd

import (
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rbarbey/tarws/tarws"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup data",
		RunE:  backup,
	}
	region, bucket, key string
	compress            bool
)

func backup(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("No path for backing up was specified")
	}

	path := args[0]
	if _, err := os.Stat(path); err != nil {
		return err
	}

	// create a pipe so that s3uploader can read from tar's writer
	reader, writer := io.Pipe()

	go func() {
		tarwsWriter := tarws.NewWriter(writer)
		defer writer.Close()
		defer tarwsWriter.Close()

		err := tarwsWriter.WriteRecursively(path)
		handle(err)
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

	return nil
}

func prepareWriter(writer io.WriteCloser) io.WriteCloser {
	if compress {
		return gzip.NewWriter(writer)
	}

	return writer
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
	backupCmd.Flags().BoolVarP(&compress, "compress", "c", false, "compress before uploading")

	TarwsCmd.AddCommand(backupCmd)
}

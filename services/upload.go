package services

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"../config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/s3"
)

//Upload files
func Upload(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter,
		r *http.Request) {
		if !alreadyLoggedIn(w, r, env.Sessions) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		file, err := os.Open("test/test.jpg")
		if err != nil {
			fmt.Printf("err opening file: %s", err)
		}
		defer file.Close()
		fileInfo, _ := file.Stat()
		size := fileInfo.Size()
		buffer := make([]byte, size) // read file content to buffer

		file.Read(buffer)
		fileBytes := bytes.NewReader(buffer)
		fileType := http.DetectContentType(buffer)
		path := "/media/" + file.Name()

		params := &s3.PutObjectInput{
			Bucket:        aws.String("go-backend-s3-jimmy"),
			Key:           aws.String(path),
			Body:          fileBytes,
			ContentLength: aws.Int64(size),
			ContentType:   aws.String(fileType),
			ACL:           aws.String("public-read"),
		}

		resp, err := env.S3.PutObject(params)
		if err != nil {
			fmt.Printf("bad response: %s", err)
		}
		fmt.Printf("response %s", awsutil.StringValue(resp))
	}
}

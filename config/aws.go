package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//NewAWS returns aws configuration
func newAWSConfig() *aws.Config {
	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, token)

	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
	}

	return aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)
}

//NewS3 returns S3 bucket
func NewS3() *s3.S3 {
	return s3.New(session.New(), newAWSConfig())
}

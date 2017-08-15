package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

//NewAWS returns aws configuration
func NewAWS() *aws.Config {
	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, token)

	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
	}

	return aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)
}

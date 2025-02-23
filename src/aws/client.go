package aws

import (
	"context"
	"dumbledore/settings"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var AwsSsmClient *ssm.Client

func InitAwsClient() aws.Config {
	client, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(settings.GetEnv().AwsRegion))
	if err != nil {
		log.Fatalf("failed to load aws credentials, %v", err)
	}
	return client
}

func InitAwsSsmClient() error {
	client := ssm.NewFromConfig(InitAwsClient())
	AwsSsmClient = client
	log.Println("initialized AWS SSM Client")
	return nil
}

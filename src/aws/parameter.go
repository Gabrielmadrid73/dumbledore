package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetParameter(parameter string) *string {

	value, err := AwsSsmClient.GetParameter(context.TODO(), &ssm.GetParameterInput{Name: aws.String(parameter), WithDecryption: aws.Bool(true)})

	if err != nil {
		log.Printf("failed to get parameter. %s", err)
		return nil
	}

	return value.Parameter.Value
}

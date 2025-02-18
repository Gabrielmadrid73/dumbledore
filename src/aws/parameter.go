package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetParameter(c context.Context, parameter string) error {
	value, err := AwsSsmClient.GetParameter(context.TODO(), &ssm.GetParameterInput{Name: aws.String(parameter), WithDecryption: aws.Bool(true)})
	if err != nil {
		return err
	}
	log.Println(*value.Parameter.Value)
	return nil
}

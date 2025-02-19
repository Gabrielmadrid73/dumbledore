package k8s

import (
	"context"
	"dumbledore/aws"
	"dumbledore/types"
	"log"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ssmAnnotationParamName = "aws-ssm/aws-param-name"
	ssmAnnotationParamType = "aws-ssm/aws-param-type"
)

func GetSecret(namespace string, name string) *v1.Secret {
	secret, err := K8sClient.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		log.Println(err)
		return nil
	}

	return secret
}

func CheckSecretAnnotation(obj *v1.Secret) *types.SecretAnnotations {
	response := &types.SecretAnnotations{
		ParamName: "",
		ParamType: "",
	}
	for key, value := range obj.Annotations {
		if strings.Contains(key, ssmAnnotationParamName) {
			response.ParamName = value
		}

		if strings.Contains(key, ssmAnnotationParamType) {
			response.ParamType = value
		}
	}
	if response.ParamName == "" || response.ParamType == "" {
		log.Printf("secret: %s missing annotations %s or %s", obj.Name, ssmAnnotationParamName, ssmAnnotationParamType)
		return nil
	}
	return response
}

func UpdateSecret(secret *v1.Secret, metadata *types.SecretAnnotations) bool {
	if value := aws.GetParameter(metadata.ParamName); value != nil {
		body := SecretBody(secret, metadata.ParamType, *value)
		if _, err := K8sClient.CoreV1().Secrets(secret.Namespace).Update(context.TODO(), body, metav1.UpdateOptions{}); err != nil {
			log.Println(err)
			return false
		} else {
			log.Printf("synchronized secret: %s with parameter store: %s", secret.Name, metadata.ParamName)
		}
	}
	return true
}

func SecretBody(secret *v1.Secret, paramType string, data string) *v1.Secret {
	value := make(map[string]string)
	value[paramType] = data
	secret.StringData = value

	return secret
}

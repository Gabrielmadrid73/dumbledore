package k8s

import (
	"dumbledore/types"
	"log"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ssmAnnotationParamName = "aws-ssm/aws-param-name"
	ssmAnnotationParamType = "aws-ssm/aws-param-type"
	ssmAnnotationParamKey  = "aws-ssm/aws-param-key"
)

func GetSecret(c *gin.Context, namespace string, name string) error {
	if secret, err := K8sClient.CoreV1().Secrets(namespace).Get(c, name, metav1.GetOptions{}); err != nil || len(secret.Annotations[ssmAnnotationParamName]) > 0 {
		log.Printf("Skipped secret not found: %s, namespace: %s", name, namespace)

	} else {
		log.Printf("Secret: %s doesn't contain the annotation %s", name, ssmAnnotationParamName)
	}
	return nil
}

func SyncSecret(c *gin.Context, obj types.StructSecrets) error {
	for secret := range obj.Secrets {
		if err := GetSecret(c, obj.Namespace, obj.Secrets[secret]); err != nil {
			return c.Error(err)
		}
	}
	return nil
}

// if err := aws.GetParameter(c, secret.Annotations["aws-ssm/aws-param-name"]); err != nil {
// 	log.Printf("Error to update secret %s, caused by: %s", name, err)
// }

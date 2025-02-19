package controller

import (
	"dumbledore/types"
	"net/http"

	"dumbledore/k8s"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
)

func SecretList(obj types.StructSecrets) []*v1.Secret {
	var secretList []*v1.Secret
	for item := range obj.Secrets {
		err := k8s.GetSecret(obj.Namespace, obj.Secrets[item])
		if err != nil {
			secretList = append(secretList, err)
		}
	}
	return secretList
}

func SyncSecret(obj types.StructSecrets) error {
	secrets := SecretList(obj)
	for item := range secrets {
		if err := k8s.CheckSecretAnnotation(secrets[item]); err != nil {
			k8s.UpdateSecret(secrets[item], err)

		}

	}

	return nil
}

func SecretController(c *gin.Context) {
	var json types.StructSecrets

	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	SyncSecret(json)

	c.JSON(http.StatusOK, gin.H{"status": "Secrets synchronized"})
}

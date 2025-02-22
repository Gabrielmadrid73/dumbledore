package controller

import (
	"net/http"

	"dumbledore/k8s"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
)

type RequestBody struct {
	Namespace string   `json:"namespace" binding:"required"`
	Secrets   []string `json:"secrets" binding:"required"`
}

func SecretList(obj *RequestBody) []*v1.Secret {
	var secretList []*v1.Secret
	for item := range obj.Secrets {
		err := k8s.GetSecret(obj.Namespace, obj.Secrets[item])
		if err != nil {
			secretList = append(secretList, err)
		}
	}
	return secretList
}

func SyncSecret(obj *RequestBody) bool {
	if secrets := SecretList(obj); secrets != nil {
		for item := range secrets {
			if err := k8s.CheckSecretAnnotation(secrets[item]); err != nil {
				k8s.UpdateSecret(secrets[item], err)
			}
		}
		return true
	} else {
		return false
	}
}

func SecretController(c *gin.Context) {
	var json *RequestBody

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"expected fields": "namespace: string, secrets: []string",
		})
		return
	}

	if err := SyncSecret(json); !err {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error during sync"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "received"})
		return
	}
}

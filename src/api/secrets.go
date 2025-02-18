package api

import (
	"dumbledore/types"
	"net/http"

	"dumbledore/k8s"

	"github.com/gin-gonic/gin"
)

func Secret(c *gin.Context) {
	var json types.StructSecrets

	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	k8s.SyncSecret(c, json)

	c.JSON(http.StatusOK, gin.H{"status": "Secrets synchronized"})
}

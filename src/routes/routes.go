package routes

import (
	"dumbledore/aws"
	"dumbledore/controller"
	"dumbledore/k8s"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	routes := gin.Default()

	k8sInit := k8s.InitK8sClient()
	if k8sInit != nil {
		panic(k8sInit.Error())
	}
	ssmInit := aws.InitAwsSsmClient()
	if ssmInit != nil {
		panic(ssmInit.Error())
	}

	// TODO use BASE_PATH environment variable
	v1 := routes.Group("/api/v1")
	{
		go v1.POST("/secrets/sync", controller.SecretController)
	}
	// TODO use BASE_PATH environment variable
	routes.GET("/health-check", func(c *gin.Context) {
		c.String(http.StatusOK, "UP")
	})

	return routes
}

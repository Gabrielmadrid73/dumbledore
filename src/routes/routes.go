package routes

import (
	"dumbledore/aws"
	"dumbledore/controller"
	"dumbledore/k8s"
	"dumbledore/settings"
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

	v1 := routes.Group(settings.GetEnv().BasePath + "/api/v1")
	{
		go v1.POST("/secrets/sync", controller.SecretController)
	}

	routes.GET(settings.GetEnv().BasePath+"/health-check", func(c *gin.Context) {
		c.String(http.StatusOK, "UP")
	})

	return routes
}

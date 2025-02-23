package main

import (
	"dumbledore/routes"
	"dumbledore/settings"
)

func main() {
	gin := routes.SetupRouter()
	gin.Run(":" + settings.GetEnv().Port)
}

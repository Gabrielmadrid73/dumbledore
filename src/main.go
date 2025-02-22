package main

import (
	"dumbledore/routes"
)

func main() {
	gin := routes.SetupRouter()
	gin.Run(":8080")
}

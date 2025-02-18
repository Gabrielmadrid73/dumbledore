package main

import (
	"dumbledore/routes"
)

func main() {
	r := routes.SetupRouter()
	r.Run(":8080")
}

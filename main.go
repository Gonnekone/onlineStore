package main

import (
	"fmt"
	"github.com/Gonnekone/onlineStore/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDBDocker()
}

func main() {
	router := InitRoutes()

	if err := router.Run(); err != nil {
		fmt.Println("Failed to start the server")
	}
}

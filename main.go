package main

import (
	"fmt"
	"github.com/Gonnekone/onlineStore/initializers"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := InitRoutes()
	go func() {
		if err := router.Run(); err != nil {
			fmt.Println("Failed to start the server")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Print("Shutting down server...")
	time.Sleep(10 * time.Second)
}

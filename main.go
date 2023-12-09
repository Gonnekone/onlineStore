package main

import (
	"fmt"
	"github.com/Gonnekone/onlineStore/handlers"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()
	router.Static("/static", "./static")

	// user
	userControllers := router.Group("/user")
	{
		userControllers.POST("/registration")
		userControllers.POST("/login")
		userControllers.GET("/auth")
	}

	// type
	typeControllers := router.Group("/type")
	{
		typeControllers.POST("/", handlers.CreateType)
		typeControllers.GET("/", handlers.ReadAllTypes)
		typeControllers.PUT("/", handlers.UpdateType)
		typeControllers.DELETE("/", handlers.DeleteType)
	}

	// brand
	brandControllers := router.Group("/brand")
	{
		brandControllers.POST("/", handlers.CreateBrand)
		brandControllers.GET("/", handlers.ReadAllBrands)
		brandControllers.PUT("/", handlers.UpdateBrand)
		brandControllers.DELETE("/", handlers.DeleteBrand)
	}

	// device
	deviceControllers := router.Group("/device")
	{
		deviceControllers.POST("/", handlers.CreateDevice)
		deviceControllers.GET("/", handlers.ReadAllDevices)
		deviceControllers.PUT("/", handlers.UpdateDevice)
		deviceControllers.DELETE("/", handlers.DeleteDevice)
		deviceControllers.GET("/:id", handlers.ReadOneDevice)
	}

	if err := router.Run(); err != nil {
		fmt.Println("Failed to start the server")
	}
}

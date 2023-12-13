package main

import (
	"fmt"
	"github.com/Gonnekone/onlineStore/handlers"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/Gonnekone/onlineStore/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()
	router.Static("/static", "./static")

	api := router.Group("/api")

	// user
	userControllers := api.Group("/user")
	{
		userControllers.POST("/registration", handlers.Registration)
		userControllers.POST("/login", handlers.Login)
		userControllers.GET("/auth", middleware.Auth, handlers.Validate)
	}

	// type
	typeControllers := api.Group("/type")
	{
		typeControllers.POST("/", middleware.Auth, middleware.CheckRoleAdmin, handlers.CreateType)
		typeControllers.GET("/", handlers.ReadAllTypes)
		typeControllers.PUT("/", middleware.Auth, middleware.CheckRoleAdmin, handlers.UpdateType)
		typeControllers.DELETE("/", middleware.Auth, middleware.CheckRoleAdmin, handlers.DeleteType)
	}

	// brand
	brandControllers := api.Group("/brand")
	{
		brandControllers.POST("/", middleware.Auth, middleware.CheckRoleAdmin, handlers.CreateBrand)
		brandControllers.GET("/", handlers.ReadAllBrands)
		brandControllers.PUT("/", middleware.Auth, middleware.CheckRoleAdmin, handlers.UpdateBrand)
		brandControllers.DELETE("/", middleware.Auth, middleware.CheckRoleAdmin, handlers.DeleteBrand)
	}

	// device
	deviceControllers := api.Group("/device")
	{
		deviceControllers.POST("/", middleware.Auth, middleware.CheckRoleAdmin, handlers.CreateDevice)
		deviceControllers.GET("/", handlers.ReadAllDevices)
		deviceControllers.PUT("/", middleware.Auth, middleware.CheckRoleAdmin, handlers.UpdateDevice)
		deviceControllers.DELETE("/", middleware.Auth, middleware.CheckRoleAdmin, handlers.DeleteDevice)
		deviceControllers.GET("/:id", handlers.ReadOneDevice)
	}

	if err := router.Run(); err != nil {
		fmt.Println("Failed to start the server")
	}
}

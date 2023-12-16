package main

import (
	"github.com/Gonnekone/onlineStore/handlers"
	"github.com/Gonnekone/onlineStore/middleware"
	"github.com/Gonnekone/onlineStore/util"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	router := gin.New()

	router.NoRoute(func(c *gin.Context) {
		response := make(chan util.Response)
		go func(context *gin.Context) {
			response <- util.ErrorResponse(util.RouteNotFound)
		}(c.Copy())
		util.SendResponse(c, <-response)
	})

	router.NoMethod(func(c *gin.Context) {
		response := make(chan util.Response)
		go func(context *gin.Context) {
			response <- util.ErrorResponse(util.MethodNotAllowed)
		}(c.Copy())
		util.SendResponse(c, <-response)
	})

	api := router.Group("/api")
	{
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
	}

	router.Static("/static", "./static")

	return router
}

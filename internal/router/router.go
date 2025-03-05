package router

import "github.com/gin-gonic/gin"

func New() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	// TODO change logger using log library
	router.Use(gin.Logger())

	// TODO: Add controller

	api := router.Group("api")
	{
		v1 := api.Group("v1")
		{
			v1.POST("/calculate")
			expressionsGroup := v1.Group("expressions")
			{
				expressionsGroup.GET("")
				expressionsGroup.GET("/:id")
			}
		}
	}
	internal := router.Group("internal")
	{
		internal.GET("/task")
		internal.POST("/task")
	}

	return router
}

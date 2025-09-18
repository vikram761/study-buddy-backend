package routes

import (
	"log"

	"github.com/gin-gonic/gin"
)

func InitRoutes(port string) {
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	auth_routes := router.Group("/auth")
	{
		// login
		// sign up
	}

	teacher_routes := router.Group("/teacher")
	{
		// TECHER ROUTES
		// eg: teacher_routes.GET, teacher_routes.POST etc
	}

	student_routes := router.Group("/student")
	{
		// TECHER ROUTES
		// eg: student_routes.GET, student_routes.POST etc
	}

	chapter_routes := router.Group("/chapter")
	{
		// chapter routes
		// create
		// all
		// ?chatperId
	}

	router.POST("/create-vnovel")
	router.POST("/create-quiz")
	router.POST("/create-module")

	router.NoRoute(noRoute)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error occurred", err)
	}
}

func noRoute(c *gin.Context) {
	c.JSON(404, gin.H{"message": "Page not found"})
}

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

	teacher_routes := router.Group("/teacher")
	{
		// TECHER ROUTES
		// eg: teacher_routes.GET, teacher_routes.POST etc
	}
	
	student_routes := router.Group("/teacher")
	{
		// TECHER ROUTES
		// eg: student_routes.GET, student_routes.POST etc
	}
	
	router.NoRoute(noRoute)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error occurred", err)
	}
}

func noRoute(c *gin.Context) {
	c.JSON(404, gin.H{"message": "Page not found"})
}

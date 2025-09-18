package routes

import (
	"database/sql"
	"log"
	"study-buddy-backend/controllers"
	"study-buddy-backend/controllers/auth"
	"study-buddy-backend/controllers/lesson"

	"github.com/gin-gonic/gin"
)

func InitRoutes(port string, db *sql.DB) {
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	authRoutes := router.Group("/auth")
	{
		// login
		authRoutes.POST("/signup", auth.SignUpHandler(db))
		authRoutes.POST("/login", auth.LoginHandler(db))
		// sign up
	}

	// teacher_routes := router.Group("/teacher")
	// {
	// 	// TECHER ROUTES
	// 	// eg: teacher_routes.GET, teacher_routes.POST etc
	// }

	// student_routes := router.Group("/student")
	// {
	// 	// TECHER ROUTES
	// 	// eg: student_routes.GET, student_routes.POST etc
	// }

	lesson_routes := router.Group("/lesson")
	{
		lesson_routes.POST("/create",lesson.CreateLessonHandler(db))
		lesson_routes.POST("/all",lesson.GetAllLesson(db))
		lesson_routes.POST("/:lesson_id/edit", lesson.GetModulesByLessonID(db))
		// chapter routes
		// all
		// ?chatperId
	}

	router.POST("/create-module", controllers.CreateModule(db))
	router.POST("/all-module", controllers.CreateModule(db))

	router.NoRoute(noRoute)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error occurred", err)
	}
}

func noRoute(c *gin.Context) {
	c.JSON(404, gin.H{"message": "Page not found"})
}

package routes

import (
	"database/sql"
	"log"
	"study-buddy-backend/controllers"
	"study-buddy-backend/controllers/auth"
	"study-buddy-backend/controllers/lesson"
	"study-buddy-backend/controllers/student"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes(port string, db *sql.DB) {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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

	student_routes := router.Group("/student")
	{
		student_routes.POST("/lessons", student.StudentLessonHandler(db))
		// TECHER ROUTES
		// eg: student_routes.GET, student_routes.POST etc
	}

	lesson_routes := router.Group("/lesson")
	{
		lesson_routes.POST("/create", lesson.CreateLessonHandler(db))
		lesson_routes.POST("/all", lesson.GetAllLesson(db))
		lesson_routes.POST("/:lesson_id", lesson.GetModulesByLessonID(db))
		lesson_routes.GET("/getall", lesson.GetAllJesus(db))
		// chapter routes
	}

	router.POST("/create-module", controllers.CreateModule(db))
	router.POST("/all-module", controllers.GetAllModule(db))
	router.POST("/mod/edit", controllers.EditModuleData(db))
	router.GET("/mod/:module_id",controllers.GetModulesByModuleID(db))
	router.POST("/gen-vnovel", controllers.GenerateVnovel(db))
	router.POST("/join/:lesson_id", lesson.JoinLessonByID(db))

	router.NoRoute(noRoute)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error occurred", err)
	}
}

func noRoute(c *gin.Context) {
	c.JSON(404, gin.H{"message": "Page not found"})
}

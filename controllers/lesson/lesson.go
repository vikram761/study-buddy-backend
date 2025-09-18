package lesson

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"study-buddy-backend/services/db"
)

type LessonRequest struct {
	LessonName string `json:"lesson_name"`
	Subject    string `json:"subject"`
	TeacherID  string `json:"teacher_id"`
}

type LessonAllRequest struct {
	TeacherId string `json:"teacher_id" binding:"required"`
}

// Accepts DB connection and returns the route handler
func CreateLessonHandler(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LessonRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err := db.CreateLesson(dbConn, req.LessonName, req.Subject, req.TeacherID )
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Lesson created successfully"})
	}
}

func GetAllLesson(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LessonAllRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		lessons, err := db.FindAll(dbConn, req.TeacherId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"lessons": lessons,
		})
	}
}

func GetModulesByLessonID(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		lessonID := c.Param("lesson_id")
		if lessonID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "lesson_id is required in path"})
			return
		}

		modules, err := db.GetModulesByLessonID(dbConn, lessonID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve modules", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"modules": modules,
		})
	}
}


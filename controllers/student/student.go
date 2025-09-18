package student

import (
	"database/sql"
	"net/http"
	"study-buddy-backend/models"
	"study-buddy-backend/services/db"

	"github.com/gin-gonic/gin"
)

type StudentLessonRequest struct {
	StudentID string `json:"student_id"`
}

func StudentLessonHandler(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StudentLessonRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		lessons, err := db.StudentRelation(dbConn, req.StudentID)

		var ldetail []models.Lesson

		for _, l := range lessons {
			nl, err := db.GetLesson(dbConn, l.LessonID)
			if err != nil {
				continue // or return err, depending on desired behavior
			}
			ldetail = append(ldetail, nl)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"lessons": ldetail,
		})
	}
}

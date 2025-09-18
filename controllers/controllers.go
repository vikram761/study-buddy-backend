package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http" 
	"study-buddy-backend/models"
	"study-buddy-backend/services/db"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type CreateModuleRequest struct {
	LessonID   string          `json:"lessonid"`
	ModuleType models.ModuleType `json:"moduletype"`
	ModuleData json.RawMessage `json:"moduledata"` // raw JSON, like your nested object
}

type ModuleAllRequest struct {
	TeacherId string `json:"teacher_id" binding:"required"`
}


func CreateModule(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateModuleRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ✅ Correctly pass RawMessage and handle return value
		err := db.CreateModule(dbConn, req.LessonID, string(req.ModuleType), req.ModuleData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create module in database",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Module created successfully",
		})
	}
}


func GetAllModule(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ModuleAllRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		modules, err := db.FindAllModule(dbConn, req.TeacherId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"lessons": modules,
		})
	}
}


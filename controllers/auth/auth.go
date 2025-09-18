package auth

import (
	"study-buddy-backend/services/db"
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthRequest struct {
	Role            string `json:"role"`
	Name            string `json:"name,omitempty"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
}

func LoginHandler(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Simplified logic
		id, hashedPass, err := db.GetStudentByEmail(dbConn, req.Email)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(req.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "id": id})
	}
}

func SignUpHandler(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		err := db.CreateStudent(dbConn, req.Name, req.Email, string(hashedPass), req.ProfileImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}


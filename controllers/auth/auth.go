package auth

import (
	"database/sql"
	"log"
	"net/http"
	"study-buddy-backend/services/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		var token string
		var jwttoken string
		if req.Role == "student" {
			id, hashedPass, err := db.GetStudentByEmail(dbConn, req.Email)
			if err != nil || bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(req.Password)) != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			token = id
			jwttoken = hashedPass
		} else if req.Role == "teacher" {
			id, hashedPass, err := db.GetTeacherByEmail(dbConn, req.Email)
			if err != nil || bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(req.Password)) != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			token = id
			jwttoken = hashedPass
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "id": token, "jwt": jwttoken})
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
		log.Println(req)
		if req.Role == "student" {
			err := db.CreateStudent(dbConn, req.Name, req.Email, string(hashedPass), req.ProfileImageURL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else if req.Role == "teacher" {
			err := db.CreateTeacher(dbConn, req.Name, req.Email, string(hashedPass), req.ProfileImageURL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect role"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

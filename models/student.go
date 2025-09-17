package models

type Student struct {
	StudentID       int64  `json:"student_id"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
}

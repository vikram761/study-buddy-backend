package models

type Teacher struct {
	TeacherID       int64  `json:"teacher_id"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
}


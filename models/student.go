package models

import "encoding/json"

type Student struct {
	StudentID       int64  `json:"student_id"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
}

type StudRelation struct {
	StudentID string          `json:"student_id"`
	LessonID  string          `json:"lesson_id"`
	Stats     json.RawMessage `json:"status"`
}

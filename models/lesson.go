package models

type Lesson struct {
	LessonID   string `json:"lesson_id"`
	LessonName string `json:"lesson_name"`
	TeacherID  string `json:"teacher_id"` 
	Subject    string `json:"subject"`
}

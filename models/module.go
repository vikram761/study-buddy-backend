package models

type ModuleType string

const (
	VNovel ModuleType = "VNOVEL"
	Quiz ModuleType = "QUIZ"
)

type Module struct {
	ModuleId string `json:"module_id"`
	ModuleType string `json:"module_type"`
	ModuleData string `json:"module_data,omitempty"`
	LessonId string `json:"lesson_id"`
}

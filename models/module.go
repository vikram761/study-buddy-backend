package models

import "encoding/json"

type ModuleType string

const (
	VNovel ModuleType = "VNOVEL"
	Quiz   ModuleType = "QUIZ"
)

type Module struct {
	ModuleId   string          `json:"module_id"`
	LessonId   string          `json:"lesson_id"`
	ModuleType ModuleType      `json:"module_type"`
	ModuleData json.RawMessage `json:"module_data,omitempty"` // raw JSON, like your nested object
	ModuleName string          `json:"module_name"`
	ModuleDesc string          `json:"module_description"`
}

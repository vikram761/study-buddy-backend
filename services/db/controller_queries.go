package db

import (
	"database/sql"
	"encoding/json"
	"study-buddy-backend/models"
)

func CreateModule(db *sql.DB , lessonId, moduleType string, moduleData json.RawMessage) error {
	_, err := db.Exec(`
		INSERT INTO module (lesson_id, module_type, module_data)
		VALUES ($1, $2, $3)`, lessonId,moduleType, moduleData,
	)
	return err
}


func FindAllModule(db *sql.DB, teacherID string) ([]models.Module, error) {
	rows, err := db.Query(`
		SELECT module_id, module_type, module_data, lesson_id
		FROM module
		WHERE lesson_id = $1
	`, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []models.Module

	for rows.Next() {
		var module models.Module
		if err := rows.Scan(&module.ModuleId, &module.ModuleType, &module.ModuleData, &module.LessonId); err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modules, nil
}

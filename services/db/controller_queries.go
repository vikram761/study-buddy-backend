package db

import (
	"database/sql"
	"encoding/json"
	"study-buddy-backend/models"
)

func CreateModule(db *sql.DB , mod models.Module) error {
	_, err := db.Exec(`
		INSERT INTO module (lesson_id, module_type, module_data , module_name , module_description)
		VALUES ($1, $2, $3, $4, $5)`, mod.LessonId,mod.ModuleType, mod.ModuleData, mod.ModuleName, mod.ModuleDesc,
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
		return nil, err }
	defer rows.Close()

	var modules []models.Module

	for rows.Next() {
		var module models.Module
		if err := rows.Scan(&module.ModuleId, &module.ModuleType, &module.ModuleData, &module.LessonId , &module.ModuleName, &module.ModuleDesc); err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modules, nil
}

func EditModule(db *sql.DB, module_id string , module_data json.RawMessage) error {
	_, err := db.Exec(`
		INSERT INTO module (module_data)
		VALUES ($1) WHERE module_id = $2`, module_data,module_id,
	)
	return err
}

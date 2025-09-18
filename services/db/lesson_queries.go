package db

import (
	"database/sql"
	"fmt"
	"study-buddy-backend/models"
)

func CreateLesson(db *sql.DB, name, subject, teacherID string) error {
	_, err := db.Exec(`
		INSERT INTO lesson (lesson_name, subject, teacher_id)
		VALUES ($1, $2, $3)
	`, name, subject, teacherID)
	return err
}

func FindAll(db *sql.DB, teacherID string) ([]models.Lesson, error) {
	rows, err := db.Query(`
		SELECT lesson_id, lesson_name, teacher_id, subject
		FROM lesson
		WHERE teacher_id = $1
	`, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.Lesson

	for rows.Next() {
		var lesson models.Lesson
		if err := rows.Scan(&lesson.LessonID, &lesson.LessonName, &lesson.TeacherID, &lesson.Subject); err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lessons, nil
}

func GetModulesByModuleID(db *sql.DB, moduleID string) (models.Module, error) {
	var mod models.Module

	rows, err := db.Query(`
	SELECT module_id, module_type, module_data, lesson_id
	FROM module
	WHERE module_id = $1
`, moduleID)

	if err != nil {
		return mod, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&mod.ModuleId, &mod.ModuleType, &mod.ModuleData, &mod.LessonId)
		if err != nil {
			return mod, err
		}
	} else {
		return mod, fmt.Errorf("no module found with ID %s", moduleID)
	}

	return mod, nil
}

func GetModulesByLessonID(db *sql.DB, lessonID string) ([]models.Module, error) {
	rows, err := db.Query(`
		SELECT module_id, module_type, module_data, lesson_id
		FROM module
		WHERE lesson_id = $1
	`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []models.Module

	for rows.Next() {
		var mod models.Module
		err := rows.Scan(&mod.ModuleId, &mod.ModuleType, &mod.ModuleData, &mod.LessonId)
		if err != nil {
			return nil, err
		}
		modules = append(modules, mod)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modules, nil
}

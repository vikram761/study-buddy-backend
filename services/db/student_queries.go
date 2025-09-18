package db

import (
	"database/sql"
	"study-buddy-backend/models"
)

func StudentRelation(db *sql.DB, studID string) ([]models.StudentRelation, error) {
	rows, err := db.Query(`
		SELECT student_id,lesson_id
		FROM stud_lesson
		WHERE student_id = $1
	`, studID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.StudentRelation

	for rows.Next() {
		var mod models.StudentRelation

		err := rows.Scan(&mod.StudentID, &mod.LessonID)
		if err != nil {
			return nil, err
		}

		lessons = append(lessons, mod)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lessons, nil
}

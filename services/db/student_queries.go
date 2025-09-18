package db

import (
	"database/sql"
	"study-buddy-backend/models"
)

func StudentRelation(db *sql.DB, studID string) ([]models.StudRelation ,error) {
	rows, err := db.Query(`
		SELECT student_id,lesson_id,status
		FROM stud_lessson
		WHERE student_id = $1
	`, studID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.StudRelation

	for rows.Next() {
		var mod models.StudRelation
		err := rows.Scan(&mod.StudentID, &mod.LessonID, &mod.Stats)
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

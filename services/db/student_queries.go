package db

import (
	"database/sql"
	"encoding/json"
	"study-buddy-backend/models"
)

func StudentRelation(db *sql.DB, studID string) ([]models.StudRelation, error) {
	rows, err := db.Query(`
		SELECT student_id,lesson_id,status
		FROM stud_lesson
		WHERE student_id = $1
	`, studID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.StudRelation

	for rows.Next() {
		var mod models.StudRelation
		var rawStatus sql.NullString

		err := rows.Scan(&mod.StudentID, &mod.LessonID, &rawStatus)
		if err != nil {
			return nil, err
		}

		// Normalize empty values to nil
		if !rawStatus.Valid || rawStatus.String == "" || rawStatus.String == "null" || rawStatus.String == "{}" || rawStatus.String == "[]" {
			mod.Stats = nil
		} else {
			mod.Stats = json.RawMessage(rawStatus.String)
		}

		lessons = append(lessons, mod)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lessons, nil
}

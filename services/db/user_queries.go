package db

import (
	"database/sql"
)

func CreateStudent(db *sql.DB, name, email, hashedPass, image string) error {
	_, err := db.Exec(`INSERT INTO students (name, email, pass, profile_image_url) VALUES ($1, $2, $3, $4)`,
		name, email, hashedPass, image)
	return err
}

func CreateTeacher(db *sql.DB, name, email, hashedPass, image string) error {
	_, err := db.Exec(`INSERT INTO teachers (name, email, pass, profile_image_url) VALUES ($1, $2, $3, $4)`,
		name, email, hashedPass, image)
	return err
}

func GetStudentByEmail(db *sql.DB, email string) (string, string, string, error) {
	var id, hashedPass, name string
	err := db.QueryRow(`SELECT student_id, pass, name FROM students WHERE email=$1`, email).Scan(&id, &hashedPass, &name)
	if err != nil {
		return "", "", name, err
	}
	return id, hashedPass, name, nil
}

func GetTeacherByEmail(db *sql.DB, email string) (string, string, string, error) {
	var id, hashedPass, name string
	err := db.QueryRow(`SELECT teacher_id, pass, name FROM teachers WHERE email=$1`, email).Scan(&id, &hashedPass, &name)
	if err != nil {
		return "", "", "", err
	}
	return id, hashedPass, name, nil
}


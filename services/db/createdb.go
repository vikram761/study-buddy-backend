package db

import (
	"database/sql"
	"log"
)

func CreateDB(db *sql.DB) {
	createPgcryptoExtension := `CREATE EXTENSION IF NOT EXISTS "pgcrypto";`

	createModuleTypeEnum := `DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_type WHERE typname = 'module_type_enum'
  ) THEN
    CREATE TYPE module_type_enum AS ENUM ('VNOVEL', 'QUIZ');
  END IF;
END $$;`

	createStudTableQuery := `
	CREATE TABLE IF NOT EXISTS students (
		student_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		pass VARCHAR(250) NOT NULL,
		profile_image_url VARCHAR(500)
	);`

	createTeacherTableQuery := `
	CREATE TABLE IF NOT EXISTS teachers (
		teacher_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		pass VARCHAR(250) NOT NULL,
		profile_image_url VARCHAR(500)
	);`

	createLessonTableQuery := `
	CREATE TABLE IF NOT EXISTS lesson (
		lesson_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		lesson_name VARCHAR(255) NOT NULL,
		teacher_id UUID NOT NULL,
		subject VARCHAR(255) NOT NULL,
		FOREIGN KEY (teacher_id) REFERENCES teachers(teacher_id) ON DELETE CASCADE
	);`

	createModuleTableQuery := `
	CREATE TABLE IF NOT EXISTS module (
		module_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		module_type module_type_enum NOT NULL,
		module_name VARCHAR(300) NOT NULL,
		module_description VARCHAR(500) NOT NULL,
		module_data JSONB,
		lesson_id UUID NOT NULL,
		FOREIGN KEY (lesson_id) REFERENCES lesson(lesson_id) ON DELETE CASCADE
	);`

	createStudRelationTableQuery := `
	CREATE TABLE IF NOT EXISTS stud_lesson (
		student_id UUID NOT NULL,
		lesson_id UUID NOT NULL,
		status JSONB,
		PRIMARY KEY (student_id, lesson_id),
		FOREIGN KEY (student_id) REFERENCES students(student_id) ON DELETE CASCADE,
		FOREIGN KEY (lesson_id) REFERENCES lesson(lesson_id) ON DELETE CASCADE
	);`

	queries := []string{
		createPgcryptoExtension,
		createModuleTypeEnum,
		createStudTableQuery,
		createTeacherTableQuery,
		createLessonTableQuery,
		createModuleTableQuery,
		createStudRelationTableQuery,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error executing query: %v\nQuery:\n%s", err, query)
		}
	}

	log.Println("Database setup completed successfully.")
}


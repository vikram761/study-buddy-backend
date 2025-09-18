package db

import (
    "database/sql"
    "study-buddy-backend/models"
)

func GetLessonsForStudent(db *sql.DB, studentID string) ([]models.Lesson, error) {
    rows, err := db.Query(`
        SELECT l.lesson_id, l.lesson_name, l.teacher_id, l.subject
        FROM lesson l
        INNER JOIN stud_lesson sl ON l.lesson_id = sl.lesson_id
        WHERE sl.student_id = $1
    `, studentID)
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

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return lessons, nil
}

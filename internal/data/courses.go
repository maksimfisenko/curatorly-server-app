package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
)

type Course struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	AcademicYear string `json:"academicYear"`
	ProjectID    int64  `json:"projectID"`
}

type CourseModel struct {
	DB *sql.DB
}

func ValidateCourse(v *validator.Validator, course *Course) {
	v.Check(course.Title != "", "title", "must be provided")
	v.Check(len(course.Title) <= 500, "title", "must not be more than 500 bytes long")
}

func (m CourseModel) Insert(course *Course) error {
	args := []any{course.Title, course.AcademicYear, course.ProjectID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, queryCourseInsert, args...).Scan(&course.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m CourseModel) GetAllForProject(projectID int64) ([]*Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, queryCourseGetAllForProject, projectID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	courses := []*Course{}

	for rows.Next() {
		var course Course

		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.AcademicYear,
			&course.ProjectID,
		)
		if err != nil {
			return nil, err
		}

		courses = append(courses, &course)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

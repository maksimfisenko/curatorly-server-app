package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
)

type Course struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int32     `json:"version"`
}

func ValidateCourse(v *validator.Validator, course *Course) {
	v.Check(course.Title != "", "title", "must be provided")
	v.Check(len(course.Title) <= 500, "title", "must be less or equal than 500 bytes long")
}

type CourseStorage struct {
	DB *sql.DB
}

func (s CourseStorage) Insert(course *Course) error {
	query := `
		INSERT INTO core.courses (title)
		VALUES ($1)
		RETURNING id, created_at, updated_at, version;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, course.Title).Scan(
		&course.ID,
		&course.CreatedAt,
		&course.UpdatedAt,
		&course.Version,
	)
}

func (s CourseStorage) Get(id int64) (*Course, error) {
	return nil, nil
}

func (s CourseStorage) Update(course *Course) error {
	return nil
}

func (s CourseStorage) Delete(id int64) error {
	return nil
}

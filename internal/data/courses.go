package data

import (
	"context"
	"database/sql"
	"errors"
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
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, title, created_at, updated_at, version
		FROM core.courses
		WHERE id = $1;
	`

	var course Course

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.Title,
		&course.CreatedAt,
		&course.UpdatedAt,
		&course.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &course, nil
}

func (s CourseStorage) Update(course *Course) error {
	query := `
		UPDATE core.courses
		SET title = $1, updated_at = $2, version = version + 1
		WHERE id = $3
		RETURNING updated_at, version;
	`

	args := []interface{}{course.Title, time.Now(), course.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(
		&course.UpdatedAt,
		&course.Version,
	)
}

func (s CourseStorage) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM core.courses
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

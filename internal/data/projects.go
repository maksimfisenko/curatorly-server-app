package data

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
)

var (
	ErrDuplicateUserInProject = errors.New("user already in project")
)

type Project struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	CreatorID  int64     `json:"creator_id"`
	AccessCode string    `json:"accessCode"`
}

func GenerateAccessCode() string {
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func AccessCodeExists(accessCode string) (bool, error) {
	return true, nil
}

func ValidateProject(v *validator.Validator, project *Project) {
	v.Check(project.Title != "", "title", "must be provided")
	v.Check(len(project.Title) <= 500, "title", "must not be more than 500 bytes long")
}

type ProjectModel struct {
	DB *sql.DB
}

func (m ProjectModel) Insert(project *Project) error {
	args := []any{project.Title, GenerateAccessCode(), project.CreatorID}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, queryProjectInsert, args...).Scan(&project.ID, &project.AccessCode, &project.CreatedAt)
	if err != nil {
		return err
	}

	args = []any{project.ID, project.CreatorID}

	result, err := m.DB.ExecContext(ctx, queryProjectUserInsert, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrFailedToAddRecord
	}

	return nil
}

func (m ProjectModel) GetByAccessCode(accessCode string) (*Project, error) {
	var project Project

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, queryProjectGetByAccessCode, accessCode).Scan(
		&project.ID,
		&project.Title,
		&project.AccessCode,
		&project.CreatorID,
		&project.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &project, nil
}

func (m ProjectModel) InsertUser(projectID, userID int64) error {
	args := []any{projectID, userID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, queryProjectUserInsert, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "projects_users_pkey"`:
			return ErrDuplicateUserInProject
		default:
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrFailedToAddRecord
	}

	return nil
}

func (m ProjectModel) GetAllForUser(userID int64) ([]*Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, queryProjectGetAllForUser, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projects := []*Project{}

	for rows.Next() {
		var project Project

		err := rows.Scan(
			&project.ID,
			&project.Title,
			&project.AccessCode,
			&project.CreatorID,
			&project.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (m ProjectModel) Get(id int64) (*Project, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	var project Project

	err := m.DB.QueryRow(queryProjectGet, id).Scan(
		&project.ID,
		&project.Title,
		&project.AccessCode,
		&project.CreatorID,
		&project.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &project, nil
}

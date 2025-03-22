package data

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
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
	query := `
	INSERT INTO content.projects (title, access_code, creator_id)
	VALUES ($1, $2, $3)
	RETURNING id, access_code, created_at
	`
	args := []any{project.Title, GenerateAccessCode(), project.CreatorID}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&project.ID, &project.AccessCode, &project.CreatedAt)
	if err != nil {
		fmt.Println(project.Title)
		fmt.Println(GenerateAccessCode())
		fmt.Println(project.CreatorID)
		return err
	}

	query = `
	INSERT INTO content.projects_users (project_id, user_id)
	VALUES ($1, $2)
	`

	args = []any{project.ID, project.CreatorID}

	result, err := m.DB.ExecContext(ctx, query, args...)
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

func (m ProjectModel) GetAllForUser(userID int64) ([]*Project, error) {
	return make([]*Project, 10), nil
}

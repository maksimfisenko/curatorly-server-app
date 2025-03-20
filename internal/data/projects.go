package data

import (
	"database/sql"
	"time"

	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
	"github.com/xyproto/randomstring"
)

type Project struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	CreatorID  int64     `json:"creator_id"`
	AccessCode string    `json:"accessCode"`
}

func GenerateAccessCode() string {
	return randomstring.String(10)
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

func (m ProjectModel) GetAllForUser(userID int64) ([]*Project, error) {
	return make([]*Project, 10), nil
}

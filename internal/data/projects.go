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

func (m ProjectModel) Insert(project *Project) error {
	return nil
}

func (m ProjectModel) Get(id int64) (*Project, error) {
	return nil, nil
}

func (m ProjectModel) GerMembers(id int64) ([]*User, error) {
	return nil, nil
}

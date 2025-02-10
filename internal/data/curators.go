package data

import (
	"database/sql"
	"time"

	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
)

type Curator struct {
	ID         int64     `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	BirthDate  time.Time `json:"birth_date"`
	City       string    `json:"city"`
	University string    `json:"university"`
	Profile    string    `json:"profile"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Version    int32     `json:"version"`
}

func ValidateCurator(v *validator.Validator, curator *Curator) {
	v.Check(curator.FirstName != "", "first_name", "must be provided")
	v.Check(len(curator.FirstName) <= 500, "first_name", "must be less or equal than 500 bytes long")

	v.Check(curator.LastName != "", "last_name", "must be provided")
	v.Check(len(curator.LastName) <= 500, "last_name", "must be less or equal than 500 bytes long")

	v.Check(len(curator.MiddleName) <= 500, "middle_name", "must be less or equal than 500 bytes long")

	// TODO: phone validation

	v.Check(validator.Matches(curator.Email, validator.EmailRX) || len(curator.Email) == 0, "email", "should be valid")

	v.Check(curator.BirthDate.Before(time.Now()), "birth_date", "must be before today")

	v.Check(len(curator.City) <= 500, "city", "must be less or equal than 500 bytes long")

	v.Check(len(curator.University) <= 500, "university", "must be less or equal than 500 bytes long")

	v.Check(len(curator.Profile) <= 500, "profile", "must be less or equal than 500 bytes long")
}

type CuratorStorage struct {
	DB *sql.DB
}

func (s CuratorStorage) Insert(curator *Curator) error {
	return nil
}

func (s CuratorStorage) Get(id int64) (*Curator, error) {
	return nil, nil
}

func (s CuratorStorage) Update(curator *Curator) error {
	return nil
}

func (s CuratorStorage) Delete(id int64) error {
	return nil
}

func (s CuratorStorage) GetAll() ([]*Curator, error) {
	return nil, nil
}

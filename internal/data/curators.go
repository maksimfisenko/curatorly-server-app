package data

import (
	"context"
	"database/sql"
	"errors"
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
	query := `
		INSERT INTO core.curators 
			(first_name, last_name, middle_name, phone, email, birth_date, city, university, profile)
		VALUES 
			($1, $2, $3, $4, $5,$6, $7, $8, $9)
		RETURNING
			id, created_at, updated_at, version;
	`

	args := []interface{}{
		curator.FirstName,
		curator.LastName,
		curator.MiddleName,
		curator.Phone,
		curator.Email,
		curator.BirthDate,
		curator.City,
		curator.University,
		curator.Profile,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(
		&curator.ID,
		&curator.CreatedAt,
		&curator.UpdatedAt,
		&curator.Version,
	)
}

func (s CuratorStorage) Get(id int64) (*Curator, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT 
			id, first_name, last_name, middle_name, phone, email, 
			birth_date, city, university, profile, created_at, updated_at, version
		FROM 
			core.curators
		WHERE 
			id = $1;
	`

	var curator Curator

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&curator.ID,
		&curator.FirstName,
		&curator.LastName,
		&curator.MiddleName,
		&curator.Phone,
		&curator.Email,
		&curator.BirthDate,
		&curator.City,
		&curator.University,
		&curator.Profile,
		&curator.CreatedAt,
		&curator.UpdatedAt,
		&curator.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &curator, nil
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

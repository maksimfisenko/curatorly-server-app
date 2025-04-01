package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
)

type Curator struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Birthday  time.Time `json:"birthday"`
	Status    string    `json:"status"`
	ProjectID int64     `json:"projectID"`
}

type CuratorModel struct {
	DB *sql.DB
}

func ValidateCurator(v *validator.Validator, curator *Curator) {
	v.Check(curator.Name != "", "name", "must be provided")
	v.Check(len(curator.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(curator.Surname != "", "surname", "must be provided")
	v.Check(len(curator.Surname) <= 500, "surname", "must not be more than 500 bytes long")

	v.Check(curator.Status != "", "status", "must be provided")
	v.Check(len(curator.Status) <= 500, "status", "must not be more than 500 bytes long")
}

func (m CuratorModel) Insert(curator *Curator) error {
	args := []any{curator.Name, curator.Surname, curator.Birthday, curator.Status, curator.ProjectID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, queryCuratorInsert, args...).Scan(&curator.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m CuratorModel) GetAllForProject(projectID int64) ([]*Curator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, queryCuratorsGetAllForProject, projectID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	curators := []*Curator{}

	for rows.Next() {
		var curator Curator

		err := rows.Scan(
			&curator.ID,
			&curator.Name,
			&curator.Surname,
			&curator.Birthday,
			&curator.Status,
			&curator.ProjectID,
		)
		if err != nil {
			return nil, err
		}

		curators = append(curators, &curator)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return curators, nil
}

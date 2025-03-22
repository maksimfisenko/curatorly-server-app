package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrFailedToAddRecord = errors.New("failed to add record")
)

type Models struct {
	Users    UserModel
	Projects ProjectModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Projects: ProjectModel{DB: db},
	}
}

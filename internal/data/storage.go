package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Storage struct {
	Courses  CourseStorage
	Curators CuratorStorage
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Courses:  CourseStorage{DB: db},
		Curators: CuratorStorage{DB: db},
	}
}

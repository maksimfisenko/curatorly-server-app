package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Storage struct {
	Courses CourseStorage
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Courses: CourseStorage{DB: db},
	}
}

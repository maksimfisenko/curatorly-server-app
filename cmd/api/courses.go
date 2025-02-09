package main

import (
	"net/http"
	"time"

	"github.com/maksimfisenko/curatorly-server-app/internal/data"
)

func (app *application) showCourseHandler(w http.ResponseWriter, r *http.Request) {
	course := data.Course{
		ID:        1,
		Title:     "Course Name 1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}

	err := app.writeJSON(w, http.StatusOK, course, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "the server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

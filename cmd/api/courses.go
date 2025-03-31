package main

import (
	"errors"
	"net/http"

	"github.com/maksimfisenko/curatorly-server-app/internal/data"
	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
)

func (app *application) createCourseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title        string `json:"title"`
		AcademicYear string `json:"academicYear"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	projectID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	project, err := app.models.Projects.Get(projectID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	course := &data.Course{
		Title:        input.Title,
		AcademicYear: input.AcademicYear,
		ProjectID:    project.ID,
	}

	v := validator.New()

	if data.ValidateCourse(v, course); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Courses.Insert(course)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"course": course}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listProjectCoursesHandler(w http.ResponseWriter, r *http.Request) {
	projectID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	project, err := app.models.Projects.Get(projectID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	courses, err := app.models.Courses.GetAllForProject(project.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"courses": courses}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

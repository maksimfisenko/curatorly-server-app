package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/maksimfisenko/curatorly-server-app/internal/data"
	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
)

func (app *application) createCuratorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Birthday string `json:"birthday"`
		Status   string `json:"status"`
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

	birthDate, _ := time.Parse(time.DateOnly, input.Birthday)

	curator := &data.Curator{
		Name:      input.Name,
		Surname:   input.Surname,
		Birthday:  birthDate,
		Status:    input.Status,
		ProjectID: project.ID,
	}

	v := validator.New()

	if data.ValidateCurator(v, curator); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Curators.Insert(curator)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"curator": curator}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listProjectCuratorsHandler(w http.ResponseWriter, r *http.Request) {
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

	curators, err := app.models.Curators.GetAllForProject(project.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"curators": curators}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

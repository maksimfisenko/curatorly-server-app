package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/maksimfisenko/curatorly-server-app/internal/data"
	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
)

func (app *application) createCuratorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName  string    `json:"first_name"`
		LastName   string    `json:"last_name"`
		MiddleName string    `json:"middle_name"`
		Phone      string    `json:"phone"`
		Email      string    `json:"email"`
		BirthDate  time.Time `json:"birth_date"`
		City       string    `json:"city"`
		University string    `json:"university"`
		Profile    string    `json:"profile"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	curator := &data.Curator{
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		MiddleName: input.MiddleName,
		Phone:      input.Phone,
		Email:      input.Email,
		BirthDate:  input.BirthDate,
		City:       input.City,
		University: input.University,
		Profile:    input.Profile,
	}

	v := validator.New()

	if data.ValidateCurator(v, curator); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.storage.Curators.Insert(curator)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/curators/%d", curator.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"curator": curator}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showCuratorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	curator, err := app.storage.Curators.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"curator": curator}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCuratorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.storage.Curators.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "curator successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

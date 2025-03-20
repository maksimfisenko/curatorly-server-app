package main

import (
	"net/http"
)

func (app *application) listUserProjectsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	projects, err := app.models.Projects.GetAllForUser(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"projects": projects}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

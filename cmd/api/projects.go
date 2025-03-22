package main

import (
	"net/http"

	"github.com/maksimfisenko/curatorly-server-app/internal/data"
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

func (app *application) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	project := &data.Project{
		Title:     input.Title,
		CreatorID: user.ID,
	}

	err = app.models.Projects.Insert(project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"project": project}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

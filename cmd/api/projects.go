package main

import (
	"errors"
	"net/http"

	_ "github.com/maksimfisenko/curatorly-server-app/docs"
	"github.com/maksimfisenko/curatorly-server-app/internal/data"
	"github.com/maksimfisenko/curatorly-server-app/internal/validator"
)

//	@Summary		List user projects
//	@Description	List all the projects in which user is a creator or an ordinory member.
//	@Tags			project
//	@Accept			json
//	@Produce		json
//	@Router			/projects [get]
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

	v := validator.New()

	if data.ValidateProject(v, project); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Projects.Insert(project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"project": project}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addUserToProject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AccessCode string `json:"access_code"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	project, err := app.models.Projects.GetByAccessCode(input.AccessCode)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("project", "project with this code doesnt exist")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Projects.InsertUser(project.ID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateUserInProject):
			v.AddError("user", "user already is a member of this project")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

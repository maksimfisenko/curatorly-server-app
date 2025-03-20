package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/api/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/users/login", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/api/v1/projects", app.requireAuthenticatedUser(app.listUserProjectsHandler))

	return app.authenticate(router)
}

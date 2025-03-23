package main

import (
	"expvar"
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

	router.HandlerFunc(http.MethodPost, "/api/v1/projects", app.requireAuthenticatedUser(app.createProjectHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/projects", app.requireAuthenticatedUser(app.listUserProjectsHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/projects/add-user", app.requireAuthenticatedUser(app.addUserToProject))

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.rateLimit(app.authenticate(router))))
}

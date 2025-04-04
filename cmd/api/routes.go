package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/maksimfisenko/curatorly-server-app/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/api/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/users/login", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/users/current", app.showCurrentUserHandler)

	router.HandlerFunc(http.MethodPost, "/api/v1/projects/:id/curators/add", app.requireAuthenticatedUser(app.createCuratorHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/projects/:id/curators", app.requireAuthenticatedUser(app.listProjectCuratorsHandler))

	router.HandlerFunc(http.MethodPost, "/api/v1/projects/:id/courses", app.requireAuthenticatedUser(app.createCourseHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/projects/:id/courses", app.requireAuthenticatedUser(app.listProjectCoursesHandler))

	router.HandlerFunc(http.MethodPost, "/api/v1/add-users", app.requireAuthenticatedUser(app.addUserToProject))

	router.HandlerFunc(http.MethodPost, "/api/v1/projects", app.requireAuthenticatedUser(app.createProjectHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/projects", app.requireAuthenticatedUser(app.listUserProjectsHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/projects/:id", app.requireAuthenticatedUser(app.showProjectHandler))

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	router.HandlerFunc(http.MethodGet, "/swagger/:any", httpSwagger.WrapHandler)

	return app.metrics(app.recoverPanic(app.rateLimit(app.authenticate(router))))
}

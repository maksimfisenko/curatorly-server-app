package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/courses", app.listCoursesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/courses/:id", app.showCourseHandler)
	router.HandlerFunc(http.MethodPost, "/v1/courses", app.createCourseHandler)
	router.HandlerFunc(http.MethodPut, "/v1/courses/:id", app.updateCourseHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/courses/:id", app.deleteCourseHandler)

	router.HandlerFunc(http.MethodPost, "/v1/curators", app.createCuratorHandler)
	router.HandlerFunc(http.MethodGet, "/v1/curators/:id", app.showCuratorHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/curators/:id", app.deleteCuratorHandler)

	return app.recoverPanic(router)
}

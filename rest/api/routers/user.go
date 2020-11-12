package routers

import (
	"vl_sa/rest/api"
	"vl_sa/rest/authentication"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// SetUserRoutes sets routes for user operations
func SetUserRoutes(router *mux.Router, api *api.UserAPI, tokenMiddleware authentication.TokenMiddlewareInterface) *mux.Router {

	router.Handle("/user", negroni.New(
		negroni.HandlerFunc(tokenMiddleware.RequireTokenAuthentication),
		negroni.HandlerFunc(api.Update),
	)).Methods("PUT")

	return router
}

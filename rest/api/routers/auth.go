package routers

import (
	"vl_sa/rest/api"

	"github.com/gorilla/mux"
)

// SetAuthRoutes sets the routes user for authentification and authorization
func SetAuthRoutes(router *mux.Router, userAPI *api.UserAPI) *mux.Router {

	router.HandleFunc(
		"/signup",
		userAPI.SignUp,
	).Methods("POST")

	router.HandleFunc(
		"/signin",
		userAPI.SignIn,
	).Methods("POST")

	router.HandleFunc(
		"/signWithGoogle",
		userAPI.SignWithGoogle,
	).Methods("POST")

	router.HandleFunc(
		"/forgetpass",
		userAPI.ForgetPassRequest,
	).Methods("POST")

	router.HandleFunc(
		"/forgetpass",
		userAPI.ForgetPassVerify,
	).Methods("GET")

	router.HandleFunc(
		"/forgetpass",
		userAPI.ForgetPassReset,
	).Methods("PUT")

	return router
}

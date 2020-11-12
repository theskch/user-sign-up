package routers

import (
	"vl_sa/html/handlers"

	"github.com/gorilla/mux"
)

// SetUserHTMLRoutes sets user related html routes
func SetUserHTMLRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/", handlers.LoginHandler)
	router.HandleFunc("/home", handlers.HomeHandler)
	router.HandleFunc("/edit", handlers.EditHandler)
	router.HandleFunc("/usersignup", handlers.SignupHandler)
	router.HandleFunc("/logout", handlers.LogoutHandler)
	router.HandleFunc("/forgetpassword", handlers.ForgotPassword)
	router.HandleFunc("/resetpass", handlers.ResetPassword)
	return router
}

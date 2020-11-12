package authentication

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//CorsMiddlewareInterface ..
type CorsMiddlewareInterface interface {
	RegisterCors(router *mux.Router) *mux.Router
}

//CorsMiddleware ..
type CorsMiddleware struct {
}

//RegisterCors ...
func (corsMiddleware *CorsMiddleware) RegisterCors(router *mux.Router) http.Handler {
	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*",
		},
	})
	return corsOpts.Handler(router)
}

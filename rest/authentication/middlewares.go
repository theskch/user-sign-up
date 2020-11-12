package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"vl_sa/models"
	"vl_sa/rest/persistance"

	"github.com/dgrijalva/jwt-go"
)

// TokenMiddlewareInterface used for token validation on requests
type TokenMiddlewareInterface interface {
	RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc)
}

// TokenMiddleware implementation
type TokenMiddleware struct {
	JWTHandler JWTHandlerInterface
	TokenRepo  persistance.BlacklistTokenRepositoryInterface
}

// RequireTokenAuthentication implementation
func (tokenMiddleware *TokenMiddleware) RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	jwtString := req.Header.Get("Authorization")
	if jwtString == "" {
		unauthorizedResponse(rw, "Bearer token in missing")
		return
	}

	jwtStrings := strings.Split(jwtString, "Bearer ")

	if len(jwtStrings) < 2 {
		unauthorizedResponse(rw, fmt.Sprintf("Bearer prefix is missing [%s]", jwtString))
		return
	}

	claims := jwt.MapClaims{}
	token, err := tokenMiddleware.JWTHandler.ParseToken(jwtStrings[1], claims)
	if err == nil && token.Valid {
		if exists, err := tokenMiddleware.TokenRepo.TokenExists(jwtStrings[1]); err != nil || exists {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusTemporaryRedirect)
			_ = json.NewEncoder(rw).Encode(models.Unauthorized{
				Code:    http.StatusTemporaryRedirect,
				Message: "please login again",
			})
			return
		}

		req.Header.Add("user", claims["sub"].(string))
		next(rw, req)
	} else {
		unauthorizedResponse(rw, "Invalid token")
	}
}

func unauthorizedResponse(rw http.ResponseWriter, errorMessage string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(rw).Encode(models.Unauthorized{
		Code:    http.StatusUnauthorized,
		Message: errorMessage,
	})
}

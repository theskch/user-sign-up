package authentication

import (
	"context"
	"net/http"
	"vl_sa/rest/domain"
	"vl_sa/utils"

	"google.golang.org/api/idtoken"
)

// GoogleAuthInterface used for google authentification
type GoogleAuthInterface interface {
	GetGoogleUser(token string) (domain.GoogleUser, error)
}

// GoogleAuth implementation
type GoogleAuth struct {
	Client http.Client
}

// GetTokenInfo implementation
func (googleAuth *GoogleAuth) GetGoogleUser(token string) (domain.GoogleUser, error) {
	payload, err := idtoken.Validate(context.Background(), token, utils.GetEnvParameter("CID"))
	if err != nil {
		return domain.GoogleUser{}, err
	}

	user := domain.GoogleUser{}

	if id, ok := payload.Claims["id"]; ok {
		user.ID = id.(string)
	}
	if email, ok := payload.Claims["email"]; ok {
		user.Email = email.(string)
	}
	if verifiedEmail, ok := payload.Claims["verified_email"]; ok {
		user.VerifiedEmail = verifiedEmail.(bool)
	}
	if name, ok := payload.Claims["name"]; ok {
		user.Name = name.(string)
	}
	if givenName, ok := payload.Claims["given_name"]; ok {
		user.GivenName = givenName.(string)
	}
	if familyName, ok := payload.Claims["family_name"]; ok {
		user.FamilyName = familyName.(string)
	}
	if link, ok := payload.Claims["link"]; ok {
		user.Link = link.(string)
	}
	if picture, ok := payload.Claims["picture"]; ok {
		user.Picture = picture.(string)
	}
	if gender, ok := payload.Claims["gender"]; ok {
		user.Gender = gender.(string)
	}
	if local, ok := payload.Claims["locale"]; ok {
		user.Locale = local.(string)
	}
	if err != nil {
		return user, err
	}
	return user, nil
}

package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"vl_sa/rest/domain"
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
	response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token))
	if err != nil {
		return domain.GoogleUser{}, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	var user domain.GoogleUser
	err = json.Unmarshal(contents, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

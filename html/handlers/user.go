package handlers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"vl_sa/html/assets"
	"vl_sa/logger"
	"vl_sa/models"
	"vl_sa/utils"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Store for cookies
var Store *sessions.CookieStore
var client http.Client

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	Store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	Store.Options = &sessions.Options{
		Domain:   utils.GetEnvParameter("DOMAIN"),
		Path:     "/",
		MaxAge:   3600 * 8,
		HttpOnly: true,
	}

	gob.Register(models.UserResponse{})
	client = http.Client{}
}

// LoginHandler renders login template
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	data := make(map[string]interface{})
	data["error"] = false

	if r.Method == http.MethodPost {
		request := new(models.SigninRequest)
		request.Email = &email
		request.Password = &password

		requestBody, err := json.Marshal(request)
		if err != nil {
			fmt.Fprintln(w, "Failed creating request to api")
			return
		}

		response, err := http.Post(
			fmt.Sprintf("%s/signin", utils.GetEnvParameter("BASE_URL")),
			"application/json",
			bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("error getting response from api %s", err.Error()))
		}

		defer response.Body.Close()
		if response.StatusCode == http.StatusOK {
			signinResponse := new(models.SigninResponse)
			decoder := json.NewDecoder(response.Body)
			if err := decoder.Decode(&signinResponse); err != nil {
				fmt.Fprintln(w, fmt.Sprintf("failed decoding response from api %s", err.Error()))
			}

			sessions, _ := Store.Get(r, "auth")
			sessions.Values["user"] = signinResponse.User
			sessions.Values["token"] = signinResponse.Token
			sessions.Save(r, w)
			http.Redirect(w, r, "/home", http.StatusFound)
			return
		}

		badRequestResponse := new(models.BadRequest)
		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&badRequestResponse); err != nil {
			fmt.Fprintln(w, fmt.Sprintf("failed decoding response from api %s", err.Error()))
		}

		data["error"] = true
		data["message"] = badRequestResponse.Message
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	loginHTML := assets.MustAssetString("html/templates/login.html")
	loginTemplate := template.Must(template.New("login_view").Parse(loginHTML))

	render(w, r, loginTemplate, "login_view", data)
}

// SignupHandler renders login template
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	data := make(map[string]interface{})
	data["error"] = false

	if r.Method == http.MethodPost {
		request := new(models.SigninRequest)
		request.Email = &email
		request.Password = &password

		requestBody, err := json.Marshal(request)
		if err != nil {
			fmt.Fprintln(w, "Failed creating request to api")
			return
		}

		response, err := http.Post(
			fmt.Sprintf("%s/signup", utils.GetEnvParameter("BASE_URL")),
			"application/json",
			bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("error getting response from api %s", err.Error()))
		}

		defer response.Body.Close()
		if response.StatusCode == http.StatusOK {
			signinResponse := new(models.SigninResponse)
			decoder := json.NewDecoder(response.Body)
			if err := decoder.Decode(&signinResponse); err != nil {
				fmt.Fprintln(w, fmt.Sprintf("failed decoding response from api %s", err.Error()))
			}

			sessions, _ := Store.Get(r, "auth")
			sessions.Values["user"] = signinResponse.User
			sessions.Values["token"] = signinResponse.Token
			sessions.Save(r, w)
			http.Redirect(w, r, "/edit", http.StatusFound)
			return
		}

		badRequestResponse := new(models.BadRequest)
		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&badRequestResponse); err != nil {
			fmt.Fprintln(w, fmt.Sprintf("failed decoding response from api %s", err.Error()))
		}

		data["error"] = true
		data["message"] = badRequestResponse.Message
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := assets.MustAssetString("html/templates/signup.html")
	template := template.Must(template.New("signup_view").Parse(html))
	render(w, r, template, "signup_view", data)
}

// HomeHandler renders home template
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "auth")
	user, err := getUser(session)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := assets.MustAssetString("html/templates/homepage.html")
	template := template.Must(template.New("homepage_view").Parse(html))

	render(w, r, template, "homepage_view", user)
}

// EditHandler renders edit template
func EditHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "auth")
	user, err := getUser(session)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := make(map[string]interface{})
	data["error"] = false
	data["user"] = user

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		if user.GoogleAuth {
			email = user.Email
		}

		fullname := r.FormValue("fullName")
		address := r.FormValue("address")
		telephone := r.FormValue("telephone")

		request := models.UpdateUserRequest{
			Email:     email,
			FullName:  fullname,
			Address:   address,
			Telephone: telephone,
		}

		requestBody, err := json.Marshal(request)
		if err != nil {
			fmt.Fprintln(w, "Failed creating request to api")
			return
		}

		var token string
		token, err = getToken(session)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		updateReq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/user", utils.GetEnvParameter("BASE_URL")), bytes.NewBuffer(requestBody))
		updateReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("error creating update request %s", err.Error()))
			return
		}

		response, err := client.Do(updateReq)
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("failed executing update %s", err.Error()))
			return
		}

		if response.StatusCode == http.StatusOK {
			user.Email = email
			user.Address = address
			user.Telephone = telephone
			user.FullName = fullname
			session.Values["user"] = &user
			if err := session.Save(r, w); err != nil {
				fmt.Printf("failed saving user in session: %s", err)
			}

			http.Redirect(w, r, "/home", http.StatusFound)
			return
		}

		badRequestResponse := new(models.BadRequest)
		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&badRequestResponse); err != nil {
			fmt.Fprintln(w, fmt.Sprintf("failed decoding response from api %s", err.Error()))
			return
		}

		data["error"] = true
		data["message"] = badRequestResponse.Message
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := assets.MustAssetString("html/templates/edit.html")
	template := template.Must(template.New("edit_view").Parse(html))

	render(w, r, template, "edit_view", data)
	return
}

// LogoutHandler clears the session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "auth")

	session.Values["user"] = &models.UserResponse{}
	session.Values["token"] = ""
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// ForgotPassword renders forgot password template
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["success"] = false
	data["failed"] = false

	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		request := new(models.PasswordResetEmailRequest)
		request.Email = &email

		requestBody, err := json.Marshal(request)
		if err != nil {
			fmt.Fprintln(w, "Failed creating request to api")
			return
		}

		response, err := http.Post(
			fmt.Sprintf("%s/forgetpass", utils.GetEnvParameter("BASE_URL")),
			"application/json",
			bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("error getting response from api %s", err.Error()))
			return
		}

		if response.StatusCode == http.StatusBadRequest {
			badRequestResponse := new(models.BadRequest)
			decoder := json.NewDecoder(response.Body)
			if err := decoder.Decode(&badRequestResponse); err != nil {
				fmt.Fprintln(w, fmt.Sprintf("failed decoding response from api %s", err.Error()))
				return
			}

			data["failed"] = true
			data["message"] = badRequestResponse.Message
		} else {
			data["success"] = true
			data["message"] = "Password reset link sent"
		}

	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := assets.MustAssetString("html/templates/forgotpass.html")
	template := template.Must(template.New("forgotpass_view").Parse(html))

	render(w, r, template, "forgotpass_view", data)
}

// ResetPassword ..
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	params, ok := r.URL.Query()["token"]
	if !ok || len(params) < 1 {
		fmt.Fprintln(w, "reset token missing")
	}

	token := params[0]
	validateRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/forgetpass", utils.GetEnvParameter("BASE_URL")), nil)

	query := validateRequest.URL.Query()
	query.Add("token", token)
	validateRequest.URL.RawQuery = query.Encode()

	if err != nil {
		fmt.Fprintln(w, fmt.Sprintf("error creating token validation request %s", err.Error()))
		return
	}

	validateResponse, err := client.Do(validateRequest)
	if err != nil {
		fmt.Fprintln(w, fmt.Sprintf("failed executing token validation %s", err.Error()))
		return
	}

	data := make(map[string]interface{})
	if validateResponse.StatusCode == http.StatusOK {
		data["failed"] = false
		tokenVerify := new(models.ForgetPasswordTokenVerify)
		decoder := json.NewDecoder(validateResponse.Body)
		if err := decoder.Decode(&tokenVerify); err != nil {
			fmt.Fprintln(w, fmt.Sprintf("failed decoding response from api %s", err.Error()))
			return
		}

		if r.Method == http.MethodPost {
			password := r.FormValue("password")
			passwordReserRequest := new(models.PasswordResetRequest)
			passwordReserRequest.Password = &password
			passwordReserRequest.Token = &token

			requestBody, err := json.Marshal(passwordReserRequest)
			if err != nil {
				fmt.Fprintln(w, "Failed creating request to api")
				return
			}

			resetRequest, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/forgetpass", utils.GetEnvParameter("BASE_URL")), bytes.NewBuffer(requestBody))
			if err != nil {
				fmt.Fprintln(w, fmt.Sprintf("error creating password reset request %s", err.Error()))
				return
			}

			_, err = client.Do(resetRequest)
			if err != nil {
				fmt.Fprintln(w, fmt.Sprintf("failed executing password reset %s", err.Error()))
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		data["failed"] = true
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := assets.MustAssetString("html/templates/resetpass.html")
	template := template.Must(template.New("resetpass_view").Parse(html))
	render(w, r, template, "resetpass_view", data)
}

// Render renders a template
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		logger.Error.Printf("Render Error: %v", err)
		return
	}
	w.Write(buf.Bytes())
}

func getUser(s *sessions.Session) (models.UserResponse, error) {
	val := s.Values["user"]
	var user = models.UserResponse{}
	user, ok := val.(models.UserResponse)
	if !ok {
		logger.Error.Printf("failed to get user from session")
		return user, fmt.Errorf("failed to get user from session")
	}
	return user, nil
}

func getToken(s *sessions.Session) (string, error) {
	val := s.Values["token"]
	token, ok := val.(string)
	if !ok {
		logger.Error.Printf("failed to get token from session")
		return token, fmt.Errorf("failed to get token from session")
	}
	return token, nil
}

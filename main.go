package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"vl_sa/html/handlers"
	htmlRouters "vl_sa/html/routers"
	"vl_sa/logger"
	"vl_sa/models"
	"vl_sa/rest/api"
	apiRouters "vl_sa/rest/api/routers"
	"vl_sa/rest/authentication"
	"vl_sa/rest/email"
	"vl_sa/rest/persistance"
	services "vl_sa/rest/service"
	_ "vl_sa/rest/swaggerui"
	"vl_sa/utils"

	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/rakyll/statik/fs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  fmt.Sprintf("%s/callback", utils.GetEnvParameter("BASE_URL")),
	ClientID:     utils.GetEnvParameter("CID"),
	ClientSecret: utils.GetEnvParameter("CSECRET"),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: google.Endpoint,
}

// Store for cookies
var store *sessions.CookieStore

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   3600 * 8,
		HttpOnly: true,
	}

	gob.Register(models.UserResponse{})
}

func main() {
	port := utils.GetEnvParameter("PORT")
	router := mux.NewRouter().StrictSlash(true)
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	staticServer := http.FileServer(statikFS)
	sh := http.StripPrefix("/rest/swaggerui/", staticServer)
	router.PathPrefix("/rest/swaggerui").Handler(sh)
	router.PathPrefix("/html/static/").Handler(http.StripPrefix("/html/static", http.FileServer(http.Dir(utils.GetEnvParameter("STATIC_DIR")))))

	router.HandleFunc("/googlelogin", googleLoginHandler)
	router.HandleFunc("/callback", callbackHandler)

	router = registerRESTRoutes(router)
	router = registerHTMLRoutes(router)
	handler := registerCors(router)

	logger.Info.Printf("Starting up the server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	oauthStateString := uniuri.New()
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintf(w, "Code exchange failed with error %s\n", err.Error())
		return
	}

	if !token.Valid() {
		fmt.Fprintln(w, "Retreived invalid token")
		return
	}

	idToken := token.Extra("id_token").(string)
	request := new(models.SignWithGoogleRequest)
	request.Token = &idToken
	requestBody, err := json.Marshal(request)
	if err != nil {
		fmt.Fprintln(w, "Failed creating request to api")
		return
	}

	response, err := http.Post(
		fmt.Sprintf("%s/signWithGoogle", utils.GetEnvParameter("BASE_URL")),
		"application/json",
		bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Fprintln(w, fmt.Sprintf("error getting response from api %s", err.Error()))
	}

	defer response.Body.Close()
	signupResponse := new(models.SignupResponse)
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&signupResponse); err != nil {
		fmt.Fprintln(w, fmt.Sprintf("failed decoding response from api %s", err.Error()))
	}

	sessions, _ := handlers.Store.Get(r, "auth")
	sessions.Values["user"] = signupResponse.User
	sessions.Values["token"] = signupResponse.Token
	sessions.Save(r, w)
	http.Redirect(w, r, "/home", http.StatusFound)
}

func registerCors(router *mux.Router) http.Handler {
	corsMiddleware := authentication.CorsMiddleware{}
	return corsMiddleware.RegisterCors(router)
}

func registerHTMLRoutes(router *mux.Router) *mux.Router {
	htmlRouters.SetUserHTMLRoutes(router)
	return router
}

func registerRESTRoutes(router *mux.Router) *mux.Router {
	dbInit := persistance.DBInitializer{
		Settings: persistance.GetDBSettings(),
	}

	database, err := dbInit.ConfigureDB()
	if err != nil {
		log.Fatalf("failed to initialize database")
	}

	userRepo := persistance.UserRepository{Database: database}
	tokenRepo := persistance.BlacklistTokenRepository{Database: database}
	jwtHandler := authentication.JWTHandler{Settings: authentication.GetJWTSettings()}
	passwordHandler := authentication.PasswordResetHandler{Settings: authentication.GetPasswordResetSettings()}
	middleware := authentication.TokenMiddleware{JWTHandler: &jwtHandler, TokenRepo: &tokenRepo}
	emailClient := email.SendgridClient{Settings: email.GetSettings()}

	googleAuth := authentication.GoogleAuth{Client: http.Client{}}

	userAPI := buildUserAPI(&userRepo, &jwtHandler, &passwordHandler, &googleAuth, &emailClient)
	router = apiRouters.SetAuthRoutes(router, &userAPI)
	router = apiRouters.SetUserRoutes(router, &userAPI, &middleware)

	return router
}

func buildUserAPI(
	userRepo persistance.UserRepositoryInterface,
	jwtHandler authentication.JWTHandlerInterface,
	passwordHandler authentication.PasswordResetHandlerInterface,
	googleAuth authentication.GoogleAuthInterface,
	emaiLClient email.ClientInterface,
) api.UserAPI {

	authService := services.AuthService{
		UserRepo:        userRepo,
		JWTHandler:      jwtHandler,
		PasswordHandler: passwordHandler,
		GoogleAuth:      googleAuth,
		EmailClient:     emaiLClient,
	}
	userService := services.UserService{
		UserRepo: userRepo,
	}

	return api.UserAPI{
		AuthService: &authService,
		UserService: &userService,
	}
}

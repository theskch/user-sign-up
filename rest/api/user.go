package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vl_sa/logger"
	"vl_sa/models"
	"vl_sa/rest/api/utils"
	"vl_sa/rest/constants"
	services "vl_sa/rest/service"
)

// UserAPI ..
type UserAPI struct {
	AuthService services.AuthServiceInterface
	UserService services.UserServiceInterface
}

// SignUp implemntation
func (api *UserAPI) SignUp(w http.ResponseWriter, r *http.Request) {
	request := new(models.SignupRequest)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		utils.BuildInternalServerError(w, err)
		return
	}

	if err := request.Validate(nil); err != nil {
		utils.BuildBadRequest(w, err)
		return
	}

	logger.Info.Printf("Sign up request: %+v\n", request)
	status, response := api.AuthService.SignUp(request)
	utils.BuildResponse(w, status, response)
}

// SignIn implementation
func (api *UserAPI) SignIn(w http.ResponseWriter, r *http.Request) {
	request := new(models.SigninRequest)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		utils.BuildInternalServerError(w, err)
		return
	}

	if err := request.Validate(nil); err != nil {
		utils.BuildBadRequest(w, err)
		return
	}

	logger.Info.Printf("Sign in request: %+v\n", request)
	status, response := api.AuthService.SignIn(request)
	utils.BuildResponse(w, status, response)
}

// SignWithGoogle implementation
func (api *UserAPI) SignWithGoogle(w http.ResponseWriter, r *http.Request) {
	request := new(models.SignWithGoogleRequest)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		utils.BuildInternalServerError(w, err)
		return
	}

	if err := request.Validate(nil); err != nil {
		utils.BuildBadRequest(w, err)
		return
	}

	logger.Info.Printf("Sign in request: %+v\n", request)
	status, response := api.AuthService.SignWithGoogle(request)
	utils.BuildResponse(w, status, response)
}

// Update implementation
func (api *UserAPI) Update(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
	request := new(models.UpdateUserRequest)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		utils.BuildInternalServerError(w, err)
		return
	}
	if err := request.Validate(nil); err != nil {
		utils.BuildBadRequest(w, err)
		return
	}
	logger.Info.Printf("update user request: %+v\n", request)
	userID := r.Header.Get("user")
	status, response := api.UserService.Update(request, userID)
	utils.BuildResponse(w, status, response)
}

// ForgetPassRequest implementation
func (api *UserAPI) ForgetPassRequest(w http.ResponseWriter, r *http.Request) {
	forgetPassRequest := new(models.PasswordResetEmailRequest)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&forgetPassRequest); err != nil {
		utils.BuildInternalServerError(w, err)
		return
	}
	if err := forgetPassRequest.Validate(nil); err != nil {
		utils.BuildBadRequest(w, err)
		return
	}
	logger.Info.Printf("Forget password request. User email: %s\n", *forgetPassRequest.Email)
	responseStatus, response := api.AuthService.ForgetPassRequest(forgetPassRequest)
	if responseStatus == constants.StatusOK {
		response = models.Sucessful{Code: http.StatusOK,
			Message: fmt.Sprintf("Successfully send reset token to [%s]", *forgetPassRequest.Email)}
	}

	utils.BuildResponse(w, responseStatus, response)
}

// ForgetPassVerify verifies password token
func (api *UserAPI) ForgetPassVerify(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	token := queryValues.Get("token")
	logger.Info.Printf("Password reset token: %s\n", token)

	responseStatus, response := api.AuthService.ForgetPassVerify(token)
	utils.BuildResponse(w, responseStatus, response)
}

// ForgetPassReset sets the new password
func (api *UserAPI) ForgetPassReset(w http.ResponseWriter, r *http.Request) {
	resetPassword := new(models.PasswordResetRequest)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&resetPassword); err != nil {
		utils.BuildInternalServerError(w, err)
		return
	}
	if err := resetPassword.Validate(nil); err != nil {
		utils.BuildBadRequest(w, err)
		return
	}
	logger.Info.Printf("Forget password reset token: %s\n", *resetPassword.Token)
	responseStatus, response := api.AuthService.ForgetPassReset(resetPassword)
	if responseStatus == http.StatusOK {
		response = models.Sucessful{
			Code: http.StatusOK,
		}
	}
	utils.BuildResponse(w, responseStatus, response)
}

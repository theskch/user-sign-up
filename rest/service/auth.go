package service

import (
	"fmt"
	"vl_sa/logger"
	"vl_sa/models"
	"vl_sa/rest/authentication"
	"vl_sa/rest/constants"
	"vl_sa/rest/domain"
	"vl_sa/rest/email"
	"vl_sa/rest/persistance"

	"github.com/google/uuid"
)

// AuthServiceInterface used for authentification and authorization of users
type AuthServiceInterface interface {
	SignUp(request *models.SignupRequest) (int, interface{})
	SignIn(request *models.SigninRequest) (int, interface{})
	SignWithGoogle(request *models.SignWithGoogleRequest) (int, interface{})

	ForgetPassRequest(emailModel *models.PasswordResetEmailRequest) (int, interface{})
	ForgetPassVerify(token string) (int, interface{})
	ForgetPassReset(resetPassword *models.PasswordResetRequest) (int, interface{})
}

// AuthService implementation
type AuthService struct {
	UserRepo        persistance.UserRepositoryInterface
	PasswordHandler authentication.PasswordResetHandlerInterface
	JWTHandler      authentication.JWTHandlerInterface
	GoogleAuth      authentication.GoogleAuthInterface
	EmailClient     email.ClientInterface
}

// SignUp implementation
func (authSvc *AuthService) SignUp(request *models.SignupRequest) (int, interface{}) {
	var user domain.UserAccount
	var err error

	if user, err = authSvc.UserRepo.FindByEmail(*request.Email); err != nil {
		return constants.StatusInternalServerError, err
	}

	if user.ID != "" {
		return constants.StatusBadRequest, fmt.Errorf("user with email %s already exists", *request.Email)
	}

	var hashedPassword string
	if hashedPassword, err = authSvc.PasswordHandler.HashPassword(*request.Password); err != nil {
		return constants.StatusInternalServerError, err
	}

	user = domain.UserAccount{
		ID:         uuid.New().String(),
		Email:      *request.Email,
		Password:   hashedPassword,
		GoogleAuth: false,
	}

	if err = authSvc.UserRepo.Insert(user); err != nil {
		return constants.StatusInternalServerError, err
	}

	token, err := authSvc.JWTHandler.GenerateToken(user.ID)
	if err != nil {
		return constants.StatusInternalServerError, err
	}

	return constants.StatusOK, models.SignupResponse{
		Token: token,
		User: &models.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			GoogleAuth: user.GoogleAuth,
		},
	}
}

// SignIn implementation
func (authSvc *AuthService) SignIn(request *models.SigninRequest) (int, interface{}) {
	var user domain.UserAccount
	var err error

	if user, err = authSvc.UserRepo.FindByEmail(*request.Email); err != nil {
		return constants.StatusInternalServerError, err
	}

	if user.ID == "" {
		return constants.StatusBadRequest, fmt.Errorf("user with email %s doesn't exist", *request.Email)
	}

	if ok := authSvc.PasswordHandler.CheckPasswordHash(*request.Password, user.Password); !ok {
		return constants.StatusBadRequest, fmt.Errorf("wrong password for user %s", user.Email)
	}

	token, err := authSvc.JWTHandler.GenerateToken(user.ID)
	if err != nil {
		return constants.StatusInternalServerError, err
	}

	return constants.StatusOK, models.SignupResponse{
		Token: token,
		User: &models.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			Address:    user.Address,
			FullName:   user.FullName,
			Telephone:  user.Telephone,
			GoogleAuth: user.GoogleAuth,
		},
	}
}

// SignWithGoogle implementation
func (authSvc *AuthService) SignWithGoogle(request *models.SignWithGoogleRequest) (int, interface{}) {
	var googleUser domain.GoogleUser
	var err error

	if googleUser, err = authSvc.GoogleAuth.GetGoogleUser(*request.Token); err != nil {
		return constants.StatusInternalServerError, err
	}

	var user domain.UserAccount
	if user, err = authSvc.UserRepo.FindByEmail(googleUser.Email); err != nil {
		return constants.StatusInternalServerError, err
	}

	if user.ID == "" {
		user.Email = googleUser.Email
		user.FullName = googleUser.Name
		user.GoogleAuth = true
		user.ID = uuid.New().String()
		if err = authSvc.UserRepo.Insert(user); err != nil {
			return constants.StatusInternalServerError, err
		}
	}

	token, err := authSvc.JWTHandler.GenerateToken(user.ID)
	if err != nil {
		return constants.StatusInternalServerError, err
	}

	return constants.StatusOK, models.SignupResponse{
		Token: token,
		User: &models.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			Address:    user.Address,
			FullName:   user.FullName,
			Telephone:  user.Telephone,
			GoogleAuth: user.GoogleAuth,
		},
	}
}

// ForgetPassRequest implemntation
func (authSvc *AuthService) ForgetPassRequest(emailModel *models.PasswordResetEmailRequest) (int, interface{}) {
	user, err := authSvc.UserRepo.FindByEmail(*emailModel.Email)
	if err != nil {
		return constants.StatusInternalServerError, err
	}

	if user.ID == "" {
		return constants.StatusBadRequest, fmt.Errorf("user with email %s doesn't exist", *emailModel.Email)
	}

	token, passwordResetLink := authSvc.PasswordHandler.GenerateResetPasswordToken(user.Email, user.Password)
	logger.Debug.Printf("Reset password token [%s]", token)

	content := struct {
		FullName          string
		PasswordResetLink string
	}{
		FullName:          user.FullName,
		PasswordResetLink: passwordResetLink,
	}

	if err = authSvc.EmailClient.SendMail(content, constants.EmailPasswordReset, user.Email, "Password reset"); err != nil {
		return constants.StatusInternalServerError, err
	}

	return constants.StatusOK, nil
}

// ForgetPassVerify implemntation
func (authSvc *AuthService) ForgetPassVerify(token string) (int, interface{}) {
	emailAddress, err := authSvc.PasswordHandler.VerifyResetPasswordToken(token, authSvc.getPasswordByEmail)

	if err != nil {
		return constants.StatusBadRequest, err
	}

	return constants.StatusOK, models.ForgetPasswordTokenVerify{Token: token, Email: emailAddress}
}

// ForgetPassReset implementation
func (authSvc *AuthService) ForgetPassReset(resetPassword *models.PasswordResetRequest) (int, interface{}) {

	emailAddress, err := authSvc.PasswordHandler.VerifyResetPasswordToken(*resetPassword.Token, authSvc.getPasswordByEmail)
	if err != nil {
		return constants.StatusBadRequest, err
	}

	hashedPassword, err := authSvc.PasswordHandler.HashPassword(*resetPassword.Password)

	if err != nil {
		return constants.StatusInternalServerError, err
	}

	if err = authSvc.UserRepo.UpdatePasswordByEmail(emailAddress, hashedPassword); err != nil {
		return constants.StatusInternalServerError, err
	}

	return constants.StatusOK, nil
}

func (authSvc *AuthService) getPasswordByEmail(email string) ([]byte, error) {
	user, err := authSvc.UserRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return []byte(user.Password), nil
}

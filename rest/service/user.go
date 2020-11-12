package service

import (
	"fmt"
	"vl_sa/models"
	"vl_sa/rest/constants"
	"vl_sa/rest/domain"
	"vl_sa/rest/persistance"
)

// UserServiceInterface used for user type operations
type UserServiceInterface interface {
	Update(request *models.UpdateUserRequest, id string) (int, interface{})
}

// UserService implementation
type UserService struct {
	UserRepo persistance.UserRepositoryInterface
}

// Update implementation
func (userSvc *UserService) Update(request *models.UpdateUserRequest, id string) (int, interface{}) {
	var user domain.UserAccount
	var err error

	if user, err = userSvc.UserRepo.FindByID(id); err != nil {
		return constants.StatusInternalServerError, err
	}

	if user.ID == "" {
		return constants.StatusBadRequest, fmt.Errorf("user not found")
	}

	if user.Email != request.Email {
		if user.GoogleAuth {
			return constants.StatusBadRequest, fmt.Errorf("user signed with google can't change email")
		}

		userWithEmail, err := userSvc.UserRepo.FindByEmail(request.Email)
		if err != nil {
			return constants.StatusInternalServerError, err
		}

		if userWithEmail.ID != "" && userWithEmail.ID != id {
			return constants.StatusBadRequest, fmt.Errorf("user with email %s already exists", request.Email)
		}
	}

	user.Email = request.Email
	user.FullName = request.FullName
	user.Address = request.Address
	user.Telephone = request.Telephone

	if err = userSvc.UserRepo.Update(user); err != nil {
		return constants.StatusInternalServerError, err
	}

	return constants.StatusOK, models.UserResponse{
		ID:        user.ID,
		Address:   user.Address,
		Email:     user.Email,
		FullName:  user.FullName,
		Telephone: user.Telephone,
	}
}

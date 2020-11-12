package service

import (
	"fmt"
	"net/http"
	"testing"
	"vl_sa/models"
	"vl_sa/rest/domain"

	"github.com/stretchr/testify/assert"
)

func TestUserUpdate(t *testing.T) {

	t.Run("valid user update", func(t *testing.T) {
		service := UserService{UserRepo: &ValidUserRepo{}}
		updateUserRequest := models.UpdateUserRequest{
			Email: "new_email",
		}

		status, response := service.Update(&updateUserRequest, "valid_id")
		assert.Equal(t, http.StatusOK, status)
		assert.IsType(t, models.UserResponse{}, response)
	})

	t.Run("database error while fetching user", func(t *testing.T) {
		service := UserService{UserRepo: &UserFindByIDErrorRepo{}}
		updateUserRequest := models.UpdateUserRequest{
			Email: "new_email",
		}

		status, response := service.Update(&updateUserRequest, "valid_id")
		assert.Equal(t, http.StatusInternalServerError, status)
		err, _ := response.(error)
		assert.Error(t, err)
	})

	t.Run("user not found error", func(t *testing.T) {
		service := UserService{UserRepo: &UserNotFoundRepo{}}
		updateUserRequest := models.UpdateUserRequest{
			Email: "new_email",
		}

		status, response := service.Update(&updateUserRequest, "valid_id")
		assert.Equal(t, http.StatusBadRequest, status)
		err, _ := response.(error)
		assert.Error(t, err)
	})

	t.Run("update email of google user error", func(t *testing.T) {
		service := UserService{UserRepo: &UpdateEmailForGoogleUserRepo{}}
		updateUserRequest := models.UpdateUserRequest{
			Email: "new_email",
		}

		status, response := service.Update(&updateUserRequest, "valid_id")
		assert.Equal(t, http.StatusBadRequest, status)
		err, _ := response.(error)
		assert.Error(t, err)
	})

	t.Run("database error while updating user", func(t *testing.T) {
		service := UserService{UserRepo: &UpdateUserDatabaseErrorRepo{}}
		updateUserRequest := models.UpdateUserRequest{
			Email: "new_email",
		}

		status, response := service.Update(&updateUserRequest, "valid_id")
		assert.Equal(t, http.StatusInternalServerError, status)
		err, _ := response.(error)
		assert.Error(t, err)
	})

	t.Run("find user by email error", func(t *testing.T) {
		service := UserService{UserRepo: &UserFindByEmailErrorRepo{}}
		updateUserRequest := models.UpdateUserRequest{
			Email: "new_email",
		}

		status, response := service.Update(&updateUserRequest, "valid_id")
		assert.Equal(t, http.StatusInternalServerError, status)
		err, _ := response.(error)
		assert.Error(t, err)
	})

	t.Run("user id missmatch", func(t *testing.T) {
		service := UserService{UserRepo: &UserIDMissmatchRepo{}}
		updateUserRequest := models.UpdateUserRequest{
			Email: "new_email",
		}

		status, response := service.Update(&updateUserRequest, "valid_id")
		assert.Equal(t, http.StatusBadRequest, status)
		err, _ := response.(error)
		assert.Error(t, err)
	})
}

type ValidUserRepo struct{}

func (repo *ValidUserRepo) FindByEmail(email string) (domain.UserAccount, error) {
	return domain.UserAccount{ID: "valid_id", Email: "valid_email", GoogleAuth: false}, nil
}

func (repo *ValidUserRepo) FindByID(email string) (domain.UserAccount, error) {
	return domain.UserAccount{ID: "valid_id", Email: "valid_email", GoogleAuth: false}, nil
}

func (repo *ValidUserRepo) Insert(user domain.UserAccount) error {
	return nil
}

func (repo *ValidUserRepo) Update(user domain.UserAccount) error {
	return nil
}

func (repo *ValidUserRepo) UpdatePasswordByEmail(email string, password string) error {
	return nil
}

type UserFindByIDErrorRepo struct{}

func (repo *UserFindByIDErrorRepo) FindByEmail(email string) (domain.UserAccount, error) {
	return domain.UserAccount{}, nil
}

func (repo *UserFindByIDErrorRepo) FindByID(email string) (domain.UserAccount, error) {
	return domain.UserAccount{}, fmt.Errorf("some error")
}

func (repo *UserFindByIDErrorRepo) Insert(user domain.UserAccount) error {
	return nil
}

func (repo *UserFindByIDErrorRepo) Update(user domain.UserAccount) error {
	return nil
}

func (repo *UserFindByIDErrorRepo) UpdatePasswordByEmail(email string, password string) error {
	return nil
}

type UserNotFoundRepo struct{}

func (repo *UserNotFoundRepo) FindByEmail(email string) (domain.UserAccount, error) {
	return domain.UserAccount{}, nil
}

func (repo *UserNotFoundRepo) FindByID(email string) (domain.UserAccount, error) {
	return domain.UserAccount{}, nil
}

func (repo *UserNotFoundRepo) Insert(user domain.UserAccount) error {
	return nil
}

func (repo *UserNotFoundRepo) Update(user domain.UserAccount) error {
	return nil
}

func (repo *UserNotFoundRepo) UpdatePasswordByEmail(email string, password string) error {
	return nil
}

type UpdateEmailForGoogleUserRepo struct{}

func (repo *UpdateEmailForGoogleUserRepo) FindByEmail(email string) (domain.UserAccount, error) {
	return domain.UserAccount{}, nil
}

func (repo *UpdateEmailForGoogleUserRepo) FindByID(email string) (domain.UserAccount, error) {
	return domain.UserAccount{ID: "valid_id", GoogleAuth: true, Email: "old_email"}, nil
}

func (repo *UpdateEmailForGoogleUserRepo) Insert(user domain.UserAccount) error {
	return nil
}

func (repo *UpdateEmailForGoogleUserRepo) Update(user domain.UserAccount) error {
	return nil
}

func (repo *UpdateEmailForGoogleUserRepo) UpdatePasswordByEmail(email string, password string) error {
	return nil
}

type UpdateUserDatabaseErrorRepo struct{}

func (repo *UpdateUserDatabaseErrorRepo) FindByEmail(email string) (domain.UserAccount, error) {
	return domain.UserAccount{ID: "valid_id"}, nil
}

func (repo *UpdateUserDatabaseErrorRepo) FindByID(email string) (domain.UserAccount, error) {
	return domain.UserAccount{ID: "valid_id", Email: "old_email"}, nil
}

func (repo *UpdateUserDatabaseErrorRepo) Insert(user domain.UserAccount) error {
	return nil
}

func (repo *UpdateUserDatabaseErrorRepo) Update(user domain.UserAccount) error {
	return fmt.Errorf("some error")
}

func (repo *UpdateUserDatabaseErrorRepo) UpdatePasswordByEmail(email string, password string) error {
	return nil
}

type UserFindByEmailErrorRepo struct{}

func (repo *UserFindByEmailErrorRepo) FindByEmail(email string) (domain.UserAccount, error) {
	return domain.UserAccount{}, fmt.Errorf("some error")
}

func (repo *UserFindByEmailErrorRepo) FindByID(email string) (domain.UserAccount, error) {
	return domain.UserAccount{ID: "valid_id", Email: "old_email"}, nil
}

func (repo *UserFindByEmailErrorRepo) Insert(user domain.UserAccount) error {
	return nil
}

func (repo *UserFindByEmailErrorRepo) Update(user domain.UserAccount) error {
	return fmt.Errorf("some error")
}

func (repo *UserFindByEmailErrorRepo) UpdatePasswordByEmail(email string, password string) error {
	return nil
}

type UserIDMissmatchRepo struct{}

func (repo *UserIDMissmatchRepo) FindByEmail(email string) (domain.UserAccount, error) {
	return domain.UserAccount{ID: "valid_other_id"}, nil
}

func (repo *UserIDMissmatchRepo) FindByID(email string) (domain.UserAccount, error) {
	return domain.UserAccount{ID: "valid_id", Email: "old_email"}, nil
}

func (repo *UserIDMissmatchRepo) Insert(user domain.UserAccount) error {
	return nil
}

func (repo *UserIDMissmatchRepo) Update(user domain.UserAccount) error {
	return fmt.Errorf("some error")
}

func (repo *UserIDMissmatchRepo) UpdatePasswordByEmail(email string, password string) error {
	return nil
}

package persistance

import (
	"vl_sa/rest/domain"

	"github.com/jinzhu/gorm"
)

// UserRepositoryInterface gives access to user repository
type UserRepositoryInterface interface {
	FindByEmail(email string) (domain.UserAccount, error)
	FindByID(id string) (domain.UserAccount, error)

	Insert(user domain.UserAccount) error
	Update(user domain.UserAccount) error

	UpdatePasswordByEmail(email string, password string) error
}

// UserRepository implementation
type UserRepository struct {
	Database *gorm.DB
}

// FindByEmail implementation
func (userRepo *UserRepository) FindByEmail(email string) (domain.UserAccount, error) {
	user := domain.UserAccount{}
	err := userRepo.Database.Where("email = ?", email).First(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		return user, nil
	}

	return user, err
}

// FindByID implementation
func (userRepo *UserRepository) FindByID(id string) (domain.UserAccount, error) {
	user := domain.UserAccount{}
	err := userRepo.Database.Where("id = ?", id).First(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		return user, nil
	}

	return user, err
}

// Insert implementation
func (userRepo *UserRepository) Insert(user domain.UserAccount) error {
	return userRepo.Database.Create(&user).Error
}

// Update implementation
func (userRepo *UserRepository) Update(user domain.UserAccount) error {
	return userRepo.Database.Save(&user).Error
}

// UpdatePasswordByEmail implementation
func (userRepo *UserRepository) UpdatePasswordByEmail(email string, password string) error {
	return userRepo.Database.Model(&domain.UserAccount{}).Where("email = ?", email).Update("password", password).Error
}

package persistance

import (
	"vl_sa/rest/domain"

	"github.com/jinzhu/gorm"
)

// BlacklistTokenRepositoryInterface used for checking if the token is still valid
type BlacklistTokenRepositoryInterface interface {
	TokenExists(token string) (bool, error)
}

// BlacklistTokenRepository implementation
type BlacklistTokenRepository struct {
	Database *gorm.DB
}

// TokenExists implementation
func (blacklistTokenRepo *BlacklistTokenRepository) TokenExists(tokenString string) (bool, error) {
	token := domain.BlacklistToken{}
	err := blacklistTokenRepo.Database.Where("token = ?", tokenString).First(&token).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}

	return true, err
}

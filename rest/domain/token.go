package domain

// BlacklistToken db model
type BlacklistToken struct {
	ID    int    `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Token string `gorm:"type:varchar(256);UNIQUE"`
}

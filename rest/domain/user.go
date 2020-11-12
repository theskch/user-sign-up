package domain

// UserAccount db model
type UserAccount struct {
	ID         string `gorm:"type:varchar(36);PRIMARY_KEY"`
	Password   string `gorm:"type:varchar(1024)"`
	FullName   string `gorm:"type:varchar(1024)"`
	Address    string `gorm:"type:varchar(1024)"`
	Telephone  string `gorm:"type:varchar(1024)"`
	Email      string `gorm:"type:varchar(256);UNIQUE"`
	GoogleAuth bool   `gorm:"default:false"`
}

// GoogleUser account
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

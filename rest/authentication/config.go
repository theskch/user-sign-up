package authentication

import (
	"log"
	"strconv"
	"vl_sa/utils"
)

// PasswordResetSettings settings for password reset related operations
type PasswordResetSettings struct {
	TokenExpirationTime int
	TokenSecretKey      string
	WebsiteBaseURL      string
}

// GetPasswordResetSettings returns settings for password reset token
func GetPasswordResetSettings() PasswordResetSettings {
	tokenExpirationTimeString := utils.GetEnvParameter("PASSWORD_RESET_TOKEN_EXPIRATION_TIME")
	tokenExpirationTime, err := strconv.Atoi(tokenExpirationTimeString)
	if err != nil {
		log.Fatalf("failed to convert password expiration time %s. Int value required", tokenExpirationTimeString)
	}

	tokenSecretKey := utils.GetEnvParameter("PASSWORD_RESET_TOKEN_SECRET_KEY")
	websiteBaseURL := utils.GetEnvParameter("BASE_URL")
	return PasswordResetSettings{
		TokenExpirationTime: tokenExpirationTime,
		TokenSecretKey:      tokenSecretKey,
		WebsiteBaseURL:      websiteBaseURL,
	}
}

// JWTSettings settings for jwt token operations
type JWTSettings struct {
	JWTExpirationTime int
	JWTSecretKey      string
}

// GetJWTSettings returns settings for jwt tokens
func GetJWTSettings() JWTSettings {
	jwtExpirationTimeString := utils.GetEnvParameter("JWT_EXPIRATION_TIME")
	jwtExpirationTime, err := strconv.Atoi(jwtExpirationTimeString)
	if err != nil {
		log.Fatalf("failed to convert jwt expiration time %s. Int value required", jwtExpirationTimeString)
	}

	jwtSecretKey := utils.GetEnvParameter("JWT_SECRET_KEY")
	return JWTSettings{
		JWTExpirationTime: jwtExpirationTime,
		JWTSecretKey:      jwtSecretKey,
	}
}

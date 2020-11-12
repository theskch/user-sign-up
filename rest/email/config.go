package email

import "vl_sa/utils"

// Settings for email client
type Settings struct {
	FromAddress string
	APIKey      string
}

//GetSettings get settings from the env variables
func GetSettings() Settings {
	fromAddress := utils.GetEnvParameter("SEND_GRID_FROM_ADDRESS")
	apiKey := utils.GetEnvParameter("SEND_GRID_API_KEY")
	return Settings{
		FromAddress: fromAddress,
		APIKey:      apiKey,
	}
}

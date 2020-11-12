package persistance

import "vl_sa/utils"

// DBSettings ...
type DBSettings struct {
	URL string
}

// GetDBSettings returnes database settings
func GetDBSettings() DBSettings {
	return DBSettings{
		URL: utils.GetEnvParameter("DATABASE_URL"),
	}
}

package persistance

import (
	"fmt"
	"vl_sa/rest/domain"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DBInitializer used to initialize database
type DBInitializer struct {
	Settings DBSettings
}

// ConfigureDB configures database connection and does a migration
func (dbInitializer *DBInitializer) ConfigureDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dbInitializer.Settings.URL)
	if err != nil {
		fmt.Printf("ERROR %+v", err)
		return nil, err
	}
	db.AutoMigrate(
		&domain.UserAccount{},
		&domain.BlacklistToken{},
	)
	return db, nil
}

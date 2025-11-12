package config

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Luggage{},
		&Country{},
		&Parcel{},
	)
}

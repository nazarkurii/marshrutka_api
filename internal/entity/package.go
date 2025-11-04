package entity

import (
	"database/sql"
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type Package struct {
	ID           uuid.UUID `gorm:"type:binary(16);primaryKey"         json:"id"`
	UserID       uuid.UUID `gorm:"type:binary(16);not null"           json:"userId"`
	ConnectionID uuid.UUID `gorm:"type:binary(16);not null"           json:"connectionID"`
	PhoneNumber  string    `gorm:"type:varchar(15);not null"                                                  json:"phoneNumber"`
	Email        string    `gorm:"type:varchar(255);not null"                          json:"email"`

	PickUpAdressID  uuid.UUID      `gorm:"type:binary(16);not null"           json:"-"`
	PickUpAdress    Address        `gorm:"foreignKey:PickUpAdressID"    json:"pickUpAddress"`
	DropOffAdressID uuid.UUID      `gorm:"type:binary(16);not null"           json:"-"`
	DropOffAdress   Address        `gorm:"foreignKey:DropOffAdressID"   json:"dropOffAddress"`
	CreatedAt       time.Time      `gorm:"not null"                     json:"createdAt"`
	CompletedAt     sql.NullTime   `                                    json:"completedAt"`
	Payment         PackagePayment `gorm:"foreignKey:PackageID"    `
	DeletedAt       gorm.DeletedAt `                                    json:"deletedAt"`
	LuggageVolume   luggage        `gorm:"type:MEDIUMINT UNSIGNED;not null"`
	QRCode          []byte         `gorm:"type:blob;not null" json:"qrCode"`
}

type PackagePayment struct {
	PackageID uuid.UUID     `gorm:"type:binary(16);not null"                                                json:"packadeId"`
	Price     int           `gorm:"type:MEDIUMINT;not null"                                           json:"price"`
	Method    paymentMethod `gorm:"type:enum('Apple Pay','Card','Cash','Google Pay');not null"        json:"method"`
	CreatedAt time.Time     `gorm:"not null"                                                          json:"createdAt"`
	SessionID string        `gorm:"type:varchar(500);not null"                                                          json:"sessionID"`
	Succeeded bool          `gorm:"not null"                                                          json:"succeeded"`
}

func MigratePackage(db *gorm.DB) error {
	return db.AutoMigrate(
		&Package{},
		&PackagePayment{},
	)

}

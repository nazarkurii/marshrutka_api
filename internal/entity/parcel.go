package entity

import (
	"database/sql"
	"fmt"
	"maryan_api/config"
	rfc7807 "maryan_api/pkg/problem"
	"strconv"
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type Parcel struct {
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
	Payment         ParcelPayment  `gorm:"foreignKey:PackageID"    `
	DeletedAt       gorm.DeletedAt `                                    json:"deletedAt"`
	LuggageVolume   luggage        `gorm:"type:MEDIUMINT UNSIGNED;not null"`
	QRCode          []byte         `gorm:"type:blob;not null" json:"qrCode"`
}

type ParcelPayment struct {
	PackageID uuid.UUID     `gorm:"type:binary(16);not null"                                                json:"packadeId"`
	Price     int           `gorm:"type:MEDIUMINT;not null"                                           json:"price"`
	Method    paymentMethod `gorm:"type:enum('Apple Pay','Card','Cash','Google Pay');not null"        json:"method"`
	CreatedAt time.Time     `gorm:"not null"                                                          json:"createdAt"`
	SessionID string        `gorm:"type:varchar(500);not null"                                                          json:"sessionID"`
	Succeeded bool          `gorm:"not null"                                                          json:"succeeded"`
}

func MigratePackage(db *gorm.DB) error {
	return db.AutoMigrate(
		&Parcel{},
		&ParcelPayment{},
	)

}

type FindParcelConnectionsRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Length string `json:"length"`
}

type FindParcelConnectionsRequestParsed struct {
	From   uuid.UUID
	To     uuid.UUID
	Width  int
	Height int
	Length int
}

func (fpcr FindParcelConnectionsRequest) Parse() (FindParcelConnectionsRequestParsed, error) {
	var params rfc7807.InvalidParams

	width, err := strconv.Atoi(fpcr.Width)
	if err != nil {
		params.SetInvalidParam("width", err.Error())
	} else if width < config.MinParcelWidth {
		params.SetInvalidParam("width", fmt.Sprintf("Has to be greater that or equal to %d.", config.MinParcelWidth))
	}

	height, err := strconv.Atoi(fpcr.Height)
	if err != nil {
		params.SetInvalidParam("height", err.Error())
	} else if height < config.MinParcelHeight {
		params.SetInvalidParam("height", fmt.Sprintf("Has to be greater that or equal to %d.", config.MinParcelHeight))
	}

	length, err := strconv.Atoi(fpcr.Length)
	if err != nil {
		params.SetInvalidParam("length", err.Error())
	} else if height < config.MinParcelLength {
		params.SetInvalidParam("length", fmt.Sprintf("Has to be greater that or equal to %d.", config.MinParcelLength))
	}

	from, _, err := config.ParseCountry(fpcr.From)
	if err != nil {
		params.SetInvalidParam("from", err.Error())
	}

	to, _, err := config.ParseCountry(fpcr.To)
	if err != nil {
		params.SetInvalidParam("to", err.Error())
	}

	if params != nil {
		return FindParcelConnectionsRequestParsed{}, rfc7807.BadRequest("invalid-params", "Invalid Params Error", "Provided params are not valid.", params...)
	}

	return FindParcelConnectionsRequestParsed{
		from,
		to,
		width,
		height,
		length,
	}, nil

}

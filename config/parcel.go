package config

import (
	"errors"
	"slices"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type Parcel struct {
	ID         uuid.UUID `gorm:"type:binary(16);primaryKey" json:"-"`
	IsOverSize bool      `gorm:"not null"`
	Height     int       `gorm:"type:SMALLINT UNSIGNED;not null"  json:"height"`
	Width      int       `gorm:"type:SMALLINT UNSIGNED;not null"  json:"width"`
	Length     int       `gorm:"type:SMALLINT UNSIGNED;not null"  json:"length"`
	Price      uint      `gorm:"type:INT UNSIGNED;not null"  json:"price"`
}

type parcelType string

const (
	UndefinedSizeParcelType parcelType = "Undefined Size"
	UssualParcelType        parcelType = "Ussual"
)

var parcelsConfig []Parcel
var undefinedParcleSizePrice uint

func LoadParcelsConfig(db *gorm.DB) {
	var allParcels []Parcel

	response := db.Find(&allParcels)
	if response.Error != nil {
		panic(response.Error)
	} else if response.RowsAffected == 0 {
		panic(errors.New("There are no parcel types provided for config."))
	}

	parcelsConfig = make([]Parcel, 0, len(allParcels)-1)

	for _, parcel := range allParcels {
		if parcel.IsOverSize {
			undefinedParcleSizePrice = parcel.Price
		} else {
			parcelsConfig = append(parcelsConfig, parcel)
		}
	}

	if undefinedParcleSizePrice == 0 {
		panic("There is no 'Undefined Size' parcel config provided or its price value is 0.")
	}

	slices.SortFunc(parcelsConfig, func(a, b Parcel) int {
		switch {
		case a.Price < b.Price:
			return -1
		case a.Price > b.Price:
			return 1
		default:
			return 0
		}
	})
}

type ParcelPposition struct {
	Height, Width, Length int
}

func CalulatePrice(positions []ParcelPposition) uint {
	for _, parcel := range parcelsConfig {
		for _, position := range positions {
			if parcel.Width >= position.Width && parcel.Height >= position.Height && parcel.Length >= position.Length {
				return parcel.Price
			}
		}
	}

	return undefinedParcleSizePrice
}

func ParcelConfigTestData() []Parcel {
	return []Parcel{
		{
			ID:         uuid.New(),
			IsOverSize: false,
			Height:     10,   // cm
			Width:      15,   // cm
			Length:     20,   // cm
			Price:      5000, // €50.00
		},
		{
			ID:         uuid.New(),
			IsOverSize: false,
			Height:     20,
			Width:      25,
			Length:     35,
			Price:      7500, // €75.00
		},
		{
			ID:         uuid.New(),
			IsOverSize: false,
			Height:     30,
			Width:      40,
			Length:     50,
			Price:      10000, // €100.00
		},
		{
			ID:         uuid.New(),
			IsOverSize: false,
			Height:     40,
			Width:      50,
			Length:     70,
			Price:      13000, // €130.00
		},
		{
			ID:         uuid.New(),
			IsOverSize: false,
			Height:     60,
			Width:      60,
			Length:     100,
			Price:      18000, // €180.00
		},
		{
			ID:         uuid.New(),
			IsOverSize: true,
			Height:     0,
			Width:      0,
			Length:     0,
			Price:      25000, // €250.00 fallback price
		},
	}
}

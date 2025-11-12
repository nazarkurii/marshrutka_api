package config

import (
	"fmt"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type Luggage struct {
	ID     uuid.UUID   `gorm:"type:binary(16);primaryKey" json:"-"`
	Type   luggageType `gorm:"type:enum('Small','Medium','Large'); not null"  json:"type"`
	Height uint        `gorm:"type:SMALLINT UNSIGNED;not null"  json:"height"`
	Width  uint        `gorm:"type:SMALLINT UNSIGNED;not null"  json:"width"`
	Length uint        `gorm:"type:SMALLINT UNSIGNED;not null"  json:"length"`
	Price  uint        `gorm:"type:INT UNSIGNED;not null"  json:"price"`
}

var luggageConfig = struct {
	Small  Luggage
	Medium Luggage
	Large  Luggage
}{}

type luggageType string

const (
	SmallLuggageType  luggageType = "Small"
	MediumLuggageType luggageType = "Medium"
	LargeLuggageType  luggageType = "Large"
)

func getLuggagePointerByType(luggageType luggageType) *Luggage {
	switch luggageType {
	case SmallLuggageType:
		return &luggageConfig.Small
	case MediumLuggageType:
		return &luggageConfig.Medium
	case LargeLuggageType:
		return &luggageConfig.Large
	default:
		return nil
	}
}

func LoadLuggageConfig(db *gorm.DB) {
	var luggageDB []Luggage

	response := db.Find(&luggageDB)
	if response.Error != nil {
		panic(response.Error)
	} else if response.RowsAffected < 3 {
		panic(fmt.Errorf("Not enough luggage types: %v", response.RowsAffected))
	}

	for _, luggage := range luggageDB {
		luggageCongigPointer := getLuggagePointerByType(luggage.Type)
		*luggageCongigPointer = luggage
	}

}

func GetLuggageConfig() (Luggage, Luggage, Luggage) {
	return luggageConfig.Small, luggageConfig.Medium, luggageConfig.Large
}

func LuggageConfigTestData() []Luggage {
	return []Luggage{
		{
			ID:     uuid.New(),
			Type:   SmallLuggageType,
			Height: 55, // cm
			Width:  40,
			Length: 20,
			Price:  3000, // €30.00
		},
		{
			ID:     uuid.New(),
			Type:   MediumLuggageType,
			Height: 65,
			Width:  45,
			Length: 25,
			Price:  5000, // €50.00
		},
		{
			ID:     uuid.New(),
			Type:   LargeLuggageType,
			Height: 75,
			Width:  55,
			Length: 30,
			Price:  7000, // €70.00
		},
	}
}

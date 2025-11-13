package config

import (
	"github.com/d3code/uuid"
	"golang.org/x/exp/rand"
)

type Luggage struct {
	ID     uuid.UUID `gorm:"type:binary(16);primaryKey" json:"-"`
	Height uint      `gorm:"type:SMALLINT UNSIGNED;not null"`
	Width  uint      `gorm:"type:SMALLINT UNSIGNED;not null"`
	Length uint      `gorm:"type:SMALLINT UNSIGNED;not null"`
	Price  uint      `gorm:"type:INT UNSIGNED;not null"`
}

type LuggageConfig struct {
	ID        uuid.UUID `gorm:"type:binary(16);primaryKey" json:"-"`
	CountryID uuid.UUID `gorm:"type:binary(16);uniqueIndex" json:"-"`

	SmallID uuid.UUID `gorm:"type:binary(16)" json:"-"`
	Small   Luggage   `gorm:"foreignKey:SmallID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	MediumID uuid.UUID `gorm:"type:binary(16)" json:"-"`
	Medium   Luggage   `gorm:"foreignKey:MediumID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	LargeID uuid.UUID `gorm:"type:binary(16)" json:"-"`
	Large   Luggage   `gorm:"foreignKey:LargeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func createTestLuggageConfigs(countries []Country) []LuggageConfig {

	var luggageConfigs []LuggageConfig

	for _, country := range countries {
		// Randomized price base per country
		basePrice := uint(2000 + rand.Intn(3000)) // 20–50 EUR (in cents)

		small := Luggage{
			ID:     uuid.New(),
			Height: 40,
			Width:  25,
			Length: 20,
			Price:  basePrice,
		}

		medium := Luggage{
			ID:     uuid.New(),
			Height: 60,
			Width:  40,
			Length: 30,
			Price:  basePrice + 2000,
		}

		large := Luggage{
			ID:     uuid.New(),
			Height: 80,
			Width:  50,
			Length: 40,
			Price:  basePrice + 4000,
		}

		config := LuggageConfig{
			ID:        uuid.New(),
			CountryID: country.ID,
			SmallID:   small.ID,
			MediumID:  medium.ID,
			LargeID:   large.ID,
			Small:     small,
			Medium:    medium,
			Large:     large,
		}
		luggageConfigs = append(luggageConfigs, config)
	}

	return luggageConfigs
}

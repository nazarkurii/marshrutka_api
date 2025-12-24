package config

import (
	"github.com/d3code/uuid"
)

type Luggage struct {
	ID     uuid.UUID
	Height uint
	Width  uint
	Length uint
	Volume int
	Price  uint
}

type LuggageConfig struct {
	Small  Luggage
	Medium Luggage
	Large  Luggage
}

var luggageConfig = createTestLuggageConfig()

func createTestLuggageConfig() LuggageConfig {

	small := Luggage{
		ID:     uuid.New(),
		Height: 40,
		Width:  20,
		Length: 30,
		Volume: 40 * 20 * 30,
		Price:  3000,
	}

	medium := Luggage{
		ID:     uuid.New(),
		Height: 50,
		Width:  40,
		Length: 30,
		Volume: 40 * 50 * 30,
		Price:  5000,
	}

	large := Luggage{
		ID:     uuid.New(),
		Height: 100,
		Width:  50,
		Length: 40,
		Volume: 40 * 100 * 50,
		Price:  7000,
	}

	return LuggageConfig{
		Small:  small,
		Medium: medium,
		Large:  large,
	}
}

func GetLoggageConfig() LuggageConfig {
	return luggageConfig
}

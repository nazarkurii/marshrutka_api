package config

import (
	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type Country struct {
	ID            uuid.UUID      `gorm:"type:binary(16);primaryKey"`
	Name          string         `gorm:"type:varchar(50);not null;unique"`
	IsMaster      bool           `gorm:"not null"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	LuggageConfig LuggageConfig  `gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParcelsConfig ParcelsConfig  `gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

var countriesConfig = map[uuid.UUID]Country{}
var masterCountry Country

func LoadCountriesConfig(db *gorm.DB) {
	var countriesDB []Country

	response := db.Table("countries").Find(&countriesDB)
	if response.Error != nil {
		panic(response.Error)
	} else if response.RowsAffected == 0 {
		panic("There are no countriesConfig provitded for the config.")
	}

	for _, country := range countriesDB {
		if country.IsMaster {
			if masterCountry.ID != uuid.Nil {
				panic("There has to be only one master country (Dublicate Master Country Error).")
			} else {
				masterCountry = country
			}
		}
		countriesConfig[country.ID] = country
	}
}

func countriesConfigTestData() []Country {
	return []Country{
		{ID: uuid.New(), Name: "Poland", IsMaster: false},
		{ID: uuid.New(), Name: "Germany", IsMaster: false},
		{ID: uuid.New(), Name: "Czechia", IsMaster: false},
		{ID: uuid.New(), Name: "Estonia", IsMaster: false},
		{ID: uuid.New(), Name: "Latvia", IsMaster: false},
		{ID: uuid.New(), Name: "Lithuania", IsMaster: false},
		{ID: uuid.New(), Name: "Slovakia", IsMaster: false},
		{ID: uuid.New(), Name: "Hungary", IsMaster: false},
		{ID: uuid.New(), Name: "Ukraine", IsMaster: true},
	}
}

func CreateTestData() []Country {
	countries := countriesConfigTestData()
	luggages := createTestLuggageConfigs(countries)
	parcels := createTestParcelConfigs(countries)
	for i := range countries {
		countries[i].LuggageConfig = luggages[i]
		countries[i].ParcelsConfig = parcels[i]
	}

	return countries
}

func GetConfig() (map[uuid.UUID]Country, Country) {
	return countriesConfig, masterCountry
}

func GetCountryByID(countryID uuid.UUID) (Country, bool) {
	country, ok := countriesConfig[countryID]
	return country, ok
}

func MustGetCountryByID(countryID uuid.UUID) Country {
	return countriesConfig[countryID]
}

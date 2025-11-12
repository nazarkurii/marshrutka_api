package config

import (
	"fmt"
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type Country struct {
	ID        uuid.UUID `gorm:"type:binary(16);primaryKey"`
	Name      string    `gorm:"type:varchar(50);not null; UNIQUE"`
	IsMaster  bool      `gorm:"not null"`
	DeletedAt gorm.DeletedAt
}

type countryType string

var countriesConfig = map[string]Country{}
var masterCountry Country

func LoadCountriesConfig(db *gorm.DB) {

	var CountriesConfigSlice []Country

	response := db.Table("countriesConfig").Select("id", "name").Scan(&CountriesConfigSlice)
	if response.Error != nil {
		panic(response.Error)
	} else if response.RowsAffected == 0 {
		panic("There are no countriesConfig provitded for the config.")
	}

	for _, country := range CountriesConfigSlice {
		if country.IsMaster {
			if masterCountry.ID != uuid.Nil {
				panic("There has to be only one master country (Dublicate Master Country Error).")
			} else {
				masterCountry = country
			}
		}
		countriesConfig[country.Name] = country
	}
}

func ParseCountry(countryName string) (uuid.UUID, *time.Location, error) {
	country, ok := countriesConfig[countryName]
	if !ok {
		return uuid.Nil, nil, fmt.Errorf("Non-existing country")
	}
	location, err := countryToLocation(country.Name)

	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("Non-existing time location")
	}

	return country.ID, location, nil
}

func MustParseTimeToLocalByCountryName(date time.Time, country string) time.Time {
	location, _ := countryToLocation(country)

	return date.In(location)
}

func MustParseTimeToLocalByCountryUUID(date time.Time, countryID uuid.UUID) time.Time {
	var name string
	for _, country := range countriesConfig {
		if countryID == country.ID {
			name = country.Name
			break
		}
	}
	location, _ := countryToLocation(name)

	return date.In(location)
}

func countryToLocation(countryName string) (*time.Location, error) {
	mapping := map[string]string{
		"Germany":   "Europe/Berlin",
		"Poland":    "Europe/Warsaw",
		"Czechia":   "Europe/Prague",
		"Estonia":   "Europe/Tallinn",
		"Latvia":    "Europe/Riga",
		"Lithuania": "Europe/Vilnius",
		"Slovakia":  "Europe/Bratislava",
		"Hungary":   "Europe/Budapest",
		"Ukraine":   "Europe/Kyiv",
	}

	if tz, ok := mapping[countryName]; ok {
		return time.LoadLocation(tz)
	}
	return nil, fmt.Errorf("unsupported country: %s", countryName)
}

func MustGetLocationByCountryID(id uuid.UUID) *time.Location {
	var name string
	for _, country := range countriesConfig {
		if country.ID == id {
			name = country.Name
			break
		}
	}
	location, _ := countryToLocation(name)

	return location
}

func GetCountriesConfig() ([]uuid.UUID, uuid.UUID) {
	var countriesConfigIDs []uuid.UUID

	for _, country := range countriesConfig {
		if country.ID != masterCountry.ID {
			countriesConfigIDs = append(countriesConfigIDs, country.ID)
		}
	}

	return countriesConfigIDs, masterCountry.ID
}

func CountriesConfigTestData() []Country {
	return []Country{
		{
			ID:       uuid.New(),
			Name:     "Poland",
			IsMaster: false,
		},
		{
			ID:       uuid.New(),
			Name:     "Germany",
			IsMaster: false,
		},
		{
			ID:       uuid.New(),
			Name:     "Czechia",
			IsMaster: false,
		},
		{
			ID:       uuid.New(),
			Name:     "Estonia",
			IsMaster: false,
		},
		{
			ID:       uuid.New(),
			Name:     "Latvia",
			IsMaster: false,
		},
		{
			ID:       uuid.New(),
			Name:     "Lithuania",
			IsMaster: false,
		},
		{
			ID:       uuid.New(),
			Name:     "Slovakia",
			IsMaster: false,
		},
		{
			ID:       uuid.New(),
			Name:     "Hungary",
			IsMaster: false,
		},
		{
			ID:       uuid.New(),
			Name:     "Ukraine", // optional: pick one as master country
			IsMaster: true,
		},
		{
			ID:       uuid.New(),
			Name:     "Germany",
			IsMaster: true,
		},
	}
}

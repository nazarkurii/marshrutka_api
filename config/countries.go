package config

import (
	"fmt"
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

var countries = map[string]uuid.UUID{}

func LoadCountries(db *gorm.DB) {

	var countriesSlice []struct {
		ID   uuid.UUID
		Name string
	}

	err := db.Table("countries").Select("id", "name").Scan(&countriesSlice).Error
	if err != nil {
		panic(err)
	}

	for _, country := range countriesSlice {
		countries[country.Name] = country.ID
	}
}

func ParseCountry(country string) (uuid.UUID, *time.Location, error) {
	id, ok := countries[country]
	if !ok {
		return uuid.Nil, nil, fmt.Errorf("Non-existing country")
	}
	location, err := countryToLocation(country)

	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("Non-existing time location")
	}

	return id, location, nil
}

func MustParseToLocal(date time.Time, country string) time.Time {
	location, _ := countryToLocation(country)

	return date.In(location)
}

func MustParseToLocalByUUID(date time.Time, country uuid.UUID) time.Time {
	var name string
	for coutryName, countryID := range countries {
		if countryID == country {
			name = coutryName
			break
		}
	}
	location, _ := countryToLocation(name)

	return date.In(location)
}

func countryToLocation(country string) (*time.Location, error) {
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

	if tz, ok := mapping[country]; ok {
		return time.LoadLocation(tz)
	}
	return nil, fmt.Errorf("unsupported country: %s", country)
}

func MustGetLocationFromCountryID(id uuid.UUID) *time.Location {
	var name string
	for coutryName, countryID := range countries {
		if countryID == id {
			name = coutryName
			break
		}
	}
	location, _ := countryToLocation(name)

	return location
}

func GetCountries() (map[string]uuid.UUID, uuid.UUID) {
	var countriesCopy = map[string]uuid.UUID{}
	for name, id := range countries {
		countriesCopy[name] = id
	}

	return countriesCopy, countriesCopy["Ukraine"]
}

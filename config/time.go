package config

import (
	"fmt"
	"time"

	"github.com/d3code/uuid"
)

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

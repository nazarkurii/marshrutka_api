package config

import (
	"slices"
)

type Parcel struct {
	Type       parcelType
	Height     uint
	Width      uint
	Length     uint
	Price      uint
	SizeParams []uint
}

type parcelType string

const (
	UndefinedSizeParcelType parcelType = "Oversized"
	UssualParcelType        parcelType = "Ussual"
	Documents               parcelType = "Documents"
)

type ParcelsConfig struct {
	ParcelTypes   []Parcel
	OverSizePrice uint
}

var parcelConfig = GenerateTestParcelsConfig()

func CalulateParcelPrice(widht, height, length uint) uint {
	var position = []uint{widht, height, length}
	slices.Sort(position)

	for _, parcelConfig := range parcelConfig.ParcelTypes {
		if parcelConfig.SizeParams[0] >= position[0] && parcelConfig.SizeParams[1] >= position[1] && parcelConfig.SizeParams[2] >= position[2] {
			return parcelConfig.Price
		}
	}

	return parcelConfig.OverSizePrice
}

func GenerateTestParcelsConfig() ParcelsConfig {
	parcels := []Parcel{
		// 1. Documents
		{
			Type:       Documents,
			Height:     10,
			Width:      10,
			Length:     40,
			Price:      7000,
			SizeParams: []uint{10, 10, 40},
		},
		// 2–7. Ussual parcel types
		{
			Type:       UssualParcelType,
			Height:     20,
			Width:      20,
			Length:     20,
			Price:      8000,
			SizeParams: []uint{20, 20, 20},
		},
		{
			Type:       UssualParcelType,
			Height:     30,
			Width:      30,
			Length:     30,
			Price:      9000,
			SizeParams: []uint{30, 30, 30},
		},
		{
			Type:       UssualParcelType,
			Height:     40,
			Width:      40,
			Length:     40,
			Price:      10000,
			SizeParams: []uint{40, 40, 40},
		},
		{
			Type:       UssualParcelType,
			Height:     60,
			Width:      60,
			Length:     60,
			Price:      11000,
			SizeParams: []uint{60, 60, 60},
		},
		{
			Type:       UssualParcelType,
			Height:     80,
			Width:      80,
			Length:     80,
			Price:      12000,
			SizeParams: []uint{80, 80, 80},
		},
		{
			Type:       UssualParcelType,
			Height:     100,
			Width:      100,
			Length:     100,
			Price:      13000,
			SizeParams: []uint{100, 100, 100},
		},
	}

	return ParcelsConfig{
		ParcelTypes:   parcels,
		OverSizePrice: 15000, // Fallback for parcels larger than 100x100x100
	}
}

package config

import (
	"slices"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type Parcel struct {
	ID             uuid.UUID  `gorm:"type:binary(16);primaryKey"`
	ParcelConfigID uuid.UUID  `gorm:"type:binary(16);index"`
	Type           parcelType `gorm:"type:enum('Oversized','Ussual','Documents,)"`
	Height         uint       `gorm:"type:SMALLINT UNSIGNED;not null"`
	Width          uint       `gorm:"type:SMALLINT UNSIGNED;not null"`
	Length         uint       `gorm:"type:SMALLINT UNSIGNED;not null"`
	Price          uint       `gorm:"type:INT UNSIGNED;not null"`
	SizeParams     []uint     `gorm:"-"`
}

type parcelType string

const (
	UndefinedSizeParcelType parcelType = "Oversized"
	UssualParcelType        parcelType = "Ussual"
	Documents               parcelType = "Documents"
)

type ParcelsConfig struct {
	ID            uuid.UUID `gorm:"type:binary(16);primaryKey"`
	CountryID     uuid.UUID `gorm:"type:binary(16);uniqueIndex"`
	ParcelTypes   []Parcel  `gorm:"foreignKey:ParcelConfigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OverSizePrice uint      `gorm:"type:INT UNSIGNED;not null"`
}

func (p *Parcel) AfterFind(tx *gorm.DB) (err error) {
	p.SizeParams = []uint{p.Width, p.Height, p.Length}
	slices.Sort(p.SizeParams)

	return
}

func (pc ParcelsConfig) CalulatePrice(widht, height, length uint) uint {
	var position = []uint{widht, height, length}
	slices.Sort(position)

	for _, parcelConfig := range pc.ParcelTypes {
		if parcelConfig.SizeParams[0] >= position[0] && parcelConfig.SizeParams[1] >= position[1] && parcelConfig.SizeParams[2] >= position[2] {
			return parcelConfig.Price
		}
	}

	return pc.OverSizePrice
}

func createTestParcelConfigs(countries []Country) []ParcelsConfig {
	var configs []ParcelsConfig

	basePriceCents := uint(7000) // Base price for 20x20x20 cm
	documentsPrice := uint(6000) // Price for documents

	// Define six parcel sizes, increasing by 15 cm each step
	sizes := []struct {
		height uint
		width  uint
		length uint
	}{
		{20, 20, 20}, // Small / Usual
		{35, 35, 35}, // Medium
		{50, 50, 50}, // Large
		{65, 65, 65}, // XL
		{80, 80, 80}, // XXL
		{2, 25, 35},  // Documents
	}

	for i, country := range countries {
		// Price multiplier based on “distance” (index difference from Ukraine)
		distanceFactor := 1.0 + (float64(i) * 0.1)

		var parcels []Parcel
		for j, size := range sizes {
			price := uint(float64(basePriceCents) * distanceFactor)
			pType := UssualParcelType

			if j == len(sizes)-1 { // Last size is documents
				price = documentsPrice
				pType = Documents
			} else {
				// Optional: Increase price for larger parcels proportionally
				price = uint(float64(price) * (1.0 + float64(j)*0.2))
			}

			parcels = append(parcels, Parcel{
				ID:             uuid.New(),
				Type:           pType,
				Height:         size.height,
				Width:          size.width,
				Length:         size.length,
				Price:          price,
				ParcelConfigID: uuid.New(),
			})
		}

		config := ParcelsConfig{
			ID:            uuid.New(),
			CountryID:     country.ID,
			ParcelTypes:   parcels,
			OverSizePrice: uint(float64(basePriceCents)*distanceFactor) + 3000,
		}

		configs = append(configs, config)
	}

	return configs
}

package agrigate

import (
	"maryan_api/internal/models/entity"
	"slices"
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type Bus struct {
	ID                       uuid.UUID         `gorm:"type:binary(16);primaryKey"                       `
	Model                    string            `gorm:"type:varchar(255);not null"                 `
	Images                   []entity.BusImage `gorm:"foreignKey:BusID"                                    `
	RegistrationNumber       string            `gorm:"type:varchar(8);not null;unique"            `
	Year                     int               `gorm:"type:smallint;not null"                     `
	GpsTrackerID             string            `gorm:"type:varchar(255);not null"                 `
	Seats                    []entity.Seat     `gorm:"foreignKey:BusID"                           `
	Structure                []Row             `gorm:"foreignKey:BusID"                                   `
	CreatedAt                time.Time         `gorm:"not null"                                   `
	UpdatedAt                time.Time         `gorm:"not null"                                   `
	DeletedAt                gorm.DeletedAt    `gorm:"index"                                      `
	LuggageCompartmentWidth  uint              `gorm:"type:SMALLINT UNSIGNED;not null"`
	LuggageCompartmentHeight uint              `gorm:"type:SMALLINT UNSIGNED;not null"`
	LuggageCompartmentLength uint              `gorm:"type:INT UNSIGNED;not null"`
	BusAvailability          []entity.BusAvailability
}

type Row struct {
	ID        uuid.UUID             `gorm:"type:binary(16);primaryKey"    `
	BusID     uuid.UUID             `gorm:"type:binary(16), not null"     `
	Number    int                   `gorm:"type:TINYINT; not null"  `
	Positions []entity.SeatPosition `gorm:"foreignKey:RowID"        `
}

func MigrateBus(db *gorm.DB) error {
	return db.AutoMigrate(
		&Bus{},
		&Row{},
	)
}

type CustomerBus struct {
	Model              string                   `json:"model"`
	Images             []string                 `json:"imageURLs"`
	RegistrationNumber string                   `json:"registrationNumber"`
	Year               int                      `json:"year"`
	Structure          [][]ResponseCustomerSeat `json:"structure"`
}

type ResponseCustomerSeat struct {
	ResponseSeat
	Taken bool `json:"taken"`
}

type ResponseSeat struct {
	ID        uuid.UUID `json:"id"`
	Number    int       `json:"number"`
	Type      string    `json:"type"`
	Direction string    `json:"direction"`
}

func (b Bus) ToCustomerBus(takenSeatsIDs []uuid.UUID) CustomerBus {
	var imageUrls = make([]string, len(b.Images))
	for i, image := range b.Images {
		imageUrls[i] = image.Url
	}

	return CustomerBus{
		Model:              b.Model,
		Images:             imageUrls,
		RegistrationNumber: b.RegistrationNumber,
		Year:               b.Year,
		Structure:          b.responseCustomerStructure(takenSeatsIDs),
	}
}

func (b Bus) responseCustomerStructure(takenSeatsIDs []uuid.UUID) [][]ResponseCustomerSeat {
	var structure = make([][]ResponseCustomerSeat, len(b.Structure))

	for _, row := range b.Structure {
		structure[row.Number] = make([]ResponseCustomerSeat, len(row.Positions))
		for _, rowObject := range row.Positions {

			if rowObject.Type == entity.SeatPossitionTypeTable || rowObject.Type == entity.SeatPossitionTypeSpace {
				structure[row.Number][rowObject.Position] = ResponseCustomerSeat{
					ResponseSeat: ResponseSeat{Type: string(rowObject.Type)},
				}
			} else {
				seatIndex := slices.IndexFunc(b.Seats, func(seat entity.Seat) bool {
					return seat.Number == rowObject.SeatNumber
				})

				structure[row.Number][rowObject.Position] = ResponseCustomerSeat{
					ResponseSeat: ResponseSeat{Type: string(b.Seats[seatIndex].Type),
						Number:    b.Seats[seatIndex].Number,
						ID:        b.Seats[seatIndex].ID,
						Direction: string(b.Seats[seatIndex].Direction)},
					Taken: slices.ContainsFunc(takenSeatsIDs, func(id uuid.UUID) bool { return id == b.Seats[seatIndex].ID }),
				}
			}

		}
	}
	return structure
}

func (b Bus) responseStructure() [][]ResponseSeat {
	var structure = make([][]ResponseSeat, len(b.Structure))

	for _, row := range b.Structure {
		structure[row.Number] = make([]ResponseSeat, len(row.Positions))
		for _, rowObject := range row.Positions {

			if rowObject.Type == entity.SeatPossitionTypeTable || rowObject.Type == entity.SeatPossitionTypeSpace {
				structure[row.Number][rowObject.Position] = ResponseSeat{
					Type: string(rowObject.Type),
				}
			} else {
				seatIndex := slices.IndexFunc(b.Seats, func(seat entity.Seat) bool {
					return seat.Number == rowObject.SeatNumber
				})

				structure[row.Number][rowObject.Position] = ResponseSeat{
					Type:      string(b.Seats[seatIndex].Type),
					Number:    b.Seats[seatIndex].Number,
					ID:        b.Seats[seatIndex].ID,
					Direction: string(b.Seats[seatIndex].Direction),
				}
			}

		}
	}
	return structure
}

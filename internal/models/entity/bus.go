package entity

import (
	"encoding/json"
	"time"

	"github.com/d3code/uuid"
)

type Seat struct {
	ID        uuid.UUID     `gorm:"type:binary(16);primaryKey;"                                                       `
	BusID     uuid.UUID     `gorm:"type:binary(16)"                                                                   `
	Number    int           `gorm:"type:tinyint;not null"                                                       `
	Type      seatType      `gorm:"type:enum('Window', 'Single', 'Single-Window', 'Aisle', 'Middle');not null"  `
	Direction seatDirection `gorm:"type:enum('Forward', 'Backward');not null"                                   `
}

type seatType string

const (
	SeatTypeSingle       seatType = "Single"
	SeatTypeSingleWindow seatType = "Single-Window"
	SeatTypeWindow       seatType = "Window"
	SeatTypeAisle        seatType = "Aisle"
	SeatTypeMidlle       seatType = "Middle"
	seatTypeError        seatType = "Error"
)

func defineSeatType(typeStr string) (seatType, bool) {
	switch seatType(typeStr) {
	case SeatTypeSingle:
		return SeatTypeSingle, true
	case SeatTypeSingleWindow:
		return SeatTypeSingleWindow, true
	case SeatTypeWindow:
		return SeatTypeWindow, true
	case SeatTypeAisle:
		return SeatTypeAisle, true
	case SeatTypeMidlle:
		return SeatTypeMidlle, true
	default:
		return seatTypeError, false
	}
}

type seatDirection string

const (
	SeatDirectionForward  seatDirection = "Forward"
	SeatDirectionBackward seatDirection = "Backward"
	seatDirectionError    seatDirection = "Error"
)

func defineSeatDirection(directionStr string) (seatDirection, bool) {
	switch seatDirection(directionStr) {
	case SeatDirectionForward:
		return SeatDirectionForward, true
	case SeatDirectionBackward:
		return SeatDirectionBackward, true
	default:
		return seatDirectionError, false

	}
}

type BusImage struct {
	BusID uuid.UUID `gorm:"type:binary(16);not null"            `
	Url   string    `gorm:"type:varchar(255);not null"    `
}

func (b BusImage) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Url)
}

type SeatPosition struct {
	RowID      uuid.UUID         `gorm:"type:binary(16), not null"                          `
	SeatNumber int               `gorm:"type:TINYINT; not null"                       `
	Type       seatPossitionType `gorm:"type:enum('Space', 'Table', 'Seat');not null" `
	Position   int               `gorm:"type:TINYINT; not null"                       `
}

type seatPossitionType string

const (
	SeatPossitionTypeSpace seatPossitionType = "Space"
	SeatPossitionTypeTable seatPossitionType = "Table"
	SeatPossitionTypeSeat  seatPossitionType = "Seat"
	seatPossitionTypeError seatPossitionType = "Error"
)

func defineSeatPositionType(typeStr string) (seatPossitionType, bool) {
	switch seatPossitionType(typeStr) {
	case SeatPossitionTypeSeat:
		return SeatPossitionTypeSeat, true
	case SeatPossitionTypeTable:
		return SeatPossitionTypeTable, true
	case SeatPossitionTypeSpace:
		return SeatPossitionTypeSpace, true
	default:
		return seatPossitionTypeError, false
	}
}

type BusAvailability struct {
	BusID   uuid.UUID             `gorm:"type:binary(16); not null"                                                  json:"-"`
	Status  busAvailabilityStatus `gorm:"type:enum('Other','Broken','Busy'); not null"                         json:"status"`
	Date    time.Time             `gorm:"not null;default:CURRENT_TIME_STAMP"                                  json:"date"`
	Comment string                `gorm:"type:varchar(500)"                                                    json:"comment;omitempty"`
}

type busAvailabilityStatus string

const (
	BusAvailabilityStatusBroken busAvailabilityStatus = "Broken"
	BusAvailabilityStatusOther  busAvailabilityStatus = "Other"
	BusAvailabilityStatusBusy   busAvailabilityStatus = "Busy"
)

func (ba busAvailabilityStatus) IsValid() bool {
	switch ba {
	case BusAvailabilityStatusBroken, BusAvailabilityStatusOther, BusAvailabilityStatusBusy:
		return true
	default:
		return false
	}
}

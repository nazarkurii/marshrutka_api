package entity

import (
	rfc7807 "maryan_api/pkg/problem"
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type tripStatus string

const (
	TripStatusRegistered        tripStatus = "Registered"
	TripStatusCanceled          tripStatus = "Canceled"
	TripStatusChangedBus        tripStatus = "Changed Bus"
	TripStatusStarted           tripStatus = "Started"
	TripStatusOutboundDone      tripStatus = "Outbound Done"
	TripStatusBreakDown         tripStatus = "Break Down"
	TripStatusBrokenBusFixed    tripStatus = "Broken Bus Fixed"
	TripStatusBrokenBusReplaced tripStatus = "Broken Bus Replaced"
	TripStatusFinished          tripStatus = "Finished"
)

type TripUpdate struct {
	TripID    uuid.UUID  `json:"-"          gorm:"type:binary(16);not null"`
	Status    tripStatus `json:"status"     gorm:"type:enum('Registered','Canceled','Changed Bus','Started','Outbound Done','Break Down','Broken Bus Fixed','Broken Bus Replaced','Finished');not null"`
	CreatedAt time.Time  `json:"createdAt"  gorm:"not null"`
	Comment   string     `json:"comment"    gorm:"type:varchar(500)"`
}

func (tu *TripUpdate) Validate() error {

	switch tu.Status {
	case TripStatusRegistered,
		TripStatusCanceled,
		TripStatusChangedBus,
		TripStatusStarted,
		TripStatusOutboundDone,
		TripStatusBreakDown,
		TripStatusBrokenBusFixed,
		TripStatusBrokenBusReplaced,
		TripStatusFinished:
	default:
		return rfc7807.BadRequest("invalid-trip-status", "Invalid Trip Status Error", "Trip status provided is not valid.")
	}

	return nil
}

func MigrateTrip(db *gorm.DB) error {
	return db.AutoMigrate(

		&TripUpdate{},
	)
}

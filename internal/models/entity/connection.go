package entity

import (
	rfc7807 "maryan_api/pkg/problem"
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type ConnectionUpdate struct {
	ConnectionID uuid.UUID        `json:"-"          gorm:"type:binary(16); not null"`
	CreatedAt    time.Time        `json:"createAt"   gorm:"not null" json:"createAt"`
	Status       connectionStatus `json:"status"     gorm:"type:enum('Registered','Canceled','Sold','Started','Finished','Stopped','Renewed','Could Not Be Finished','Departure Time Changed');not null" `
	Coment       string           `json:"coment"     gorm:"type:varchar(500)"`
}

type connectionStatus string

type StopUpdate struct {
	StopID    uuid.UUID  `gorm:"type:binary(16); not null"                       json:"-" `
	Status    stopStatus `gorm:"type:enum('Confirmed','Missed','Completed')"     json:"status"`
	Comment   string     `gorm:"type:varchar(500)"                               json:"comment"`
	CreatedAt time.Time  `gorm:"not null"                                        json:"createdAt"`
}

type stopStatus string

const (
	RegisteredConnectionStatus           connectionStatus = "Registered"
	ChangedDepartureTimeConnectionStatus connectionStatus = "Departure Time Changed"
	CanceledConnectionStatus             connectionStatus = "Canceled"
	SoldConnectionStatus                 connectionStatus = "Sold"
	StartedConnectionStatus              connectionStatus = "Started"
	FinishedConnectionStatus             connectionStatus = "Finished"
	StoppedConnectionStatus              connectionStatus = "Stopped"
	RenewedConnectionStatus              connectionStatus = "Renewed"
	CouldNotBeFinishConnectionStatus     connectionStatus = "Could Not Be Finished"

	ConfirmedStopStatus stopStatus = "Confirmed"
	MissedStopStatus    stopStatus = "Missed"
	CompletedStopStatus stopStatus = "Completed"
)

func (tu ConnectionUpdate) Validate() error {
	switch tu.Status {
	case RegisteredConnectionStatus,
		ChangedDepartureTimeConnectionStatus,
		CanceledConnectionStatus,
		SoldConnectionStatus,
		StartedConnectionStatus,
		FinishedConnectionStatus,
		StoppedConnectionStatus,
		RenewedConnectionStatus,
		CouldNotBeFinishConnectionStatus:
	default:
		return rfc7807.BadRequest("invalid-connection-status", "Invalid Connection Status Error", "Connection status provided is not valid.")
	}

	return nil
}

func MigrateConnection(db *gorm.DB) error {
	return db.AutoMigrate(
		&ConnectionUpdate{},
		&StopUpdate{},
	)

}

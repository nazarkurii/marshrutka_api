package agrigate

import (
	"maryan_api/internal/models/entity"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Trip struct {
	ID uuid.UUID `gorm:"type:binary(16);primaryKey"                         json:"id"`

	LeadDriver        entity.User   `gorm:"foreignKey:LeadDriverID;references:ID"      `
	LeadDriverID      uuid.NullUUID `gorm:"type:binary(16);unique"                `
	AssistantDriver   entity.User   `gorm:"foreignKey:AssistantDriverID;references:ID" `
	AssistantDriverID uuid.NullUUID `gorm:"type:binary(16);unique"                  `

	BusID uuid.UUID  `gorm:"type:binary(16);not null"`
	Bus   entity.Bus `gorm:"foreignKey:BusID"`

	OutboundConnectionID uuid.UUID  `gorm:"type:binary(16);not null"                           json:"-"`
	OutboundConnection   Connection `gorm:"foreignKey:OutboundConnectionID;references:ID"      json:"outboundConnection"`
	ReturnConnectionID   uuid.UUID  `gorm:"type:binary(16);not null"                           json:"-"`
	ReturnConnection     Connection `gorm:"foreignKey:ReturnConnectionID;references:ID"        json:"returnConnection"`

	Updates []entity.TripUpdate `                                                    json:"updates"`
}

type TripSimplified struct {
	ID                 uuid.UUID            `json:"id"`
	OutboundConnection ConnectionSimplified `json:"outboundConnection"`
	ReturnConnection   ConnectionSimplified `json:"returnConnection"`
}

func MigrateTrip(db *gorm.DB) error {
	return db.AutoMigrate(
		&Trip{},
	)
}

func PreloadTrip() []string {
	return []string{
		clause.Associations,

		"OutboundConnection.Bus",
		"OutboundConnection.Bus.Images",
		"OutboundConnection.Bus.LeadDriver",
		"OutboundConnection.Bus.AssistantUser",
		"OutboundConnection.Bus.Seats",
		"OutboundConnection.Bus.Structure",
		"OutboundConnection.Bus.Structure.Positions",

		"OutboundConnection.ReplacedBus",
		"OutboundConnection.ReplacedBus.Images",
		"OutboundConnection.ReplacedBus.LeadDriver",
		"OutboundConnection.ReplacedBus.AssistantUser",
		"OutboundConnection.ReplacedBus.Seats",
		"OutboundConnection.ReplacedBus.Structure",
		"OutboundConnection.ReplacedBus.Structure.Positions",

		"OutboundConnection.Stops",
		"OutboundConnection.Stops.Ticket",
		"OutboundConnection.Stops.Updates",
		"OutboundConnection.Updates",

		"ReturnConnection.Bus",
		"ReturnConnection.Bus.Images",
		"ReturnConnection.Bus.LeadDriver",
		"ReturnConnection.Bus.AssistantUser",
		"ReturnConnection.Bus.Seats",
		"ReturnConnection.Bus.Structure",
		"ReturnConnection.Bus.Structure.Positions",

		"ReturnConnection.ReplacedBus",
		"ReturnConnection.ReplacedBus.Images",
		"ReturnConnection.ReplacedBus.LeadDriver",
		"ReturnConnection.ReplacedBus.AssistantUser",
		"ReturnConnection.ReplacedBus.Seats",
		"ReturnConnection.ReplacedBus.Structure",
		"ReturnConnection.ReplacedBus.Structure.Positions",

		"ReturnConnection.Stops",
		"ReturnConnection.Stops.Ticket",
		"ReturnConnection.Stops.Updates",
		"ReturnConnection.Updates",
	}
}

package dataStore

import (
	"context"
	"maryan_api/config"
	"maryan_api/internal/entity"
	"maryan_api/pkg/dbutil"
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Parsel interface {
	GetAvailableConnections(ctx context.Context, pagination dbutil.Pagination, luggageVolume uint, fromCountry, toCountry uuid.UUID) ([]entity.Connection, int, error, bool)
}

type parselMysql struct {
	db *gorm.DB
}

func (pds *parselMysql) GetAvailableConnections(ctx context.Context, pagination dbutil.Pagination, luggageVolume uint, fromCountry, toCountry uuid.UUID) ([]entity.Connection, int, error, bool) {

	pagination.Where(
		`(buses.luggage_volume 
       - COALESCE(parcels.luggage_volume, 0) 
       - COALESCE(tickets.luggage_volume, 0)
	- ((bus_seats.number - (COALESCE(passengers.number, 0))) * (? + ?))
     ) >= ? 
     AND connections.departure_country_id = ? 
     AND connections.destination_country_id = ? 
     AND connections.sell_before > ?`,
		config.BackpackVolume,
		config.LargeLuggageVolume,
		luggageVolume,
		fromCountry,
		toCountry,
		config.MustParseToLocalByUUID(time.Now(), fromCountry).UTC(),
	)

	pagination.Order = "connections.departure_time ASC"

	return dbutil.Paginate[entity.Connection](
		ctx,
		pds.db.Table("connections").Preload(clause.Associations).
			Joins("LEFT JOIN buses ON connections.bus_id = buses.id").
			Joins("LEFT JOIN (SELECT SUM(luggage_volume) AS luggage_volume, connection_id FROM parcels GROUP BY connection_id) parcels ON parcels.connection_id = connections.id").
			Joins("LEFT JOIN (SELECT SUM(luggage_volume) AS luggage_volume, connection_id FROM tickets GROUP BY connection_id) tickets ON tickets.connection_id = connections.id").
			Joins("LEFT JOIN (SELECT COUNT(id) AS number, bus_id FROM seats GROUP BY bus_id) bus_seats ON bus_seats.bus_id = buses.id").
			Joins(`LEFT JOIN (
                    SELECT connections.id AS connection_id, COUNT(ticket_seats.seat_id) AS number
                    FROM ticket_seats
                    JOIN tickets ON tickets.id = ticket_seats.ticket_id
                    JOIN connections ON connections.id = tickets.connection_id
                    GROUP BY connections.id
               ) passengers ON passengers.connection_id = connections.id`),
		pagination,
	)
}

func NewParsel(db *gorm.DB) Parsel {
	return &parselMysql{db}
}

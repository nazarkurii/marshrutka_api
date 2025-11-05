package repo

import (
	"context"
	"maryan_api/internal/entity"
	dataStore "maryan_api/internal/infrastructure/persistence"
	"maryan_api/pkg/dbutil"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type Parcel interface {
	GetConnectionByID(ctx context.Context, id uuid.UUID, passengerNumber int) (entity.Connection, []uuid.UUID, uint, error)
	GetAvailableConnections(ctx context.Context, pagination dbutil.Pagination, luggageVolume uint, fromCountry, toCountry uuid.UUID) ([]entity.Connection, int, error, bool)
	// SavePackage(ctx context.Context, ticket *entity.Ticket) error
	// CreatePackageStops(ctx context.Context, paymentSessionID string) error
	// PaymentSucceeded(ctx context.Context, paymentSessionID string) error
	// RemovePassengerStops(ctx context.Context, paymentSessionID string) error
	// GetTickets(ctx context.Context, pagination dbutil.Pagination) ([]entity.Ticket, []entity.Connection, int, error, bool)
}

type parcelRepo struct {
	parsel     dataStore.Parsel
	connection dataStore.Connection
}

func (r *parcelRepo) GetConnectionByID(ctx context.Context, id uuid.UUID, passengerNumber int) (entity.Connection, []uuid.UUID, uint, error) {
	return r.connection.GetByID(ctx, id, passengerNumber)
}

func (r *parcelRepo) GetAvailableConnections(ctx context.Context, pagination dbutil.Pagination, luggageVolume uint, fromCountry, toCountry uuid.UUID) ([]entity.Connection, int, error, bool) {
	return r.parsel.GetAvailableConnections(ctx, pagination, luggageVolume, fromCountry, toCountry)
}

// func (r *ticketRepo) GetTickets(ctx context.Context, pagination dbutil.Pagination) ([]entity.Ticket, []entity.Connection, int, error, bool) {
// 	return r.ticket.GetTickets(ctx, pagination)
// }

// func (r *ticketRepo) DeleteTickets(ctx context.Context, paymentSessionID string) error {
// 	return r.ticket.DeleteTickets(ctx, paymentSessionID)
// }

// func (r *ticketRepo) PaymentSucceeded(ctx context.Context, paymentSessionID string) error {
// 	return r.ticket.PaymentSucceeded(ctx, paymentSessionID)
// }

// func (r *ticketRepo) CreatePassengerStops(ctx context.Context, paymentSessionID string) error {
// 	return r.ticket.CreatePassengerStops(ctx, paymentSessionID)
// }

// func (r *ticketRepo) RemovePassengerStops(ctx context.Context, paymentSessionID string) error {
// 	return r.ticket.RemovePassengerStops(ctx, paymentSessionID)
// }

// func (r *ticketRepo) SaveTicket(ctx context.Context, ticket *entity.Ticket) error {
// 	return r.ticket.Create(ctx, ticket)
// }

func NewParcelRepo(db *gorm.DB) Parcel {
	return &parcelRepo{
		dataStore.NewParsel(db), dataStore.NewConnection(db),
	}
}

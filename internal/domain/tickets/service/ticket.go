package service

import (
	"context"
	"maryan_api/config"
	"maryan_api/internal/domain/tickets/repo"
	"maryan_api/internal/entity"
	"maryan_api/internal/infrastructure/clients/stripe"
	"maryan_api/pkg/auth"
	"maryan_api/pkg/dbutil"
	"maryan_api/pkg/hypermedia"
	rfc7807 "maryan_api/pkg/problem"
	"net/http"
	"slices"
	"time"

	"github.com/d3code/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type Ticket interface {
	Purchase(ctx context.Context, userID uuid.UUID, newTicket entity.NewTicketJSON) (string, error)
	PurchaseFailed(ctx context.Context, sessionID, token string) error
	PurchaseSucceded(ctx context.Context, sessionID, token string) error
	GetTickets(ctx context.Context, paginationStr dbutil.PaginationStr, userID uuid.UUID) ([]entity.CustomerTicket, hypermedia.Links, error)
}

type serviceImpl struct {
	repo repo.Ticket
}

func (s *serviceImpl) GetTickets(ctx context.Context, paginationStr dbutil.PaginationStr, userID uuid.UUID) ([]entity.CustomerTicket, hypermedia.Links, error) {
	pagination, err := paginationStr.ParseWithCondition(dbutil.Condition{"user_id = ?", []any{userID}}, []string{}, "created_at")
	if err != nil {

		return nil, nil, err
	}

	tickets, connections, total, err, empty := s.repo.GetTickets(ctx, pagination)
	if err != nil || empty {
		return nil, nil, err
	}

	var respose = make([]entity.CustomerTicket, len(tickets))

	for i, ticket := range tickets {
		connectionIndex := slices.IndexFunc(connections, func(connection entity.Connection) bool {
			return connection.ID == ticket.ConnectionID
		})

		if connectionIndex == -1 {
			return nil, nil, rfc7807.DB("internal")
		}

		respose[i] = entity.CustomerTicket{
			Ticket:     ticket,
			Connection: connections[connectionIndex].Simplify(),
			Expired:    connections[connectionIndex].DepartureTime.Before(time.Now().UTC()),
		}
	}

	return respose, hypermedia.Pagination(paginationStr, total), nil
}

func (s *serviceImpl) PurchaseFailed(ctx context.Context, sessionID, token string) error {
	_, err := auth.VerifyAccessToken(token, config.PaymentSecretKey())
	if err != nil {
		return rfc7807.Unauthorized("unauthorized", "Unauthorized Error", "Unauthorized Error")
	}
	err = stripe.CancelPaymentIntent(sessionID)
	if err != nil {
		return rfc7807.BadGateway("payment-cancelation",
			"Payment Cancelation Error", err.Error())
	}
	err = s.repo.RemoveStopsAccordingToTickets(ctx, sessionID)
	if err != nil {
		return err
	}
	return s.repo.DeleteTickets(ctx, sessionID)
}

func (s *serviceImpl) PurchaseSucceded(ctx context.Context, sessionID, token string) error {
	_, err := auth.VerifyAccessToken(token, config.PaymentSecretKey())
	if err != nil {
		return rfc7807.Unauthorized("unauthorized", "Unauthorized Error", "Unauthorized Error")
	}
	return s.repo.PaymentSucceeded(ctx, sessionID)
}

func (s *serviceImpl) Purchase(ctx context.Context, userID uuid.UUID, newTicket entity.NewTicketJSON) (string, error) {
	email, phoneNumber, err := newTicket.ParseContaanctInfo()

	connection, takenSeats, err := s.repo.GetConnectionByID(ctx, newTicket.ConnectionID)
	if err != nil {
		return "", err
	} else if connection.DepartureTime.Before(time.Now().UTC()) {
		return "", rfc7807.BadRequest("unavailable-connection0", "Unavailavble Connection Error", "Connection has alredy departed.")
	}

	if len(newTicket.SeatIDs) != len(newTicket.Passengers) {
		return "", rfc7807.BadRequest("seats-passengers", "Seats Passengers Error", "The seats number and the passengers number have to be equal.")
	}

	var totalSeats int
	for _, seat := range connection.Bus.Seats {
		if seat.Number != 0 {
			totalSeats++
		}
	}

	if totalSeats-len(takenSeats) < len(newTicket.SeatIDs) {
		return "", rfc7807.New(http.StatusConflict, "available-seats", "Available Seats Error", "There are no anough available seats for this many passengers")
	}

	for _, seat := range newTicket.SeatIDs {
		if slices.Contains(takenSeats, seat) {
			return "", rfc7807.New(http.StatusConflict, "taken-seat", "Taken Seat Error", seat.String()+" is already taken.")
		}
	}

	pickUpAdressID, err := s.CreateAdress(ctx, newTicket.PickUpAdress, userID, connection.DepartureCountryID)
	if err != nil {
		return "", err
	}

	dropOffAdressID, err := s.CreateAdress(ctx, newTicket.DropOffAdress, userID, connection.DestinationCountryID)
	if err != nil {
		return "", err
	}

	var passengerIDs []uuid.UUID
	for _, passenger := range newTicket.Passengers {
		id, err := s.CreatePassenger(ctx, passenger, userID)
		if err != nil {
			return "", err
		}
		passengerIDs = append(passengerIDs, id)
	}
	token, err := auth.GenerateAccessToken(config.PaymentSecretKey(), jwt.MapClaims{
		"expires": time.Now().Add(time.Minute * 15).Unix(),
	})
	redirectURL, sessionID, err := stripe.CreateStripeCheckoutSession(int64(connection.Price)*int64(len(newTicket.SeatIDs)), token)
	if err != nil {
		return "", rfc7807.BadGateway("payment", "Payment Error", err.Error())
	}

	var tickets = make([]*entity.Ticket, len(passengerIDs))
	for i, passengerID := range passengerIDs {
		ticketID := uuid.New()
		tickets[i] = &entity.Ticket{
			ID:              ticketID,
			UserID:          userID,
			PhoneNumber:     phoneNumber,
			Email:           email,
			ConnectionID:    connection.ID,
			SeatID:          newTicket.SeatIDs[i],
			PassengerID:     passengerID,
			PickUpAdressID:  pickUpAdressID,
			DropOffAdressID: dropOffAdressID,
			TicketPayment: entity.TicketPayment{
				TicketID:  ticketID,
				Price:     connection.Price,
				Method:    entity.PaymentMethodCard,
				SessionID: sessionID,
				Succeeded: false,
			},
		}
	}

	err = s.repo.SaveTickets(ctx, tickets)
	if err != nil {
		return "", err
	}

	err = s.repo.CreateStopsAccordingToTickets(ctx, sessionID)
	if err != nil {
		return "", err
	}
	return redirectURL, nil
}

func (s *serviceImpl) CreateAdress(ctx context.Context, newAdress entity.NewAddress, userID uuid.UUID, countryID uuid.UUID) (uuid.UUID, error) {
	adress := newAdress.ToAddress(countryID)
	err := adress.Prepare(userID)

	if err != nil {
		return uuid.Nil, err
	}

	return adress.ID, s.repo.CreateAdress(ctx, &adress)
}
func (s *serviceImpl) CreatePassenger(ctx context.Context, newPassenger entity.NewPassenger, userID uuid.UUID) (uuid.UUID, error) {
	passenger := newPassenger.Parse()
	params := passenger.Prepare(userID)
	if params != nil {
		return uuid.Nil, rfc7807.BadRequest("passenger-invalid-data", "Passenger Data Error", "Provided data is not valid.", params...)
	}

	return passenger.ID, s.repo.CreatePassenger(ctx, &passenger)
}

func NewTicketService(repo repo.Ticket) Ticket {
	return &serviceImpl{
		repo,
	}
}

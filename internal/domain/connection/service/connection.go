package service

import (
	"context"
	"maryan_api/config"
	"maryan_api/internal/domain/connection/repo"
	"maryan_api/internal/entity"
	dataStore "maryan_api/internal/infrastructure/persistence"
	"maryan_api/pkg/dbutil"
	"maryan_api/pkg/hypermedia"
	rfc7807 "maryan_api/pkg/problem"
	"slices"
	"strconv"
	"time"

	"github.com/d3code/uuid"
)

type AdminConnection interface {
	GetByID(ctx context.Context, id string, passengerNumber string) (entity.Connection, error)
	GetConnections(ctx context.Context, pagination dbutil.PaginationStr, complete string) ([]entity.ConnectionSimplified, hypermedia.Links, error)
	RegisterUpdate(ctx context.Context, update entity.ConnectionUpdate) error
}

type CustomerConnection interface {
	GetByID(ctx context.Context, id string, passengerNumber string) (entity.CustomerConnection, error)
	GetConnections(ctx context.Context, userID uuid.UUID, pagination dbutil.PaginationStr, complete string) ([]entity.CustomerConnection, hypermedia.Links, error)
	FindConnections(ctx context.Context, request entity.FindConnectionsRequestJSON) (entity.FindConnectionsResponse, error)
}

type connectionService struct {
	repo repo.Connection
}

func (c *connectionService) FindConnections(ctx context.Context, requestJSON entity.FindConnectionsRequestJSON) (entity.FindConnectionsResponse, error) {
	request, invalidParams := requestJSON.Parse()

	if invalidParams != nil {
		return entity.FindConnectionsResponse{}, rfc7807.BadRequest("request-data", "Request Data Error", "Provied data is not valid.", invalidParams...)
	}

	found, err := c.repo.FindConnections(ctx, request)
	if err != nil {
		return entity.FindConnectionsResponse{}, err
	}

	var response = entity.FindConnectionsResponse{
		Connections: make([]entity.FoundConnection, len(found.Connections)),
	}

	for i, connection := range found.Connections {
		ticketsLeft := found.TicketsLeft[slices.IndexFunc(found.TicketsLeft, func(ticketsLeft dataStore.TicketsLeft) bool {
			return ticketsLeft.ID == connection.ID
		})]

		response.Connections[i] = entity.FoundConnection{
			ConnectionSimplified: connection.Simplify(),
			TicketsLeft:          int(ticketsLeft.Number),
			Fits:                 int(ticketsLeft.Number)-request.Adults-request.Children-request.Teenagers >= 0,
			Available:            config.MustParseToLocal(time.Now(), connection.DepartureCountry.Name).UTC().Before(connection.SellBefore),
		}
	}

	response.LeftRange = make([]entity.ConnectionsRange, request.Range)
	length := len(found.LeftRange)
	for i := 0; i < request.Range; i++ {
		if i < length {
			response.LeftRange[i] = found.LeftRange[i]
			response.LeftRange[i].Available = !response.LeftRange[i].SellBefore.Before(config.MustParseToLocalByUUID(time.Now(), request.From).UTC())
			response.LeftRange[i].Date = response.LeftRange[i].Date.In(config.MustGetLocationFromCountryID(request.From))
		} else if i == 0 {
			response.LeftRange[i] = entity.ConnectionsRange{
				Date: request.Date.Add(-24 * time.Hour),
			}
		} else {
			response.LeftRange[i] = entity.ConnectionsRange{
				Date: response.LeftRange[i-1].Date.Add(-24 * time.Hour),
			}
		}
	}
	slices.Reverse(response.LeftRange)
	response.RightRange = make([]entity.ConnectionsRange, request.Range)
	length = len(found.RightRange)
	for i := 0; i < request.Range; i++ {
		if i < length {
			response.RightRange[i] = found.RightRange[i]
			response.RightRange[i].Available = !response.RightRange[i].SellBefore.Before(config.MustParseToLocalByUUID(time.Now(), request.From).UTC())
			response.RightRange[i].Date = response.RightRange[i].Date.In(config.MustGetLocationFromCountryID(request.From))
		} else if i == 0 {
			response.RightRange[i] = entity.ConnectionsRange{
				Date: request.Date.Add(24 * time.Hour),
			}
		} else {
			response.RightRange[i] = entity.ConnectionsRange{
				Date: response.RightRange[i-1].Date.Add(24 * time.Hour),
			}
		}
	}
	return response, nil
}

func (c *connectionService) getByID(ctx context.Context, idStr string, passengerNumber string) (entity.Connection, []uuid.UUID, error) {
	passengers, err := strconv.Atoi(passengerNumber)
	if err != nil {
		return entity.Connection{}, nil, rfc7807.BadRequest("parsing", "Parsing Error", err.Error())
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return entity.Connection{}, nil, rfc7807.UUID(err.Error())
	}

	return c.repo.GetByID(ctx, id, passengers)
}

func (c *connectionService) getConnections(ctx context.Context, paginationStr dbutil.PaginationStr, completed string, condition dbutil.Condition) ([]entity.Connection, hypermedia.Links, error) {
	if completed == "false" {
		condition.Where += " connections.arrival_time IS NULL"
	}

	pagination, err := paginationStr.ParseWithCondition(condition, []string{
		"connections.id", "connections.line", "connections.bus_id",
		"buses.model", "buses.registration_number", "buses.driver_id",
		"passenger.name", "passenger.surname",
		"users.first_name", "users.last_name",
	}, "connections.departure_date")

	if err != nil {
		return nil, nil, err
	}

	connections, total, err, empty := c.repo.GetConnections(ctx, pagination)
	if err != nil || empty {
		return nil, nil, err
	}

	return connections, hypermedia.Pagination(paginationStr, total, hypermedia.DefaultParam{"completed", "false", completed}), nil
}

type adminService struct {
	connectionService
	repo repo.Connection
}

type customerService struct {
	connectionService
	repo repo.Connection
}

//-------------------------Interface implementation--------------------------------

func (c *adminService) GetByID(ctx context.Context, idStr string, passengerNumber string) (entity.Connection, error) {

	connection, _, err := c.getByID(ctx, idStr, passengerNumber)
	return connection, err
}

func (c *adminService) GetConnections(ctx context.Context, paginationStr dbutil.PaginationStr, completed string) ([]entity.ConnectionSimplified, hypermedia.Links, error) {

	connections, urls, err := c.getConnections(ctx, paginationStr, completed, dbutil.Condition{})
	if err != nil {
		return nil, nil, err
	}

	var connectionsSimplified = make([]entity.ConnectionSimplified, len(connections))
	for i, connection := range connections {
		connectionsSimplified[i] = connection.Simplify()
	}

	return connectionsSimplified, urls, nil
}

func (c *adminService) RegisterUpdate(ctx context.Context, update entity.ConnectionUpdate) error {
	err := update.Validate()
	if err != nil {
		return err
	}

	return c.repo.RegisterUpdate(ctx, &update)
}

func (c *customerService) GetByID(ctx context.Context, connectionIDStr string, passengersNumber string) (entity.CustomerConnection, error) {

	connection, takedSeatsIDs, err := c.getByID(ctx, connectionIDStr, passengersNumber)

	if err != nil {
		return entity.CustomerConnection{}, err
	}

	luggage := config.GetLoggageConfig()
	return connection.ToCustomer(takedSeatsIDs, int(luggage.Small.Price), int(luggage.Medium.Price), int(luggage.Large.Price)), nil
}

func (c *customerService) GetConnections(ctx context.Context, userID uuid.UUID, paginationStr dbutil.PaginationStr, completed string) ([]entity.CustomerConnection, hypermedia.Links, error) {
	connections, urls, err := c.getConnections(ctx, paginationStr, completed, dbutil.Condition{Where: "users.id = ?", Values: []any{userID}})
	if err != nil {
		return nil, nil, err
	}

	var connectionsCustomer = make([]entity.CustomerConnection, len(connections))
	for i, connection := range connections {
		connectionsCustomer[i] = connection.ToCustomer(nil, 0, 0, 0)
	}

	return connectionsCustomer, urls, nil
}

//Declaration functions

func NewAdminConnection(repo repo.Connection) AdminConnection {
	return &adminService{connectionService{repo}, repo}
}

func NewCustomerConnection(repo repo.Connection) CustomerConnection {
	return &customerService{connectionService{repo}, repo}
}

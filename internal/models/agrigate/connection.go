package agrigate

import (
	"maryan_api/config"
	entity "maryan_api/internal/models/agrigates"
	rfc7807 "maryan_api/pkg/problem"
	"maryan_api/pkg/uuidutil"
	"strconv"
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// =================== CONNECTION DECLARATION ==============
// ============================================================
type Connection struct {
	ID    uuid.UUID `gorm:"type:binary(16);primaryKey"`
	Line  int       `gorm:"type:SMALLINT;not null"`
	Price int       `gorm:"type:MEDIUMINT UNSIGNED;not null"`

	DepartureCountryID   uuid.UUID `gorm:"type:binary(16);not null"`
	DestinationCountryID uuid.UUID `gorm:"type:binary(16);not null"`

	DepartureTime     time.Time `gorm:"not null"`
	ArrivalTime       time.Time `gorm:"not null"`
	EstimatedDuration int       `gorm:"-"`
	GoogleMapsTripURL string    `gorm:"type:varchar(1000);not null"`

	Stops     []Stop
	CreatedAt time.Time
	Updates   []entity.ConnectionUpdate `gorm:"not null"`

	Type connectionType `gorm:"type:enum('Comertial','Special Assignment','Break Down Return','Break Down Replacement');not null"`

	SellBefore             time.Time `gorm:"not null"`
	MaxLuggageParcelWidth  int       `gorm:"type:SMALLINT UNSIGNED;not null"`
	MaxLuggageParcelHeight int       `gorm:"type:SMALLINT UNSIGNED;not null"`
	MaxLuggageParcelLength int       `gorm:"type:SMALLINT UNSIGNED;not null"`
}

func (c *Connection) AfterFind(tx *gorm.DB) (err error) {
	c.EstimatedDuration = int(c.ArrivalTime.Sub(c.DepartureTime).Minutes())
	return
}

type connectionType string

type Stop struct {
	ID           uuid.UUID           `gorm:"type:binary(16);primaryKey"                   json:"id"`
	TicketID     uuid.NullUUID       `gorm:"type:binary(16)"                              json:"-"`
	Ticket       Ticket              `gorm:"foreignKey:TicketID"                          json:"ticket"`
	ParcelID     uuid.NullUUID       `gorm:"type:binary(16)"                              json:"-"`
	Parcel       Parcel              `gorm:"foreignKey:ParcelID"                          json:"package"`
	ConnectionID uuid.UUID           `gorm:"type:binary(16);not null"                     json:"-"`
	Type         stopType            `gorm:"type:enum('Passenger','Parcel')"              json:"type"`
	LocationType stopLocationType    `gorm:"type:enum('Pick-up','Drop-off')"              json:"locationType"`
	Updates      []entity.StopUpdate `gorm:"constraint:OnDelete:CASCADE"                  json:"updates"`
}

type stopType string
type stopLocationType string

const (
	ComertialConnectionType            connectionType = "Comertial"
	SpecialAsignmentConnectionType     connectionType = "Special Asignment"
	BreakDownRetunConnectionType       connectionType = "Break down Return"
	BreakDownReplacementConnectionType connectionType = "Break Down Replacement"

	PassengerStopType stopType = "Passenger"
	ParcelStopType    stopType = "Parcel"

	PickUpStopType  stopLocationType = "Pick-up"
	DropOffStopType stopLocationType = "Drop-off"
)

func PreloadConnection() []string {
	return []string{
		clause.Associations,

		"Bus.Images",
		"Bus.LeadDriver",
		"Bus.AssistantDriver",
		"Bus.Seats",
		"Bus.Structure",
		"Bus.Structure.Positions",
	}
}

// ============================================================
// ============================================================
// /
// /
// /
// /
// /
// /
// =================== CONNECTION SIMPLIFICATION ==============
// ============================================================
type ConnectionSimplified struct {
	ID                   uuid.UUID `json:"id"`
	Price                int       `json:"price"`
	Line                 int       `json:"line"`
	DepartureCountryID   uuid.UUID `json:"departureCountry"`
	DestinationCountryID uuid.UUID `json:"destinationCountry"`
	DepartureTime        time.Time `json:"departureTime"`
	ArrivalTime          time.Time `json:"arrivalTime"`
	EstimatedDuration    int       `json:"estimatedDuration"`
	SellBefore           time.Time `json:"sellBefore"`
}

func (c *Connection) Simplify(departureCountryID, destinationCountryID uuid.UUID) ConnectionSimplified {

	return ConnectionSimplified{
		ID:                   c.ID,
		Price:                c.Price,
		DepartureCountryID:   departureCountryID,
		DestinationCountryID: destinationCountryID,
		DepartureTime:        config.MustParseTimeToLocalByCountryUUID(c.DepartureTime, c.DepartureCountryID),
		ArrivalTime:          config.MustParseTimeToLocalByCountryUUID(c.ArrivalTime, c.DestinationCountryID),
		Line:                 c.Line,
		EstimatedDuration:    c.EstimatedDuration,
		SellBefore:           c.SellBefore,
	}
}

// ============================================================
// ============================================================
// /
// /
// /
// /
// /
// /
// =============== PASSENGER CONNECTIONS REQUEST ==============
// ============================================================
type PassengerConnection struct {
	ConnectionSimplified
	GoogleMapsTripURL string               `json:"googleMapsTripUrl"`
	Bus               CustomerBus          `json:"bus"`
	Stops             []Stop               `json:"stops"`
	Luggage           config.LuggageConfig `json:"luggage"`
}

func (c *Connection) ToPassenger(departureCountryID, destinationCountryID uuid.UUID, bus CustomerBus) PassengerConnection {
	country := config.MustGetCountryByID(c.DepartureCountryID)
	if country.IsMaster {
		country = config.MustGetCountryByID(c.DestinationCountryID)
	}

	return PassengerConnection{
		ConnectionSimplified: c.Simplify(departureCountryID, destinationCountryID),
		GoogleMapsTripURL:    c.GoogleMapsTripURL,
		Bus:                  bus,
		Stops:                c.Stops,
		Luggage:              country.LuggageConfig,
	}
}

type FindPassengerConnectionsRequestJSON struct {
	DepartureCountry   string `json:"from"`
	DestinationCountry string `json:"to"`
	Date               string `json:"date"`
	Adults             string `json:"adults"`
	Children           string `json:"children"`
	Teenagers          string `json:"teenagers"`
	Range              string `json:"range"`
}

type FindPassengerConnectionsRequest struct {
	DepartureCountry   config.Country
	DestinationCountry config.Country
	Date               time.Time
	Adults             int
	Children           int
	Teenagers          int
	Range              int
}

func (r FindPassengerConnectionsRequestJSON) Parse() (FindPassengerConnectionsRequest, error) {
	var invalidParams rfc7807.InvalidParams

	adults, err := strconv.Atoi(r.Adults)
	if err != nil {
		invalidParams.SetInvalidParam("adults", err.Error())
	} else if adults < 0 {
		invalidParams.SetInvalidParam("adults", "cannot be less that 0.")
	}
	teenagers, err := strconv.Atoi(r.Teenagers)
	if err != nil {
		invalidParams.SetInvalidParam("infants", err.Error())
	} else if teenagers < 0 {
		invalidParams.SetInvalidParam("infants", "cannot be less that 0.")
	} else if teenagers < 1 && adults < 1 {
		invalidParams.SetInvalidParam("passengers", "There has to be at least one adult or one teenager.")
	}

	children, err := strconv.Atoi(r.Children)
	if err != nil {
		invalidParams.SetInvalidParam("children", err.Error())
	} else if children < 0 {
		invalidParams.SetInvalidParam("children", "cannot be less that 0.")
	} else if children > 0 && adults < 1 {
		invalidParams.SetInvalidParam("children", "cannot be more than 0 if there is no adult.")
	}

	connectionsRange, err := strconv.Atoi(r.Range)
	if err != nil {
		invalidParams.SetInvalidParam("range", err.Error())
	} else if connectionsRange < 1 {
		invalidParams.SetInvalidParam("range", "cannot be less that 1.")
	}

	fromID, err := uuidutil.Parse(r.DepartureCountry)
	if err != nil {
		invalidParams.SetInvalidParam("from", err.Error())
	}

	toID, err := uuidutil.Parse(r.DestinationCountry)
	if err != nil {
		invalidParams.SetInvalidParam("to", err.Error())
	}

	departureCountry, ok := config.GetCountryByID(fromID)
	if !ok {
		invalidParams.SetInvalidParam("from", "Non-existing country.")
	}

	destinationCountry, ok := config.GetCountryByID(toID)
	if !ok {
		invalidParams.SetInvalidParam("from", "Non-existing country.")
	}

	if invalidParams != nil {
		return FindPassengerConnectionsRequest{}, rfc7807.BadRequest("find-connection", "Find Connection Request Error", "Provided request data is not valid.", invalidParams...)
	}

	date, err := time.ParseInLocation("2006-01-02", r.Date, config.MustGetLocationByCountryID(departureCountry.ID))
	if err != nil {
		invalidParams.SetInvalidParam("date", err.Error())
	}

	return FindPassengerConnectionsRequest{
		DepartureCountry:   departureCountry,
		DestinationCountry: destinationCountry,
		Date:               date,
		Adults:             adults,
		Children:           children,
		Teenagers:          teenagers,
		Range:              connectionsRange,
	}, nil

}

type FindPassengerConnectionsResponse struct {
	Connections []FoundConnection  `json:"connections"`
	LeftRange   []ConnectionsRange `json:"leftRange"`
	RightRange  []ConnectionsRange `json:"rightRange"`
}

type FoundConnection struct {
	ConnectionSimplified
	TicketsLeft int  `json:"ticketsLeft"`
	Fits        bool `json:"fits"`
	Available   bool `json:"available"`
}

type ConnectionsRange struct {
	Date       time.Time `gorm:"column:date" json:"date"`
	Available  bool      `json:"available"`
	SellBefore time.Time `gorm:"column:sellBefore" json:"-"`
	Number     int       `gorm:"column:number" json:"number"`
	MinPrice   int       `gorm:"column:minPrice" json:"minPrice"`
}

// ============================================================
// ============================================================
// /
// /
// /
// /
// /
// /
// =============== PARCEL CONNECTIONS REQUEST =================
// ============================================================
type ConnectionParcel struct {
	ConnectionSimplified
	Price     uint `json:"price"`
	Fits      bool `json:"fits"`
	DayNumber int  `json:"dayNumber"`
	DayMonth  int  `json:"dayMonth"`
}

func (c *Connection) ToParcel(departureCountryID, destinationCountryID uuid.UUID, fits bool, dayNumber, dayMonth int, price uint) ConnectionParcel {
	return ConnectionParcel{
		ConnectionSimplified: c.Simplify(departureCountryID, destinationCountryID),
		Fits:                 fits,
		DayNumber:            dayNumber,
		DayMonth:             dayMonth,
		Price:                price,
	}
}

// ============================================================
// ============================================================
// /
// /
// /
// /
// /
// /
// ======================== MIGRATIONS ========================
// ============================================================
func MigrateConnection(db *gorm.DB) error {
	return db.AutoMigrate(
		&Connection{},
		&Stop{},
	)

}

// ============================================================
// ============================================================

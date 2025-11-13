package entity

import (
	"time"

	"github.com/d3code/uuid"
	"gorm.io/gorm"
)

type TicketSeat struct {
	TicketID uuid.UUID `gorm:"type:binary(16);not null"           json:"-"`
	SeatID   uuid.UUID `gorm:"type:binary(16);not null"           json:"-"`
	Seat     Seat      `gorm:"foreignKey:SeatID"   json:"seat"`
}

type TicketPayment struct {
	TicketID  uuid.UUID     `gorm:"type:binary(16);not null"                                                json:"ticketId"`
	Price     int           `gorm:"type:MEDIUMINT;not null"                                           json:"price"`
	Method    paymentMethod `gorm:"type:enum('Apple Pay','Card','Cash','Google Pay');not null"        json:"method"`
	CreatedAt time.Time     `gorm:"not null"                                                          json:"createdAt"`
	SessionID string        `gorm:"type:varchar(500);not null"                                                          json:"sessionID"`
	Succeeded bool          `gorm:"not null"                                                          json:"succeeded"`
}

type paymentMethod string

const (
	PaymentMethodApplePay  = "Apple Pay"
	PaymentMethodCard      = "Card"
	PaymentMethodCash      = "Cash"
	PaymentMethodGooglePay = "Google Pay"
)

func DefinePaymentMethod(v string) (paymentMethod, bool) {
	switch paymentMethod(v) {
	case PaymentMethodApplePay, PaymentMethodCard, PaymentMethodCash, PaymentMethodGooglePay:
		return paymentMethod(v), true
	default:
		return "", false
	}
}

func MigrateTicket(db *gorm.DB) error {
	return db.AutoMigrate(

		&TicketPayment{},
		&TicketSeat{},
	)

}

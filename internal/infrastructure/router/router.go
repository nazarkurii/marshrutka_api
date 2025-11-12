package router

import (
	adress "maryan_api/internal/domain/adress/transport/http"
	bus "maryan_api/internal/domain/bus/transport/http"
	connection "maryan_api/internal/domain/connection/transport/http"
	"maryan_api/internal/domain/documents"
	parcel "maryan_api/internal/domain/parcel/transport/http"
	passenger "maryan_api/internal/domain/passenger/transport/http"
	ticket "maryan_api/internal/domain/tickets/transport/http"
	trip "maryan_api/internal/domain/trip/transport/http"
	user "maryan_api/internal/domain/user/transport/http"
	"maryan_api/internal/infrastructure/clients/payment"

	ginutil "maryan_api/pkg/ginutils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(s *gin.Engine, db *gorm.DB, client *http.Client, payment payment.Payment) {
	s.Use(ginutil.LogMiddlewear(db))

	passenger.RegisterRoutes(db, s, client, payment)
	user.RegisterRoutes(db, s, client, payment)
	bus.RegisterRoutes(db, s, client, payment)
	adress.RegisterRoutes(db, s, client, payment)
	connection.RegisterRoutes(db, s, client, payment)
	trip.RegisterRoutes(db, s, client, payment)
	ticket.RegisterRoutes(db, s, client, payment)
	documents.RegisterRoutes(db, s, client, payment)
	parcel.RegisterRoutes(db, s, client, payment)
}

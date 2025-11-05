package http

import (
	"maryan_api/internal/domain/parcel/repo"
	"maryan_api/internal/domain/parcel/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, s *gin.Engine, client *http.Client) {
	// customerRouter := ginutil.CreateAuthRouter("/customer", auth.Customer.SecretKey(), s)

	customerHandler := newHandler(service.NewParcelServie(repo.NewParcelRepo(db)))

	//-----------------------Ticket Routes---------------------------------------

	// customerRouter.POST("/connection/purchase-parsel", customerHandler.purchase)
	// customerRouter.GET("/parsels", customerHandler.getTickets)
	s.GET("/connection/available-parsel-dates/:from/:to/:width/:height/:length", customerHandler.FindConnections)
	// s.GET("/connection/purchase-parcel/failed/:id/:token", customerHandler.purchaseFailed)
	// s.GET("/connection/purchase-parcel/succeded/:id/:token", customerHandler.purchaseSucceded)
}

package http

import (
	"maryan_api/internal/domain/tickets/repo"
	"maryan_api/internal/domain/tickets/service"
	"maryan_api/pkg/auth"
	ginutil "maryan_api/pkg/ginutils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, s *gin.Engine, client *http.Client) {
	customerRouter := ginutil.CreateAuthRouter("/customer", auth.Customer.SecretKey(), s)

	customerHandler := newHandler(service.NewTicketService(repo.NewTicketRepo(db)))

	//-----------------------Ticket Routes---------------------------------------

	customerRouter.POST("/connection/purchase-ticket", customerHandler.purchase)
	customerRouter.GET("/tickets", customerHandler.getTickets)
	s.GET("/connection/purchase-ticket/failed/:id/:token", customerHandler.purchaseFailed)
	s.GET("/connection/purchase-ticket/succeded/:id/:token", customerHandler.purchaseSucceded)
}

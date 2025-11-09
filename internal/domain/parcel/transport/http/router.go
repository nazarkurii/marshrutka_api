package http

import (
	"maryan_api/internal/domain/parcel/repo"
	"maryan_api/internal/domain/parcel/service"
	"maryan_api/pkg/auth"
	ginutil "maryan_api/pkg/ginutils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, s *gin.Engine, client *http.Client) {
	customerRouter := ginutil.CreateAuthRouter("/customer", auth.Customer.SecretKey(), s)

	customerHandler := newHandler(service.NewParcelService(repo.NewParcelRepo(db), client))

	s.GET("/connection/available-parcel-dates/:from/:to/:year/:month", customerHandler.findConnections)
	customerRouter.POST("/connection/:id/purchase-parcel", customerHandler.purchase)
	customerRouter.GET("/parcels", customerHandler.getParcels)
	s.GET("/connection/purchase-parcel/failed/:id/:token", customerHandler.purchaseFailed)
	s.GET("/connection/purchase-parcel/succeded/:id/:token", customerHandler.purchaseSucceded)
}

package http

import (
	"context"
	"fmt"
	"maryan_api/internal/domain/parcel/service"
	"maryan_api/internal/entity"
	"maryan_api/pkg/dbutil"
	ginutil "maryan_api/pkg/ginutils"
	"maryan_api/pkg/hypermedia"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type parcelHandler struct {
	service service.Parcel
}

func newHandler(service service.Parcel) *parcelHandler {
	return &parcelHandler{service}
}

func (ch *parcelHandler) FindConnections(ctx *gin.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*10)
	defer cancel()

	req := entity.FindParcelConnectionsRequest{
		From:   ctx.Param("from"),
		To:     ctx.Param("to"),
		Width:  ctx.Param("width"),
		Length: ctx.Param("length"),
		Height: ctx.Param("height"),
	}
	connections, links, err := ch.service.FindConnections(
		ctxWithTimeout,
		req,
		dbutil.PaginationStr{
			Path:     fmt.Sprintf("/connection/available-parsel-dates/%s/%s/%s/%s/%s", req.From, req.To, req.Width, req.Height, req.Length),
			Page:     ctx.DefaultQuery("page", "1"),
			Size:     ctx.DefaultQuery("size", "9"),
			OrderBy:  ctx.DefaultQuery("order_by", "created_at"),
			OrderWay: ctx.DefaultQuery("order_way", "desc"),
			Search:   ctx.DefaultQuery("search", ""),
		},
	)

	if err != nil {
		ginutil.ServiceErrorAbort(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, struct {
		Connections []entity.ConnectionSimplified `json:"connections"`
		Links       hypermedia.Links              `json:"links"`
		ginutil.Response
	}{
		connections,
		links,
		ginutil.Response{
			Message: "The connections has successfuly been found.",
		},
	})
}

// func (p *passengerHandler) purchase(ctx *gin.Context) {
// 	var request entity.NewTicketJSON

// 	err := ctx.ShouldBindJSON(&request)
// 	if err != nil {

// 		ginutil.HandlerProblemAbort(ctx, rfc7807.JSON(err.Error()))
// 		return
// 	}

// 	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*10)
// 	defer cancel()

// 	redirectURL, err := p.service.Purchase(ctxWithTimeout, ctx.MustGet("userID").(uuid.UUID), request)
// 	if err != nil {
// 		ginutil.ServiceErrorAbort(ctx, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, ginutil.Response{
// 		"The purchase procces has started",
// 		hypermedia.Links{
// 			{"redirect", hypermedia.LinkData{
// 				Href:   redirectURL,
// 				Method: "",
// 			}},
// 		},
// 	})
// }

// func (p *passengerHandler) purchaseSucceded(ctx *gin.Context) {
// 	var sessionID = ctx.Param("id")
// 	if sessionID == "" {
// 		ginutil.HandlerProblemAbort(ctx, rfc7807.Unauthorized("unauthorized", "Unauthorized Error", "Unauthorized."))
// 		return
// 	}

// 	var token = ctx.Param("token")
// 	if token == "" {
// 		ginutil.HandlerProblemAbort(ctx, rfc7807.Unauthorized("unauthorized", "Unauthorized Error", "Unauthorized."))
// 		return
// 	}

// 	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*10)
// 	defer cancel()

// 	err := p.service.PurchaseSucceded(ctxWithTimeout, sessionID, token)
// 	if err != nil {
// 		ginutil.ServiceErrorAbort(ctx, err)
// 		return
// 	}

// 	ctx.Redirect(http.StatusFound, config.FrontendURL()+"/profile/tickets")
// }

// func (p *passengerHandler) getTickets(ctx *gin.Context) {

// 	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*10)
// 	defer cancel()

// 	tickets, links, err := p.service.GetTickets(ctxWithTimeout, dbutil.PaginationStr{
// 		"/customer/tickets",
// 		ctx.DefaultQuery("page", "1"),
// 		ctx.DefaultQuery("size", "9"),
// 		ctx.DefaultQuery("order_by", "created_at"),
// 		ctx.DefaultQuery("order_way", "desc"),
// 		ctx.DefaultQuery("search", ""),
// 	}, ctx.MustGet("userID").(uuid.UUID))

// 	if err != nil {
// 		ginutil.ServiceErrorAbort(ctx, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, struct {
// 		ginutil.Response
// 		Links   hypermedia.Links        `json:"links"`
// 		Tickets []entity.CustomerTicket `json:"tickets"`
// 	}{
// 		ginutil.Response{
// 			"The adresses have successfuly beeen found.",
// 			hypermedia.Links{},
// 		},
// 		links,
// 		tickets,
// 	})

// }

// func (p *passengerHandler) purchaseFailed(ctx *gin.Context) {
// 	var sessionID = ctx.Param("id")
// 	if sessionID == "" {
// 		ginutil.HandlerProblemAbort(ctx, rfc7807.Unauthorized("unauthorized", "Unauthorized Error", "Unauthorized."))
// 		return
// 	}

// 	var token = ctx.Param("token")
// 	if token == "" {
// 		ginutil.HandlerProblemAbort(ctx, rfc7807.Unauthorized("unauthorized", "Unauthorized Error", "Unauthorized."))
// 		return
// 	}

// 	err := p.service.PurchaseFailed(ctx, sessionID, token)
// 	if err != nil {
// 		ctx.Redirect(http.StatusFound, config.FrontendURL()+"/internal-server-error")
// 		return
// 	}

// 	ctx.Redirect(http.StatusFound, config.FrontendURL()+"/connection/failed/0/0/0/true")
// }

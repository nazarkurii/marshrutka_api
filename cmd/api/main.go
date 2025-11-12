package main

import (
	"maryan_api/config"
	"maryan_api/internal/infrastructure/clients/payment"
	dataStore "maryan_api/internal/infrastructure/persistence"
	"maryan_api/internal/infrastructure/router"
	"maryan_api/pkg/languages"
	"maryan_api/pkg/timezone"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig("../../.env")
	timezone.Load()

	db := dataStore.Init()
	dataStore.Migrate(db)
	config.LoadCountriesConfig(db)
	config.LoadLuggageConfig(db)
	config.LoadLuggageConfig(db)

	payment := payment.Init()
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Authorization", "Content-Type", "X-Email-Access-Token", "X-Customer-Update-Token", "Content-Language"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	server.Use(languages.GinMiddlewear)
	client := http.DefaultClient
	router.RegisterRoutes(server, db, client, payment)
	server.Static("/imgs", "../../static/images")
	server.GET("", func(ctx *gin.Context) {
		ctx.JSON(
			http.StatusOK, struct {
				Message string
			}{
				"Hello, World!",
			},
		)
	})
	gin.SetMode(gin.ReleaseMode)

	server.Run(":9990")
}

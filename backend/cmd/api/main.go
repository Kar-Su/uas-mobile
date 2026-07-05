package main

import (
	"log"
	"time"

	"github.com/Kar-Su/uas-mobile.git/internal/modules/auth"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/barang"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/satuan_barang"
	sseModule "github.com/Kar-Su/uas-mobile.git/internal/modules/sse"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/tipe_barang"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/user"
	"github.com/Kar-Su/uas-mobile.git/internal/package/env"
	"github.com/Kar-Su/uas-mobile.git/internal/providers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func run(server *gin.Engine) {
	port := env.GetWithDefault[string]("GO_PORT", "8080")

	var serve string
	app := env.GetWithDefault[string]("GO_APP", "localhost")
	if app == "localhost" {
		serve = "0.0.0.0:" + port
	} else {
		serve = ":" + port
	}
	log.Printf("server is running on %s", app)
	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %s", err)
	}
}

func main() {
	injector := do.New()

	server := gin.Default()
	configCors := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cache-Control", "Connection"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	configCors.AllowAllOrigins = true
	server.Use(cors.New(configCors))

	providers.RegisterProviders(injector)

	auth.RegisterRoutes(server, injector)
	user.RegisterRoutes(server, injector)
	tipe_barang.RegisterRoutes(server, injector)
	satuan_barang.RegisterRoutes(server, injector)
	barang.RegisterRoutes(server, injector)
	sseModule.RegisterRoutes(server, injector)

	totalRoutes := len(server.Routes())
	log.Println("Total routes:", totalRoutes)

	run(server)
}

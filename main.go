package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/snipextt/dayer/handler"
	"github.com/snipextt/dayer/internal/cron"
	"github.com/snipextt/dayer/middleware"
	"github.com/snipextt/dayer/storage"
	"github.com/snipextt/dayer/utils"
	clerk_utils "github.com/snipextt/dayer/utils/clerk"
)

func main() {
	godotenv.Load()
	app := fiber.New()

	utils.SetGoogleAuthConfig()
	storage.Init()
	clerk_utils.SetClerk()
	cron.Init()

	app.Use(cors.New())

	api := app.Group("/v1")
	api.Use(middleware.AuthMiddleware)

	api.Post("/onboarding", handler.Onboarding)

	// Get all calendar connections
	api.Get("/calendar/connections", handler.GetConnectedCalendars)

	// Google calendar routes
	api.Get("/calendar/auth/google", handler.GoogleAuthUrl)
	api.Get("/calendar/google", handler.ListAllGoogleCalendarsForConnection)
	api.Post("/calendar/sync/google", handler.SyncGoogleCalendars)

	// Microsoft calendar routes
	api.Get("/calendar/auth/microsoft", handler.MsAuthUrl)

	api.Get("/extension/all", handler.GetExtensions)

	workspace := api.Group("/workspace")

  workspace.Get("/reports", handler.Reports)

	workspace.Get("/", handler.GetCurrentWorkspace)
	workspace.Post("/", handler.CreateWorkspace)

	workspace.Get("/team/:id", handler.GetTeam)
	workspace.Post("/team", handler.CreateTeam)

	workspace.Post("/timedoctor/connect", handler.ConnectTimeDoctor)
	workspace.Post("/timedoctor/company", handler.ConnectTimeDoctorCompany)

	// Callback routes
	callack := api.Group("/callback")
	callack.Get("/oauth/google", handler.GoogleAuthCallback)
	callack.Post("/oauth/microsoft", handler.MsAuthCallback)

	port := os.Getenv("PORT")

	app.Listen("0.0.0.0:" + port)
}

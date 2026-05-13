package main

import (
	"log"
	"os"
	"revoked/cmd/revoked/hooks"
	_ "revoked/migrations"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, falling back to system environment variables")
	}

	app := pocketbase.New()

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASSWORD")

	if adminEmail != "" && adminPass != "" {
		hooks.BindCreateSuperuserAccount(app, adminEmail, adminPass)
	}

	// Run migrations
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: true,
	})

	// Binds hooks and routes
	hooks.BindHooksAndRoutes(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"gitlab.yurtal.tech/company/blitz/business-card/back/cmd/logo"
	"gitlab.yurtal.tech/company/blitz/business-card/back/cmd/parse"
	application "gitlab.yurtal.tech/company/blitz/business-card/back/internal/app"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/config"
)

func main() {
	log.Print("config initializing")
	cfg := config.GetConfig()

	app := application.NewApp(cfg)

	// Register parse command
	app.RootCmd.AddCommand(parse.ParseCommand(app))

	// Register logo command
	app.RootCmd.AddCommand(logo.LogoCommand(app))

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

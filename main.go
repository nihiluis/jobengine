package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/nihiluis/jobengine/api"
	"github.com/nihiluis/jobengine/database"
	"github.com/rs/zerolog/log"
)

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file")
	}

	db, err := database.New(context.Background())
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.RunMigrations()
	if err != nil {
		return err
	}

	queries := database.NewQueries(db)
	api := api.NewAPI(queries)

	return api.Start()
}

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run")
	}
}

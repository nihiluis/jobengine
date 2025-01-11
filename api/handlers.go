package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/nihiluis/jobengine/database"
	"github.com/rs/zerolog/log"
)

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

// API struct holds dependencies for API handlers
type API struct {
	app    huma.API
	router *chi.Mux
}

type internalAPI struct {
	queries database.Queries
}

// NewAPI creates a new instance of API
func NewAPI(queries database.Queries) *API {
	internalAPI := &internalAPI{
		queries: queries,
	}

	router := chi.NewMux()
	humaApp := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

	internalAPI.registerRoutes(humaApp)

	return &API{
		app:    humaApp,
		router: router,
	}
}

// registerRoutes registers all API routes
func (api *internalAPI) registerRoutes(humaApi huma.API) {
	huma.Register(humaApi, huma.Operation{
		OperationID: "get-ping",
		Method:      http.MethodGet,
		Path:        "/api/v1/ping",
		Summary:     "Health check endpoint",
		Description: "Health check endpoint",
		Tags:        []string{"Health"},
	}, api.pingHandler)

	huma.Register(humaApi, huma.Operation{
		OperationID: "get-job",
		Method:      http.MethodGet,
		Path:        "/api/v1/jobs/{id}",
		Summary:     "Get a job",
		Description: "Get a job by ID",
		Tags:        []string{"Jobs"},
	}, api.getJobHandler)

	huma.Register(humaApi, huma.Operation{
		OperationID: "get-jobs",
		Method:      http.MethodGet,
		Path:        "/jobs",
		Summary:     "Get jobs",
		Description: "Get jobs",
		Tags:        []string{"Jobs"},
	}, api.getJobsHandler)

	huma.Register(humaApi, huma.Operation{
		OperationID: "create-job",
		Method:      http.MethodPost,
		Path:        "/api/v1/jobs",
		Summary:     "Create a job",
		Description: "Create a job",
		Tags:        []string{"Jobs"},
	}, api.createJobHandler)

	huma.Register(humaApi, huma.Operation{
		OperationID: "finish-job",
		Method:      http.MethodPost,
		Path:        "/api/v1/jobs/{id}/finish",
		Summary:     "Finish a job",
		Description: "Finish a job",
		Tags:        []string{"Jobs"},
	}, api.finishJobHandler)
}

func (api *API) WriteOpenAPISpec() error {
	b, err := api.app.OpenAPI().DowngradeYAML()
	if err != nil {
		return fmt.Errorf("failed to generate OpenAPI spec: %w", err)
	}

	err = os.WriteFile("openapi.yaml", b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write OpenAPI spec: %w", err)
	}

	return nil

}

// Start starts the HTTP server
func (api *API) Start() error {
	apiAddress := os.Getenv("ADDRESS")
	if apiAddress == "" {
		return errors.New("ADDRESS is not set")
	}

	log.Info().Msgf("Starting server on %s", apiAddress)
	return http.ListenAndServe(apiAddress, api.router)
}

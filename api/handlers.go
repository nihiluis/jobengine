package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nihiluis/jobengine/database"
	"github.com/nihiluis/jobengine/database/queries"
)

type M map[string]any

// API struct holds dependencies for API handlers
type API struct {
	e       *echo.Echo
	queries *database.Queries
}

// NewAPI creates a new instance of API
func NewAPI(queries *database.Queries) *API {
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := &API{
		e:       e,
		queries: queries,
	}

	// Register routes
	api.registerRoutes()

	return api
}

func (api *API) ping(c echo.Context) error {
	return c.JSON(http.StatusOK, M{
		"message": "OK",
	})
}

func (api *API) getJob(c echo.Context) error {
	// Parse job ID from URL parameter
	jobID := c.Param("id")

	// Get job from database
	job, err := api.queries.GetJobByID(c.Request().Context(), jobID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{
			"error": "Failed to fetch job",
		})
	}

	return c.JSON(http.StatusOK, job)
}

func (api *API) getJobs(c echo.Context) error {
	jobStatusStr := c.Param("status")

	var jobStatus queries.JobStatus
	err := jobStatus.Scan(jobStatusStr)

	if err != nil {
		return c.JSON(http.StatusBadRequest, M{
			"error": "Invalid job status",
		})
	}

	jobs, err := api.queries.GetJobsByStatus(c.Request().Context(), jobStatus)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{
			"error": "Failed to fetch jobs",
		})
	}

	return c.JSON(http.StatusOK, jobs)
}

type CreateJobRequest struct {
	JobType string         `json:"job_type"`
	Payload map[string]any `json:"payload"`
	Process bool           `json:"process"`
}

func (api *API) createJob(c echo.Context) error {
	var req CreateJobRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, M{"error": "Invalid request"})
	}

	// Validate required fields
	if req.JobType == "" {
		return c.JSON(http.StatusBadRequest, M{"error": "job_type is required"})
	}

	ctx := c.Request().Context()

	var (
		job *queries.Job
		err error
	)

	if req.Process {
		job, err = api.queries.CreateJobAndProcess(ctx, req.JobType, req.Payload)
	} else {
		job, err = api.queries.CreateJob(ctx, req.JobType, req.Payload)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"error": "Failed to create job"})
	}
	return c.JSON(http.StatusOK, job)
}

type FinishJobRequest struct {
	JobID  string         `json:"job_id"`
	Result map[string]any `json:"result"`
	Status string         `json:"status"`
}

func (api *API) finishJob(c echo.Context) error {
	var req FinishJobRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, M{"error": "Invalid request"})
	}

	err := api.queries.FinishJob(c.Request().Context(), req.JobID, req.Status, req.Result)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"error": "Failed to finish job"})
	}

	return c.JSON(http.StatusOK, M{"message": "Job finished"})
}

// registerRoutes registers all API routes
func (api *API) registerRoutes() {
	api.e.GET("/api/v1/ping", api.ping)
	api.e.GET("/api/v1/job/:id", api.getJob)
	api.e.GET("/api/v1/job/:status", api.getJobs)
	api.e.POST("/api/v1/job", api.createJob)
	api.e.POST("/api/v1/job/finish", api.finishJob)
}

// Start starts the HTTP server
func (api *API) Start(address string) error {
	return api.e.Start(address)
}

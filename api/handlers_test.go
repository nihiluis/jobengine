package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/nihiluis/jobengine/database"
	"github.com/nihiluis/jobengine/database/queries"

	"github.com/stretchr/testify/assert"
)

// StubQueries is a mock implementation of the database.Queries interface
type StubQueries struct {
}

func (s *StubQueries) CreateJob(ctx context.Context, jobType string, payload map[string]any) (*queries.Job, error) {
	return &queries.Job{
		JobType: jobType,
		Status:  queries.JobStatusPending,
	}, nil
}

func (m *StubQueries) CreateJobAndProcess(ctx context.Context, jobType string, payload map[string]any) (*queries.Job, error) {
	return &queries.Job{
		JobType: jobType,
		Status:  queries.JobStatusProcessing,
	}, nil
}

func (m *StubQueries) GetJobByID(ctx context.Context, id string) (*queries.Job, error) {
	return &queries.Job{
		JobType: "test-type",
		Status:  queries.JobStatusProcessing,
	}, nil
}

func (m *StubQueries) GetJobsByStatus(ctx context.Context, status queries.JobStatus) ([]queries.Job, error) {
	return []queries.Job{
		{
			JobType: "test-type",
			Status:  status,
		},
	}, nil
}

func (m *StubQueries) FinishJob(ctx context.Context, jobIDStr string, status string, result map[string]any) error {
	return nil
}

func setupMockApi(t *testing.T, queries database.Queries) humatest.TestAPI {
	internalAPI := &internalAPI{
		queries: queries,
	}

	_, api := humatest.New(t)
	internalAPI.registerRoutes(api)

	return api
}

func TestCreateJobProcessTrue(t *testing.T) {
	humaApi := setupMockApi(t, &StubQueries{})

	req := CreateJobRequestBody{}
	req.JobType = "test-type"
	req.Payload = map[string]any{}
	req.Process = true

	resp := humaApi.Post("/api/v1/jobs", req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var output CreateJobResponseBody
	err := json.NewDecoder(resp.Body).Decode(&output)
	assert.NoError(t, err)
	assert.Equal(t, queries.JobStatusProcessing, output.Job.Status)
}

func TestCreateJobProcessFalse(t *testing.T) {
	humaApi := setupMockApi(t, &StubQueries{})

	req := CreateJobRequestBody{}
	req.JobType = "test-type"
	req.Payload = map[string]any{}
	req.Process = false
	resp := humaApi.Post("/api/v1/jobs", req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var output CreateJobResponseBody
	err := json.NewDecoder(resp.Body).Decode(&output)
	assert.NoError(t, err)
	assert.Equal(t, queries.JobStatusPending, output.Job.Status)
}

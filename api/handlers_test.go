package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/nihiluis/jobengine/database/queries"
	"github.com/nihiluis/jobengine/job"

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

func (m *StubQueries) FinishJob(ctx context.Context, jobIDStr string, status string, message string, result map[string]any) error {
	return nil
}

func setupMockApi(t *testing.T, jobService job.JobService) humatest.TestAPI {
	internalAPI := &internalAPI{
		jobService: jobService,
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
	assert.Equal(t, string(queries.JobStatusProcessing), output.Job.Status)
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
	assert.Equal(t, string(queries.JobStatusPending), output.Job.Status)
}

func TestFinishJob(t *testing.T) {
	humaApi := setupMockApi(t, &StubQueries{})

	req := FinishJobRequestBody{
		JobID:   "test-id",
		Message: "Test completed",
		Result:  map[string]any{"key": "value"},
		Status:  string(queries.JobStatusCompleted),
	}

	resp := humaApi.Post("/api/v1/jobs/finish", req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var output FinishJobResponseBody
	err := json.NewDecoder(resp.Body).Decode(&output)
	assert.NoError(t, err)
	assert.Equal(t, "OK", output.Message)
}

func TestGetJob(t *testing.T) {
	humaApi := setupMockApi(t, &StubQueries{})

	resp := humaApi.Get("/api/v1/jobs/test-id")

	assert.Equal(t, http.StatusOK, resp.Code)

	var output GetJobResponseBody
	err := json.NewDecoder(resp.Body).Decode(&output)
	assert.NoError(t, err)
	assert.Equal(t, "test-type", output.Job.JobType)
	assert.Equal(t, string(queries.JobStatusProcessing), output.Job.Status)
}

func TestGetJobs(t *testing.T) {
	humaApi := setupMockApi(t, &StubQueries{})

	resp := humaApi.Get("/api/v1/jobs/status/processing")

	assert.Equal(t, http.StatusOK, resp.Code)

	var output GetJobsResponseBody
	t.Logf("Response body: %s", resp.Body)
	err := json.NewDecoder(resp.Body).Decode(&output)
	assert.NoError(t, err)
	assert.Len(t, *output.Jobs, 1)
	assert.Equal(t, "test-type", (*output.Jobs)[0].JobType)
	assert.Equal(t, string(queries.JobStatusProcessing), (*output.Jobs)[0].Status)
}

func TestGetJobsInvalidStatus(t *testing.T) {
	humaApi := setupMockApi(t, &StubQueries{})

	resp := humaApi.Get("/api/v1/jobs/status/invalid-status")

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

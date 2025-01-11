package api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nihiluis/jobengine/database/queries"
)

type GetJobsOutput struct {
	Body struct {
		Jobs *[]queries.Job `json:"jobs" doc:"The jobs"`
	}
}

func (api *internalAPI) getJobsHandler(ctx context.Context, input *struct {
	Status string `path:"status" enum:"pending,processing,completed,failed,cancelled,retrying" doc:"The status of the jobs to get"`
}) (*GetJobsOutput, error) {
	var jobStatus queries.JobStatus
	err := jobStatus.Scan(input.Status)

	if err != nil {
		return nil, huma.Error400BadRequest("invalid job status")
	}

	jobs, err := api.queries.GetJobsByStatus(ctx, jobStatus)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to fetch jobs")
	}

	resp := &GetJobsOutput{}
	resp.Body.Jobs = &jobs

	return resp, nil
}

package api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nihiluis/jobengine/database/queries"
)

type CreateJobRequestBody struct {
	JobType string         `json:"jobType" validate:"required" doc:"The type of the job to create"`
	Payload map[string]any `json:"payload" doc:"The payload of the job"`
	Process bool           `json:"process" doc:"Whether the job will process the job immediately"`
}

type CreateJobOutput struct {
	Body CreateJobResponseBody `json:"body"`
}

type CreateJobResponseBody struct {
	Job *queries.Job `json:"job" doc:"The created job" validate:"required"`
}

func (api *internalAPI) createJobHandler(ctx context.Context, input *struct{ Body CreateJobRequestBody }) (*CreateJobOutput, error) {
	var (
		job *queries.Job
		err error
	)

	if input.Body.Process {
		job, err = api.queries.CreateJobAndProcess(ctx, input.Body.JobType, input.Body.Payload)
	} else {
		job, err = api.queries.CreateJob(ctx, input.Body.JobType, input.Body.Payload)
	}

	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create job")
	}

	resp := &CreateJobOutput{}
	resp.Body.Job = job

	return resp, nil
}

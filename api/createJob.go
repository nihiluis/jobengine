package api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nihiluis/jobengine/database/queries"
	"github.com/rs/zerolog/log"
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
	Job *JobOutput `json:"job" doc:"The created job" validate:"required"`
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
		log.Error().Err(err).Msg("failed to create job")
		return nil, huma.Error500InternalServerError("failed to create job")
	}

	mappedJob := &JobOutput{}
	mappedJob.FromQueries(job)

	resp := &CreateJobOutput{}
	resp.Body.Job = mappedJob

	return resp, nil
}

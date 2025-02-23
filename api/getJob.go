package api

import (
	"context"

	"github.com/rs/zerolog/log"
)

type GetJobOutput struct {
	Body GetJobResponseBody
}

type GetJobResponseBody struct {
	Job *JobOutput `json:"job" doc:"The job"`
}

func (api *internalAPI) getJobHandler(ctx context.Context, input *struct {
	JobID string `path:"id" doc:"The ID of the job to get"`
}) (*GetJobOutput, error) {
	// Parse job ID from URL parameter
	jobID := input.JobID

	// Get job from database
	job, err := api.jobService.GetJobByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	resp := &GetJobOutput{}
	resp.Body.Job = &JobOutput{}
	resp.Body.Job.FromQueries(job)

	log.Info().
		Str("id", resp.Body.Job.ID).
		Str("type", resp.Body.Job.JobType).
		Str("status", resp.Body.Job.Status).
		Msg("getJob")

	return resp, nil
}

package api

import (
	"context"
)

type GetJobOutput struct {
	Body struct {
		Job *JobOutput `json:"job" doc:"The job"`
	}
}

func (api *internalAPI) getJobHandler(ctx context.Context, input *struct {
	JobID string `path:"id" doc:"The ID of the job to get"`
}) (*GetJobOutput, error) {
	// Parse job ID from URL parameter
	jobID := input.JobID

	// Get job from database
	job, err := api.queries.GetJobByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	resp := &GetJobOutput{}
	resp.Body.Job = &JobOutput{}
	resp.Body.Job.FromQueries(job)

	return resp, nil
}

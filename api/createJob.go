package api

import (
	"context"
	"time"

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

type JobOutput struct {
	ID        string `json:"id" doc:"The ID of the job"`
	JobType   string `json:"jobType" doc:"The type of the job"`
	Status    string `json:"status" doc:"The status of the job"`
	Payload   string `json:"payload" doc:"The payload of the job"`
	Result    string `json:"result" doc:"The result of the job"`
	CreatedAt string `json:"createdAt" doc:"The creation time of the job"`
}

func (j *JobOutput) FromQueries(q *queries.Job) {
	j.ID = q.ID.String()
	j.JobType = q.JobType
	j.Status = string(q.Status)
	j.Payload = string(q.Payload)
	j.Result = string(q.Result)
	j.CreatedAt = q.CreatedAt.Time.Format(time.RFC3339)
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
		return nil, huma.Error500InternalServerError("failed to create job")
	}

	mappedJob := &JobOutput{}
	mappedJob.FromQueries(job)

	resp := &CreateJobOutput{}
	resp.Body.Job = mappedJob

	return resp, nil
}

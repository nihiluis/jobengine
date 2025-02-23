package api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type FinishJobInput struct {
	Body FinishJobRequestBody
}

type FinishJobRequestBody struct {
	JobID   string         `json:"jobId" validate:"required" doc:"The ID of the job that was finished"`
	Message string         `json:"message" doc:"The message of the job"`
	Result  map[string]any `json:"result" doc:"The result payload of the job"`
	Status  string         `json:"status" validate:"required" doc:"The new status of the job"`
}

type FinishJobOutput struct {
	Body FinishJobResponseBody
}

type FinishJobResponseBody struct {
	Message string `json:"message" example:"OK" doc:"The message"`
}

func (api *internalAPI) finishJobHandler(ctx context.Context, input *FinishJobInput) (*FinishJobOutput, error) {
	err := api.jobService.FinishJob(ctx, input.Body.JobID, input.Body.Status, input.Body.Message, input.Body.Result)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to finish job")
	}

	resp := &FinishJobOutput{}
	resp.Body.Message = "OK"

	return resp, nil
}

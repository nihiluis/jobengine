package api

import (
	"context"
)

type PingOutput struct {
	Body struct {
		Message string `json:"message" example:"OK" doc:"The message to return"`
	}
}

func (api *internalAPI) pingHandler(ctx context.Context, input *struct{}) (*PingOutput, error) {
	resp := &PingOutput{}
	resp.Body.Message = "OK"
	return resp, nil
}

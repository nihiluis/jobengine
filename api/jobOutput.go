package api

import (
	"time"

	"github.com/nihiluis/jobengine/database/queries"
)

type JobOutput struct {
	ID         string `json:"id" doc:"The ID of the job"`
	JobType    string `json:"jobType" doc:"The type of the job"`
	Status     string `json:"status" doc:"The status of the job"`
	Payload    string `json:"payload" doc:"The payload of the job"`
	Result     string `json:"result" doc:"The result of the job"`
	OutMessage string `json:"outMessage" doc:"The message of the job"`
	CreatedAt  string `json:"createdAt" doc:"The creation time of the job"`
}

func (j *JobOutput) FromQueries(q *queries.Job) {
	j.ID = q.ID.String()
	j.JobType = q.JobType
	j.Status = string(q.Status)
	j.Payload = string(q.Payload)
	j.Result = string(q.Result)
	j.OutMessage = q.OutMessage.String
	j.CreatedAt = q.CreatedAt.Time.Format(time.RFC3339)
}

package circleci

import "time"

type PipelineTriggerRequest struct {
	Branch     string                    `json:"branch"`
	Parameters PipelineTriggerParameters `json:"parameters"`
}
type PipelineTriggerParameters struct {
	Action string `json:"GHA_Action"`
}

type PipelineTriggerResponse struct {
	Number    int       `json:"number"`
	State     string    `json:"state"`
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type PipelinesStatusResponse struct {
	NextPageToken interface{}      `json:"next_page_token"`
	Items         []PipelineStatus `json:"items"`
}

type PipelineStatus struct {
	PipelineID string `json:"pipeline_id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
}

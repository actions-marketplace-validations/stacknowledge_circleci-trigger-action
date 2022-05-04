package circleci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	circleCIAPIURL     = "https://circleci.com/api/v2"
	triggerPipelineURL = circleCIAPIURL + "/project/gh/%s/pipeline"
	pipelineStatusURL  = circleCIAPIURL + "/pipeline/%s/workflow"

	StatusFailed   = "failed"
	StatusSuccess  = "success"
	StatusCanceled = "canceled"
)

type CircleCIAPI struct {
	token  string
	client http.Client
}

func NewCircleCIAPI(token string) *CircleCIAPI {
	return &CircleCIAPI{
		token: token,
		client: http.Client{
			Timeout: time.Duration(3 * time.Second),
		},
	}
}

func (api *CircleCIAPI) TriggerPipeline(project, branch, action string) (string, error) {
	payload, err := json.Marshal(PipelineTriggerRequest{
		Branch:     branch,
		Parameters: PipelineTriggerParameters{Action: action},
	})

	if err != nil {
		return "", ErrSerializingPayload
	}

	response, err := api.request(http.MethodPost, fmt.Sprintf(triggerPipelineURL, project), payload)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusCreated {
		return "", ErrTriggeringPipeline
	}

	pipelineTriggerResponse := PipelineTriggerResponse{}
	json.NewDecoder(response.Body).Decode(&pipelineTriggerResponse)

	defer response.Body.Close()
	return pipelineTriggerResponse.ID, nil
}

func (api *CircleCIAPI) GetPipelineStatus(pipelineID string) (*PipelineStatus, error) {
	response, err := api.request(http.MethodGet, fmt.Sprintf(pipelineStatusURL, pipelineID), nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrFetchingPipelineStatus
	}

	pipelineStatus := PipelinesStatusResponse{}
	json.NewDecoder(response.Body).Decode(&pipelineStatus)

	defer response.Body.Close()
	return &pipelineStatus.Items[0], nil
}

func (api *CircleCIAPI) request(method, url string, body []byte) (*http.Response, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header = map[string][]string{
		"Circle-Token": {api.token},
		"Content-Type": {"application/json"},
	}

	response, err := api.client.Do(request)
	if err != nil {
		return nil, ErrRequestingAPIResource
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	return response, nil
}

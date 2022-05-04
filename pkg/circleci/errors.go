package circleci

import "errors"

var (
	ErrSerializingPayload     = errors.New("error serializing api request payload")
	ErrRequestingAPIResource  = errors.New("error requesting circleci api resource")
	ErrTriggeringPipeline     = errors.New("error triggering pipeline on requested project")
	ErrFetchingPipelineStatus = errors.New("error fetching pipeline status on requested project")
	ErrNotFound               = errors.New("provided token is invalid or resource is does not exist")
)

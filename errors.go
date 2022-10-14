package takeout

import "errors"

var (
	FailedToVerifyToken      = errors.New("failed to verify token")
	FailedToGetCloudTemplate = errors.New("failed to get cloud template")
)

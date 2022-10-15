package takeout

import "errors"

var (
	FailedToVerifyToken      = errors.New("failed to verify token")
	FailedToSendEmail        = errors.New("failed to send email")
	FailedToGetCloudTemplate = errors.New("failed to get cloud template")
)

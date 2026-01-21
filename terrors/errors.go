package terrors

import (
	"errors"
)

var (
	MaxRetryReached = errors.New("the maximum number of retries has been reached")
)

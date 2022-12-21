package apperror

import "errors"

var (
	InvalidPrefixLength = errors.New("invalid prefix length")
	InvalidPrefix       = errors.New("invalid prefix")
	InvalidCommandBody  = errors.New("invalid command body")

	CommandNotInterpreted = errors.New("command not interpreted")
)

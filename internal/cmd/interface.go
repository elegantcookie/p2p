package cmd

import (
	"context"
)

type Command interface {
	Type() int
	Execute(ctx context.Context) ([]byte, error)
}

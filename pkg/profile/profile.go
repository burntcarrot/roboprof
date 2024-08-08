package profile

import (
	"context"
	"io"
)

type Profile interface {
	Name() string
	Start(ctx context.Context, f io.Writer) error
	Stop() error
}

var AllProfiles = []Profile{
	&CPUProfile{},
	&MemProfile{},
}

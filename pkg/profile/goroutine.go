package profile

import (
	"context"
	"io"
	"runtime/pprof"
)

type GoroutineProfile struct{}

func (p *GoroutineProfile) Name() string {
	return "goroutineprof"
}

func (p *GoroutineProfile) Start(ctx context.Context, f io.Writer) error {
	err := pprof.Lookup("goroutine").WriteTo(f, 0)
	if err != nil {
		return err
	}
	return nil
}

func (p *GoroutineProfile) Stop() error {
	return nil
}

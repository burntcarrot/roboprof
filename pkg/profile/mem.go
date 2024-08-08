package profile

import (
	"context"
	"io"
	"runtime/pprof"
)

type MemProfile struct{}

func (p *MemProfile) Name() string {
	return "memprof"
}

func (p *MemProfile) Start(ctx context.Context, f io.Writer) error {
	err := pprof.WriteHeapProfile(f)
	if err != nil {
		return err
	}
	return nil
}

func (p *MemProfile) Stop() error {
	return nil
}

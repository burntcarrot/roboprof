package profile

import (
	"context"
	"io"
	"runtime/pprof"
	"time"

	"github.com/burntcarrot/roboprof/pkg/utils"
)

type CPUProfile struct {
	ProfileDuration time.Duration
}

func (p *CPUProfile) Name() string {
	return "cpuprof"
}

func (p *CPUProfile) Start(ctx context.Context, f io.Writer) error {
	err := pprof.StartCPUProfile(f)
	if err != nil {
		return err
	}

	utils.Sleep(p.ProfileDuration, ctx.Done())

	pprof.StopCPUProfile()
	return nil
}

func (p *CPUProfile) Stop() error {
	return nil
}

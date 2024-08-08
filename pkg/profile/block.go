package profile

import (
	"context"
	"io"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/burntcarrot/roboprof/pkg/utils"
)

type BlockProfile struct {
	ProfileDuration time.Duration
}

func (p *BlockProfile) Name() string {
	return "blockprof"
}

func (p *BlockProfile) Start(ctx context.Context, f io.Writer) error {
	runtime.SetBlockProfileRate(1)

	utils.Sleep(p.ProfileDuration, ctx.Done())

	err := pprof.Lookup("block").WriteTo(f, 0)
	if err != nil {
		return err
	}
	runtime.SetBlockProfileRate(0)
	return nil
}

func (p *BlockProfile) Stop() error {
	return nil
}

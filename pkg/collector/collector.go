package collector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/burntcarrot/roboprof/pkg/profile"
	"github.com/burntcarrot/roboprof/pkg/storage/fs"
	"golang.org/x/sync/errgroup"
)

type Collector struct {
	conf Config
	log  *log.Logger

	done chan struct{}
	stop chan struct{}
}

func New(opts ...Option) *Collector {
	collector := &Collector{
		conf: Config{},
		done: make(chan struct{}),
		stop: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(collector)
	}

	if len(collector.conf.profConf.enabledProfilers) == 0 {
		collector.conf.profConf.enabledProfilers = append(collector.conf.profConf.enabledProfilers, &profile.CPUProfile{})
	}

	collector.log = log.New(os.Stdout, "roboprof: ", log.Ldate|log.Ltime)

	if collector.conf.Quiet {
		collector.log = log.New(io.Discard, "", log.Ldate|log.Ltime)
	}

	return collector
}

func Start(opts ...Option) (*Collector, error) {
	collector := New(opts...)

	seen := make(map[string]bool)
	var uniqueProfilers []profile.Profile
	for _, prof := range collector.conf.profConf.enabledProfilers {
		if _, ok := seen[prof.Name()]; !ok {
			seen[prof.Name()] = true
			uniqueProfilers = append(uniqueProfilers, prof)
		}
	}
	collector.conf.profConf.enabledProfilers = uniqueProfilers

	if err := collector.Start(context.Background()); err != nil {
		return nil, err
	}

	return collector, nil
}

func (c *Collector) Start(ctx context.Context) error {
	go c.collect(ctx)
	return nil
}

func (c *Collector) Stop() {
	close(c.stop)
	<-c.done
}

func (c *Collector) serialCollect(ctx context.Context) error {
	for _, prof := range c.conf.profConf.enabledProfilers {
		var buf bytes.Buffer

		err := prof.Start(ctx, &buf)
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("%s_%s.pprof", prof.Name(), time.Now().Format("2006-01-02 15:04:05"))
		store := fs.NewFSStorage(fs.WithDir(c.conf.storageConf.FSStorageConfig.Dir))

		err = store.Write(buf.Bytes(), filename)
		if err != nil {
			return err
		}

		c.log.Printf("wrote profile to %s\n", filepath.Join(c.conf.storageConf.FSStorageConfig.Dir, filename))

		// err = shipper.Ship(ctx, &buf)
		// if err != nil {
		// 	return err
		// }

		buf.Reset()
	}
	return nil
}

func (c *Collector) concurrentCollect(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, prof := range c.conf.profConf.enabledProfilers {
		var buf bytes.Buffer
		prof := prof

		g.Go(func() error {
			err := prof.Start(ctx, &buf)
			if err != nil {
				return err
			}

			filename := fmt.Sprintf("%s_%s.pprof", prof.Name(), time.Now().Format("2006-01-02 15:04:05"))
			store := fs.NewFSStorage(fs.WithDir(c.conf.storageConf.FSStorageConfig.Dir))

			err = store.Write(buf.Bytes(), filename)
			if err != nil {
				return err
			}

			c.log.Printf("wrote profile to %s\n", filepath.Join(c.conf.storageConf.FSStorageConfig.Dir, filename))

			// err = shipper.Ship(ctx, &buf)
			// if err != nil {
			// 	return err
			// }

			return nil
		})
	}

	return nil
}

func (c *Collector) collect(ctx context.Context) {
	defer close(c.done)

	timer := time.NewTimer(0)

	for {
		select {
		case <-c.stop:
			if !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
			var err error
			if c.conf.CollectionType == CollectionSerial {
				err = c.serialCollect(ctx)
			} else {
				err = c.concurrentCollect(ctx)
			}
			if err != nil {
				c.log.Printf("could not collect profiles, err: %v", err)
			}

			timer.Reset(c.conf.tickInterval)
		}
	}
}

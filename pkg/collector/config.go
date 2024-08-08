package collector

import (
	"time"

	"github.com/burntcarrot/roboprof/pkg/profile"
)

type Config struct {
	Quiet          bool
	CollectionType CollectionType
	tickInterval   time.Duration

	profConf    ProfileConfig
	storageConf StorageConfig
}

type CollectionType int

const (
	CollectionSerial CollectionType = iota
	CollectionConcurrent
)

type LogType int

const (
	StdoutLog LogType = iota
	QuietLog
)

type ProfileConfig struct {
	enabledProfilers []profile.Profile
}

type StorageConfig struct {
	FSStorageConfig FSStorageConfig
}

type FSStorageConfig struct {
	Dir string
}

type Option func(c *Collector)

func LogMode(mode LogType) Option {
	return func(c *Collector) {
		if mode == QuietLog {
			c.conf.Quiet = true
		}
	}
}

func CollectionMode(mode CollectionType) Option {
	return func(c *Collector) {
		if mode == CollectionConcurrent {
			c.conf.CollectionType = CollectionConcurrent
		} else {
			c.conf.CollectionType = CollectionSerial
		}
	}
}

func WithStorageConf(sc StorageConfig) Option {
	return func(c *Collector) {
		c.conf.storageConf = sc
	}
}

func WithTickInterval(t time.Duration) Option {
	return func(c *Collector) {
		c.conf.tickInterval = t
	}
}

func WithCPUProfile(d time.Duration) Option {
	return func(c *Collector) {
		profileDuration := 10 * time.Second
		if d != 0 {
			profileDuration = d
		}
		c.conf.profConf.enabledProfilers = append(c.conf.profConf.enabledProfilers, &profile.CPUProfile{ProfileDuration: profileDuration})
	}
}

func WithMemProfile() Option {
	return func(c *Collector) {
		c.conf.profConf.enabledProfilers = append(c.conf.profConf.enabledProfilers, &profile.MemProfile{})
	}
}

func WithBlockProfile(d time.Duration) Option {
	return func(c *Collector) {
		profileDuration := 10 * time.Second
		if d != 0 {
			profileDuration = d
		}
		c.conf.profConf.enabledProfilers = append(c.conf.profConf.enabledProfilers, &profile.BlockProfile{ProfileDuration: profileDuration})
	}
}

func WithGoroutineProfile() Option {
	return func(c *Collector) {
		c.conf.profConf.enabledProfilers = append(c.conf.profConf.enabledProfilers, &profile.GoroutineProfile{})
	}
}

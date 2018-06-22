package conf

import (
	"errors"
	"github.com/gigalixir/eventarbiter/common/flag"
	"github.com/golang/glog"
	"k8s.io/heapster/common/flags"
	"runtime"
	"time"
)

var version string
var commitHash string
var Conf Config

const (
	defaultFlushFrequency  = 5 * time.Second
	defaultMaxProcs        = 4
	defaultEventChanLength = 100000
)

type Config struct {
	LogFlushFrequency time.Duration
	MaxProcs          uint

	Version    string
	CommitHash string

	Source flags.Uri

	Sinks flags.Uris

	FilteredAlertEventReasons flag.List

	EventChanLength int

	Environment flag.Map
}

func (conf Config) Validate() error {
	if err := conf.validateMaxProcs(); err != nil {
		return err
	}

	if err := conf.validateSink(); err != nil {
		return err
	}

	if err := conf.validateSource(); err != nil {
		return err
	}

	return nil
}

func (conf Config) validateMaxProcs() error {
	cpuCapacity := runtime.NumCPU()

	if conf.MaxProcs > uint(cpuCapacity) {
		glog.Errorf(
			"max procs specified exceeds max available cpu number. max: %d, specified: %d",
			cpuCapacity, conf.MaxProcs,
		)
		return errors.New("max procs exceeds available")
	}

	return nil
}

func (conf Config) validateSource() error {
	if conf.Source.Key != "kubernetes" {
		glog.Errorf("source must be kubernetes. supplied %q", conf.Source.Key)
		return errors.New("source must be kubernetes")
	}

	return nil
}

func (conf Config) validateSink() error {
	if len(conf.Sinks) == 0 {
		glog.Errorf("at least one sink must be specified. supported stdout and callback")
		return errors.New("sink is empty")
	}

	for _, sink := range conf.Sinks {
		if sink.Key != "stdout" && sink.Key != "callback" {
			glog.Errorf("sink must be one of stdout and callback. supplied: %q", sink.Key)
			return errors.New("supplied sink is not supported")
		}
	}

	for _, sink := range conf.Sinks {
		if sink.Key == "callback" {
			if sink.Val.String() == "" {
				glog.Errorf("callback sink url can not be empty")
				return errors.New("callback sink callback url is empty")
			}
		}
	}

	return nil
}

func (conf Config) SetMaxProcs() {
	glog.Infof("set max procs to %d", conf.MaxProcs)
	runtime.GOMAXPROCS(int(conf.MaxProcs))
}

func init() {
	Conf = Config{
		LogFlushFrequency: defaultFlushFrequency,
		MaxProcs:          defaultMaxProcs,
		EventChanLength:   defaultEventChanLength,
		Version:           version,
		CommitHash:        commitHash,
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/gigalixir/eventarbiter/cmd/eventarbiter/conf"
	"github.com/gigalixir/eventarbiter/cmd/eventarbiter/signal"
	"github.com/gigalixir/eventarbiter/handler"
	_ "github.com/gigalixir/eventarbiter/handler/node"
	_ "github.com/gigalixir/eventarbiter/handler/npd"
	_ "github.com/gigalixir/eventarbiter/handler/pod"
	"github.com/gigalixir/eventarbiter/models"
	"github.com/gigalixir/eventarbiter/sink/callback"
	"github.com/gigalixir/eventarbiter/sink/stdout"
	"github.com/gigalixir/eventarbiter/source/kubernetes"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"os"
	"strings"
	"sync"
	"time"
)

func mergeCmdFlags(conf *conf.Config) {
	var versionFlag bool

	flag.BoolVar(&versionFlag, "version", false, "version info")
	flag.UintVar(&conf.MaxProcs, "max_procs", conf.MaxProcs,
		"max cpu number arbiter will use. Can not exceeds the number of max cpu cores.")
	flag.DurationVar(&conf.LogFlushFrequency, "log_flush_frequency", conf.LogFlushFrequency,
		"Maximum number of seconds between glog flushes")

	flag.Var(&conf.Source, "source", "source to read event from. Currently only support Kubernetes")
	flag.Var(&conf.Sinks, "sink", "sink(s) that receive event alert")
	flag.Var(&conf.FilteredAlertEventReasons, "event_filter",
		"event alert reasons that should be filtered. Seperated by comma")
	flag.Var(&conf.Environment, "environment", "comma seperated key-value pairs. This map will be"+
		" as the 'Environment' field in callback http request body along side event alert.")

	flag.Parse()

	if versionFlag {
		fmt.Printf("version: %s\ncommit: %s\n", conf.Version, conf.CommitHash)
		os.Exit(0)
	}
}

func flushLogPeriodically(duration time.Duration) {
	glog.Infof("set glog flush interval to %v", duration)

	go func() {
		flusher := time.NewTicker(duration)
		for {
			select {
			case <-flusher.C:
				glog.Flush()
			}
		}
	}()
}

func filterAlertEvent() {
	for _, filteredEvent := range conf.Conf.FilteredAlertEventReasons.Value {
		glog.Infof("filter out %s event alert", filteredEvent)
		delete(handler.EventAlertHandlers, filteredEvent)
	}

	var eventAlertReasonEnabled []string
	for eventAlertReason := range handler.EventAlertHandlers {
		eventAlertReasonEnabled = append(eventAlertReasonEnabled, eventAlertReason)
	}

	glog.Infof("enable event alert for %s", strings.Join(eventAlertReasonEnabled, ","))

	var eventReasons []string
	for eventReason := range handler.EventReasons {
		eventReasons = append(eventReasons, eventReason)
	}

	glog.V(2).Infof("listening Kubernetes event: %s", strings.Join(eventReasons, ","))
}

func StartMain(sinks []models.Sink, eventChan <-chan *api.Event, stopWG *sync.WaitGroup) {
	go func() {
		for event := range eventChan {
			if _, exists := handler.EventReasons[event.Reason]; exists {
				glog.Infof("got %s. %#v", event.Reason, event)

				for alertEventReason, eventHandler := range handler.EventAlertHandlers {
					glog.V(3).Infof("sending event to %s", alertEventReason)
					stopWG.Add(1)

					// range variable eventHandler can not be captured by func literal
					go func(eventHandler models.EventHandler, event *api.Event) {
						defer stopWG.Done()

						eventHandler.HandleEvent(sinks, event)
					}(eventHandler, event)
				}
			}
		}
	}()
}

func main() {
	mergeCmdFlags(&conf.Conf)

	err := conf.Conf.Validate()
	if err != nil {
		glog.Exit("arbiter encounters something that can not be handled internally. Exit")
	}

	conf.Conf.SetMaxProcs()

	flushLogPeriodically(conf.Conf.LogFlushFrequency)

	filterAlertEvent()

	var sinks []models.Sink
	for _, sinker := range conf.Conf.Sinks {
		if sinker.Key == "stdout" {
			sinks = append(sinks, stdout.NewStdout())
			glog.Info("start json sink")

			continue
		}

		if sinker.Key == "callback" {
			sinks = append(sinks, callback.NewCallback(sinker.Val.String()))
			glog.Infof("start callback sink on %s", sinker.Val.String())

			continue
		}

		glog.Exit("unrecognized sink %q. only support stdout and callback", sinker.Key)
	}

	source := kubernetes.MustNewKubernetes(&conf.Conf.Source.Val)

	eventChan := make(chan *api.Event, conf.Conf.EventChanLength)
	source.Start(eventChan)

	stopWG := sync.WaitGroup{}
	StartMain(sinks, eventChan, &stopWG)

	select {
	case exitSignal := <-signal.ExitChan:
		glog.Infof("receive signal %s", exitSignal)
		glog.Flush()

		source.Stop()

		glog.Info("stop sink")
		for {
			if len(eventChan) != 0 {
				time.Sleep(200 * time.Millisecond)
				glog.Infof("waiting for event chan to be empty")
				continue
			}

			break
		}
		stopWG.Wait()
		glog.Info("sink stopped")

		glog.Infoln("stop gracefully")
		glog.Flush()
		os.Exit(0)
	}
}

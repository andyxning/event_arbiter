package handler

import (
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"strings"
)

type failedSync struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewFailedSync() models.EventHandler {
	return failedSync{
		kind:             "POD",
		reason:           events.FailedSync,
		alertEventReason: "pod_failedsync",
	}
}

func (fs failedSync) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == fs.kind && event.Reason == fs.reason {
		var eventAlert = models.PodEventAlert{
			Kind:          event.InvolvedObject.Kind,
			Name:          event.InvolvedObject.Name,
			Namespace:     event.ObjectMeta.Namespace,
			Host:          event.Source.Host,
			Reason:        event.Reason,
			Message:       event.Message,
			LastTimestamp: event.LastTimestamp.Local().String(),
		}

		for _, sink := range sinks {
			sink.Sink(fs.kind, eventAlert)
		}
	}
}

func (fs failedSync) AlertEventReason() string {
	return fs.alertEventReason
}

func (fs failedSync) Reason() string {
	return fs.reason
}

func init() {
	fs := NewFailedSync()
	handler.MustRegisterEventAlertReason(fs.AlertEventReason(), fs)
	handler.RegisterEventReason(fs.Reason())
}
package pod

import (
	"github.com/gigalixir/eventarbiter/cmd/eventarbiter/conf"
	"github.com/gigalixir/eventarbiter/handler"
	"github.com/gigalixir/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"strings"
)

const (
	PodFailedSyncReason = events.FailedSync
)

type failedSync struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewFailedSync() models.EventHandler {
	return failedSync{
		kind:             "POD",
		reason:           PodFailedSyncReason,
		alertEventReason: "pod_failedsync",
	}
}

func (fs failedSync) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == fs.kind && event.Reason == fs.reason {
		var eventAlert = models.PodEventAlert{
			Kind:          strings.ToUpper(event.InvolvedObject.Kind),
			Name:          event.InvolvedObject.Name,
			Namespace:     event.ObjectMeta.Namespace,
			Host:          event.Source.Host,
			Reason:        event.Reason,
			Message:       event.Message,
			LastTimestamp: event.LastTimestamp.Local().String(),
			Environment:   conf.Conf.Environment.Value,
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

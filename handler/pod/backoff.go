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
	PodBackOffReason = events.BackOffStartContainer
)

type backOff struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewBackOff() models.EventHandler {
	return backOff{
		kind:             "POD",
		reason:           PodBackOffReason,
		alertEventReason: "pod_backoff",
	}
}

func (bf backOff) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == bf.kind && event.Reason == bf.reason {
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
			sink.Sink(bf.kind, eventAlert)
		}
	}
}

func (bf backOff) AlertEventReason() string {
	return bf.alertEventReason
}

func (bf backOff) Reason() string {
	return bf.reason
}

func init() {
	bf := NewBackOff()
	handler.MustRegisterEventAlertReason(bf.AlertEventReason(), bf)
	handler.RegisterEventReason(bf.Reason())
}

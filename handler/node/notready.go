package node

import (
	"github.com/gigalixir/eventarbiter/cmd/eventarbiter/conf"
	"github.com/gigalixir/eventarbiter/handler"
	"github.com/gigalixir/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"strings"
)

const (
	NodeNotReadyReason = events.NodeNotReady
)

type notReady struct {
	kind             string
	reason           string
	alertEventReason string
}

func newNotReady() models.EventHandler {
	return notReady{
		kind:             "NODE",
		reason:           NodeNotReadyReason,
		alertEventReason: "node_notready",
	}
}

func (nr notReady) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == nr.kind && event.Reason == nr.reason {
		var eventAlert = models.NodeEventAlert{
			Kind:          strings.ToUpper(event.InvolvedObject.Kind),
			Name:          event.InvolvedObject.Name,
			Namespace:     event.ObjectMeta.Namespace,
			Reason:        event.Reason,
			Message:       event.Message,
			LastTimestamp: event.LastTimestamp.Local().String(),
			Environment:   conf.Conf.Environment.Value,
		}

		for _, sink := range sinks {
			sink.Sink(nr.kind, eventAlert)
		}
	}
}

func (nr notReady) AlertEventReason() string {
	return nr.alertEventReason
}

func (nr notReady) Reason() string {
	return nr.reason
}

func init() {
	nr := newNotReady()
	handler.MustRegisterEventAlertReason(nr.AlertEventReason(), nr)
	handler.RegisterEventReason(nr.Reason())
}

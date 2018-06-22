package npd

import (
	"github.com/gigalixir/eventarbiter/cmd/eventarbiter/conf"
	"github.com/gigalixir/eventarbiter/handler"
	"github.com/gigalixir/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"strings"
)

const (
	ContainerOOMReason = "OOMKilling"
)

type oom struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewOOM() models.EventHandler {
	return oom{
		kind:             "NODE",
		reason:           ContainerOOMReason,
		alertEventReason: "npd_oomkilling",
	}
}

func (om oom) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == om.kind && event.Reason == om.reason {
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
			sink.Sink(om.kind, eventAlert)
		}
	}
}

func (om oom) AlertEventReason() string {
	return om.alertEventReason
}

func (om oom) Reason() string {
	return om.reason
}

func init() {
	om := NewOOM()
	handler.MustRegisterEventAlertReason(om.AlertEventReason(), om)
	handler.RegisterEventReason(om.Reason())
}

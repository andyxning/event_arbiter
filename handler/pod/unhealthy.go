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
	PodUnhealthyReason = events.ContainerUnhealthy
)

type unhealthy struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewUnhealthy() models.EventHandler {
	return unhealthy{
		kind:             "POD",
		reason:           PodUnhealthyReason,
		alertEventReason: "pod_unhealthy",
	}
}

func (uh unhealthy) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == uh.kind && event.Reason == uh.reason {
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
			sink.Sink(uh.kind, eventAlert)
		}
	}
}

func (uh unhealthy) AlertEventReason() string {
	return uh.alertEventReason
}

func (uh unhealthy) Reason() string {
	return uh.reason
}

func init() {
	uh := NewUnhealthy()
	handler.MustRegisterEventAlertReason(uh.AlertEventReason(), uh)
	handler.RegisterEventReason(uh.Reason())
}

package handler

import (
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"strings"
)

type unhealthy struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewUnhealthy() models.EventHandler {
	return unhealthy{
		kind:             "POD",
		reason:           events.ContainerUnhealthy,
		alertEventReason: "pod_unhealthy",
	}
}

func (uh unhealthy) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == uh.kind && event.Reason == uh.reason {
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
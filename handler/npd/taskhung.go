package npd

import (
	"github.com/gigalixir/eventarbiter/cmd/eventarbiter/conf"
	"github.com/gigalixir/eventarbiter/handler"
	"github.com/gigalixir/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"strings"
)

const (
	TaskHungReason = "TaskHung"
)

type taskHung struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewTaskHung() models.EventHandler {
	return taskHung{
		kind:             "NODE",
		reason:           TaskHungReason,
		alertEventReason: "npd_taskhung",
	}
}

func (th taskHung) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == th.kind && event.Reason == th.reason {
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
			sink.Sink(th.kind, eventAlert)
		}
	}
}

func (th taskHung) AlertEventReason() string {
	return th.alertEventReason
}

func (th taskHung) Reason() string {
	return th.reason
}

func init() {
	th := NewTaskHung()
	handler.MustRegisterEventAlertReason(th.AlertEventReason(), th)
	handler.RegisterEventReason(th.Reason())
}

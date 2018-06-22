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
	NodeRebootedReason = events.NodeRebooted
)

type rebooted struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewRebooted() models.EventHandler {
	return rebooted{
		kind:             "NODE",
		reason:           NodeRebootedReason,
		alertEventReason: "node_rebooted",
	}
}

func (rbt rebooted) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == rbt.kind && event.Reason == rbt.reason {
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
			sink.Sink(rbt.kind, eventAlert)
		}
	}
}

func (rbt rebooted) AlertEventReason() string {
	return rbt.alertEventReason
}

func (rbt rebooted) Reason() string {
	return rbt.reason
}

func init() {
	rbt := NewRebooted()
	handler.MustRegisterEventAlertReason(rbt.AlertEventReason(), rbt)
	handler.RegisterEventReason(rbt.Reason())
}

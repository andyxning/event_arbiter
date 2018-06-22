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
	PodKillingReason = events.KillingContainer
)

type killing struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewKilling() models.EventHandler {
	return killing{
		kind:             "POD",
		reason:           PodKillingReason,
		alertEventReason: "pod_killing",
	}
}

func (bf killing) HandleEvent(sinks []models.Sink, event *api.Event) {
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

func (bf killing) AlertEventReason() string {
	return bf.alertEventReason
}

func (bf killing) Reason() string {
	return bf.reason
}

func init() {
	bf := NewKilling()
	handler.MustRegisterEventAlertReason(bf.AlertEventReason(), bf)
	handler.RegisterEventReason(bf.Reason())
}

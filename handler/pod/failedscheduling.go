package pod

import (
	"github.com/gigalixir/eventarbiter/cmd/eventarbiter/conf"
	"github.com/gigalixir/eventarbiter/handler"
	"github.com/gigalixir/eventarbiter/models"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"strings"
)

const (
	PodFailedSchedulingReason = "FailedScheduling"
)

type failedScheduling struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewFailedScheduling() models.EventHandler {
	return failedScheduling{
		kind:             "POD",
		reason:           PodFailedSchedulingReason,
		alertEventReason: "pod_failedscheduling",
	}
}

func (icm failedScheduling) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == icm.kind && event.Reason == icm.reason {
		// Currently in Kubernetes 1.4.0, message for insufficient cpu is proportional with
		// the number of minion machines. This may be very large and longer information is useless.
		// So, we can just truncate it to maxMessageLength.
		// See https://gist.github.com/andyxning/8065dc35889f07073e129bb75a6e57fe for an example.
		if len(event.Message) >= maxMessageLength {
			glog.Warningf("truncate event message. %s", event.Message)
			event.Message = event.Message[:maxMessageLength]
		}

		var eventAlert = models.PodEventAlert{
			Kind:          strings.ToUpper(event.InvolvedObject.Kind),
			Name:          event.InvolvedObject.Name,
			Namespace:     event.ObjectMeta.Namespace,
			Reason:        icm.reason,
			Message:       event.Message,
			LastTimestamp: event.LastTimestamp.Local().String(),
			Environment:   conf.Conf.Environment.Value,
		}

		for _, sink := range sinks {
			sink.Sink(icm.kind, eventAlert)
		}
	}
}

func (icm failedScheduling) AlertEventReason() string {
	return icm.alertEventReason
}

func (icm failedScheduling) Reason() string {
	return icm.reason
}

func init() {
	icm := NewFailedScheduling()
	handler.MustRegisterEventAlertReason(icm.AlertEventReason(), icm)
	handler.RegisterEventReason(icm.Reason())
}

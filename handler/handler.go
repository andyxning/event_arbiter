package handler

import (
	"fmt"
	"github.com/gigalixir/eventarbiter/models"
)

var (
	EventAlertHandlers map[string]models.EventHandler = make(map[string]models.EventHandler)
	EventReasons       map[string]struct{}            = make(map[string]struct{})
)

func MustRegisterEventAlertReason(eventAlertReason string, handler models.EventHandler) {
	if _, exists := EventAlertHandlers[eventAlertReason]; exists {
		panic(fmt.Sprintf("duplicate event alert reason %s", eventAlertReason))
	}
	EventAlertHandlers[eventAlertReason] = handler
}

func RegisterEventReason(eventReason string) {
	if _, exists := EventReasons[eventReason]; exists {
		return
	}
	EventReasons[eventReason] = struct{}{}
}

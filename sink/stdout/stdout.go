package stdout

import (
	"bytes"
	"encoding/json"
	"github.com/gigalixir/eventarbiter/models"
	"github.com/golang/glog"
	"strings"
)

type stdout struct {
	name string
}

func NewStdout() models.Sink {
	return &stdout{
		name: "stdout",
	}
}

func (so *stdout) Name() string {
	return so.name
}

func (so *stdout) Sink(kind string, alert models.EventAlert) {
	switch strings.ToUpper(kind) {
	case "POD":
		if v, ok := alert.(models.PodEventAlert); ok {
			so.sinkPodEvent(v)
			return
		}
		glog.Errorf("associate kind pod with a none pod event alert. %v", alert)
	case "NODE":
		if v, ok := alert.(models.NodeEventAlert); ok {
			so.sinkNodeEvent(v)
			return
		}
		glog.Errorf("associate kind node with a none node event alert. %v", alert)
	default:
		glog.Errorf("unknown event alert. %#v", alert)
	}
}

func (js *stdout) sinkPodEvent(alert models.PodEventAlert) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "    ")

	err := encoder.Encode(alert)
	if err != nil {
		glog.Errorf("encode pod alert event error. %s. %v", err, alert)
	}

	glog.Info(buf.String())
}

func (ss *stdout) sinkNodeEvent(alert models.NodeEventAlert) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "    ")

	err := encoder.Encode(alert)
	if err != nil {
		glog.Errorf("encode node alert event error. %s. %v", err, alert)
	}

	glog.Info(buf.String())
}

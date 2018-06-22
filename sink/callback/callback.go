package callback

import (
	"bytes"
	"encoding/json"
	"github.com/gigalixir/eventarbiter/common"
	"github.com/gigalixir/eventarbiter/models"
	"github.com/golang/glog"
	"strings"
)

type callback struct {
	name string
	url  string
}

func NewCallback(url string) models.Sink {
	return &callback{
		name: "callback",
		url:  url,
	}
}

func (cb *callback) Name() string {
	return cb.name
}

func (cb *callback) Sink(kind string, alert models.EventAlert) {
	switch strings.ToUpper(kind) {
	case "POD":
		if v, ok := alert.(models.PodEventAlert); ok {
			cb.sinkPodEvent(v)
			return
		}
		glog.Errorf("associate kind pod with a none pod event alert. %v", alert)
	case "NODE":
		if v, ok := alert.(models.NodeEventAlert); ok {
			cb.sinkNodeEvent(v)
			return
		}
		glog.Errorf("associate kind node with a none node event alert. %v", alert)
	default:
		glog.Errorf("unknown event alert. %#v", alert)
	}
}

func (cb *callback) sinkPodEvent(alert models.PodEventAlert) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	err := encoder.Encode(alert)
	if err != nil {
		glog.Errorf("encode pod alert event error. %s. %v", err, alert)
	}

	err = common.SendAlert(&buf, cb.url)
	if err != nil {
		glog.Errorf("send pod event alert error. %v", alert)
		return
	}

	glog.Infof("pod event alert sent. %v", alert)
}

func (cb *callback) sinkNodeEvent(alert models.NodeEventAlert) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	err := encoder.Encode(alert)
	if err != nil {
		glog.Errorf("encode node alert event error. %s. %v", err, alert)
	}

	err = common.SendAlert(&buf, cb.url)
	if err != nil {
		glog.Errorf("send node event alert error. %v", alert)
		return
	}

	glog.Infof("node event alert sent. %v", alert)
}

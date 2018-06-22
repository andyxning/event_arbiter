package common

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
)

const (
	defaultTimeout = 12 * time.Second
)

func SendAlert(alert io.Reader, url string) error {
	client := http.DefaultClient
	client.Timeout = defaultTimeout

	req, err := http.NewRequest("POST", url, alert)
	if err != nil {
		glog.Errorf("new http reqeust error. %s", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		glog.Errorf("perform http request error. %s", err)
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("read http response error. %s", err)
		return err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		glog.Errorf("http respond with status code %d", resp.StatusCode)
		return err
	}

	return err
}

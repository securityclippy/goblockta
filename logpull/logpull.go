package logpull

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/securityclippy/goblockta/events"
	"github.com/securityclippy/goblockta/events/user"
	"encoding/json"
	"net/url"
//	"strings"
	"time"
	"github.com/sirupsen/logrus"
	"strings"
	"errors"
	"io"
)



type Manager struct {
	OrgURL string
	APIToken string
	Client http.Client
}

func NewManager (orgURL, apiToken string) (Manager) {
	c := http.Client{}
	m := Manager{
		OrgURL: orgURL,
		APIToken: apiToken,
		Client: c,
	}
	return m
}

func (m Manager) AddRequestHeaders(req *http.Request){
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", m.APIToken))
	req.Header.Add("Accept", "application/json")
}

func (m Manager) AddRequestParam(params *url.Values, key, value string) (*url.Values) {
	if params == nil {
		params := url.Values{}
		params.Add(key, value)
		return &params
	} else {
		params.Add(key, value)
		return params
	}
}

func (m Manager) NewOktaApiRequest(method, apiPath string, params url.Values, postBody io.Reader) (*http.Request, error) {
	//url := path.Join(m.OrgURL, apiPath)
	requestUrl := fmt.Sprintf("%s%s?%s", m.OrgURL, apiPath, params.Encode())
	req, err := http.NewRequest(method, requestUrl, postBody)
	if err != nil {
		return &http.Request{}, err
	}
	m.AddRequestHeaders(req)
	return req, nil
}

func (m Manager) GetSyslogs(since time.Time) ([]events.LogEvent, error) {
	params := m.AddRequestParam(nil, "filter", fmt.Sprintf("eventType eq %s", user.SessionStart))
	params = m.AddRequestParam(params, "since", OktaTime(since))
	//logrus.Infof("Since: %s", since.Format(time.RFC3339Nano))
	req, err := m.NewOktaApiRequest("GET", "/api/v1/logs", *params, nil)
	if err != nil {
		return []events.LogEvent{}, err
	}
	//req, err := m.NewOktaApiRequest("GET", "/api/v1/logs", nil)
	res, err := m.Client.Do(req)
	if err != nil {
		return []events.LogEvent{}, err
	}
	body, err := ReadResponseBody(res)
	if err != nil {
		logrus.Warnf("%+v", string(body))
		return []events.LogEvent{}, err
	}

	logs := []events.LogEvent{}
	err = json.Unmarshal(body, &logs)

	if err != nil {
		if strings.Contains(string(body), "You do not have permission to perform the requested action") {
			return []events.LogEvent{}, errors.New("Invalid Credentials")
		}
		logrus.Warningf("body: %+v", string(body))
		return []events.LogEvent{}, err
	}

	return logs, nil
}

func ReadResponseBody(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return body, nil
}

func OktaTime(t time.Time) (string) {
	//2018-05-01T21:39:15.384Z
	return t.Format("2006-01-02T15:04:05.999Z")
}


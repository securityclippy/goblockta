package conf

import (
	"github.com/securityclippy/goblockta/ipzonemanager"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"encoding/json"
)

var lg = logrus.WithField("service", "config")

type Config struct {
	IPZoneWhitelist []ipzonemanager.IPAddress `json:"ip_zone_whitelist"`
	IPWhitelist []string `json:"ip_whitelist"`
	CountryWhitelist []string `json:"country_whitelist"`
	PerMinuteWarning float64 `json:"per_minute_warning"`
	BanDuration int `json:"ban_duration"`
	OrgURL string `json:"org_url"`
	OktaAPIKey string `json:"okta_api_key"`
	WarnThreshold float64 `json:"warn_threshold"`
	LogToSlack bool `json:"log_to_slack"`
	BlockThreshold float64 `json:"block_threshold"`
	SlackUserName string `json:"slack_user"`
	SlackWebhookURL string `json:"slack_webhook_url"`
	BlockingMode bool `json:"blocking_mode"`
	BlockZoneNames []string `json:"block_zone_names"`

}


func ReadConfig(filePath string) (Config, error) {
	config := Config{}
	infile, err := ioutil.ReadFile(filePath)
	if err != nil {
		lg.Fatal(err)
	}
	err = json.Unmarshal(infile, &config)
	if err != nil {
		lg.Fatal(err)
	}
	return config, nil
}


package ipzonemanager

import (
	"github.com/securityclippy/goblockta/logpull"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"fmt"
	"bytes"
	"strings"
	"time"
	"github.com/pkg/errors"
)

var lg = logrus.WithField("service", "IPZoneManager")


type IPZoneManager struct {
	Manager logpull.Manager
	BlockZones map[string]IPZone
	BlockedIPS map[string]time.Time
	WhitelistedIPS map[string]bool
}

func NewIPZoneManager(manager logpull.Manager, ipwhitelist []string) (IPZoneManager){
	IPWhitelist := make(map[string]bool)
	for _, ip := range ipwhitelist {
		IPWhitelist[ip] = true
	}
	return IPZoneManager{
		Manager: manager,
		BlockedIPS: make(map[string]time.Time),
		WhitelistedIPS: IPWhitelist,
		BlockZones: make(map[string]IPZone),
	}

}


func (ipzm IPZoneManager) GetIPZones() ([]IPZone, error) {
	path := "/api/v1/zones"
	req, err := ipzm.Manager.NewOktaApiRequest("GET", path, nil, nil)
	if err != nil {
		lg.Error(err)
		return []IPZone{}, err
	}
	res, err := ipzm.Manager.Client.Do(req)
	if err != nil {
		lg.Error(err)
		return []IPZone{}, err
	}
	body, err := logpull.ReadResponseBody(res)
	if err != nil {
		lg.Error(err)
		return []IPZone{}, err
	}
	zones := []IPZone{}
	if strings.Contains(string(body), "errorCode") {
		return []IPZone{}, errors.New(string(body))
	}
	err = json.Unmarshal(body, &zones)
	if err != nil {
		return []IPZone{}, err
	}
	return zones, nil
}

func (ipzm IPZoneManager) GetZone(zoneID string) (IPZone, error) {
	apiPath := fmt.Sprintf("/api/v1/zones/%s", zoneID)
	req, err := ipzm.Manager.NewOktaApiRequest("GET", apiPath, nil, nil)
	if err != nil {
		return IPZone{}, err
	}
	res, err := ipzm.Manager.Client.Do(req)
	if err != nil {
		return IPZone{}, err
	}
	ipz := IPZone{}
	body, err := logpull.ReadResponseBody(res)
	if err != nil {
		return IPZone{}, err
	}
	err = json.Unmarshal(body, &ipz)
	if err != nil {
		return IPZone{}, err
	}
	return ipz, nil
}

func (ipzm IPZoneManager) UpdateIPzone(zoneID string, zoneObject IPZone) (error) {
	apiPath := fmt.Sprintf("/api/v1/zones/%s", zoneID)
	js, err := json.Marshal(zoneObject)
	if err != nil {
		lg.Error(err)
		return err
	}
	req, err := ipzm.Manager.NewOktaApiRequest("PUT", apiPath, nil, bytes.NewReader(js))
	if err != nil {
		return err
	}
	res, err := ipzm.Manager.Client.Do(req)
	if err != nil {
		return err
	}
	body, err := logpull.ReadResponseBody(res)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		lg.Errorf("Return code: %d", res.StatusCode)
	}
	lg.Info(string(body))
	return nil
}

func (ipzm IPZoneManager) BlockCIDR(cidrIP string, zone IPZone) (error) {
	if s := strings.Split(cidrIP, "/"); len(s) <= 1 {
		cidrIP = fmt.Sprintf("%s/32", cidrIP)
		lg.Infof("Adding cidr notation, you bum")
	}
	zone.Gateways = append(zone.Gateways, IPAddress{
		Type: "CIDR",
		Value: cidrIP,
	})
	err := ipzm.UpdateIPzone(zone.ID, zone)
	if err != nil {
		return err
	}
	return nil
}

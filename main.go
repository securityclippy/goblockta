package main

import (
	"github.com/securityclippy/goblockta/logpull"
	"github.com/sirupsen/logrus"
	"github.com/securityclippy/goblockta/ipmanager"
	"time"
	//"github.com/securityclippy/badgerwrapper"
	"github.com/securityclippy/goblockta/ruleprocessor"
	"github.com/securityclippy/goblockta/ipzonemanager"
	"flag"
	"github.com/securityclippy/goblockta/conf"
	"github.com/securityclippy/goslacker"
	"github.com/securityclippy/goslacker/slackclient"
		"github.com/gin-gonic/gin/json"
	"fmt"
)


var lg = logrus.WithField("service", "main")
//import "flag"

//func GetLogs (db badgerwrapper.Badger, m logpull.Manager, ipm ipmanager.IPManager) {
func GetLogs (db *ipmanager.IpDB, m logpull.Manager, ipm ipmanager.IPManager) {
	for {
		timeCursor := db.TimeCursor
		/*if err != nil {
			t1 := time.Now().AddDate(0, 0, -1)
			logrus.Infof("setting time: %+v", logpull.OktaTime(t1))
			err = ipm.UpdateTimeCursor(t1)
			if err != nil {
				logrus.Error(err)
			}
		}*/
		//timeCursor, err = ipm.GetTimeCursor()
		var backOff bool
		//logrus.Warningf("starting from %+v", timeCursor)
		logs, err := m.GetSyslogs(timeCursor.Add(time.Millisecond * 1))
		if err != nil {
			logrus.Error(err)
		}

		if len(logs) <= 20 {
			backOff = true
		}

		//logrus.Infof("parsing %d logs", len(logs))
		for _, log := range logs {
			t1, err := time.Parse(time.RFC3339, log.Published)
			if t1.Sub(timeCursor) < 0 {
				//logrus.Info(t1.Sub(timeCursor))
				continue
			} else {
				timeCursor = t1
				//logrus.Warningf("setting new time to: %+v", t1)
				//err = ipm.DB.UpdateTimeCursor(t1)
				db.TimeCursor = t1
				if err != nil {
					logrus.Error(err)
				}
				//time.Sleep(time.Second * 2)
				//tc, err := ipm.DB.GetTimeCursor()
				//logrus.Infof("TimeCursor: %+v", tc)
				//time.Sleep(time.Second * 3)
				err = ipm.ParseLog(log)
				if err != nil {
					logrus.Error(err)
				}
			}

			//logrus.Info(log.EventType)
			//logrus.Infof(log.Published)
		}
		if backOff {
			//logrus.Infof("backing off to reduce load...")
			time.Sleep(time.Second * 45)
			backOff = false
		}
		time.Sleep(time.Second * 1)

	}

}

//func ProcessLogs (db badgerwrapper.Badger, m logpull.Manager, ipm ipmanager.IPManager) {
func ProcessLogs (db *ipmanager.IpDB, ipzm ipzonemanager.IPZoneManager, sc slackclient.SlackClient, config conf.Config) {
	ll := logrus.WithField("Func", "ProcessLogs")
	for {
		rp := ruleprocessor.NewRuleProcessor(db)
		vals, err := rp.ReadDB()
		if err != nil {
			ll.Error(err)
		}
		for _, j := range vals {
			//rate, err := rp.TotalFailureRate(j.IP)
			if err != nil {
				ll.Error(err)
				continue
			}
			whitelisted, ok := ipzm.WhitelistedIPS[j.IP]
			if !whitelisted && !ok && j.Failure > 0  {
				_, blocked := ipzm.BlockedIPS[j.IP]
				perMinute := rp.FailureRatePerMinute(j, time.Minute * 15)
				if perMinute > config.WarnThreshold && !blocked {
					timeFailureRate, duration := rp.FailureRateOverTime(j, time.Minute * 15)
					logrus.Infof("%s time failure rate: %f, over %+v", j.IP, timeFailureRate, duration)
					//logrus.Infof("%s failures per minute: %+v", j.IP, perMinute)
					ans := rp.FailurePerInterval(j, 5, 5)
					for k, v := range ans {
						ll.Infof("Failures per minute for %+v:  %f", k, v)
					}
					if perMinute >= config.BlockThreshold {
						if !blocked {
							// if actual blocking enabled
							if config.BlockingMode {
								lg.Warningf("Blocking %s", j.IP)
								s := fmt.Sprintf("blocking %s\n, Failure Rate: %f", j.IP, perMinute)
								if config.LogToSlack {
									msg := goslacker.NewAttachmentMessage(config.SlackUserName, s, "New IP Block", "")
									err = sc.PostToIncomingWebHook(msg)
									if err != nil {
										lg.Error(err)
									}

								}
								for _, z := range ipzm.BlockZones {
									ipzm.BlockCIDR(j.IP, z)
								}
							} else {
								// testing mode, don't actually block things
								lg.Warningf("Blocking %s", j.IP)
								s := fmt.Sprintf("blocking %s\n, Failure Rate: %f", j.IP, perMinute)
								if config.LogToSlack {
									msg := goslacker.NewAttachmentMessage(config.SlackUserName, s, "New IP Block", "")
									err = sc.PostToIncomingWebHook(msg)
									if err != nil {
										lg.Error(err)
									}

								}

							}
							ipzm.BlockedIPS[j.IP] = time.Now()
						}
					}
				}
			}
		}
		time.Sleep(time.Second * 20)
		if len(ipzm.BlockedIPS) > 0 {
			lg.Infof("blocked IPs: %+v", ipzm.BlockedIPS)
		}

	}
}

//func Run(db badgerwrapper.Badger, m logpull.Manager, ipm ipmanager.IPManager) {
func Run(db *ipmanager.IpDB, m logpull.Manager, ipm ipmanager.IPManager, ipzm ipzonemanager.IPZoneManager, sc slackclient.SlackClient, config conf.Config) {
		go GetLogs(db, m, ipm)
		go ProcessLogs(db, ipzm, sc, config)
		for {
			time.Sleep(time.Second * 10)
		}
}

func main() {

	confFile := flag.String("c", "/var/config.json", "config file")


	flag.Parse()
	config, err := conf.ReadConfig(*confFile)
	if err != nil {
		lg.Fatal(err)
	}
	sc := slackclient.NewSlackClient(config.SlackWebhookURL)
	startMsg := goslacker.NewAttachmentMessage(config.SlackUserName, "Starting GoBlockta", "", "")
	err =sc.PostToIncomingWebHook(startMsg)
	if err != nil {
		lg.Fatal(err)
	}
	db := ipmanager.NewIpDB()
	m := logpull.NewManager(config.OrgURL, config.OktaAPIKey)
	ipm := ipmanager.NewIPManager(db)
	ipzm := ipzonemanager.NewIPZoneManager(m, config.IPWhitelist)
	zones, err := ipzm.GetIPZones()
	if err != nil {
		logrus.Fatal(err)
	}
	for _, zone := range zones {
		lg.Infof("Found IP Zone: %s", zone.Name)
		js, _ := json.MarshalIndent(zone.Gateways, "", "  ")
		lg.Debug("IP Gateways for: %s \n, %s", zone.Name, string(js))
	}
	for _, zone := range zones {
		// get our block zones
		if stringInList(zone.Name, config.BlockZoneNames) {
			ipzm.BlockZones[zone.Name] = zone
			lg.Infof("adding %s to list of zones used for blocking", zone.Name)
		}
	}
	Run(db, m, ipm, ipzm, sc, config)
}

func stringInList(s string, sList []string) bool {
	for _, i := range sList {
		if s == i {
			return true
		}
	}
	return false
}
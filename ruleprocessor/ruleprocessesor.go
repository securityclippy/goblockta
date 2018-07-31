package ruleprocessor

import (
	"github.com/securityclippy/badgerwrapper"
	"github.com/sirupsen/logrus"
	"github.com/securityclippy/goblockta/ipmanager"
	//"encoding/json"
	"time"
)

func derp() {
	//db, _ := badgerwrapper.NewBadgerDB()
	//db.Get()

}

type ProcessorDB interface {
	SearchPrefix(prefix string) ([]badgerwrapper.KVP, error)
	Get(key []byte) ([]byte, error)
}

type RuleProcessor struct {
	DB ipmanager.IpDB
	//DB ProcessorDB
}

func NewRuleProcessor(db *ipmanager.IpDB) (RuleProcessor) {
	return RuleProcessor{
		DB: *db,
	}
}


func (rp RuleProcessor) ReadDB() ([]ipmanager.IPLog, error) {
	entries, err := rp.DB.SearchPrefix("")
	if err != nil {
		logrus.Error(err)
		return []ipmanager.IPLog{}, err
	}
	logs := make([]ipmanager.IPLog, len(entries))
	for _, e := range entries {
		logs = append(logs, e)
	}
	/*for i, e := range entries {
		l := ipmanager.IPLog{}
		if (string(e.Key) == "timecursor") {
			continue
		}
		err := json.Unmarshal(e.Value, &l)
		if err != nil {
			logrus.Errorf("Erroring key: %s, erroring val: %s", string(e.Key), string(e.Value))
			return []ipmanager.IPLog{}, err
		}
		logs[i] = l

		//logrus.Infof("ip: %s\n failures: %d\n success: %d\n", l.IP, l.Failure, l.Success)
	}*/
	return logs, nil
}

func (rp RuleProcessor) TotalFailureRate(ipAddress string) (float64, error) {
	//val, err := rp.DB.Get([]byte(ipAddress))
	logs, err := rp.DB.Get(ipAddress)
	if err != nil {
		return 0.0, err
	}
	/*logs := ipmanager.IPLog{}
	err = json.Unmarshal(val, &logs)
	if err != nil {
		return 0.0, err
	}
	*/
	var rate float64
	if logs.Failure > 0 {
		total := logs.Failure + logs.Success
		rate = float64(logs.Failure)/float64(total)
	} else {
		rate = 0
	}
	return rate, nil
}

func (rp RuleProcessor) FailureRateOverTime(iplog ipmanager.IPLog, lastDuration time.Duration) (float64, time.Duration) {
	startTime := time.Now().Add(-lastDuration)
	failures := 0.0
	successes := 0.0
	for _, l := range iplog.Logs {
		if startTime.Before(l.Timestamp) {
			if l.Result == "FAILURE" {
				failures ++
			} else {
				successes ++
			}
		}
	}
	failureRate := failures / (failures + successes)
	return failureRate, lastDuration
}

func (rp RuleProcessor) FailureRatePerMinute(ipLog ipmanager.IPLog, d time.Duration) (float64) {
	startTime := time.Now().Add(-d)
	failures := 0.0
	successes := 0.0
	for _, l := range ipLog.Logs {
		if startTime.Before(l.Timestamp) {
			if l.Result == "FAILURE" {
				failures ++
			} else {
				successes ++
			}
		}
	}
	minutes := d.Minutes()
	return failures / float64(minutes)
}

func (rp RuleProcessor) FailurePerInterval(log ipmanager.IPLog, interval int, numIntervals int) (map[time.Duration]float64) {
	ans := map[time.Duration]float64{}
	for n := 1; n <= numIntervals; n++ {
		dur := time.Minute * time.Duration(n*interval)
		r := rp.FailureRatePerMinute(log, dur)
		ans[dur] = r
	}
	return ans
}

//func

//func (rp RuleProcessor) CalcFailurePerMinute(iplog ipmanager.IPLog) (float64, error) {
//
//}

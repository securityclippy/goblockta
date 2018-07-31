package ipmanager

import (
	"github.com/securityclippy/goblockta/events"
	//"github.com/securityclippy/badgerwrapper"
	"time"
	"sort"
	"github.com/sirupsen/logrus"
	"errors"
)

type IPMDB interface {
	//Update(key, value []byte) (error)
	//Get(key []byte) ([]byte, error)
	Update(key string, value IPLog) (error)
	Get(key string) (IPLog, error)
	UpdateTimeCursor(newTime time.Time) (error)
	GetTimeCursor() (time.Time, error)
}

type IPManager struct {
	DB IPMDB
}

func NewIPManager(db IPMDB) (IPManager) {
	//db, err := badgerwrapper.NewBadgerDB()
	ipm := IPManager{
		DB: db,
	}
	return ipm
}

type IPLog struct {
	IP string `json:"ip"`
	Success int `json:"success"`
	Failure int `json:"failure"`
	Logs []TimeResult `json:"logs"`
}

func (ipl *IPLog) SortAscending() {
	sort.Slice(ipl.Logs, func(i, j int) bool {return ipl.Logs[i].Timestamp.Before(ipl.Logs[j].Timestamp)})
	logrus.Infof("sorted times: %+v", ipl.Logs)
}

type TimeResult struct {
	Timestamp time.Time `json:"timestamp"`
	Result string `json:"result"`
}

func (i IPManager) ParseLog(event events.LogEvent) (error) {
	t1, err := time.Parse(time.RFC3339, event.Published)
	if err != nil {
		return err
	}
	tr := TimeResult{
		Timestamp: t1,
		Result: event.Outcome.Result,
	}

	//rec, err := i.DB.Get([]byte(event.Client.IPAddress))
	rec, err := i.DB.Get(event.Client.IPAddress)
	if err != nil {
		iplog := IPLog{
			IP: event.Client.IPAddress,
			Logs: []TimeResult{tr},

		}
		if event.Outcome.Result == "SUCCESS" {
			iplog.Success = 1
		} else {
			iplog.Failure = 1
		}
		/*js, err := json.Marshal(iplog)
		if err != nil {
			return err
		}*/
		//err = i.DB.Update([]byte(iplog.IP), js)
		err = i.DB.Update(iplog.IP, iplog)
		if err != nil {
			return err
		}
		return nil
	}

	//iplog := IPLog{}

	/*err = json.Unmarshal(rec, &iplog)
	if err != nil {
		return err
	}*/

	rec.Logs = append(rec.Logs, tr)
	if event.Outcome.Result == "SUCCESS" {
		rec.Success ++
	} else {
		rec.Failure ++
	}
	//iplog.Logs = append(iplog.Logs, tr)

	//if event.Outcome.Result == "SUCCESS" {
		//iplog.Success += 1
	//} else {
		//iplog.Failure += 1
	//}
	//js, err := json.Marshal(iplog)
	//if err != nil {
		//return err
	//}

	//err = i.DB.Update([]byte(iplog.IP), js)
	err = i.DB.Update(rec.IP, rec)
	if err != nil {
		return err
	}
	return nil

}

func (i IPManager) UpdateTimeCursor(newTime time.Time) (error) {
	//err := i.DB.Update([]byte("timecursor"), []byte(logpull.OktaTime(newTime)))
	err := i.DB.UpdateTimeCursor(newTime)
	if err != nil {
		return err
	}
	return nil
}

func (i IPManager) GetTimeCursor() (time.Time, error) {
	//val, err := i.DB.Get([]byte("timecursor"))
	t, err := i.DB.GetTimeCursor()
	if err != nil {
		return time.Time{}, err
	}
	//t, err := time.Parse(time.RFC3339Nano, string(val))
	//if err != nil {
		//return time.Time{}, err
	//}
	return t, nil
}


type IpDB struct {
	IpMap map[string]IPLog
	TimeCursor time.Time
}

//drop in replacement for badgerdb.Update
func (ipdb IpDB) Update(key string, value IPLog) (error) {
	ipdb.IpMap[key] = value
	return nil
}

func (ipdb IpDB) Get(key string) (IPLog, error) {
	val, ok := ipdb.IpMap[key]
	if !ok {
		return IPLog{}, errors.New("Key Does not Exist")
	}
	return val, nil
}

func (ipdb IpDB) SearchPrefix(prefix string) (map[string]IPLog, error){
	return ipdb.IpMap, nil
}

func (ipdb IpDB) UpdateTimeCursor(newTime time.Time) (error) {
	//logrus.Infof("recieved time : %+v", newTime)
	ipdb.TimeCursor = newTime
	//logrus.Infof("my time now: %+v", ipdb.TimeCursor)
	return nil
}

func (ipdb IpDB) GetTimeCursor() (time.Time, error) {
	return ipdb.TimeCursor, nil
}

func NewIpDB() (*IpDB) {
	tc := time.Now().AddDate(0, 0, -1).UTC()
	newTC := tc.Add(time.Hour * 23)
	return &IpDB{
		IpMap: make(map[string]IPLog),
		TimeCursor: newTC,
	}
}
package funcs

import (
"log"
"strings"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
"time"
"strconv"
"github.com/open-falcon/falcon-plus/common/model"
"fmt"
"errors"
)

func NewMongo(host string, port string, user string, passwd string, dbName string)  *mgo.Session{
	var (
		err error
		session *mgo.Session
	)
	if len(user)>0 && len(passwd)>0 {
		dialInfo := &mgo.DialInfo{
			Addrs:     []string{host+":"+port},
			Direct:    true,
			Timeout:   time.Second * 3,
			Database:  dbName,
			Source:    dbName,
			Username:  user,
			Password:  passwd,
		}
		session, err = mgo.DialWithInfo(dialInfo)
	} else {
		session, err = mgo.Dial(strings.Join([]string{host,port},":"))
	}

	if err != nil {
		log.Println("Connect mongodb server error:",err)
		return nil
	}
	session.SetMode(mgo.PrimaryPreferred, true)
	return session
}

func MongoStatInfo(host string, port string, user string, passwd string, dbName string)(L []*model.MetricValue) {
	session := NewMongo(host, port, user, passwd, dbName)
	defer session.Close()

	result := bson.M{}
	if err := session.DB("admin").Run(bson.D{{"serverStatus", 1}}, &result); err != nil {
		log.Println("Mongodb run command serverStatus error!",err)
		return
	} else {
		for k,v := range result {
			L = append(L, ParseMapKV(port, strings.Join([]string{"mongodb",k},"."), v)...)
		}
	}
	L = append(L,MongoReplStatInfo(session,port)...)
	L = append(L,MongoReplTimeDiffHours(session,port)...)

	return
}

func ParseMapKV( port string, key string, data interface{})(L []*model.MetricValue){
	monitorKeys := map[string]string{
		"mongodb.mem.resident":                      "GAUGE",
		"mongodb.mem.virtual":                       "GAUGE",
		"mongodb.mem.mapped":                        "GAUGE",
		"mongodb.connections.current":               "GAUGE",
		"mongodb.connections.available":             "GAUGE",
		"mongodb.asserts.regular":                   "COUNTER",
		"mongodb.asserts.warning":                   "COUNTER",
		"mongodb.asserts.msg":                       "COUNTER",
		"mongodb.asserts.user":                      "COUNTER",
		"mongodb.opcounters.insert":                 "COUNTER",
		"mongodb.opcounters.query":                  "COUNTER",
		"mongodb.opcounters.update":                 "COUNTER",
		"mongodb.opcounters.delete":                 "COUNTER",
		"mongodb.opcounters.command":                "COUNTER",
		"mongodb.opcounters.getmore":                "COUNTER",
		"mongodb.network.bytesIn":                   "COUNTER",
		"mongodb.network.bytesOut":                  "COUNTER",
		"mongodb.network.numRequests":               "COUNTER",
		"mongodb.dur.journaledMB":                   "COUNTER",
		"mongodb.dur.writeToDataFilesMB":            "COUNTER",
		"mongodb.globalLock.currentQueue.readers":   "GAUGE",
		"mongodb.globalLock.currentQueue.writers":   "GAUGE",
		"mongodb.globalLock.activeClients.readers":  "GAUGE",
		"mongodb.globalLock.activeClients.writers":  "GAUGE",
		"mongodb.backgroundFlushing.last_ms":        "GAUGE",
		"mongodb.extra_info.page_faults":            "COUNTER",
	}

	switch value := data.(type) {
	default:
		if  t,ok := monitorKeys[key];ok {
			if t == "GAUGE" {
				L = append(L, GaugeValue(key+"."+port, fmt.Sprintf("%T",value)))
			} else if t == "COUNTER" {
				L = append(L, CounterValue(key+"."+port, fmt.Sprintf("%T",value)))
			}
		}
	case string:
		if  t,ok := monitorKeys[key];ok {
			if t == "GAUGE" {
				L = append(L, GaugeValue(key+"."+port, value))
			} else if t == "COUNTER" {
				L = append(L, CounterValue(key+"."+port, value))
			}
		}
	case int:
		if  t,ok := monitorKeys[key];ok {
			if t == "GAUGE" {
				L = append(L, GaugeValue(key+"."+port, strconv.Itoa(value)))
			} else if t == "COUNTER" {
				L = append(L, CounterValue(key+"."+port, strconv.Itoa(value)))
			}
		}
	case int64:
		if  t,ok := monitorKeys[key];ok {
			if t == "GAUGE" {
				L = append(L, GaugeValue(key+"."+port, strconv.FormatInt(value, 10)))
			} else if t == "COUNTER" {
				L = append(L, CounterValue(key+"."+port, strconv.FormatInt(value, 10)))
			}
		}
	case float64:
		if  t,ok := monitorKeys[key];ok {
			if t == "GAUGE" {
				L = append(L, GaugeValue(key+"."+port, strconv.FormatFloat(value, 'f', 2, 64)))
			} else if t == "COUNTER" {
				L = append(L, CounterValue(key+"."+port, strconv.FormatFloat(value, 'f', 2, 64)))
			}
		}
	case bson.M:
		for k,v := range value {
			L = append(L, ParseMapKV(port, strings.Join([]string{key,k},"."), v)...)
		}
	}
	return
}

func MongoReplStatInfo(session *mgo.Session, port string)(L []*model.MetricValue) {
	result := bson.M{}
	if err := session.DB("admin").Run("replSetGetStatus", &result); err != nil {
		log.Println("Mongodb cluster error: ",err)
		return
	} else {
		var (
			SecondaryOptime uint32
			PrimayOptime  uint32
			IsPrimary bool
		)

		if  v,ok := result["ok"];ok {
			status,err := IsNumber(v)
			if (err != nil) || (status != 1) {
				log.Println("Mongodb cluster is not ok",status)
			}
			L = append(L, CounterValue("mongodb.replSet.ok."+port, status))
		}

		if  v,ok := result["myState"];ok {
			mySate,err := IsNumber(v)
			if err != nil {
				log.Println("Mongodb cluster parse mySate error:",err)
			}
			if mySate == 1 {
				//log.Println("Mongodb cluster mySate is Primary")
				IsPrimary = true
			} else {
				//log.Println("Mongodb cluster mySate is Secondary")
				IsPrimary = false
			}
		}

		if  v,vok := result["members"];vok {
			if vv,vvok := v.([]interface{});vvok {
				for _,vvv := range vv {
					if vvvv,vvvok := vvv.(bson.M);vvvok{
						if vvvvv,vvvvvok := vvvv["self"];vvvvvok{
							x,err:=IsBool(vvvvv)
							if err !=nil {
								log.Println(err)
								continue
							}
							if x {
								y,err := IsMongTimestamp(vvvv["optime"])
								if err != nil {
									log.Println(err)
									continue
								}else {
									SecondaryOptime = y
									//log.Println("Self optime:",y)
								}
							}
						} else if vvvvv,vvvvvok := vvvv["stateStr"];vvvvvok {
							if "PRIMARY" == GetBsonValue(vvvvv) {
								z,err := IsMongTimestamp(vvvv["optime"])
								if err != nil {
									log.Println(err)
									continue
								}else {
									PrimayOptime = z
									//log.Println("Primary optime:",z)
								}
							}
						}
					}
				}
			}
		}
		if IsPrimary {
			//log.Println("Mongodb.replSet.lag."+port, "0.00")
			L = append(L, GaugeValue("mongodb.replSet.lag."+port, "0.00"))
		} else {
			lag := strconv.FormatFloat(float64(PrimayOptime-SecondaryOptime), 'f', 2, 64)
			L = append(L, GaugeValue("mongodb.replSet.lag."+port, lag))
			//log.Println("Mongodb.replSet.lag."+port,lag)
		}
	}
	return
}

func MongoReplTimeDiffHours(session *mgo.Session, port string)(L []*model.MetricValue) {
	var (
		FirstOpTime uint32
		LastOpTime uint32
		timeDiffHours string
	)
	result := bson.M{}

	db := session.DB("local")
	collection := db.C("oplog.rs")

	countNum, err := collection.Count()
	if err != nil {
		panic(err)
	}
	if countNum < 2 {
		return
	}

	err = collection.Find(nil).Sort("$natural").One(&result)
	if err != nil {
		log.Println(err)
		return
	} else {
		if v, ok := result["ts"]; ok {
			vv,err := IsMongTimestamp(v)
			if err != nil {
				log.Println(err)
			}else {
				FirstOpTime = vv
				//log.Println("FirstOpTime:",vv)
			}
		}
	}

	err = collection.Find(nil).Sort("-$natural").One(&result)
	if err != nil {
		log.Println(err)
		return
	} else {
		if v, ok := result["ts"]; ok {
			vv,err := IsMongTimestamp(v)
			if err != nil {
				log.Println(err)
			}else {
				LastOpTime = vv
				//log.Println("LastOpTime:",vv)
			}
		}
	}

	timeDiffHours = strconv.FormatFloat((float64(LastOpTime-FirstOpTime))/3600.0, 'f', 2, 64)
	L = append(L, GaugeValue("mongodb.replSet.timeDiffHours."+port, timeDiffHours))
	//log.Println("timeDiffHours:",timeDiffHours)
	return
}

func IsNumber(data interface{})(r  int, err error) {
	switch value := data.(type) {
	default:
		r =  0
		err = errors.New("not int type")
	case int:
		r=  int(value)
	case int64:
		r=int(value)
	case float64:
		r=int(value)
	}
	return
}

func IsBool(data interface{})(r  bool, err error) {
	if res,ok:=data.(bool);ok{
		r = res
	} else {
		err = errors.New("not bool type")
	}
	return
}

func IsMongTimestamp(data interface{})(r  uint32, err error) {
	if res,ok := data.(bson.MongoTimestamp);ok{
		timestamp := int64(res)
		r = uint32(timestamp >> 32)
	} else {
		err = errors.New("no bson.MongoTimestamp type")
	}
	return
}
func GetBsonValue(data interface{}) (r string){
	switch value := data.(type) {
	default:
		r =  fmt.Sprintf("%T",value)
	case string:
		r=value
	case int:
		r=  strconv.Itoa(value)
	case int64:
		r=strconv.FormatInt(value, 10)
	case float64:
		r=strconv.FormatFloat(value, 'f', 2, 64)
	}
	return
}

package funcs

import (
	"log"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/open-falcon/falcon-plus/common/model"
)

func NewRedis(arg []string) redis.Conn {
	c, err := redis.Dial("tcp", arg[0])
	if err != nil {
		log.Println(err)
		return nil
	}
	//密码授权
	_, err = c.Do("AUTH", arg[1])
	if err != nil {
		return nil
	}
	return c
}

func RedisStatInfo(host string,port string,passwd string) (L []*model.MetricValue) {
	conn := NewRedis([]string{strings.Join([]string{host,port},":"),passwd})
        if conn == nil {
		log.Println("Error Redis connect failed,port=",port)
		return
	} 
	defer conn.Close()
	monitorKeys := map[string]string{
		"connected_clients":        "GAUGE",
		"blocked_clients":          "GAUGE",
		"used_memory":              "GAUGE",
		"used_memory_rss":          "GAUGE",
		"mem_fragmentation_ratio":  "GAUGE",
		"total_commands_processed": "COUNTER",
		"rejected_connections":     "COUNTER",
		"expired_keys":             "COUNTER",
		"evicted_keys":             "COUNTER",
		"keyspace_hits":            "COUNTER",
		"keyspace_misses":          "COUNTER",
		"keys_num":                 "GAUGE",
		"role":                     "COUNTER"}

	var (
		keyspace_hits float64 
		keyspace_misses float64
	)

	r, _ := redis.String(conn.Do("info"))
	r = strings.Replace(string(r),"\r\n", "\n", -1)
	for _, item := range strings.Split(r, "\n") {
		rkey := strings.Split(item, ":")
		if t, ok := monitorKeys[rkey[0]]; ok {
			if rkey[0] == "keyspace_hits" {
				keyspace_hits,_ = strconv.ParseFloat(rkey[1],64)
				L = append(L, CounterValue("redis."+rkey[0]+"."+port, rkey[1]))
			} else if rkey[0] == "keyspace_misses" {
				keyspace_misses,_ = strconv.ParseFloat(rkey[1],64)
				L = append(L, CounterValue("redis."+rkey[0]+"."+port, rkey[1]))
			} else if rkey[0] == "role" {
				if rkey[1] == "master"{
					L = append(L, CounterValue("redis."+rkey[0]+"."+port, "0"))
				} else {
					L = append(L, CounterValue("redis."+rkey[0]+"."+port, "1"))
				}
			} else if t == "GAUGE" {
				L = append(L, GaugeValue("redis."+rkey[0]+"."+port, rkey[1]))
			} else if t == "COUNTER" {
				L = append(L, CounterValue("redis."+rkey[0]+"."+port, rkey[1]))
			}
		}
	}

	if keyspace_hits > 1.0 {
		keyspace_hit_ratio := keyspace_hits / (keyspace_hits + keyspace_misses)
		L = append(L, GaugeValue("redis.keyspace_hit_ratio."+port, strconv.FormatFloat(keyspace_hit_ratio, 'f', 2, 64)))
	} else {
		L = append(L, GaugeValue("redis.keyspace_hit_ratio."+port, "0.00"))
	}

	return
}

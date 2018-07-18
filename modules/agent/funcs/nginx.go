package funcs

import (
	"log"
	"fmt"
	"strings"
	"net/http"
	"github.com/open-falcon/falcon-plus/common/model"

	"io/ioutil"
)

func NginxStatInfo(host string,port string)(L []*model.MetricValue){
	L = append(L,NginxServerInfo(host,port)...)
	L = append(L,NginxRequestInfo(host,port)...)
	return
}

func NginxServerInfo(host string,port string)(L []*model.MetricValue){
	resp, err := http.Get("http://"+host+":"+port+"/nxstatus/server")
	if err != nil {
		log.Println("Error! Nginx server request failed!",err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error! Nginx server status read error!",err.Error())
		return
	}
	for _, row := range strings.Split(strings.TrimSpace(string(body)), "\n") {
		if strings.Contains(row, "Active") {
			item := strings.Split(strings.TrimSpace(row), " ")
			if len(item) == 3 {
				L = append(L, GaugeValue("nginx.server.active_connection", item[2]))
	        		//log.Println("nginx.server.active_connection == ", item[2])
			} else {
				log.Println("Error! Nginx server stats error! ",item)
				continue
			}
		} else if strings.Contains(row, "Reading") {
			item := strings.Split(strings.TrimSpace(row), " ")
			if len(item) == 6 {
				L = append(L, GaugeValue("nginx.server.reading", item[1]))
				L = append(L, GaugeValue("nginx.server.writing", item[3]))
				L = append(L, GaugeValue("nginx.server.waiting", item[5]))
	        		//log.Println("nginx.server.reading == ", item[1])
	        		//log.Println("nginx.server.writing == ", item[3])
	        		//log.Println("nginx.server.waiting == ", item[5])
			} else {
				log.Println("Error! Nginx server stats error! ",item)
				continue
			}
		}
	}
	return
}

func NginxRequestInfo(host string,port string)(L []*model.MetricValue){
	resp, err := http.Get("http://"+host+":"+port+"/nxstatus/request")
	if err != nil {
		log.Println("Error! Nginx request failed!",err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error! Nginx request stats read error!",err.Error())
		return
	}
	for _, row := range strings.Split(strings.TrimSpace(string(body)), "\n") {
		item := strings.Split(strings.TrimSpace(row), "\t")
		if len(item) != 2 {
			log.Println("Error! Nginx request stats error! ",item)
			continue
		}
		ks := strings.TrimSpace(item[0])
		ks = strings.Replace(ks, ".", "_", -1)
		ks = strings.Replace(ks, "-", ".", -1)
		tmp := strings.Split(ks, ".")
		k := fmt.Sprintf("%s.%s", tmp[1], tmp[0])
		if strings.Contains(k, "Time") && strings.Contains(item[1], "|") {
			v := strings.Split(strings.TrimSpace(item[1]), "|")
			if len(v) == 6 {
				L = append(L, GaugeValue("nginx.request."+tmp[1]+".p995."+tmp[0], v[0]))
				L = append(L, GaugeValue("nginx.request."+tmp[1]+".p99."+tmp[0], v[1]))
				L = append(L, GaugeValue("nginx.request."+tmp[1]+".p98."+tmp[0], v[2]))
				L = append(L, GaugeValue("nginx.request."+tmp[1]+".p95."+tmp[0], v[3]))
				L = append(L, GaugeValue("nginx.request."+tmp[1]+".p90."+tmp[0], v[4]))
				L = append(L, GaugeValue("nginx.request."+tmp[1]+".p80."+tmp[0], v[5]))
	        		//log.Println("nginx.request."+tmp[1]+".p995."+tmp[0]+" == ", v[0])
	        		//log.Println("nginx.request."+tmp[1]+".p99."+tmp[0]+" == ", v[1])
	        		//log.Println("nginx.request."+tmp[1]+".p98."+tmp[0]+" == ", v[2])
	        		//log.Println("nginx.request."+tmp[1]+".p95."+tmp[0]+" == ", v[3])
	        		//log.Println("nginx.request."+tmp[1]+".p90."+tmp[0]+" == ", v[4])
	        		//log.Println("nginx.request."+tmp[1]+".p80."+tmp[0]+" == ", v[5])
			} else {
				log.Println("Error! Nginx request stats error! ",v)
				continue
			}

		} else {
			v := strings.TrimSpace(item[1])
			L = append(L, GaugeValue("nginx.request."+k, v))
	        	//log.Println("nginx.request."+k+" == ", v)
		}
	}
	return
}

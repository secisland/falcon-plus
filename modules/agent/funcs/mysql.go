package funcs

import (
	"log"
	"strconv"
	"database/sql"
	"fmt"
	"strings"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/falcon-plus/common/model"
)

func NewMySQL(arg []string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",arg[0],arg[1],arg[2],arg[3],arg[4]))
	if err != nil {
		log.Println(err)
		return nil
	}
	return db
}

func MySQLStatInfo(host string, port string, user string, passwd string, dbName string) (L []*model.MetricValue) {
	db := NewMySQL([]string{user, passwd, host, port, dbName})
	if db == nil {
		log.Println("Error MySQL connect failed,port=",port)
		return
	}
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Println("Error MySQL connect error,port=",port,err.Error())
		return
	}
	L = append(L,MySQLStatusStatInfo(db,port)...)
	L = append(L,MySQLEngineStatInfo(db,port)...)
	L = append(L,MySQLSlaveInfo(db,port)...)
	return
}


func MySQLStatusStatInfo(db *sql.DB, port string)(L []*model.MetricValue){
	rows, err := db.Query("SHOW GLOBAL STATUS")
	defer rows.Close()

	monitorKeys := map[string]string{
		"Com_select":"COUNTER",
		"Qcache_hits":"COUNTER",
		"Com_insert":"COUNTER",
		"Com_update":"COUNTER",
		"Com_delete":"COUNTER",
		"Com_replace":"COUNTER",
		"MySQL_QPS":"COUNTER",
		"MySQL_TPS":"COUNTER",
		"ReadWrite_ratio":"GAUGE",
		"Innodb_buffer_pool_read_requests":"COUNTER",
		"Innodb_buffer_pool_reads":"COUNTER",
		"Innodb_buffer_read_hit_ratio":"GAUGE",
		"Innodb_buffer_pool_pages_flushed":"COUNTER",
		"Innodb_buffer_pool_pages_free":"GAUGE",
		"Innodb_buffer_pool_pages_dirty":"GAUGE",
		"Innodb_buffer_pool_pages_data":"GAUGE",
		"Bytes_received":"COUNTER",
		"Bytes_sent":"COUNTER",
		"Innodb_rows_deleted":"COUNTER",
		"Innodb_rows_inserted":"COUNTER",
		"Innodb_rows_read":"COUNTER",
		"Innodb_rows_updated":"COUNTER",
		"Innodb_os_log_fsyncs":"COUNTER",
		"Innodb_os_log_written":"COUNTER",
		"Created_tmp_disk_tables":"COUNTER",
		"Created_tmp_tables":"COUNTER",
		"Connections":"COUNTER",
		"Innodb_log_waits":"COUNTER",
		"Slow_queries":"COUNTER",
		"Binlog_cache_disk_use":"COUNTER"}
	var (
		Com_select float64
		Qcache_hits float64
		Com_insert float64
		Com_update float64
		Com_delete float64
		Com_replace float64
		Innodb_buffer_pool_read_requests float64
		Innodb_buffer_pool_reads float64
		MySQL_QPS float64
		MySQL_TPS float64
	)
	for rows.Next() {
		var Variable_name string
		var Value string
		err = rows.Scan(&Variable_name, &Value)
		if t, ok := monitorKeys[Variable_name]; ok {
			//fmt.Println(Variable_name,"=>", Value)
			if Variable_name == "Com_select" {
				Com_select,_ = strconv.ParseFloat(Value,64)
				L = append(L, CounterValue("mysql.Com_select."+port, Value))
			} else if Variable_name == "Qcache_hits" {
				Qcache_hits,_ = strconv.ParseFloat(Value,64)
				L = append(L, CounterValue("mysql.Qcache_hits."+port, Value))
			} else if Variable_name == "Com_insert" {
				Com_insert,_ = strconv.ParseFloat(Value,64)
				L = append(L, CounterValue("mysql.Com_insert."+port, Value))
			} else if Variable_name == "Com_update" {
				Com_update,_ = strconv.ParseFloat(Value,64)
				L = append(L, CounterValue("mysql.Com_update."+port, Value))
			} else if Variable_name == "Com_delete" {
				Com_delete,_ = strconv.ParseFloat(Value,64)
				L = append(L, CounterValue("mysql.Com_delete."+port, Value))
			} else if Variable_name == "Com_replace" {
				Com_replace,_ = strconv.ParseFloat(Value,64)
				L = append(L, CounterValue("mysql.Com_replace."+port, Value))
			} else if Variable_name == "Innodb_buffer_pool_read_requests" {
				Innodb_buffer_pool_read_requests,_ = strconv.ParseFloat(Value,64)
				L = append(L, CounterValue("mysql.Innodb_buffer_pool_read_requests."+port, Value))
			} else if Variable_name == "Innodb_buffer_pool_reads" {
				Innodb_buffer_pool_reads,_ = strconv.ParseFloat(Value,64)
				L = append(L, CounterValue("mysql.Innodb_buffer_pool_reads."+port, Value))
			} else if t == "GAUGE" {
				L = append(L, GaugeValue("mysql."+Variable_name+"."+port, Value))
			} else if t == "COUNTER" {
				L = append(L, CounterValue("mysql."+Variable_name+"."+port, Value))
			}
		}
	}

	MySQL_QPS = Com_select + Qcache_hits
	MySQL_TPS = Com_insert + Com_update + Com_delete + Com_replace

	L = append(L, CounterValue("mysql.MySQL_QPS."+port, strconv.FormatFloat(MySQL_QPS, 'f', 2, 64)))
	L = append(L, CounterValue("mysql.MySQL_TPS."+port, strconv.FormatFloat(MySQL_TPS, 'f', 2, 64)))
	//fmt.Println("MySQL_QPS","=>", MySQL_QPS)
	//fmt.Println("MySQL_TPS","=>", MySQL_TPS)

	if Innodb_buffer_pool_read_requests > 1.0 {
		Innodb_buffer_read_hit_ratio := (Innodb_buffer_pool_read_requests - Innodb_buffer_pool_reads) / Innodb_buffer_pool_read_requests
		//fmt.Println("Innodb_buffer_read_hit_ratio","=>", Innodb_buffer_read_hit_ratio)
		L = append(L, GaugeValue("mysql.Innodb_buffer_read_hit_ratio."+port, strconv.FormatFloat(Innodb_buffer_read_hit_ratio, 'f', 2, 64)))
	} else {
		//fmt.Println("Innodb_buffer_read_hit_ratio","=>", "0.00")
		L = append(L, GaugeValue("mysql.Innodb_buffer_read_hit_ratio."+port, "0.00"))
	}

	if MySQL_TPS > 1.0 {
		//fmt.Println("ReadWrite_ratio","=>", MySQL_QPS/MySQL_TPS)
		L = append(L, GaugeValue("mysql.ReadWrite_ratio."+port, strconv.FormatFloat(MySQL_QPS/MySQL_TPS, 'f', 2, 64)))
	} else {
		//fmt.Println("ReadWrite_ratio","=>", "0.00")
		L = append(L, GaugeValue("mysql.ReadWrite_ratio."+port, "0.00"))
	}

	err = rows.Err()
	if err != nil {
		log.Println("Error MySQL status query error,port=",port,err.Error())
		return
	}
	return
}

func MySQLEngineStatInfo(db *sql.DB, port string)(L []*model.MetricValue){
	rows, err := db.Query("SHOW ENGINE INNODB STATUS")
	if err != nil {
		log.Println("Error MySQL engin innodb status query error,port=",port,err.Error())
		return
	}
	defer rows.Close()
	reg, _ := regexp.Compile("History list length ([0-9]+)\n")

	for rows.Next() {
		var Type string
		var Name string
		var Status string
		err = rows.Scan(&Type, &Name, &Status)
		if strings.Contains(Status,"History list length") {
			v := reg.FindStringSubmatch(Status)
			if nil != v {
				//fmt.Println("Undo_Log_Length","=>", v[1])
				L = append(L, GaugeValue("mysql.Undo_Log_Length."+port, v[1]))
			} else {
				//fmt.Println("Undo_Log_Length","=>", "0.00")
				L = append(L, GaugeValue("mysql.Undo_Log_Length."+port, "0.00"))
			}
		}
	}
	err = rows.Err()
	if err != nil {
		log.Println("Error MySQL engin innodb status query error,port=",port,err.Error())
		return
	}
	return
}

func MySQLSlaveInfo(db *sql.DB, port string)(L []*model.MetricValue){
	rows, err := db.Query("SHOW SLAVE STATUS")
	if err != nil {
		log.Println("MySQL slave status query error,port=",port,err)
		return
	}
	defer rows.Close()

	status, err := ScanMap(rows)

	if err != nil {
		log.Println("MySQL slave status parse error,port=",port,err)
		return
	}
	if v,ok := status["Slave_IO_Running"];ok {
		if v.Valid && v.String == "YES"{
			L = append(L, CounterValue("mysql.Slave_IO_Running."+port, "0.00"))
		} else {
			L = append(L, CounterValue("mysql.Slave_IO_Running."+port, "1.00"))
		}
		//log.Println("Slave_IO_Running =>",v.String)
	}
	if v,ok := status["Slave_SQL_Running"];ok {
		if v.Valid && v.String == "YES"{
			L = append(L, CounterValue("mysql.Slave_SQL_Running."+port, "0.00"))
		} else {
			L = append(L, CounterValue("mysql.Slave_SQL_Running."+port, "1.00"))
		}
		//log.Println("Slave_SQL_Running =>",v.String)
	}
	if v,ok := status["Seconds_Behind_Master"];ok {
		if v.Valid {
			L = append(L, CounterValue("mysql.Seconds_Behind_Master."+port,v.String))
		} else {
			L = append(L, CounterValue("mysql.Seconds_Behind_Master."+port,"-1.00"))
		}
		//log.Println("Seconds_Behind_Master =>",v.String)
	}
	//log.Println("status =>",status)
	return
}

func ScanMap(rows *sql.Rows) (map[string]sql.NullString, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, err
		} else {
			return nil, nil
		}
	}

	values := make([]interface{}, len(columns))

	for index := range values {
		values[index] = new(sql.NullString)
	}

	err = rows.Scan(values...)

	if err != nil {
		return nil, err
	}

	result := make(map[string]sql.NullString)

	for index, columnName := range columns {
		result[columnName] = *values[index].(*sql.NullString)
	}

	return result, nil
}

falcon-agent
===

This is a linux monitor agent. Just like zabbix-agent and tcollector.
新增内置mysql/redis/mongodb监控项，支持数据库多实例端口的监控，只需在配置文件文增加配置即可开启，具体配置可参考 cfg.example.json 。

## Installation

It is a golang classic project

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/secisland/falcon-plus.git
cd falcon-plus/modules/agent
go get
./control build
./control start

# goto http://localhost:1988
```

I use [linux-dash](https://github.com/afaqurk/linux-dash) as the page theme.

## Configuration

- heartbeat: heartbeat server rpc address
- transfer: transfer rpc address
- ignore: the metrics should ignore


## Auto deployment

Just look at https://github.com/open-falcon/ops-updater

## Metrics description

### MySQL监控metrics如下：
```
"mysql.Com_select.{{PORT}}":                          "COUNTER",
"mysql.Qcache_hits.{{PORT}}":                         "COUNTER",
"mysql.Com_insert.{{PORT}}":                          "COUNTER",
"mysql.Com_update.{{PORT}}":                          "COUNTER",
"mysql.Com_delete.{{PORT}}":                          "COUNTER",
"mysql.Com_replace.{{PORT}}":                         "COUNTER",
"mysql.MySQL_QPS.{{PORT}}":                           "COUNTER",
"mysql.MySQL_TPS.{{PORT}}":                           "COUNTER",
"mysql.ReadWrite_ratio.{{PORT}}":                     "GAUGE",
"mysql.Innodb_buffer_pool_read_requests.{{PORT}}":    "COUNTER",
"mysql.Innodb_buffer_pool_reads.{{PORT}}":            "COUNTER",
"mysql.Innodb_buffer_read_hit_ratio.{{PORT}}":        "GAUGE",
"mysql.Innodb_buffer_pool_pages_flushed.{{PORT}}":    "COUNTER",
"mysql.Innodb_buffer_pool_pages_free.{{PORT}}":       "GAUGE",
"mysql.Innodb_buffer_pool_pages_dirty.{{PORT}}":      "GAUGE",
"mysql.Innodb_buffer_pool_pages_data.{{PORT}}":       "GAUGE",
"mysql.Bytes_received.{{PORT}}":                      "COUNTER",
"mysql.Bytes_sent.{{PORT}}":                          "COUNTER",
"mysql.Innodb_rows_deleted.{{PORT}}":                 "COUNTER",
"mysql.Innodb_rows_inserted.{{PORT}}":                "COUNTER",
"mysql.Innodb_rows_read.{{PORT}}":                    "COUNTER",
"mysql.Innodb_rows_updated.{{PORT}}":                 "COUNTER",
"mysql.Innodb_os_log_fsyncs.{{PORT}}":                "COUNTER",
"mysql.Innodb_os_log_written.{{PORT}}":               "COUNTER",
"mysql.Created_tmp_disk_tables.{{PORT}}":             "COUNTER",
"mysql.Created_tmp_tables.{{PORT}}":                  "COUNTER",
"mysql.Connections.{{PORT}}":                         "COUNTER",
"mysql.Innodb_log_waits.{{PORT}}":                    "COUNTER",
"mysql.Slow_queries.{{PORT}}":                        "COUNTER",
"mysql.Binlog_cache_disk_use.{{PORT}}":               "COUNTER",
"mysql.Undo_Log_Length.{{PORT}}":                     "GAUGE",
``` 
MySQL从库复制状态的监控项(主库不会产生下列metric):
```
"mysql.Slave_IO_Running.{{PORT}}":                    "GAUGE",
"mysql.Slave_SQL_Running.{{PORT}}":                   "GAUGE",
"mysql.Seconds_Behind_Master.{{PORT}}":               "GAUGE",
```

### Redis监控metrics如下：
```
"redis.connected_clients.{{PORT}}":             "GAUGE",
"redis.blocked_clients.{{PORT}}":               "GAUGE",
"redis.used_memory.{{PORT}}":                   "GAUGE",
"redis.used_memory_rss.{{PORT}}":               "GAUGE",
"redis.mem_fragmentation_ratio.{{PORT}}":       "GAUGE",
"redis.total_connections_received.{{PORT}}":    "COUNTER",
"redis.total_commands_processed.{{PORT}}":      "COUNTER",
"redis.rejected_connections.{{PORT}}":          "COUNTER",
"redis.total_net_input_bytes.{{PORT}}":         "COUNTER",
"redis.total_net_output_bytes.{{PORT}}":        "COUNTER",
"redis.instantaneous_ops_per_sec.{{PORT}}":     "GAUGE",
"redis.expired_keys.{{PORT}}":                  "COUNTER",
"redis.evicted_keys.{{PORT}}":                  "COUNTER",
"redis.keyspace_hits.{{PORT}}":                 "COUNTER",
"redis.keyspace_misses.{{PORT}}":               "COUNTER",
"redis.keys_num.{{PORT}}":                      "GAUGE",
"redis.role.{{PORT}}":                          "COUNTER",
"redis.keyspace_hit_ratio.{{PORT}}":            "GAUGE"
```
RedisCmdStatInfo如下(redis-cli> info Commandstats 命令所有结果)：
```
"redis.cmdstat_auth.{{PORT}}":              "COUNTER",    
"redis.cmdstat_config.{{PORT}}":            "COUNTER",
"redis.cmdstat_expire.{{PORT}}":            "COUNTER",
"redis.cmdstat_get.{{PORT}}":               "COUNTER",
"redis.cmdstat_info.{{PORT}}":              "COUNTER",
"redis.cmdstat_pexpire.{{PORT}}":           "COUNTER",
"redis.cmdstat_ping.{{PORT}}":              "COUNTER",
"redis.cmdstat_scan.{{PORT}}":              "COUNTER",
"redis.cmdstat_select.{{PORT}}":            "COUNTER",
"redis.cmdstat_ttl.{{PORT}}":               "COUNTER",
"redis.cmdstat_type.{{PORT}}":              "COUNTER",
"redis.cmdstat_zadd.{{PORT}}":              "COUNTER",
"redis.cmdstat_zcard.{{PORT}}":             "COUNTER",
"redis.cmdstat_zrange.{{PORT}}":            "COUNTER",
"redis.cmdstat_zrem.{{PORT}}":              "COUNTER",
"redis.cmdstat_zremrangebyrank.{{PORT}}":   "COUNTER",
"redis.cmdstat_zrevrange.{{PORT}}":         "COUNTER",
"redis.cmdstat_zrevrangebyscore.{{PORT}}":  "COUNTER",
```
### MongoDB监控metrics如下：
```
"mongodb.mem.resident.{{PORT}}":                        "GAUGE",
"mongodb.mem.virtual.{{PORT}}":                         "GAUGE",
"mongodb.mem.mapped.{{PORT}}":                          "GAUGE",
"mongodb.connections.current.{{PORT}}":                 "GAUGE",
"mongodb.connections.available.{{PORT}}":               "GAUGE",
"mongodb.asserts.regular.{{PORT}}":                     "COUNTER",
"mongodb.asserts.warning.{{PORT}}":                     "COUNTER",
"mongodb.asserts.msg.{{PORT}}":                         "COUNTER",
"mongodb.asserts.user.{{PORT}}":                        "COUNTER",
"mongodb.opcounters.insert.{{PORT}}":                   "COUNTER",
"mongodb.opcounters.query.{{PORT}}":                    "COUNTER",
"mongodb.opcounters.update.{{PORT}}":                   "COUNTER",
"mongodb.opcounters.delete.{{PORT}}":                   "COUNTER",
"mongodb.opcounters.command.{{PORT}}":                  "COUNTER",
"mongodb.opcounters.getmore.{{PORT}}":                  "COUNTER",
"mongodb.network.bytesIn.{{PORT}}":                     "COUNTER",
"mongodb.network.bytesOut.{{PORT}}":                    "COUNTER",
"mongodb.network.numRequests.{{PORT}}":                 "COUNTER",
"mongodb.dur.journaledMB.{{PORT}}":                     "COUNTER",
"mongodb.dur.writeToDataFilesMB.{{PORT}}":              "COUNTER",
"mongodb.globalLock.currentQueue.readers.{{PORT}}":     "GAUGE",
"mongodb.globalLock.currentQueue.writers.{{PORT}}":     "GAUGE",
"mongodb.globalLock.activeClients.readers.{{PORT}}":    "GAUGE",
"mongodb.globalLock.activeClients.writers.{{PORT}}":    "GAUGE",
"mongodb.backgroundFlushing.last_ms.{{PORT}}":          "GAUGE",
"mongodb.extra_info.page_faults.{{PORT}}":              "COUNTER",
```
副本集的状态监控metrics如下(非副本集不会产生下列metric)：
```
"mongodb.replSet.ok.{{PORT}}":                          "COUNTER",
"mongodb.replSet.lag.{{PORT}}":                         "GAUGE",
"mongodb.replSet.timeDiffHours.{{PORT}}":               "GAUGE"
```

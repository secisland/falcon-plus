// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"log"
)

type FuncsAndInterval struct {
	Fs       []func() []*model.MetricValue
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := g.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		{
			Fs: []func() []*model.MetricValue{
				AgentMetrics,
				CpuMetrics,
				NetMetrics,
				KernelMetrics,
				LoadAvgMetrics,
				MemMetrics,
				DiskIOMetrics,
				IOStatsMetrics,
				NetstatMetrics,
				ProcMetrics,
				UdpMetrics,
				ModuleMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*model.MetricValue{
				DeviceMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*model.MetricValue{
				PortMetrics,
				SocketStatSummaryMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*model.MetricValue{
				DuMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*model.MetricValue{
				UrlMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*model.MetricValue{
				GpuMetrics,
			},
			Interval: interval,
		},
	}
}

func ModuleMetrics() (L []*model.MetricValue) {
        modules := g.Config().Modules
        sz := len(modules)
	log.Printf("=> <Total=%d> Modules Config\n", sz)
        if sz == 0 {
                return
        }

        for _,module := range modules {
                if module.Name == "redis" {
                        L = append(L, RedisStatInfo(module.Host, module.Port, module.Passwd)...)
                } else if module.Name == "mysql" {
                        L = append(L, MySQLStatInfo(module.Host, module.Port, module.User, module.Passwd, module.DbName)...)
                } else if module.Name == "mongodb" {
                        L = append(L, MongoStatInfo(module.Host, module.Port, module.User, module.Passwd, module.DbName)...)
                } else if module.Name == "nginx" {
                        L = append(L, NginxStatInfo(module.Host, module.Port)...)
                }
        }
	log.Printf("=> <Total=%d> ModulesMetrics\n", len(L))

        return
}

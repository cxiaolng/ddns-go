package dns

import (
	"log"
	"time"

	"github.com/jeessy2/ddns-go/v5/config"
	"github.com/jeessy2/ddns-go/v5/util"
)

// DNS interface
type DNS interface {
	Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache)
	// 添加或更新IPv4/IPv6记录
	AddUpdateDomainRecords() (domains config.Domains)
}

var Ipcache = [][2]util.IpCache{}

// RunTimer 定时运行
func RunTimer(firstDelay time.Duration, delay time.Duration) {
	log.Printf("第一次运行将等待 %d 秒后运行 (等待网络)", int(firstDelay.Seconds()))
	time.Sleep(firstDelay)
	for {
		RunOnce()
		time.Sleep(delay)
	}
}

// RunOnce RunOnce
func RunOnce() {
	conf, err := config.GetConfigCached()
	if err != nil {
		return
	}
	// 校验是否睡眠模式
	if conf.SleepMode.SleepDisable && util.IsSleepMode(conf.SleepMode.SleepTimeStart, conf.SleepMode.SleepTimeEnd) {
		log.Printf("当前处于睡眠状态")
		return
	}
	if util.ForceCompareGlobal || len(Ipcache) != len(conf.DnsConf) {
		Ipcache = [][2]util.IpCache{}
		for range conf.DnsConf {
			Ipcache = append(Ipcache, [2]util.IpCache{{}, {}})
		}
	}

	for i, dc := range conf.DnsConf {
		var dnsSelected DNS
		switch dc.DNS.Name {
		case "alidns":
			dnsSelected = &Alidns{}
		case "dnspod":
			dnsSelected = &Dnspod{}
		case "cloudflare":
			dnsSelected = &Cloudflare{}
		case "huaweicloud":
			dnsSelected = &Huaweicloud{}
		case "callback":
			dnsSelected = &Callback{}
		case "baiducloud":
			dnsSelected = &BaiduCloud{}
		case "porkbun":
			dnsSelected = &Porkbun{}
		case "godaddy":
			dnsSelected = &GoDaddyDNS{}
		case "googledomain":
			dnsSelected = &GoogleDomain{}
		case "namecheap":
			dnsSelected = &NameCheap{}
		default:
			dnsSelected = &Alidns{}
		}
		dnsSelected.Init(&dc, &Ipcache[i][0], &Ipcache[i][1])
		domains := dnsSelected.AddUpdateDomainRecords()
		// webhook
		v4Status, v6Status := config.ExecWebhook(&domains, &conf)
		// 重置单个cache
		if v4Status == config.UpdatedFailed {
			Ipcache[i][0] = util.IpCache{}
		}
		if v6Status == config.UpdatedFailed {
			Ipcache[i][1] = util.IpCache{}
		}
	}

	util.ForceCompareGlobal = false

}

package utils

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/datasource"
	"math"
	"time"
)

const ipFrameSize = 2

func init() {
	resetGroupIpList()
}

func resetGroupIpList() {
	log.Println("ip_day_lucky.resetGroupIpList start")
	cacheObj := datasource.InstanceCache()
	for i := 0; i < ipFrameSize; i++ {
		key := fmt.Sprintf("day_ips_%d", i)
		cacheObj.Do("DEL", key)
	}
	log.Println("ip_day_lucky.resetGroupIpList end")
	// IP当天的统计数.零点的时候归零,设置定时器
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupIpList)
}

func IncrIpLuckyNum(ip string) int64 {
	ipInt := comm.Ip4toInt(ip)
	i := ipInt % ipFrameSize
	key := fmt.Sprintf("day_ips_%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, ipInt, 1)
	if err != nil {
		log.Println("ip_day_lucky redis HINCRBY error = ", err)
		return math.MaxInt32
	} else {
		return rs.(int64)
	}
}

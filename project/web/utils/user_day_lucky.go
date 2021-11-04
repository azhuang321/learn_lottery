package utils

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/datasource"
	"math"
	"time"
)

const userFrameSize = 2

func init() {
	resetGroupUserList()
}

func resetGroupUserList() {
	log.Println("user_day_lucky.resetGroupList start")
	cacheObj := datasource.InstanceCache()
	for i := 0; i < userFrameSize; i++ {
		key := fmt.Sprintf("day_users_%d", i)
		cacheObj.Do("DEL", key)
	}
	log.Println("user_day_lucky.resetGroupList end")
	// 零点归零的定时器
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupUserList)
}

func IncrUserLuckyNum(uid int) int64 {
	i := uid % userFrameSize
	key := fmt.Sprintf("day_users_%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, uid, 1)
	if err != nil {
		log.Println("user_day_lucky redis hincrby key", key, ",uid", uid, ".err", err)
		return math.MaxInt32
	} else {
		num := rs.(int64)
		return num
	}
}

func InitUserLuckyNum(uid int, num int64) {
	if num <= 1 {
		return
	}
	i := uid % userFrameSize
	key := fmt.Sprintf("day_users_%d", i)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("HSET", key, uid, num)
	if err != nil {
		log.Println("user_day_lucky redis hset key", key, ",uid", uid, ".err", err)
	}
}

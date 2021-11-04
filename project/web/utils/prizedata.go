package utils

import (
	"log"
	"lottery/comm"
	"lottery/datasource"
	"lottery/services"
)

func GetGiftPoolNum(id int) int {
	key := "gift_pool"
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HGET", key, id)
	if err != nil {
		log.Println("prizedata.GetGiftpoolnum err=", err)
		return 0
	} else {
		num := comm.GetInt64(rs, 0)
		return int(num)
	}
}

func prizeServGift(id int) bool {
	key := "gift_pool"
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, id, -1)
	if err != nil {
		log.Println("prizedata.GetGiftpoolnum err=", err)
		return false
	}
	num := comm.GetInt64(rs, -1)
	if num >= 0 {
		return true
	} else {
		return false
	}
}

func PrizeGift(id, leftNum int) bool {
	ok := false
	ok = prizeServGift(id)
	if ok {
		giftService := services.NewGiftService()
		rows, err := giftService.DecrLeftNum(id, 1)
		if rows < 1 || err != nil {
			log.Println("prizedata.PrizeGift gitservice.decrleftnum error=", err)
			return false
		}
	}
	return ok
}

func PrizeCodeDiff(id int, service services.CodeService) string {
	lockUid := 0 - id - 1000000000
	LockLucky(lockUid)
	defer UnLockLucky(lockUid)

	codeId := 0
	codeInfo := service.NextUsingCode(id, codeId)
	if codeInfo != nil && codeInfo.Id > 0 {
		codeInfo.SysStatus = 2
		codeInfo.SysUpdated = comm.NowUnix()
		service.Update(codeInfo, nil)
	} else {
		log.Println("prizedata.prizecodediff num codeinfo,gift_id=", id)
		return ""
	}
	return codeInfo.Code
}

package utils

import (
	"log"
	"lottery/comm"
	"lottery/services"
)

func PrizeGift(id, leftNum int) bool {
	giftService := services.NewGiftService()
	rows, err := giftService.DecrLeftNum(id, 1)
	if rows < 1 || err != nil {
		log.Println("prizedata.PrizeGift gitservice.decrleftnum error=", err)
		return false
	}
	return true
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

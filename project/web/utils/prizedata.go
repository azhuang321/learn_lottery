package utils

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/datasource"
	"lottery/models"
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
	return prizeServCodeDiff(id, service)
}

func PrizeLocalCodeDiff(id int, service services.CodeService) string {
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

func ImportCacheCodes(id int, code string) bool {
	// 集群版本需要放到redis中
	// [暂时]本机版本直接从数据库中处理
	// redis 中缓存 key 值
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("SADD", key, code)
	if err != nil {
		log.Println("prizedata.Recachecodes sadd error = ", err)
		return false
	}
	return true
}

// 重新整理优惠券的编码到缓存中
func RecacheCodes(id int, codeService services.CodeService) (sucNum, errNum int) {
	// 集群版本需要放入到redis中
	// [暂时]本机版本的就直接从数据库中处理吧
	list := codeService.Search(id)
	if list == nil || len(list) <= 0 {
		return 0, 0
	}
	// redis中缓存的key值
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	tmpKey := "tmp_" + key
	for _, data := range list {
		if data.SysStatus == 0 {
			code := data.Code
			_, err := cacheObj.Do("SADD", tmpKey, code)
			if err != nil {
				log.Println("prizedata.RecacheCodes SADD error=", err)
				errNum++
			} else {
				sucNum++
			}
		}
	}
	_, err := cacheObj.Do("RENAME", tmpKey, key)
	if err != nil {
		log.Println("prizedata.RecacheCodes RENAME error=", err)
	}
	return sucNum, errNum
}

func GetCacheCodeNum(id int, service services.CodeService) (int, int) {
	num := 0
	cacheNum := 0
	// 统计数据库中有效的编码数量
	list := service.Search(id)
	if len(list) > 0 {
		for _, data := range list {
			if data.SysStatus == 0 {
				num++
			}
		}
	}

	// redis 中缓存的key值
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("SCARD", key)
	if err != nil {
		log.Println("prizedata.recachecodes rename err = ", err)
	} else {
		cacheNum = int(comm.GetInt64(rs, 0))
	}
	return num, cacheNum
}

func prizeServCodeDiff(id int, service services.CodeService) string {
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("SPOP", key)
	if err != nil {
		log.Println("prizedata.prizeservcodediff err=", err)
		return ""
	}
	code := comm.GetString(rs, "")
	if code == "" {
		log.Println("prizedata.prizeservcodediff err=", err)
		return ""
	}
	service.UpdateByCode(&models.LtCode{
		Code:       code,
		SysUpdated: comm.NowUnix(),
		SysStatus:  2,
	}, nil)
	return code
}

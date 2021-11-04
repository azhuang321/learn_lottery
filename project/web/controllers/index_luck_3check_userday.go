package controllers

import (
	"fmt"
	"log"
	"lottery/conf"
	"lottery/models"
	"lottery/services"
	"lottery/web/utils"
	"strconv"
	"time"
)

func (c *IndexController) checkUserDay(uid int, num int64) bool {
	userDayService := services.NewUserdayService()
	userDayInfo := userDayService.GetUserToday(uid)
	if userDayInfo != nil && userDayInfo.Uid == uid {
		// 今天存在抽奖记录
		if userDayInfo.Num >= conf.UserPrizeMax {
			if int(num) < userDayInfo.Num {
				utils.InitUserLuckyNum(uid, int64(userDayInfo.Num))
			}
			return false
		} else {
			userDayInfo.Num++
			if int(num) < userDayInfo.Num {
				utils.InitUserLuckyNum(uid, int64(userDayInfo.Num))
			}
			err103 := userDayService.Update(userDayInfo, nil)
			if err103 != nil {
				log.Println("index_lucky_check_userday ServiceUserDay.Update "+
					"err103=", err103)
			}
		}
	} else {
		// 创建今天的用户参与记录
		y, m, d := time.Now().Date()
		strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
		day, _ := strconv.Atoi(strDay)
		userDayInfo = &models.LtUserday{
			Uid:        uid,
			Day:        day,
			Num:        1,
			SysCreated: int(time.Now().Unix()),
		}
		err103 := userDayService.Create(userDayInfo)
		if err103 != nil {
			log.Println("index_lucky_check_userday ServiceUserDay.Create "+
				"err103=", err103)
		}
		utils.InitUserLuckyNum(uid, 1)
	}
	return true
}

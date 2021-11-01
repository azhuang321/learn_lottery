package controllers

import (
	"fmt"
	"lottery/comm"
	"lottery/conf"
	"lottery/models"
	"lottery/web/utils"
)

func (c *IndexController) GetLucky() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	// 1.验证登录
	loginUser := comm.GetLoginUser(c.Ctx.Request())
	if loginUser == nil || loginUser.Uid < 1 {
		rs["code"] = 101
		rs["msg"] = "请登录"
		return rs
	}
	// 2.用户分抽奖布式锁
	ok := utils.LockLucky(loginUser.Uid)
	if !ok {
		rs["code"] = 102
		rs["msg"] = "正在抽奖中,请稍后重试"
		return rs
	}
	defer utils.UnLockLucky(loginUser.Uid)
	// 3.验证用户今日参与次数
	ok = c.checkUserDay(loginUser.Uid)
	if !ok {
		rs["code"] = 103
		rs["msg"] = "今日抽奖次数已用完,明日再来吧"
		return rs
	}
	// 4.验证今日的参与次数
	ipStr := comm.ClientIp(c.Ctx.Request())
	ipDayNum := utils.IncrIpLuckyNum(ipStr)
	if ipDayNum > conf.IpLimitMax {
		rs["code"] = 104
		rs["msg"] = "相同ip参与次数过多,明天再来参与吧"
		return rs
	}
	// 5.验证IP黑名单
	limitBlack := false
	if ipDayNum > conf.IpPrizeMax {
		limitBlack = true
	}
	var blackIpInfo *models.LtBlackip
	if !limitBlack {
		ok, blackIpInfo = c.checkBlackIp(ipStr)
		if !ok {
			fmt.Println("黑名单中的IP", ipStr, limitBlack)
			limitBlack = true
		}
	}
	// 6.验证黑名单用户
	var userInfo *models.LtUser
	if !limitBlack {
		ok, userInfo = c.checkBlackUser(loginUser.Uid)
		if !ok {
			fmt.Println("黑名单中的用户", loginUser.Uid, limitBlack)
			limitBlack = true
		}
	}
}

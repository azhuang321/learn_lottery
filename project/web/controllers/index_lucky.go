package controllers

import (
	"fmt"
	"log"
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
	userDayNum := utils.IncrUserLuckyNum(loginUser.Uid)
	if userDayNum > conf.UserPrizeMax {
		rs["code"] = 103
		rs["msg"] = "今日抽奖次数已用完,明日再来吧"
		return rs
	} else {
		ok = c.checkUserDay(loginUser.Uid, userDayNum)
		if !ok {
			rs["code"] = 103
			rs["msg"] = "今日抽奖次数已用完,明日再来吧"
			return rs
		}
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
	// 7.获得抽奖编码
	prizeCode := comm.Random(10000)
	// 8.匹配奖品是否中奖
	prizeGift := c.prize(prizeCode, limitBlack)
	if prizeGift == nil || prizeGift.PrizeNum < 0 || (prizeGift.PrizeNum > 0 && prizeGift.LeftNum <= 0) {
		rs["code"] = 205
		rs["msg"] = "很遗憾,没有中奖,请下次再试"
		return rs
	}
	// 9.有限制奖品发放
	if prizeGift.PrizeNum > 0 {
		ok = utils.PrizeGift(prizeGift.Id, prizeGift.LeftNum)
		if !ok {
			rs["code"] = 207
			rs["msg"] = "很遗憾,没有中奖,请下次再试"
			return rs
		}
	}
	// 10.不同编码的优惠券发放
	if prizeGift.Gtype == conf.GtypeCodeDiff {
		code := utils.PrizeCodeDiff(prizeGift.Id, c.ServiceCode)
		if code == "" {
			rs["code"] = 208
			rs["msg"] = "很遗憾,没有中奖,请下次再试"
			return rs
		}
		prizeGift.Gdata = code
	}
	// 11.记录中奖记录
	result := &models.LtResult{
		GiftId:     prizeGift.Id,
		GiftName:   prizeGift.Title,
		GiftType:   prizeGift.Gtype,
		Uid:        loginUser.Uid,
		Username:   loginUser.Username,
		PrizeCode:  prizeCode,
		GiftData:   prizeGift.Gdata,
		SysCreated: comm.NowUnix(),
		SysIp:      ipStr,
		SysStatus:  0,
	}
	err := c.ServiceResult.Create(result)
	if err != nil {
		log.Println("index_lucky.GetLucky serviceresult.create", result, "err", err)
		rs["code"] = 209
		rs["msg"] = "很遗憾.没有中奖,请下次再试"
		return rs
	}
	if prizeGift.Gtype == conf.GtypeGiftLarge {
		// 如果是实物大奖,需要将用户,ip设置为黑名单一段时间
		c.prizeLarge(ipStr, loginUser, userInfo, blackIpInfo)
	}

	// 12.返回抽奖结果
	rs["gift"] = prizeGift
	return rs
}

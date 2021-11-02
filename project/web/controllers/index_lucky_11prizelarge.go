package controllers

import (
	"lottery/comm"
	"lottery/models"
)

func (c *IndexController) prizeLarge(ip string, loginUser *models.ObjLoginuser, userInfo *models.LtUser, blackIpInfo *models.LtBlackip) {
	nowTime := comm.NowUnix()
	blackTime := 30 * 86400
	//更新用户黑名单信息
	if userInfo == nil || userInfo.Id <= 0 {
		userInfo = &models.LtUser{
			Id:         loginUser.Uid,
			Username:   loginUser.Username,
			Blacktime:  nowTime + blackTime,
			SysCreated: nowTime,
			SysIp:      ip,
		}
		c.ServiceUser.Create(userInfo)
	} else {
		userInfo = &models.LtUser{Id: loginUser.Uid, Blacktime: nowTime + blackTime, SysUpdated: nowTime}
		c.ServiceUser.Update(userInfo, nil)
	}
	// 更新黑名单信息
	if blackIpInfo == nil || blackIpInfo.Id <= 0 {
		blackIpInfo = &models.LtBlackip{
			Ip:         ip,
			Blacktime:  nowTime + blackTime,
			SysCreated: nowTime,
		}
		c.ServiceBlackip.Create(blackIpInfo)
	} else {
		blackIpInfo.Blacktime = nowTime + blackTime
		blackIpInfo.SysUpdated = nowTime
		c.ServiceBlackip.Update(blackIpInfo, nil)
	}
}

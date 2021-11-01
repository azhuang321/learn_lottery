package controllers

import (
	"lottery/models"
	"time"
)

func (c *IndexController) checkBlackIp(ip string) (bool, *models.LtBlackip) {
	info := c.ServiceBlackip.GetByIp(ip)
	if info == nil || info.Ip == "" {
		return true, nil
	}
	if info.Blacktime > int(time.Now().Unix()) {
		return false, info
	}
	return true, info
}

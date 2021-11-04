package controllers

import (
	"lottery/conf"
	"lottery/models"
)

func (c *IndexController) prize(prizeCode int, limitBlack bool) *models.ObjGiftPrize {
	var prizeGift *models.ObjGiftPrize
	giftList := c.ServiceGift.GetAllUse(true)
	for _, gift := range giftList {
		if gift.PrizeCodeA <= prizeCode && gift.PrizeCodeB >= prizeCode {
			// 中奖编码满足条件,可以中奖
			if !limitBlack || gift.Gtype < conf.GtypeGiftSmall {
				prizeGift = &gift
				break
			}
		}
	}
	return prizeGift
}

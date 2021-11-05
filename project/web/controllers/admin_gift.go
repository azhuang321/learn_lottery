package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"lottery/comm"
	"lottery/models"
	"lottery/web/utils"
	"lottery/web/viewmodels"
	"time"

	"lottery/services"
)

type AdminGiftController struct {
	Ctx            iris.Context
	ServiceUser    services.UserService
	ServiceGift    services.GiftService
	ServiceCode    services.CodeService
	ServiceResult  services.ResultService
	ServiceUserday services.UserdayService
	ServiceBlackip services.BlackipService
}

func (c *AdminGiftController) Get() mvc.Result {
	dataList := c.ServiceGift.GetAll(true)
	total := len(dataList)
	for i, giftInfo := range dataList {
		//奖品发放的计划数据
		prizeData := make([][2]int, 0)
		err := json.Unmarshal([]byte(giftInfo.PrizeData), &prizeData)
		if err != nil || len(prizeData) < 1 {
			dataList[i].PrizeData = "[]"
		} else {
			newPd := make([]string, len(prizeData))
			for index, pd := range prizeData {
				ct := comm.FormatFromUnixTime(int64(pd[0]))
				newPd[index] = fmt.Sprintf(" [%s] : %d ", ct, pd[1])
			}
			str, err := json.Marshal(newPd)
			if err == nil && len(str) > 0 {
				dataList[i].PrizeData = string(str)
			} else {
				dataList[i].PrizeData = "[]"
			}
		}
	}
	return mvc.View{
		Name: "admin/gift.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "gift",
			"Datalist": dataList,
			"Total":    total,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminGiftController) GetEdit() mvc.Result {
	id := c.Ctx.URLParamIntDefault("id", 0)
	giftInfo := viewmodels.ViewGift{}
	if id > 0 {
		data := c.ServiceGift.Get(id, true)
		giftInfo.Id = data.Id
		giftInfo.Title = data.Title
		giftInfo.PrizeNum = data.PrizeNum
		giftInfo.PrizeCode = data.PrizeCode
		giftInfo.PrizeTime = data.PrizeTime
		giftInfo.Img = data.Img
		giftInfo.Displayorder = data.Displayorder
		giftInfo.Gtype = data.Gtype
		giftInfo.Gdata = data.Gdata
		giftInfo.TimeBegin = comm.FormatFromUnixTime(int64(data.TimeBegin))
		giftInfo.TimeEnd = comm.FormatFromUnixTime(int64(data.TimeEnd))
	}
	return mvc.View{
		Name: "admin/giftEdit.html",
		Data: iris.Map{
			"Title":   "管理后台",
			"Channel": "gift",
			"info":    giftInfo,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminGiftController) PostSave() mvc.Result {
	data := &viewmodels.ViewGift{}
	err := c.Ctx.ReadForm(data)
	if err != nil {
		fmt.Println("admin_gift.PostSave readFrom err = ", err)
		return mvc.Response{
			Text: fmt.Sprintf("ReadFrom 转换异常,err = %s", err),
		}
	}
	giftInfo := models.LtGift{}
	giftInfo.Id = data.Id
	giftInfo.Title = data.Title
	giftInfo.PrizeNum = data.PrizeNum
	giftInfo.PrizeCode = data.PrizeCode
	giftInfo.PrizeTime = data.PrizeTime
	giftInfo.Img = data.Img
	giftInfo.Displayorder = data.Displayorder
	giftInfo.Gtype = data.Gtype
	giftInfo.Gdata = data.Gdata
	t1, err1 := comm.ParseTime(data.TimeBegin)
	t2, err2 := comm.ParseTime(data.TimeEnd)
	if err1 != nil || err2 != nil {
		return mvc.Response{
			Text: fmt.Sprintf("开始结束时间不正确,err1 = %s,err2 = %s", err1, err2),
		}
	}
	giftInfo.TimeBegin = int(t1.Unix())
	giftInfo.TimeEnd = int(t2.Unix())

	if giftInfo.Id > 0 {
		//更新
		dataInfo := c.ServiceGift.Get(giftInfo.Id, true)
		if dataInfo != nil && dataInfo.Id > 0 {
			//奖品数量发生变化
			giftInfo.LeftNum = dataInfo.LeftNum - dataInfo.PrizeNum - giftInfo.PrizeNum
			if giftInfo.LeftNum < 0 || giftInfo.PrizeNum <= 0 {
				giftInfo.LeftNum = 0
			}
			// 奖品总数发生变化
			utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
		}
		if dataInfo.PrizeTime != giftInfo.PrizeTime {
			// 发奖周期变化
			utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
		}
		giftInfo.SysUpdated = int(time.Now().Unix())
		c.ServiceGift.Update(&giftInfo, []string{""})
	} else {
		giftInfo.Id = 0
	}

	if giftInfo.Id == 0 {
		giftInfo.LeftNum = giftInfo.PrizeNum
		giftInfo.SysIp = comm.ClientIp(c.Ctx.Request())
		giftInfo.SysCreated = int(time.Now().Unix())
		c.ServiceGift.Create(&giftInfo)
		// 新的奖品,更新奖品的发奖计划
		utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
	}

	return mvc.Response{
		Path: "/admin/gift",
	}
}

func (c *AdminGiftController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceGift.Delete(id)
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

func (c *AdminGiftController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceGift.Update(&models.LtGift{Id: id, SysStatus: 0}, []string{"sys_status"})
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

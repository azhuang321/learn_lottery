package controllers

import (
	"fmt"
	"lottery/comm"
	"lottery/models"
	"lottery/services"

	"github.com/kataras/iris/v12"
)

type IndexController struct {
	Ctx            iris.Context
	ServiceUser    services.UserService
	ServiceGift    services.GiftService
	ServiceCode    services.CodeService
	ServiceResult  services.ResultService
	ServiceUserday services.UserdayService
	ServiceBlackip services.BlackipService
}

// http://localhost:8080/
func (c *IndexController) Get() string {
	c.Ctx.Header("Content-Type", "text/html")
	return "欢迎来到抽奖系统,<a href='/public/index.html'>开始抽奖</a>"
}

// http://localhost:8080/gifts
func (c *IndexController) GetGifts() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	datalist := c.ServiceGift.GetAll(false)
	list := make([]models.LtGift, 0)

	for _, data := range datalist {
		if data.SysStatus == 0 {
			list = append(list, data)
		}
	}
	rs["gifts"] = list
	return rs
}

// http://localhost:8080/newprize
func (c *IndexController) GetNewprize() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""

	return rs
}

func (c *IndexController) GetLogin() {
	uid := comm.Random(100000)
	loginUser := &models.ObjLoginuser{
		Uid:      uid,
		Username: fmt.Sprintf("admin-%d", uid),
		Now:      comm.NowUnix(),
		Ip:       comm.ClientIp(c.Ctx.Request()),
	}
	comm.SetLoginUser(c.Ctx.ResponseWriter(), loginUser)
	comm.Redirect(c.Ctx.ResponseWriter(), "/public/index.html?from=login")
}

func (c *IndexController) GetLogout() {
	comm.SetLoginUser(c.Ctx.ResponseWriter(), nil)
	comm.Redirect(c.Ctx.ResponseWriter(), "/public/index.html?from=logout")
}

package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"math/rand"
	"time"
)

var packageList map[uint32][]uint = make(map[uint32][]uint)

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

func (c *lotteryController) Get() map[uint32][2]int {
	rs := make(map[uint32][2]int)
	for id, list := range packageList {
		var money int
		for _, v := range list {
			money += int(v)
		}
		rs[id] = [2]int{len(list), money}
	}
	return rs
}

// 发红包
// http://localhost:8080/set?uid=1&money=28.88&num=19
func (c *lotteryController) GetSet() string {
	uid, errUid := c.Ctx.URLParamInt("uid")
	money, errMoney := c.Ctx.URLParamFloat64("money")
	num, errNum := c.Ctx.URLParamInt("num")
	if errUid != nil || errMoney != nil || errNum != nil {
		return fmt.Sprintf("参数异常,errUid=%d,errMoney=%d,errNum=%d\n", errUid, errMoney, errNum)
	}
	moneyTotal := int(money * 100)
	if uid < 1 || moneyTotal < num || num < 1 {
		return fmt.Sprintf("参数异常,uid=%d,money=%f,num=%d\n", uid, money, num)
	}
	// 金额分配算法
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rMax := 0.55
	if num > 1000 {
		rMax = 0.01
	} else if num >= 100 {
		rMax = 0.1
	} else if num >= 10 {
		rMax = 0.3
	}
	list := make([]uint, num)
	leftMoney := moneyTotal
	leftNum := num
	// 开始分配金额
	for leftNum > 0 {
		if leftNum == 1 { // 最后一个红包,剩余的金额给他
			list[num-1] = uint(leftMoney)
			break
		}
		if leftMoney == leftNum {
			for i := num - leftNum; i < num; i++ {
				list[i] = 1
			}
			break
		}
		rMoney := int(float64(leftMoney-leftNum) * rMax) // 保证最后剩余红包中有余额
		m := r.Intn(rMoney)
		if m < 1 {
			m = 1
		}
		list[num-leftNum] = uint(m)
		leftMoney -= m
		leftNum--
	}
	// 红包唯一ID
	id := r.Uint32()
	packageList[id] = list
	// 返回抢红包的地址
	return fmt.Sprintf("/get?id=%d&uid=%d&num=%d", id, uid, num)
}

// 抢红包
// http://localhost:8080/get?uid=1&id=1
func (c *lotteryController) GetGet() string {
	uid, errUid := c.Ctx.URLParamInt("uid")
	id, errId := c.Ctx.URLParamInt("num")
	if errUid != nil || errId != nil {
		return fmt.Sprintf("参数异常,errUid=%d,errId=%d\n", uid, id)
	}
	if uid < 1 || id < 1 {
		return fmt.Sprintf("参数异常,uid=%d,id=%d\n", uid, id)
	}
	list, ok := packageList[uint32(id)]
	if !ok || len(list) < 1 {
		return fmt.Sprintf("红包不存在:%d", id)
	}
	// 分配随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(len(list))
	money := list[i]
	// 更新红包列表中的信息
	if len(list) > 1 {
		if i == len(list)-1 {
			packageList[uint32(id)] = list[:i]
		} else if i == 0 {
			packageList[uint32(id)] = list[1:]
		} else {
			packageList[uint32(id)] = append(list[0:i], list[i+1:]...)
		}
	} else {
		delete(packageList, uint32(id))
	}
	return fmt.Sprintf("恭喜你,抢到 `%d` 的红包", money)
}

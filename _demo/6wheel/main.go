package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// 奖品中奖概率
type Prate struct {
	Rate  int    //万分之N的中奖概率
	Total int    //总数量  0表示无限
	CodeA int    //中奖概率起始编码(包含)
	CodeB int    //中奖概率终止编码(包含)
	Left  *int32 //剩余数
}

// 奖品列表
var prizeList []string = []string{
	"一等奖,火星单程票",
	"二等奖,南极之旅",
	"三等奖,iphone",
	"", //未中奖
}

// 中奖的概率设置,与上面的prizeList 对应的设置
var left int32 = 1000
var rateList []Prate = []Prate{
	{100, 1000, 0, 9999, &left},
	//{2, 2, 1, 2, 2},
	//{5, 10, 3, 5, 10},
	//{100, 0, 0, 9999, 0},
}

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

// http://localhost:8080
func (c *lotteryController) Get() string {
	c.Ctx.Header("Content-Type", "text/html")
	return fmt.Sprintf("大转盘奖品列表:<br/> %s", strings.Join(prizeList, "<br/>"))
}

func (c *lotteryController) GetDebug() string {
	return fmt.Sprintf("获奖概率: %v\n", rateList)
}

func (c *lotteryController) GetPrize() string {
	// 第一步抽奖,更具随机数匹配奖品
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	code := r.Intn(10000)

	var myPrize string
	var prizeRate *Prate
	// 奖品列表中匹配是否中奖
	for i, prize := range prizeList {
		rate := &rateList[i]
		if code >= rate.CodeA && code <= rate.CodeB {
			// 中奖
			myPrize = prize
			prizeRate = rate
			break
		}
	}
	if myPrize == "" {
		myPrize = "未中奖"
		return myPrize
	}
	// 第二部 中奖了,开始发奖
	if prizeRate.Total == 0 {
		return myPrize
	} else if *prizeRate.Left > 0 {
		left := atomic.AddInt32(prizeRate.Left, -1)
		if left >= 0 {
			return myPrize
		}
	}
	myPrize = "未中奖"
	return myPrize

}

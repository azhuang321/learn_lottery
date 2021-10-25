/**
1 即开即得型
2 双色球自选型
*/

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

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

// 1.即开即得型 http://localhost:8080/
func (c *lotteryController) Get() string {
	prize := ""
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Intn(10)

	switch {
	case code == 1:
		prize = "一等奖"
	case code > 1 && code <= 3:
		prize = "二等奖"
	case code > 3 && code <= 6:
		prize = "三等奖"
	default:
		prize = "未中奖"
	}
	return fmt.Sprintf("尾号为1获得`一等奖`\n尾号为2或者3获得`二等奖`\n尾号为4/5/6获得`三等奖`\ncode=`%d`\n中奖信息:`%s`", code, prize)
}

// 2.双色球 http://localhost:8080/prize
func (c *lotteryController) GetPrize() string {
	var prize [7]int
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	// 6 位红球
	for i := 0; i < 6; i++ {
		prize[i] = r.Intn(33) + 1
	}
	// 最后一位篮球
	prize[6] = r.Intn(16) + 1
	return fmt.Sprintf("今日开奖号码: %v", prize)
}

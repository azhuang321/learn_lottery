/**
curl http://localhost:8080/
curl --data "users=test1,test2,test3"  http://localhost:8080/import
curl http://localhost:8080/lucky
*/

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

var userList []string

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
	userList = []string{}

	app.Run(iris.Addr(":8080"))
}

func (c *lotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("当前总共参加抽奖人数: %d\n", count)
}

// POST http://localhost:8080/import
// paras: users
func (c *lotteryController) PostImport() string {
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strings.TrimSpace(strUsers), ",")
	count1 := len(userList)
	for _, v := range users {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			userList = append(userList, v)
		}
	}
	count2 := len(userList)

	return fmt.Sprintf("当前总共参与抽奖人数: %d,成功导入用户数: %d", count2, count2-count1)
}

// GET http://localhost:8080/lucky
func (c *lotteryController) GetLucky() string {
	count := len(userList)
	user := ""
	if count == 0 {
		return fmt.Sprintf("已经没有用户参与抽奖,请先通过 /import 导入用户 \n")
	}
	if count == 1 {
		user = userList[0]
	} else {
		seed := time.Now().UnixNano()
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		user = userList[index]
		userList = append(userList[0:index], userList[index+1:]...)
	}
	return fmt.Sprintf("当前中奖用户: %s ,剩余用户数: %d \n", user, count-1)

}

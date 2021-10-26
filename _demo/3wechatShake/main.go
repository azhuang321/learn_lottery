/**
 * 微信摇一摇
 * 基础功能:
 * http://localhost:8080/lucky 只有一个抽奖的接口
 * 压力测试:
 * wrk -t10 -c10 -d5 http://localhost:8080/lucky (-t:线程数 -c:连接数 -d:持续时间)
 */
package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// 奖品类型,枚举:
const (
	giftTypeCoin      = iota //虚拟币
	giftTypeCoupon           //不同券
	giftTypeCouponFix        //相同券
	giftTypeRealSmall        //实物小奖
	giftTypeRealLarge        //实物大奖
)

type gift struct {
	id       int      //奖品id
	name     string   //奖品名称
	pic      string   //奖品图片
	link     string   //奖品链接
	gType    int      //奖品类型
	data     string   //奖品数据(特定的配置信息)
	dataList []string //奖品数据集合
	total    int      //奖品总数 0 为不限量
	left     int      //奖品剩余数量
	inUse    bool     //是否使用中
	rate     int      //中奖概率  万分之N 0-9999
	rateMin  int      //大于等于最小中奖编码
	rateMax  int      //小于中奖编码
}

// 最大中奖号码
const rateMax = 10000

var logger *log.Logger
var mux sync.Mutex

// 奖品列表
var giftList []*gift

type lotteryController struct {
	Ctx iris.Context
}

// 初始化日志
func initLog() {
	f, _ := os.Create("./log/lottery_demo.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}

//初始化奖品列表
func initGift() {
	giftList = make([]*gift, 5)
	giftList[0] = &gift{
		id:       1,
		name:     "手机大奖",
		pic:      "",
		link:     "",
		gType:    giftTypeRealLarge,
		data:     "",
		dataList: nil,
		total:    20000,
		left:     20000,
		inUse:    true,
		rate:     10000,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[1] = &gift{
		id:       2,
		name:     "充电器",
		pic:      "",
		link:     "",
		gType:    giftTypeRealSmall,
		data:     "",
		dataList: nil,
		total:    5,
		left:     5,
		inUse:    false,
		rate:     10,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[2] = &gift{
		id:       3,
		name:     "优惠券满200减50",
		pic:      "",
		link:     "",
		gType:    giftTypeCouponFix,
		data:     "mall-coupon-2018",
		dataList: nil,
		total:    50,
		left:     50,
		inUse:    false,
		rate:     500,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[3] = &gift{
		id:       4,
		name:     "直降优惠券50元",
		pic:      "",
		link:     "",
		gType:    giftTypeCoupon,
		data:     "",
		dataList: []string{"c01", "c02", "c03", "c04", "c05", "c06", "c07", "c08", "c09", "c10"},
		total:    10,
		left:     10,
		inUse:    false,
		rate:     100,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[4] = &gift{
		id:       5,
		name:     "金币",
		pic:      "",
		link:     "",
		gType:    giftTypeCoin,
		data:     "10金币",
		dataList: nil,
		total:    5,
		left:     5,
		inUse:    false,
		rate:     5000,
		rateMin:  0,
		rateMax:  0,
	}
	// 整理数据,中奖区间数据
	rateStart := 0
	for _, data := range giftList {
		if !data.inUse {
			continue
		}
		data.rateMin = rateStart
		data.rateMax = rateStart + data.rate
		if data.rateMax > rateMax {
			data.rateMax = rateMax
			rateStart = 0
		} else {
			rateStart += data.rate
		}
		fmt.Printf("%+v\n", data)
	}

}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})

	initLog()
	initGift()

	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

// 奖品的数量信息 GET http://localhost:8080
func (c *lotteryController) Get() string {
	count := 0
	total := 0
	for _, data := range giftList {
		if data.inUse && (data.total == 0 || (data.total > 0 && data.left > 0)) {
			count++
			total += data.left
		}
	}
	return fmt.Sprintf("当前有效奖品种类数量:%d,限量奖品总数=%d\n", count, total)
}

// 抽奖 http://localhost:8080/lucky
func (c *lotteryController) GetLucky() map[string]interface{} {
	mux.Lock()
	defer mux.Unlock()
	code := luckyCode()

	ok := false
	result := make(map[string]interface{})
	result["success"] = ok

	for _, data := range giftList {
		if !data.inUse && (data.total == 0 || (data.total > 0 && data.left > 0)) {
			continue
		}
		if data.rateMin <= int(code) && data.rateMax > int(code) {
			// 中奖了,抽奖编码在奖品编码范围内
			// 开始发奖
			sendData := ""
			switch data.gType {
			case giftTypeCoin:
				ok, sendData = sendCoin(data)
			case giftTypeCoupon:
				ok, sendData = sendCoupon(data)
			case giftTypeCouponFix:
				ok, sendData = sendCouponFix(data)
			case giftTypeRealSmall:
				ok, sendData = sendRealSmall(data)
			case giftTypeRealLarge:
				ok, sendData = sendRealLarge(data)
			}
			if ok {
				// 中奖
				saveLuckyData(code, data.id, data.name, data.link, sendData, data.left)
				result["success"] = ok
				result["id"] = data.id
				result["name"] = data.name
				result["link"] = data.link
				result["data"] = sendData
				break
			}
		}
	}
	return result
}

func sendCoin(data *gift) (bool, string) {
	if data.total == 0 { //数量不限
		return true, data.data
	} else if data.left > 0 {
		data.left -= 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

func sendCoupon(data *gift) (bool, string) {
	if data.left > 0 {
		data.left -= 1
		return true, data.dataList[data.left]
	} else {
		return false, "奖品已发完"
	}
}

func sendCouponFix(data *gift) (bool, string) {
	if data.total == 0 { //数量不限
		return true, data.data
	} else if data.left > 0 {
		data.left -= 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

func sendRealSmall(data *gift) (bool, string) {
	if data.total == 0 { //数量不限
		return true, data.data
	} else if data.left > 0 {
		data.left -= 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

func sendRealLarge(data *gift) (bool, string) {
	if data.total == 0 { //数量不限
		return true, data.data
	} else if data.left > 0 {
		data.left -= 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

func luckyCode() int32 {
	seed := time.Now().UnixNano()
	return rand.New(rand.NewSource(seed)).Int31n(rateMax)
}
func saveLuckyData(code int32, id int, name, link, sendData string, left int) {
	logger.Printf("lucky: code=%d,gift=%d,name=%s,data=%s,link=%s,left=%d", code, id, name, sendData, name, left)
}

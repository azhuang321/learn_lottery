package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func TestMVC(t *testing.T) {
	e := httptest.New(t, newApp())

	var wg sync.WaitGroup
	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("当前总共参加抽奖人数: 0\n")

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e.POST("/import").WithFormField("users", fmt.Sprintf("test_u%d", i)).Expect().Status(httptest.StatusOK)
		}(i)
	}
	wg.Wait()

	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("当前总共参加抽奖人数: 100\n")
	e.GET("/lucky").Expect().Status(httptest.StatusOK)
	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("当前总共参加抽奖人数: 99\n")
}

package main

import (
	"fmt"
	"lottery/bootstrap"
	"lottery/web/middleware/identity"
	"lottery/web/routes"
)

var port = 8080

func newApp() *bootstrap.Bootstrapper {
	// 初始化应用
	app := bootstrap.New(
		"GO抽奖系统",
		"azhuang",
	)
	app.Bootstrap()
	app.Configure(
		identity.Configure,
		routes.Configure,
	)
	return app
}

func main() {
	app := newApp()
	app.Listen(fmt.Sprintf(":%d", port))
}

package main

import (
	"github.com/echopairs/skygo/version"
	"github.com/echopairs/skygo/gweb/web"
)

func main() {
	version.Show()
	web.Helloworld()
	//cfg := &web.Config{}
	//web.StartServer(cfg)
}
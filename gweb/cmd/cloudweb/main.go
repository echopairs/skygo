package main

import (
	"flag"
	"fmt"
	"github.com/echopairs/skygo/gweb/web"
	"github.com/echopairs/skygo/version"
	"log"
	"os"
	"os/signal"
)

func main() {
	version.Show()

	var fileName string
	flag.StringVar(&fileName, "c", "gweb.yaml", "gweb config file")
	flag.Parse()
	cfg, err := web.NewConfig(fileName)
	if err != nil {
		log.Fatal("NewConfig failed for ", err)
	}

	err, server := web.StartServer(cfg)
	if err != nil {
		log.Fatal("StartServer failed for ", err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case <- c:
		fmt.Printf("The program receives a stop signal, Waiting to stop ...\n")
		server.Close()
	}
}
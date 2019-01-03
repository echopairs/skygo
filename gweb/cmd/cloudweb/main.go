package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/echopairs/skygo/gweb/web"
	"github.com/echopairs/skygo/version"
)

func main() {
	version.Show()

	var fileName string
	flag.StringVar(&fileName, "c", "gweb.yaml", "gweb config file")
	flag.Parse()
	cfg, err := web.NewConfig(fileName)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	if err != nil {
		log.Fatal("NewConfig failed for ", err)
	}
	log.Print("begin run")

	err, server := web.StartServer(cfg)
	if err != nil {
		log.Fatal("StartServer failed for ", err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		log.Print("The program receives a stop signal, Waiting to stop ...\n")
		server.Close()
	}
}

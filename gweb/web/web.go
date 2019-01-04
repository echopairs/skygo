package web

import (

	"log"
	"io/ioutil"
	"net/http"

	"github.com/echopairs/skygo/gweb/web/auth"
	"github.com/echopairs/skygo/gweb/web/router"
	"github.com/echopairs/skygo/zsql"
	"gopkg.in/yaml.v2"

	_ "github.com/echopairs/skygo/gweb/web/book"
	"github.com/echopairs/skygo/gweb/web/ws"
)

type Config struct {
	HttpConfig    *auth.HttpConfig `yaml:"http_config"`
	SqlAddress    *zsql.SqlAddress `yaml:"sql_address"`
	ServerAddress string           `yaml:"server_addr"`
}

func NewConfig(filename string) (*Config, error) {
	// 1. filename -> []byte
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 2. []byte -> config
	c := &Config{}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func StartServer(cfg *Config) (err error, server *http.Server) {
	db, err := zsql.Connect(cfg.SqlAddress)
	if err != nil {
		log.Printf("create db failed %s\n", err.Error())
		return
	}

	ss, err := auth.NewSessionStorage(cfg.HttpConfig)
	if err != nil {
		log.Printf("create session failed %s\n", err.Error())
		return
	}
	auth.Set(ss, db)

	route := router.GetDefaultRouter()
	server = &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: route,
	}

	err = ws.Start()
	if err != nil {
		log.Printf("ws start error: %v", err)
		return
	}

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Fatalf("listen and serve failed %s \n", err.Error())
		}
	}()

	log.Printf("StartServer ok\n")
	return
}

func init() {
	router.RegisterHttpHandleFunc("GET", "/routes", "routeIndex", router.RouteIndex)
}

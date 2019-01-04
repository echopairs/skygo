package ws

import (
	"github.com/echopairs/skygo/gweb/web/router"

	"net/http"
)

var server = NewWebSocketServer()

func init() {
	router.RegisterHttpHandleFunc("GET", "/ws", "handleWebsocket", handleWebsocket)
}

func Start() error {
	return server.Start()
}

func Stop() {
	server.Stop()
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	server.HandleWebsocket(w, r)
}

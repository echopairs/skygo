package ws

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/echopairs/skygo/gweb/web/common"
)

type WebSocketServer struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcastIn chan []byte

	// outbound message from server to clients
	broadcastOut chan interface{}

	// Register requests from the clients
	register chan *Client

	// Unregister request from the clients
	unregister chan *Client

	// for cancel
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWebSocketServer() *WebSocketServer {
	ctx, cancel := context.WithCancel(context.Background())
	s := &WebSocketServer{
		clients:      make(map[*Client]bool),
		broadcastIn:  make(chan []byte),
		broadcastOut: make(chan interface{}),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		ctx:          ctx,
		cancel:       cancel,
	}
	go s.pump(ctx)
	return s
}

func (server *WebSocketServer) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("create websocket conn failed: %s ", err.Error())
		common.WriteError(w, common.ERR_CREATE_WEBSOCKET_ERROR, http.StatusInternalServerError)
		return
	}
	client := &Client{
		server: server,
		conn:   conn,
		send:   make(chan interface{}),
	}
	server.register <- client
	go client.writePump(server.ctx)
	go client.readPump(server.ctx)
	//common.WriteOk(w)
}

func (server *WebSocketServer) Start() error {
	go server.pushMockMessage(server.ctx)
	return nil
}

func (server *WebSocketServer) Stop() {
	if server.cancel != nil {
		server.cancel()
	}
}

func (server *WebSocketServer) pump(ctx context.Context) {
	for {
		select {
		case client := <-server.register:
			server.clients[client] = true
		case client := <-server.unregister:
			if _, ok := server.clients[client]; ok {
				delete(server.clients, client)
				close(client.send)
			}
		case outMsg := <-server.broadcastIn:
			log.Printf("recv message: %s", string(outMsg))
		case inMsg := <-server.broadcastOut:
			for client := range server.clients {
				select {
				case client.send <- inMsg:
				default:
					log.Println("client is not ready")
				}
			}
		case <-ctx.Done():
			log.Printf("websocket pump cancel")
			break
		}
	}
}

// for test
func (server *WebSocketServer) pushMockMessage(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// todo update data
			data := map[string]interface{}{
				"id":      1,
				"name":    "websocket",
				"mapname": "local.map",
				"x":       1,
				"y":       2,
				"angle":   3,
				"battery": 4,
				"status":  5,
			}
			server.broadcastOut <- data
		case <-ctx.Done():
			break
		}
	}
}

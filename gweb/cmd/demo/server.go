package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

type Client struct {
	id 			string
	socket 		*websocket.Conn
	send 		chan []byte
}

type Server struct {
	clients 	map[*Client]bool
	broadcast   chan []byte
	register 	chan *Client
	unregister  chan *Client
}

type Message struct {
	Sender 		string  `json:"sender,omitempty"`
	Recipient	string 	`json:"recipient,omitempty"`
	Content 	string  `json:"content,omitempty"`
}

var server = Server{
	clients: make(map[*Client]bool),
	broadcast: make(chan []byte),
	register: make(chan *Client),
	unregister: make(chan *Client),
}

func (manager *Server) Start() {
	for {
		select {
		case conn := <- manager.register:
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{
				Content: "/A new socket has connected."})
			manager.send(jsonMessage, conn)
		}
	}
}

func (manager *Server) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (c *Client) read() {
	defer func() {
		server.unregister <- c
		c.socket.Close()
	}()
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			server.unregister <- c
			c.socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{
			Sender: c.id, Content: string(message)})
		server.broadcast <- jsonMessage
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <- c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func main() {
	go server.Start()
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	client := &Client{
		socket: conn,
		send: make(chan []byte),
	}
	server.register <- client
	go client.read()
	go client.write()
}







package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	players []*Client
}

type Client struct {
	conn       *websocket.Conn
	messageBox chan []byte
}

func (c *Client) writePump(s *Server) {
	for {
		message, _ := <-c.messageBox
		_, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			// delete disconnect client
			for i, v := range s.players {
				log.Println(v, c)
				if v == c {
					s.players = append(s.players[:i], s.players[i+1:]...)
					break
				}
			}
			return
		}
		c.conn.WriteMessage(1, message)
		// w, _ := c.conn.NextWriter(websocket.TextMessage)
		// w.Write(message)
	}
}

func (c *Client) readPump(s *Server) {
	for {
		_, message, _ := c.conn.ReadMessage()
		for _, p := range s.players {
			p.messageBox <- message
		}
	}
}

func serveWs(s *Server, w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, _ := upgrader.Upgrade(w, r, nil)
	client := &Client{conn: conn, messageBox: make(chan []byte, 256)}
	s.players = append(s.players, client)
	log.Println(s.players)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.readPump(s)
	go client.writePump(s)
}

func serveHome(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	var addr = flag.String("addr", ":8080", "http service address")

	http.HandleFunc("/", serveHome)
	server := &Server{}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(server, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

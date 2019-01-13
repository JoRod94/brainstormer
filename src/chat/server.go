package chat

import (
  "log"
  "net/http"
  "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Server struct {
  messages []*Message `json:"messages"`
  clients map[int]*Member
  listener chan *Message
} 

func NewServer() {
  Messages := []*Message{}
  clients := make(map[int]*Member)
  listener := make(chan *Message)
  return &Server{
    messages,
    clients,
    listener
  } 
}

func (server *Server) handleConnections(w http.ResponseWriter, r *http.Request){
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Fatal(err)
  }

  defer ws.Close()

  server.clients[ws] = true 

  initErr := ws.WriteJSON(server)
  if initErr != nil {
        ws.Close()
        delete(server.clients, ws)
      }

  for {
    var msg Message
    msgErr := ws.ReadJSON(&msg)
    if msgErr != nil {
      delete(server.clients, ws)
    }
    server.listener <- msg
  }
}

func (server *Server) Run() {
  http.HandleFunc("/ws", server.handleConnections)
  for {
    msg := <- server.listener
    server.Messages = append(server.Messages, msg)
    for cws := range server.clients {
      err := cws.WriteJSON(msg)
      if err != nil {
        cws.Close()
        delete(server.clients, cws)
      }
    }
  }
}
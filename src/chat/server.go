package chat

import (
  "log"
  "net/http"
  "strings"
  "strconv"

  "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Server struct {
  Messages []*Message `json:"messages"`
  Ideas []*Idea `json:"ideas"`
  CurrentIdea int `json:"currentIdea"`
  clients map[*websocket.Conn]bool
  msgListener chan Message
  cmdListener chan Command
}

type IdeaUpdate struct {
  Ideas []*Idea `json:"ideas"`
  CurrentIdea int `json:"currentIdea"`
}

type Command struct {
  key string
  arg string
}

func NewServer() *Server {
  Messages := []*Message{}
  Ideas := []*Idea{}
  CurrentIdea := -1
  clients := make(map[*websocket.Conn]bool)
  msgListener := make(chan Message)
  cmdListener := make(chan Command)
  return &Server{
    Messages,
    Ideas,
    CurrentIdea,
    clients,
    msgListener,
    cmdListener,
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
      break
    }
    server.msgListener <- msg
  }
}

func (server *Server) readCommands(msg Message){
  fields := strings.Fields(msg.Text)
  if(len(fields) > 1){
    server.handleCommand(Command{fields[0], strings.Join(fields[1:], " ")})
  }
}

func (server *Server) handleCommand(command Command) {
  switch(command.key){
    case "/idea":
      ideaNumber, err := strconv.Atoi(command.arg)
      if(err == nil || ideaNumber >= 0 || ideaNumber < len(server.Ideas)){
        server.CurrentIdea = ideaNumber
      }
    case "/newidea":
      server.CurrentIdea++
      newIdea := Idea{
        command.arg, 
        []string{}, 
        []string{}, 
        0,
      }
      server.Ideas = append(server.Ideas, &newIdea)
    case "/why":
      server.Ideas[server.CurrentIdea].Why = append(server.Ideas[server.CurrentIdea].Why, command.arg)
    case "/whynot":
      server.Ideas[server.CurrentIdea].WhyNot = append(server.Ideas[server.CurrentIdea].WhyNot, command.arg)
    case "/vote":
      if(command.arg == "yes"){
        server.Ideas[server.CurrentIdea].Votes++
      } else {
        server.Ideas[server.CurrentIdea].Votes--
      }
  }

  for cws := range server.clients {
    err := cws.WriteJSON(&IdeaUpdate{server.Ideas, server.CurrentIdea})
    if err != nil {
      cws.Close()
      delete(server.clients, cws)
    }
  }
}

func (server *Server) Run() {
  http.HandleFunc("/ws", server.handleConnections)
  for {
    msg := <- server.msgListener
    server.Messages = append(server.Messages, &msg)
    server.readCommands(msg)
    for cws := range server.clients {
      err := cws.WriteJSON(msg)
      if err != nil {
        cws.Close()
        delete(server.clients, cws)
      }
    }
  }
}
package chat

import (
  "log"
  "net/http"
  "strings"
  "strconv"

  "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// JSON tagged fields will be marshalled and sent to the client. Uppercase required
type Server struct {
  Messages []*Message `json:"messages"`
  Ideas []*Idea `json:"ideas"`
  CurrentIdea int `json:"currentIdea"`
  clients map[*websocket.Conn]bool
  msgListener chan Message
  cmdListener chan Command
}

// Struct used to send updates on idea list
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
  // Initiated as -1 so that the first idea is correctly indexed at 0
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

// Accepts new connections and initiates web socket receiver loop
func (server *Server) handleConnections(w http.ResponseWriter, r *http.Request){
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Fatal(err)
  }

  defer ws.Close()

  server.clients[ws] = true 
  initErr := ws.WriteJSON(server)
  // A failed write would indicate that the client is down and should be removed from the list
  if initErr != nil {
        ws.Close()
        delete(server.clients, ws)
      }

  // Web Socket receiver loop
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

// Parses received server messages to determine possible commands
func (server *Server) readCommands(msg Message){
  fields := strings.Fields(msg.Text)
  if(len(fields) > 1){
    server.handleCommand(Command{fields[0], strings.Join(fields[1:], " ")})
  }
}

// Detects commands in messages, executing them if present
func (server *Server) handleCommand(command Command) {
  switch(command.key){
    // Switch to another idea
    case "/idea":
      ideaNumber, err := strconv.Atoi(command.arg)
      if(err == nil || ideaNumber >= 0 || ideaNumber < len(server.Ideas)){
        server.CurrentIdea = ideaNumber
      }
    // Create new idea
    case "/newidea":
      server.CurrentIdea++
      newIdea := Idea{
        command.arg, 
        []string{}, 
        []string{}, 
        0,
      }
      server.Ideas = append(server.Ideas, &newIdea)
    // Add idea pros
    case "/why":
      server.Ideas[server.CurrentIdea].Why = append(server.Ideas[server.CurrentIdea].Why, command.arg)
    // Add idea cons
    case "/whynot":
      server.Ideas[server.CurrentIdea].WhyNot = append(server.Ideas[server.CurrentIdea].WhyNot, command.arg)
    // Add vote. Only "yes" is considered a positive vote
    case "/vote":
      if(command.arg == "yes"){
        server.Ideas[server.CurrentIdea].Votes++
      } else {
        server.Ideas[server.CurrentIdea].Votes--
      }
  }

  // After changing the idea list, send new state to the client
  for cws := range server.clients {
    err := cws.WriteJSON(&IdeaUpdate{server.Ideas, server.CurrentIdea})
    if err != nil {
      cws.Close()
      delete(server.clients, cws)
    }
  }
}

// Main server exectution loop
func (server *Server) Run() {
  // Sets WebSocket handler for new connections
  http.HandleFunc("/ws", server.handleConnections)

  // Receives channel messages from connection goroutines, processing them
  // NOTE: Concurrency revision is needed
  for {
    msg := <- server.msgListener
    server.Messages = append(server.Messages, &msg)
    server.readCommands(msg)
    // After processing a new message, send it to all clients to update their state
    for cws := range server.clients {
      err := cws.WriteJSON(msg)
      if err != nil {
        cws.Close()
        delete(server.clients, cws)
      }
    }
  }
}
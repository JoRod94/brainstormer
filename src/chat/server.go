package chat

import (
  "log"
  "net/http"
  "strings"
  "strconv"
  "sync"

  "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// JSON tagged fields will be marshalled and sent to the client. Uppercase required
type Server struct {
  Messages []*Message `json:"messages"`
  ideasLock *sync.Mutex
  Ideas []*Idea `json:"ideas"`
  CurrentIdea int `json:"currentIdea"`
  clientsLock *sync.Mutex
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
  ideasLock := &sync.Mutex{}
  Ideas := []*Idea{}
  // Initiated as -1 so that the first idea is correctly indexed at 0
  CurrentIdea := -1
  clientsLock := &sync.Mutex{}
  clients := make(map[*websocket.Conn]bool)
  msgListener := make(chan Message)
  cmdListener := make(chan Command)
  return &Server{
    Messages,
    ideasLock,
    Ideas,
    CurrentIdea,
    clientsLock,
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
  server.clientsLock.Lock()
  server.clients[ws] = true
  server.clientsLock.Unlock()
  initErr := ws.WriteJSON(server)
  // A failed write would indicate that the client is down and should be removed from the list
  if initErr != nil {
        ws.Close()
        server.clientsLock.Lock()
        delete(server.clients, ws)
        server.clientsLock.Unlock()
      }

  // Web Socket receiver loop
  for {
    var msg Message
    msgErr := ws.ReadJSON(&msg)
    if msgErr != nil {
      server.clientsLock.Lock()
      delete(server.clients, ws)
      server.clientsLock.Unlock()
      break
    }
    server.msgListener <- msg
  }
}

// Parses received server messages to determine possible commands
func (server *Server) readCommand(msg Message){
  fields := strings.Fields(msg.Text)
  if(len(fields) > 1){
    server.handleCommand(Command{fields[0], strings.Join(fields[1:], " ")}, msg.Username)
  }
}

// Detects commands in messages, executing them if present
func (server *Server) handleCommand(command Command, username string) {
  noIdeas := server.CurrentIdea < 0
  switch(command.key){
    // Create new idea
    case "/newidea":
      server.ideasLock.Lock()
      server.CurrentIdea++
      newIdea := Idea{
        command.arg, 
        []string{}, 
        []string{}, 
        make(map[string]bool),
      }
      server.Ideas = append(server.Ideas, &newIdea)
      server.ideasLock.Unlock()
    // Switch to another idea
    case "/idea":
      if(noIdeas){ break }
      server.ideasLock.Lock()
      ideaNumber, err := strconv.Atoi(command.arg)
      if(err == nil || ideaNumber >= 0 || ideaNumber < len(server.Ideas)){
        server.CurrentIdea = ideaNumber
      }
      server.ideasLock.Unlock()
    // Add idea pros
    case "/why":
      if(noIdeas){ break }
      server.ideasLock.Lock()
      server.Ideas[server.CurrentIdea].Why = append(server.Ideas[server.CurrentIdea].Why, command.arg)
      server.ideasLock.Unlock()
    // Add idea cons
    case "/whynot":
      if(noIdeas){ break }
      server.ideasLock.Lock()
      server.Ideas[server.CurrentIdea].WhyNot = append(server.Ideas[server.CurrentIdea].WhyNot, command.arg)
      server.ideasLock.Unlock()
    // Add vote. Only "yes" is considered a positive vote
    case "/vote":
      if(noIdeas){ break }
      server.ideasLock.Lock()
      if(command.arg == "yes"){
        server.Ideas[server.CurrentIdea].Votes[username] = true
      } else if(command.arg == "no") {
        server.Ideas[server.CurrentIdea].Votes[username] = false
      }
      server.ideasLock.Unlock()
  }

  // After changing the idea list, send new state to the client
  for cws := range server.clients {
    err := cws.WriteJSON(&IdeaUpdate{server.Ideas, server.CurrentIdea})
    if err != nil {
      cws.Close()
      server.clientsLock.Lock()
      delete(server.clients, cws)
      server.clientsLock.Unlock()
    }
  }
}

func (server *Server) handleMessages() {
  // Receives channel messages from connection goroutines, processing them
  for {
    msg := <- server.msgListener
    server.Messages = append(server.Messages, &msg)
    server.readCommand(msg)
    // After processing a new message, send it to all clients to update their state
    for cws := range server.clients {
      err := cws.WriteJSON(msg)
      if err != nil {
        cws.Close()
        server.clientsLock.Lock()
        delete(server.clients, cws)
        server.clientsLock.Unlock()
      }
    }
  }
}

// Main server exectution loop
func (server *Server) Run() {
  // Sets WebSocket handler for new connections
  http.HandleFunc("/ws", server.handleConnections)

  // Starts listening for messages
  server.handleMessages()
}
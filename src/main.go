package main

import (
  "log"
  "net/http"

  "github.com/gorilla/websocket"
)

func main() {
  http.Handle("/", http.FileServer(http.Dir("../public")))
  server := chat.NewServer()
  go server.Run()

  log.Println("http server started on :8000")
  err := http.ListenAndServe(":8000", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
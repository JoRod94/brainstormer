package main

import (
  "log"
  "net/http"
  "./chat"
)

func main() {
  // Redirect file requests to React build folder
  http.Handle("/", http.FileServer(http.Dir("./client/build")))

  // Instantiate and run server
  server := chat.NewServer()
  go server.Run()

  log.Println("http server started on :8000")
  err := http.ListenAndServe(":8000", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
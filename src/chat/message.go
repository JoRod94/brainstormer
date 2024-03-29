package chat

// Chat message. Username is kept here for simplicity, rather than creating a User type
type Message struct {
  Username string `json:"username"`
  Text string `json:"text"`
}
package chat

type Idea struct {
  What string `json:"what"`
  Why []string `json:"why"`
  WhyNot []string `json:"whynot"`
  Votes map[string]bool `json:"votes"`
}


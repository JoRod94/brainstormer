package chat
import (
  "testing"
  "fmt"
  "time"
  "os"
)

var testServer *Server

func TestMain(m *testing.M) {
    // setup tests
    testServer = NewServer()
    go testServer.Run()

    result := m.Run() 

    close(testServer.Shutdown)
    os.Exit(result)
}

func TestMessageAdd(t *testing.T) {
  newMessage := Message{"Bob", "New Message"}
  testServer.msgListener <- newMessage

  time.Sleep(500 * time.Millisecond)
  if(len(testServer.Messages) != 1 && *testServer.Messages[0] == newMessage) {
    t.Error("Message addition test failure")
  }
  fmt.Println("Passed TestMessageAdd")
  testServer.ClearData()
}

func TestInvalidMessageAdd(t *testing.T) {
  newMessage := Message{"Bob", ""}
  testServer.msgListener <- newMessage

  time.Sleep(500 * time.Millisecond)
  if(len(testServer.Messages) != 0) {
    t.Error("Invalid message addition test failure")
  }
  fmt.Println("Passed TestInvalidMessageAdd")
  testServer.ClearData()
}

func TestMessageAddConcurrency(t *testing.T) {
  for i := 1; i <= 100; i++ {
    go func(id int) {
      testServer.msgListener <- Message{fmt.Sprintf("Bob %d", id), "New Message"}
      time.Sleep(200 * time.Millisecond)
      testServer.msgListener <- Message{fmt.Sprintf("Bob %d", id), "New Message"}
    }(i)
  }

  // Must be able to add 200 messages concurrently in 1 second without failing
  time.Sleep(1000 * time.Millisecond)
  if(len(testServer.Messages) != 200) {
    t.Error("Large concurrent message addition test failure")
  }

  fmt.Println("Passed TestMessageAddConcurrency")
  testServer.ClearData()
}

func TestCommandNewIdea(t *testing.T) {
  newCommand := Message{"Bob", "/newidea ideaName"}

  testServer.msgListener <- newCommand

  time.Sleep(500 * time.Millisecond)
  if(len(testServer.Ideas) != 1 || testServer.CurrentIdea != 0 || testServer.Ideas[0].What != "ideaName") {
    t.Error("/newidea command test failure")
  }
  
  fmt.Println("Passed TestCommandNewIdea")
  testServer.ClearData()
}


func TestCommandIdea(t *testing.T) {
  newIdeaCommand1 := Message{"Bob", "/newidea ideaName1"}
  newIdeaCommand2 := Message{"Bob", "/newidea ideaName2"}
  ideaCommand := Message{"Bob", "/idea 0"}

  testServer.msgListener <- newIdeaCommand1
  testServer.msgListener <- newIdeaCommand2
  testServer.msgListener <- ideaCommand

  time.Sleep(500 * time.Millisecond)
  if(len(testServer.Ideas) != 2 || testServer.CurrentIdea != 0 || testServer.Ideas[testServer.CurrentIdea].What != "ideaName1") {
    t.Error("/idea command test failure")
  }
  
  fmt.Println("Passed TestCommandIdea")
  testServer.ClearData()
}


func TestCommandWhy(t *testing.T) {
  newIdeaCommand := Message{"Bob", "/newidea ideaName"}
  whyCommand := Message{"Bob", "/why pro"}

  testServer.msgListener <- newIdeaCommand
  testServer.msgListener <- whyCommand

  time.Sleep(500 * time.Millisecond)
  // idea list and current idea number test should be assumed as correct, as the /newidea test should confirm
  if(len(testServer.Ideas[0].Why) != 1 || testServer.Ideas[0].Why[0] != "pro") {
    t.Error("/why command test failure")
  }
  
  fmt.Println("Passed TestCommandWhy")
  testServer.ClearData()
}

func TestCommandWhyConcurrency(t *testing.T) {
  newIdeaCommand := Message{"Bob", "/newidea ideaName"}
  testServer.msgListener <- newIdeaCommand

  for i := 1; i <= 100; i++ {
    go func(id int) {
      testServer.msgListener <- Message{fmt.Sprintf("Bob %d", id), "/why pro"}
      time.Sleep(200 * time.Millisecond)
      testServer.msgListener <- Message{fmt.Sprintf("Bob %d", id), "/why pro"}
    }(i)
  }

  time.Sleep(1000 * time.Millisecond)
  if(len(testServer.Ideas[0].Why) != 200) {
    t.Error("Large concurrent /why command test failure")
    
    fmt.Println("Passed TestCommandWhyConcurrency")}
  testServer.ClearData()
}

func TestCommandWhyNot(t *testing.T) {
  newIdeaCommand := Message{"Bob", "/newidea ideaName"}
  whyNotCommand := Message{"Bob", "/whynot con"}

  testServer.msgListener <- newIdeaCommand
  testServer.msgListener <- whyNotCommand

  time.Sleep(500 * time.Millisecond)
  // idea list and current idea number test should be assumed as correct, as the /newidea test should confirm
  if(len(testServer.Ideas[0].WhyNot) != 1 || testServer.Ideas[0].WhyNot[0] != "con") {
    t.Error("/whynot command test failure")
  }
  
  fmt.Println("Passed TestCommandWhyNot")
  testServer.ClearData()
}

func TestCommandWhyNotConcurrency(t *testing.T) {
  newIdeaCommand := Message{"Bob", "/newidea ideaName"}
  testServer.msgListener <- newIdeaCommand

  for i := 1; i <= 100; i++ {
    go func(id int) {
      testServer.msgListener <- Message{fmt.Sprintf("Bob %d", id), "/whynot con"}
      time.Sleep(200 * time.Millisecond)
      testServer.msgListener <- Message{fmt.Sprintf("Bob %d", id), "/whynot con"}
    }(i)
  }

  time.Sleep(1000 * time.Millisecond)
  if(len(testServer.Ideas[0].WhyNot) != 200) {
    t.Error("Large concurrent /whynot command test failure")
    
    fmt.Println("Passed TestCommandWhyNotConcurrency")}
  testServer.ClearData()
}

func TestCommandVote(t *testing.T) {
  newIdeaCommand := Message{"Bob", "/newidea ideaName"}
  voteCommand := Message{"Bob", "/vote yes"}

  testServer.msgListener <- newIdeaCommand
  testServer.msgListener <- voteCommand

  time.Sleep(500 * time.Millisecond)
  // idea list and current idea number test should be assumed as correct, as the /newidea test should confirm
  if(!testServer.Ideas[0].Votes["Bob"]) {
    t.Error("/vote command test failure")
  }
  
  fmt.Println("Passed TestCommandVote")
  testServer.ClearData()
}

func TestCommandRepeatedVote(t *testing.T) {
  newIdeaCommand := Message{"Bob", "/newidea ideaName"}
  voteCommand := Message{"Bob", "/vote yes"}

  testServer.msgListener <- newIdeaCommand
  testServer.msgListener <- voteCommand
  testServer.msgListener <- voteCommand

  time.Sleep(500 * time.Millisecond)
  // idea list and current idea number test should be assumed as correct, as the /newidea test should confirm
  if(!testServer.Ideas[0].Votes["Bob"]) {
    t.Error("Repeated /vote command test failure")
  }
  
  fmt.Println("Passed TestCommandRepeatedVote")
  testServer.ClearData()
}

func TestCommandInvertedVote(t *testing.T) {
  newIdeaCommand := Message{"Bob", "/newidea ideaName"}
  voteYesCommand := Message{"Bob", "/vote yes"}
  voteNoCommand := Message{"Bob", "/vote no"}

  testServer.msgListener <- newIdeaCommand
  testServer.msgListener <- voteYesCommand
  testServer.msgListener <- voteNoCommand

  time.Sleep(500 * time.Millisecond)
  // idea list and current idea number test should be assumed as correct, as the /newidea test should confirm
  if(testServer.Ideas[0].Votes["Bob"] != false) {
    t.Error("Inverted /vote command test failure")
  }
  
  fmt.Println("Passed TestCommandInvertedVote")
  testServer.ClearData()
}

func TestCommandInvalidVote(t *testing.T) {
  newIdeaCommand := Message{"Bob", "/newidea ideaName"}
  voteCommand := Message{"Bob", "/vote meh"}

  testServer.msgListener <- newIdeaCommand
  testServer.msgListener <- voteCommand

  time.Sleep(500 * time.Millisecond)
  // idea list and current idea number test should be assumed as correct, as the /newidea test should confirm
  if(len(testServer.Ideas[0].Votes) != 0) {
    t.Error("Invalid /vote command test failure")
  }
  
  fmt.Println("Passed TestCommandInvalidVote")
  testServer.ClearData()
}

func TestCommandVoteConcurrency(t *testing.T) {
  newIdeaCommand := Message{"Bob", "/newidea ideaName"}
  testServer.msgListener <- newIdeaCommand

  for i := 1; i <= 100; i++ {
    go func(id int) {
      testServer.msgListener <- Message{fmt.Sprintf("Bob %d", id), "/vote yes"}
      time.Sleep(200 * time.Millisecond)
      testServer.msgListener <- Message{fmt.Sprintf("Bob %d", id), "/vote no"}
    }(i)
  }

  time.Sleep(1000 * time.Millisecond)

  hasTrue := func(votes map[string]bool) bool{
    for _, v := range votes {
      if(v){
        return true
      }
    }
    return false
  }

  if(len(testServer.Ideas[0].Votes) != 100 || hasTrue(testServer.Ideas[0].Votes)) {
    t.Error("Large concurrent /why command test failure")
    
    fmt.Println("Passed TestCommandVoteConcurrency")}
  testServer.ClearData()
}

func TestInvalidCommand(t *testing.T) {
  newCommand := Message{"Bob", "/build time machine"}

  testServer.msgListener <- newCommand

  time.Sleep(500 * time.Millisecond)

  if(len(testServer.Ideas) != 0 || testServer.CurrentIdea != -1) {
    t.Error("Invalid command test failure")
  }
  
  fmt.Println("Passed TestInvalidCommand")
  testServer.ClearData()
}

func TestEmptyArgumentCommand(t *testing.T) {
  newCommand := Message{"Bob", "/newidea"}

  testServer.msgListener <- newCommand

  time.Sleep(500 * time.Millisecond)

  if(len(testServer.Ideas) != 0 || testServer.CurrentIdea != -1) {
    t.Error("Invalid command test failure")
  }
  
  fmt.Println("Passed TestEmptyArgumentCommand")
  testServer.ClearData()
}

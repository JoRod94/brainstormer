# Brainstormer

A chat with commands made in **golang** using WebSockets. Client made in **React**.

Use text commands to brainstorm ideas, adding pros and cons, and voting for and against ideas.

To run, use `go run main.go`.

If you edit the client, remember to create a new build with `npm run build` in `./client`. The go server serves build files.

Tests available in `src/chat/chat_test.go`. Use `go test` to run them.

PS: It's ugly.

## Commands

| Command                | Action                                      |
|------------------------|---------------------------------------------|
| `/newidea <idea_name>` | Create new idea and present it in the board |
| `/idea <idea_number>`  | Switch to another idea                      |
| `/why <reason>`        | Add a positive for the current idea         |
| `/whynot <reason>`     | Add a negative to the current idea          |
| `/vote <"yes"/"no">`   | Vote the current idea up or down            |

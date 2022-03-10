# go-repl

Golang REPL library.

Your REPL's using this library will enjoy the following features:
* session history with *reverse-search*
  * Ctrl-r: to start *reverse-search*
  * most edit commands, except the most basic ones, exit the *reverse-search* mode
  * use Up/Down to cycle through a filtered list of history entries (while in *reverse-search* mode)
* the input buffer is redrawn when a resize is detected
* status bar at bottom with current working dir and other info
* truncation of very long inputs (status bar displays info about cursor position)
* common edit and movement commands:
   * Right/Left: move cursor one character at a time
   * Up/Down: cycle through history
   * Backspace/Delete: works as expected
   * Shift-Enter: enter newline into buffer without invoking Eval
   * Ctrl-a or Home: move to start of buffer
   * Ctrl-e or End: move to end of buffer
   * Ctrl-w: delete preceding word
   * Ctrl-q: delete following phrase
   * Ctrl-Right/Left: move cursor one word at a time
   * Ctrl-c or Esc: ignore current input
   * Ctrl-d: quit
   * Ctrl-u: clear buffer to start
   * Ctrl-k: clear buffer to end
   * Ctrl-l: reset prompt at top and redraw buffer
   * Ctrl-y: insert previous deletion (from Ctrl-k, Ctrl-u, Ctrl-q or Ctrl-w)

No dependency on *ncurses*. 

Performance hasn't yet been optimized and I haven't yet tested all corner cases exhaustively.

Might not work in Windows command prompt (keystroke codes could differ, ANSI escape sequences might not be supported, the method that sets terminal to raw mode might not work crossplatform).

# Usage

Fetch this library with the following command:
```shell
$ go get -u github.com/openengineer/go-repl
```

In order to create your own REPL you have to define a type that implements the `Handler` interface:
```golang
type Handler interface {
  Prompt() string
  Tab(buffer string) string
  Eval(line string) string
}
```

Here is a complete example (can also be found in `./examples/simple/main.go`):

```golang
package main

import (
  "fmt"
  "log"
  "strconv"
  "strings"

  repl "github.com/openengineer/go-repl"
)

var helpMessage = `help              display this message
add <int> <int>   add two numbers
quit              quit this program`

// implements repl.Handler interface
type MyHandler struct {
  r *repl.Repl
}

func main() {
  fmt.Println("Welcome, type \"help\" for more info")

  h := &MyHandler{}
  h.r = repl.NewRepl(h)

  // start the terminal loop
  if err := h.r.Run(); err != nil {
    log.Fatal(err)
  }
}

func (h *MyHandler) Prompt() string {
  return "> " 
}

func (h *MyHandler) Tab(buffer string) string {
  return "" // do nothing
}

func (h *MyHandler) Eval(line string) string {
  fields := strings.Fields(line)

  if len(fields) == 0 {
    return ""
  } else {
    cmd, args := fields[0], fields[1:]

    switch cmd {
    case "help":
      return helpMessage
    case "add":
      if len(args) != 2 {
        return "\"add\" expects 2 args"
      } else {
        return add(args[0], args[1])
      }
    case "quit":
      h.r.Quit()
      return ""
    default:
      return fmt.Sprintf("unrecognized command \"%s\"", cmd)
    }
  }
}

func add(a_ string, b_ string) string {
  a, err := strconv.Atoi(a_)
  if err != nil {
    return "first arg is not an integer"
  }

  b, err := strconv.Atoi(b_)
  if err != nil {
    return "second arg is not an integer"
  }

  return strconv.Itoa(a + b)
}
```

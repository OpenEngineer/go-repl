package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	repl "github.com/openengineer/go-repl"
)

type ShellWrapper struct {
	r *repl.Repl
}

func (h *ShellWrapper) Prompt() string {
	return "> "
}

func (h *ShellWrapper) Tab(buffer string) string {
	// a tab is simply 2 spaces here
	return "  "
}

func (h *ShellWrapper) Eval(line string) string {
	// upon eval the Stdin should be unblocked
	if strings.TrimSpace(line) != "" {
		endCmd := make(chan bool)

		if line == "quit" {
			h.r.Quit()
			return ""
		}

		go func() {
			h.r.UnmakeRaw()
			defer h.r.MakeRaw()

			fields := strings.Fields(line)

			cmdName, args := fields[0], fields[1:]

			cmd := exec.Command(cmdName, args...)

			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr

			err := cmd.Start()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())

				endCmd <- false
			} else {
				cmd.Wait()
				endCmd <- true
			}
		}()

		ok := <-endCmd

		if ok {
			return ""
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func main() {
	h := &ShellWrapper{}
	h.r = repl.NewRepl(h)

	if err := h.r.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

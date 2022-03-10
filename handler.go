package repl

// Implement this interface in order to use `Repl` with your custom logic.
type Handler interface {
	Prompt() string
	Eval(line string) string
	Tab(prec string) string
}

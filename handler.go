package repl

// This interface must be implemented in order to be able to use Repl with your own logic
type Handler interface {
	Prompt() string
	Eval(line string) string
	Tab(prec string) string
}

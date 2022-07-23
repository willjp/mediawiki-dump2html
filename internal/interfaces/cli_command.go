package interfaces

// command pattern, for encapsulating commands set on commandline
type CliCommand interface {
	Call() error
}

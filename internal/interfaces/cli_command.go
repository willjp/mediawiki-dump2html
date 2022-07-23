package interfaces

// All available commands have their own CliCommand.
type CliCommand interface {
	Call() error
}

package interfaces

// abstracts a os.File
type OsFile interface {
	Name() string
	WriteString(s string) (n int, err error)
	Close() (err error)
}

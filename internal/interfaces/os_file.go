package interfaces

// abstracts a os.File
type OsFile interface {
	Name() string
	WriteString(v string) (n int, err error)
	Close() (err error)
}

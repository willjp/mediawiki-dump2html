package cli

import "fmt"

type ShowHelp struct{}

func (this *ShowHelp) Call() error {
	fmt.Println("help")
	return nil
}

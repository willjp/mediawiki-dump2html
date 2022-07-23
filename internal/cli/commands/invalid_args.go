package cli

import "fmt"

type InvalidArgs struct{}

func (this *InvalidArgs) Call() error {
	fmt.Println("invalid args")
	return nil
}

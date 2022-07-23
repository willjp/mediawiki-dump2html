package cli

import "fmt"

type Build struct{}

func (this *Build) Call() error {
	fmt.Println("build")
	return nil
}

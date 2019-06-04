package main

import (
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	cli := &CLI{
		outStream:   os.Stdout,
		errStream:   os.Stderr,
		configFlags: genericclioptions.NewConfigFlags(true),
	}
	os.Exit(cli.Run(os.Args))
}

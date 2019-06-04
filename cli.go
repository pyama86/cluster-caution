package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd/api"
)

const Name string = "cluster-caution"
const Version string = "0.1.0"

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
	rawConfig            api.Config
	configFlags          *genericclioptions.ConfigFlags
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	if err := cli.run(args); err != nil {
		fmt.Errorf("%s\n", err)
		return ExitCodeError
	}
	return ExitCodeOK
}
func (cli *CLI) run(args []string) error {
	var (
		init bool

		version bool
	)
	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&init, "init", false, "Init Repository")
	flags.BoolVar(&init, "i", false, "Init Repository(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")
	cli.configFlags.AddFlags(flags)

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	rawConfig, err := cli.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}
	currentContext, exists := rawConfig.Contexts[rawConfig.CurrentContext]
	if !exists {
		return err
	}
	top, err := getRepositoryTop()
	if err != nil {
		return err
	}

	fp := filepath.Join(strings.TrimRight(top, "\n"), ".kube-cluster-coution")
	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return nil
	}

	if init {
		if err := writeFile(currentContext, fp); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(ac *api.Context, path string) error {
	bdata, err := json.Marshal(ac)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bdata, os.ModePerm)
}
func readFile(path string) (*api.Context, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var ac api.Context
	err = json.Unmarshal(raw, &ac)
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

func getRepositoryTop() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

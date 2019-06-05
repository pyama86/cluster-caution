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
	"syscall"

	"github.com/Songmu/prompter"
	"github.com/k0kubun/pp"
	flag "github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd/api"
)

const Name string = "cluster-caution"
const Version string = "0.1.1"

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
		fmt.Println(err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func (cli *CLI) run(args []string) error {
	var (
		add bool
		del bool

		version bool
	)
	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&add, "add-context", false, "Add Context Config to Repository")
	flags.BoolVar(&del, "delete-context", false, "Delete Context Config to Repository")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")
	cli.configFlags.AddFlags(flags)

	// Parse commandline flag
	flagError := flags.Parse(args[1:])

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

	acs, err := readFile(fp)
	if err != nil {
		return err
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return nil
	}

	resultMessage := fmt.Sprintf("add cluster %s(%s)\n", currentContext.Cluster, currentContext.Namespace)
	if add || del {
		if flagError != nil {
			return flagError
		}
		if add {
			a, exist := uniqueAppendContext(acs, currentContext)
			if exist {
				resultMessage = fmt.Sprintf("cluster %s(%s) is exists\n", currentContext.Cluster, currentContext.Namespace)
			}
			acs = a
		} else if del {
			a, exist := deleteContext(acs, currentContext)
			if exist {
				resultMessage = fmt.Sprintf("delete cluster %s(%s)\n", currentContext.Cluster, currentContext.Namespace)
			} else {
				resultMessage = fmt.Sprintf("cluster %s(%s) is not exists\n", currentContext.Cluster, currentContext.Namespace)
			}
			acs = a
		}
		if err := writeFile(acs, fp); err != nil {
			return err
		}
		fmt.Println(resultMessage)
	} else {
		runKubectl(acs, currentContext)
	}

	return nil
}

func uniqueAppendContext(ac []*api.Context, newAc *api.Context) ([]*api.Context, bool) {
	exist := false
	for _, a := range ac {
		if a.Cluster == newAc.Cluster && a.Namespace == newAc.Namespace {
			exist = true
		}
	}
	if !exist {
		ac = append(ac, newAc)
	}
	return ac, exist
}
func deleteContext(ac []*api.Context, delAc *api.Context) ([]*api.Context, bool) {
	ret := []*api.Context{}
	exist := false
	for _, a := range ac {
		if a.Cluster == delAc.Cluster && a.Namespace == delAc.Namespace {
			exist = true
			continue
		}
		ret = append(ret, delAc)
	}
	return ret, exist
}

func writeFile(ac []*api.Context, path string) error {
	bdata, err := json.Marshal(ac)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bdata, os.ModePerm)
}

func readFile(path string) ([]*api.Context, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		if _, notfound := err.(*os.PathError); notfound {
			return []*api.Context{}, nil
		}
		pp.Println(err)
		return nil, err
	}
	var ac []*api.Context
	err = json.Unmarshal(raw, &ac)
	if err != nil {
		return nil, fmt.Errorf("can't parse .kube-cluster-coution %s", err.Error())
	}
	return ac, nil
}

func runKubectl(ac []*api.Context, cur *api.Context) {
	exist := false
	for _, a := range ac {
		if a.Cluster == cur.Cluster && a.Namespace == cur.Namespace {
			exist = true
		}
	}

	if !exist && len(ac) > 0 {
		if !prompter.YesNo("Repository configuration is different from cluster or namespace.\nDo you want to continue?(Y/n)", true) {
			return
		}
	}
	args := os.Args[1:]
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		status := 1
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				status = s.ExitStatus()
			}
		}
		os.Exit(status)
	}
}

func getRepositoryTop() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

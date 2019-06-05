# cluster-caution
## Description

This subcommand prevents cluster and namespace slip.

## Usage

```
# make .kube-cluster-caution file
$ kubectl cluster-caution --add
$ cat .kube-cluster-caution
[{"LocationOfOrigin":"/Users/example/.kube/config","cluster":"cluster.example.com","user":"example","namespace":"example"}]

# switch other cluster
$ kubectl config use-context other.example.com --namespace=other
$ alias kc='kubectl cluster-caution'
$ kc delete pod --all
Repository configuration is different from cluster or namespace.
Do you want to continue?(Y/n) (yes/no) [yes]:
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/pyama86/cluster-caution
```

## Contribution

1. Fork ([https://github.com/pyama86/cluster-caution/fork](https://github.com/pyama86/cluster-caution/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pyama86](https://github.com/pyama86)

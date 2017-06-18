# fint-model



## Description
Generates Java and C# models from EA XMI export. This utility is mainly for internal FINT use, but if you 
find it usefull, please use it!

## Usage

```
fint-model                                                     ✓  3268  19:46:17
NAME:
   fint-model - Generates Java and C# models from EA XMI export. This utility is mainly for internal FINT use, but if you find it usefull, please use it!

USAGE:
   fint-model [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   FINTProsjektet

COMMANDS:
     printClasses    list classes
     generate
     listPackages    list Java packages
     listNamespaces  list .Net namespaces
     listTags        list tags
     listBranches    list branches
     help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --tag value    the tag (version) of the model to generate (default: "latest")
   --help, -h     show help
   --version, -v  print the version
```


## Install

### Binaries

Precompiled binaries can be downloaded [here](https://github.com/FINTprosjektet/fint-model/releases/latest)

### Go

To install, use `go get`:

```bash
$ go get -d github.com/FINTProsjektet/fint-model
$ go install github.com/FINTProsjektet/fint-model
```



## Author

[FINTProsjektet](https://fintprosjektet.github.io)

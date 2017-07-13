# fint-model



## Description
Generates `Java` and `C#` models from EA XMI export. This utility is mainly for internal FINT use, but if you 
find it usefull, please use it!

## Usage

```
$ fint-model
NAME:
   fint-model - Generates Java and C# models from EA XMI export. This utility is mainly for internal FINT use,
   but if you find it usefull, please use it!

USAGE:
   fint-model [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   FINTProsjektet

COMMANDS:
     printClasses    list classes
     generate        generates JAVA/CS models
     listPackages    list Java packages
     listNamespaces  list CS namespaces
     listTags        list tags
     listBranches    list branches
     help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --tag value, -t value  the tag (version) of the model to generate (default: "latest")
   --force, -f            force downloading XMI for GitHub.
   --help, -h             show help
   --version, -v          print the version
```

The downloaded XMI file is put in the `$HOME/.fint-model/.cache`. If you don't use the 
`force` flag and the file exists in the temp directory it uses this one. 

## Install

### Binaries

Precompiled binaries can be downloaded [here](https://github.com/FINTprosjektet/fint-model/releases/latest)

### Go

To install, use `go get`:

```bash
go get -d github.com/FINTprosjektet/fint-model
go install github.com/FINTprosjektet/fint-model
```

## Author

[FINTProsjektet](https://fintprosjektet.github.io)

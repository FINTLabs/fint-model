# fint-model

## Description
Generates `Java` and `C#` models from EA XMI export. This utility is mainly for internal FINT use, but if you 
find it useful, please use it!

## Usage

```
$ fint-model
NAME:
   fint-model - Generates Java and C# models from EA XMI export. This utility is mainly for internal FINT use, but if you find it usefull, please use it!

USAGE:
   fint-model [global options] command [command options] [arguments...]

VERSION:
   0.0.0

AUTHOR:
   FINTLabs

COMMANDS:
     printClasses    list classes
     generate        generates JAVA/CS models
     listPackages    list Java packages
     listNamespaces  list CS namespaces
     listTags        list tags
     listBranches    list branches
     help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --owner value          Git repository containing model (default: "FINTLabs") [$GITHUB_OWNER]
   --repo value           Git repository containing model (default: "fint-informasjonsmodell") [$GITHUB_PROJECT]
   --filename value       File name containing information model (default: "FINT-informasjonsmodell.xml") [$MODEL_FILENAME]
   --tag value, -t value  the tag (version) of the model to generate (default: "latest")
   --force, -f            force downloading XMI for GitHub.
   --help, -h             show help
   --version, -v          print the version
```

The downloaded XMI file is put in the `$HOME/.fint-model/.cache`. If you don't use the 
`force` flag and the file exists in the cache directory it uses this one. 

## Install

### Binaries

Precompiled binaries are available as [Docker images](https://cloud.docker.com/u/fint/repository/docker/fint/fint-model)

Mount the directory where you want the generated source code to be written as `/src`.

Linux / MacOS:
```bash
docker run -v $(pwd):/src fint/fint-model:latest <ARGS>
```

Windows PowerShell:
```ps1
docker run -v ${pwd}:/src fint/fint-model:latest <ARGS>
```

### Source

To install, use `go get`:

```bash
go get -d github.com/FINTLabs/fint-model
go install github.com/FINTLabs/fint-model
```

## Author

[FINTLabs](https://FINTLabs.github.io)

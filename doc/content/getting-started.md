---
linkTitle: Getting Started
title: Getting Started
description: Create your first xc file
menu: main
weight: -7
---

## Installation

{{% details "Binary" %}}
Download the binary from [GitHub Releases](https://github.com/joerdav/xc/releases) and add to your `$PATH`
{{% /details %}}
{{% details "Go Install" %}}
```sh
go install github.com/joerdav/xc/cmd/xc@latest
```
{{% /details %}}
{{% details "Homebrew" %}}
```sh
brew tap joerdav/xc
brew install xc
```
{{% /details %}}
{{% details "Snap" %}}
```sh
TODO
```
{{% /details %}}
{{% details "Nix" %}}
```sh
TODO
```
{{% /details %}}
{{% details "Scoop" %}}
```sh
TODO
```
{{% /details %}}
{{% details "AUR" %}}
```sh
TODO
```
{{% /details %}}

## Verify Installation

Run `xc -version` to verify the installation.
If installed via `go install` the version will be `devel`.

## Create some tasks.

Create a file named README.md:

````
## Tasks
### hello
Prints hello
```sh
echo hello
```
### world
Prints world
Requires: hello
```sh
echo world
```
````

## List tasks.

Run `xc` to list the tasks.

```
$ xc
    hello  Prints hello
    world  Prints world
           Requires:  hello
```

## Run a task.

Run `xc world` to run the `world` task.
```
$ xc world
echo hello
hello
echo world
world
```

## RTFM!

Now that you have a simple task working, [RTFM](/task-syntax/) to learn what else `xc` is capable of.

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
{{% details "tea" %}}
With `tea` just type `xc`.
```sh
$ xc --version
# ^^ or `tea xc` if there’s no magic
```
{{% /details %}}
{{% details "Snap" %}}
```sh
TODO
```
{{% /details %}}
{{% details "Nix" %}}
There is a nix flake that can be used:
```sh
nix develop github:joerdav/xc
```
Or to create your own `xc.nix`, replace `<version>` and add the correct `sha256` and `vendorSha256` for the version:
```nix
{ config, pkgs, fetchFromGitHub, ... }:

let
  xc = pkgs.buildGoModule rec {
    pname = "xc";
    version = "<version>";
    subPackages = ["cmd/xc"];
    src = pkgs.fetchFromGitHub {
      owner = "joerdav";
      repo = "xc";
      rev = version;
      sha256 = "";
    };
    vendorSha256 = "";
  };
in
{
  environment.systemPackages = [ xc ];
}
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

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
{{% details "MacPorts" %}}
```sh
sudo port install xc
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
scoop install xc
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

## Install completion

Run `xc -complete` to install auto completion.

Run `xc -uncomplete` to uninstall auto completion.

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

Or if you prefer, create a file called README.org:

```org
** Tasks
*** hello
Prints hello
#+begin_src sh
echo hello
#+end_src
*** world
Prints world
Requires: hello
#+begin_src sh
echo world
#+end_src
```

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

## Let people know you use xc.

Add the following badge to your `README` so that people know it's xc compatible.

```
[![xc compatible](https://xcfile.dev/badge.svg)](https://xcfile.dev)
```

<a href="https://xcfile.dev" alt="xc compatible">
        <img src="/badge.svg" /></a>

## RTFM!

Now that you have a simple task working, [RTFM](/task-syntax/) to learn what else `xc` is capable of.

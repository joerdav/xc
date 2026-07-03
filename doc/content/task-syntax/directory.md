---
title: "Directory"
description:
linkTitle: "Directory"
menu: { main: { parent: 'task-syntax', weight: 10 } }
---

## Task Directory

The `directory` attribute can be used to define a file path in which the task should run.

## Syntax

Do this by adding a `dir:` or `directory:` between the task name and the script.

````markdown
## Tasks
### Build
directory: ./src
```
sh build.sh
```
````

## Using the caller's working directory

By default, tasks run in the directory where the README file is located, even if `xc` is invoked from a subdirectory. Setting `directory` to `$PWD` will instead run the task in the directory where `xc` was invoked.

````markdown
## Tasks
### List
directory: $PWD
```
ls
```
````

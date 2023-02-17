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

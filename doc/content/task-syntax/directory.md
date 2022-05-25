---
title: "Directory"
description:
linkTitle: "Directory"
menu: { main: { parent: 'task-syntax', weight: 10 } }
---

## How to use a custom working directory

You can define the working directory of the task for the command to run in.

Do this by adding a `dir:` or `directory:` between the H3 of the task name and the command.

````markdown
## Tasks
### Build
directory: ./src
```
sh build.sh
```
````
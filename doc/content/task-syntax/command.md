---
title: "Task Command"
description:
linkTitle: "Task Command"
menu: { main: { parent: 'task-syntax', weight: 10 } }
---

## How to define a command to run a task

The command of a task is what will be executed through the command line when the task is run.

To define a command you must add a code block under a particular task.

````markdown
## Tasks
### Task1
```
echo "Hello 世界!"
echo "Hello العالمية!"
echo "Hello ертөнц!"
```
````

Everything inside the code block will be run on the command line as if it is a shell script.
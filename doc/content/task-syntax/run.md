---
title: "Run"
description:
linkTitle: "Run"
menu: { main: { parent: 'task-syntax', weight: 10 } }
---

## Run attribute

By default, a task can run as many times as it appears in the requires tree.
Consider a scenario where a task is required to run only one time per `xc` invocation,
regardless of how many times it is required.

The solution would be to set the `run` attribute to `once` (defaults to `always`).

````markdown
### setup

run: once

```
echo "TASK 3"
```
````

This will result in the task only running the first time it is invoked.

The default is `always`, which can be omitted or specified.

````markdown
### setup

```
echo "TASK 3"
```
````

is the same as

````markdown
### setup

run: always

```
echo "TASK 3"
```
````

---
title: "Name"
description:
linkTitle: "Name"
menu: { main: { parent: 'task-syntax', weight: 2 } }
---

## Task Name

The name of a task should explain concisely it's purpose - it serves as documentation and the identifier for running a task.

## Syntax

Under the `## Tasks` heading, you may define as many tasks as you'd like. Each task is defined with a name by adding a H3 heading like `### TaskName`

```markdown
## Tasks

### Task1

### Task2

### Task3
```

## Constraints

You cannot use spaces in the task name.

But you may use `-` or `_`

```markdown
## Tasks

### Task_1

### Task-2
```

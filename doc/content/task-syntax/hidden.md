---
title: "Hidden tasks"
description:
linkTitle: "Hidden tasks"
menu: { main: { parent: 'task-syntax', weight: 10 } }
---

## How to hide a task in the list

When `xc` is run without any arguments, it lists all the tasks available in the project.
However, sometimes you may want to hide a task from this list, while still allowing it to be run.

To hide a task, you can set the `hidden` attribute to `true`.

```markdown
### usefull-task-but-hidden

hidden: true
```


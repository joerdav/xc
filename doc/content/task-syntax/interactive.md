---
title: "Interactive"
description:
linkTitle: "Interactive"
menu: { main: { parent: "task-syntax", weight: 12 } }
---

## Interactive attribute

By default, the logs of a task are prefixed with the task name. This does not work well for interactive tasks which usually require complete control over the terminal.
If you want to run a task interactively, you can set the `interactive` attribute to `true`.

```markdown
### configure

interactive: true
```

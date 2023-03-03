---
title: "Task List"
description:
linkTitle: "Task List"
menu: { main: { parent: 'task-syntax', weight: 1 } }
---

## Task List

The task list is the section of markdown that contains xc tasks.

## Syntax

An xc compatible markdown file needs to define a heading of any level called `Tasks`.

The tasks within a `Tasks` section will need to be one heading level lower than `Tasks`.

The xc heading can be overridden with the flags `-H` or `-heading`.

```markdown
## Tasks

### Task1

### Task2

## Next H2 Heading
```

Between `## Tasks` and `## Next H2 Heading` you can define your tasks.

or

```markdown
# Tasks

## Task1

## Task2

# Next H1 Heading
```

Please note that the word `tasks` is not case sensitive.

## Constraints

You cannot define two `Tasks` sections. If you do, the one that appears first in the markdown file will be used

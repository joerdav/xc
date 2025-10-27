---
title: "Task List"
description:
linkTitle: "Task List"
menu: { main: { parent: 'task-syntax', weight: 1 } }
---

## Task List

The task list is the section of documentation that contains xc tasks.

## Syntax

An xc compatible text file needs to define a heading of any level called `Tasks`.

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

Note `xc` cannot find tasks in two or more `Tasks` sections. If you have more than one `Tasks` section, the one that appears first will be used.

## The `<!-- xc-heading -->` Comment

In Markdown files, you can put the special comment `<!-- xc-heading -->` on the line before a heading, and `xc` will recognize that as the task heading.

```markdown
<!-- xc-heading -->
## Getting started

### Task1

### Task2

## Next H2 Heading
```

## The `:xc_heading:` Tag

In org-mode files, you can put the special tag `:xc_heading:` on a heading, and `xc` will recognize that as the task heading.

```org
** Getting started                                         :xc_heading:

*** Task1

*** Task2

** Next H2 Heading
```

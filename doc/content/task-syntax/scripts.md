---
title: "Scripts"
description:
linkTitle: "Scripts"
menu: { main: { parent: 'task-syntax', weight: 3 } }
---

## Task Scripts

A script is what defines a tasks behaviour.

## Syntax

To define a script you must add a code block under the name of a task.

In Markdown:

````markdown
## Tasks
### Task1
```
echo "Hello 世界!"
echo "Hello العالمية!"
echo "Hello ертөнц!"
```
````

Or in org-mode:

```org
** Tasks
*** Task1
#+begin_src bash
echo "Hello 世界!"
echo "Hello العالمية!"
echo "Hello ертөнц!"
#+end_src
```

## Indented code blocks

A script can also be written as an indented code block, using four spaces or a tab
instead of a ``` fence. This is handy when the surrounding document already uses
fenced blocks for something else.

````markdown
## Tasks
### Task1

    echo "Hello 世界!"
    echo "Hello العالمية!"
````

## Shebangs

To define an alternative interpreter such as python, then include a shebang, similar to the unix style.

xc will parse this and create a temporary file that will be executed using the specified interpreter.

````markdown
## Tasks
### python-task
```
#!/usr/bin/env python
print("foo")
```
````

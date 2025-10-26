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

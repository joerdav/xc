---
title: "Org-mode Features"
description:
linkTitle: "Org-mode Features"
menu: { main: {  weight: 10 } }
---

## How `xc` Finds org-mode Files

When using `xc` with org-mode, it looks in the current directory (then in parents) for a `README.org` and treats sub-headers under the "Tasks" header as tasks.

/Note: in a directory that contains both `README.org` and `README.md`, `xc` will look for tasks in the Markdown README only./
You can change this using the `type` flag, i.e. `xc -type org TASK`

## Comment headings
`xc` will skip over org-mode [comment heading](https://orgmode.org/manual/Comment-Lines.html) beginning with `COMMENT`.

## Header tags
When `xc` looks for tasks, it expects to find them under a header called "Tasks". To override this, you can put an `:xc_heading:` [tag](https://orgmode.org/manual/Tags.html) on any heading and `xc` will find tasks under the first such heading it encounters. Here's what this looks like in practice, supposing you want to use "Getting started" as your task heading:
```org
** Getting started                                              :xc_heading:
*** build
  [...etc]
```

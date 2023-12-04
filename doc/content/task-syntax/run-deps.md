---
title: "Run Dependencies"
description:
linkTitle: "RunDeps"
menu: { main: { parent: "task-syntax", weight: 11 } }
---

## RunDeps attribute

By default, the dependencies of a task are run sequentially, in the order they are listed.  
However we may prefer for all the dependencies of a task to be run in paralled.

The solution would be to set the `runDeps` attribute to `async` (defaults to `sync`).

```markdown
### build-all

requires: build-js, build-css
runDeps: async
```

This will result in both `build-js` and `build-css` being run in parallel.

The default is `sync`, which can be omitted or specified.

```markdown
### build-all

requires: build-js, build-css
runDeps: sync
```

is the same as

```markdown
### build-all

requires: build-js, build-css
runDeps: sync
```

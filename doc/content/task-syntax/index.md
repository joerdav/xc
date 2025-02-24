---
title: "Task Syntax"
weight: 1
description:
linkTitle: "Task Syntax"
---

## The anatomy of an `xc` task

### Structure

- [Task List](/task-syntax/task-list/)
  - [Name](/task-syntax/task-name/)
    - [Scripts](/task-syntax/scripts/)
    - [Requires](/task-syntax/requires/)
    - [Run](/task-syntax/run/)
    - [RunDeps](/task-syntax/run-deps/)
    - [Directory](/task-syntax/directory/)
    - [Environment Variables](/task-syntax/environment-variables/)
    - [Inputs](/task-syntax/inputs/)
    - [Interactive](/task-syntax/interactive/)
    - [Hidden](/task-syntax/hidden/)

### Example

````md
## Tasks

### deploy

Requires: test
Directory: ./deployment
Env: ENVIRONMENT=STAGING

```
sh deploy.sh
```
````

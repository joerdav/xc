---
title: "Task Syntax"
weight: 1
description:
linkTitle: "Task Syntax"
---

## The anatomy of an `xc` task.

### Structure

- Tasks Section
  - Task Name
    - Dependencies
    - Directory
    - Environment Variables
    - Command

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

---
title: "Task Syntax"
weight: 1
description: Task Syntax
linkTitle: "Task Syntax"
---

## The anatomy of an `xc` task.

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

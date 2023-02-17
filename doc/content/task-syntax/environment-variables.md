---
title: "Environment Variables"
description:
linkTitle: "Environment Variables"
menu: { main: { parent: 'task-syntax', weight: 10 } }
---
## Task Environment Variables

The `environment` attribute can be used to define environment variables that should be set before a task runs.

## Syntax

Do this by adding an `env:` or `environment:` between the task name and the script.

````markdown
## Tasks
### Task1
Env: ENVIRONMENT=PRODUCTION
```
echo $ENVIRONMENT
```
````

## Multiple environment variables

You can define multiple environment variables.

````markdown
## Tasks
### Task1
Env: ENVIRONMENT=PRODUCTION
Env: VERSION=1.2
```
echo $ENVIRONMENT
echo $VERSION
```
````

or

````markdown
## Tasks
### Task1
Env: ENVIRONMENT=PRODUCTION, VERSION=1.2
```
echo $ENVIRONMENT
echo $VERSION
```
````

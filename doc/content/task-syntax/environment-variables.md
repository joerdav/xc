---
title: "Environment Variables"
description:
linkTitle: "Environment Variables"
menu: { main: { parent: 'task-syntax', weight: 10 } }
---
## How to use environment variables

You can define environment variables that will be used during the execution of your task.

Do this by adding an `env:` or `environment:` between the H3 of the task name and the command

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
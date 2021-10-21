# xc - eXeCute project tasks from a readme file

## Installation

```
go install github.com/joe-davidson1802/xc/cmd/xc@latest
```

## Usage

To list available tasks run the following.

```
xc 
```

To run a task simply run the following replacing `task` with the task to be run.

```
xc task
```

To run a task from a file not named README.md run with the `-f` of `--file` flag.

```
xc --file OTHERFILE.md task
```

## Syntax

### Task section

To signify the start of the task definition section create a heading name "Tasks".
If a heading of the same level or greater than the "Tasks" heading is found this signifys the end of the Task section.

> ### Tasks
> - Tasks go here
> #### Subheading - Does not end the task section
> - More tasks go here
> ## Another heading - Ends the task section

### Task definition

Once in the task section a task can be defined using the following syntax:

```` md
taskname: task description
!task-dependency1, task-dependency2
```
command
```
````

#### Name

Name is everything before the `:`, and can contain md styling.
The following are all valid task names:

_task-1_
__task-2__
*task-3*

### Description

Name is everything after the `:`, and can contain any content.

### Dependencies

Other tasks can be ran by defining dependencies at the beginning.
The are signified by a `!`, they can be comma delimited or on separate lines.
The following are equivelant:

```
!!task1, task2, task3
```
```
!!task1
!!task2, task3
```

### Command

The tasks command is signified by a md codeblock.

```
command --args
```

### Subheadings

Subheadings can exist within the tasks section, as long as their level is less than the tasks title. (more # is a lower level)

## Example

### Tasks

#### Tests

__test__: test project
```
go test ./...
```

#### Development

__get__: get dependencies of the project
```
go get ./...
```

__deploy-version__: tag current commit with a version
!!test
```
sh ./push-tag.sh
```


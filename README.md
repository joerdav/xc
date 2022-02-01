# xc - eXeCute project tasks from a readme file

## Installation

Go:
```
go install github.com/joe-davidson1802/xc/cmd/xc@latest
```

Brew:
```
brew tap joe-davidson1802/xc
brew install xc
```

Install tab completion for bash:

```
xc -complete
```

When keys `xc<Tab>` are pressed, completion should be triggered:

```
[0][~/src/xc][main]$ xc
combo    echoenv  get      ls       tag      test
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
> ## Another heading - Ends the task section

### Task definition

Once in the task section a task can be defined by a subheading with a lower level:

```` md
### taskname
taskdescription
Requires: task-dependency1, task-dependency2
```
command
```
````

#### Name

The name is denoted by a heading lower than the Tasks heading.

### Description

Anything between the task name and the task command, that is not a "Requires:" section is a description.

### Dependencies

Other tasks can be ran by defining dependencies at the beginning.
They are signified by the `Requires:` prefix, they can be comma delimited or on separate lines.
The following are equivelant:

```
Requires: task1, task2, task3
```
```
Requires: task1
Requires: task2, task3
```

### Directory

Directory by default will be the current working directory. However, if you provide a "Directory:" section then it can be overridden.

### Environment Variables

Environment variables can be set with "Env:".

### Command

The tasks command is signified by a md codeblock.

```
command --args
```

## VIM Usage

I plan to introduce a vim plugin to provide this functionality, but for now you can install:

- <https://github.com/junegunn/fzf>
- <https://github.com/christoomey/vim-run-interactive>

And use the following mapping:

``` sh
:map <leader>xc :call fzf#run({'source':'xc -short', 'options': '--prompt "xc> " --preview "xc -md {}"', 'sink': 'RunInInteractiveShell xc', 'window': {'width': 0.9, 'height': 0.6}})
```

## Example

### Tasks

#### test

Test the project.

```
go test ./...
```

#### get
Get the project dependencies.

```
go get ./...
```

#### tag
Deploys a new tag for the repo.

Also runs tests

Requires: test
```
sh ./push-tag.sh
```

#### ls

Directory: ./parser
```
ls
```

#### echoenv

Env: SOME_VAR=test

```
printenv SOME_VAR
```

#### combo

Runs multiple commands but doesn't have it's own.

Requires: echoenv, ls


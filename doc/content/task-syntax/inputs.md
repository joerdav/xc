---
title: "Inputs"
description:
linkTitle: "Inputs"
menu: { main: { parent: 'task-syntax', weight: 10 } }
---

## Pass an input to a task

Task definitions are run as shell scripts, therefore you can pass arguments to them as such.

`xc [task-name] [inputs...]`

## Named Inputs

The `Inputs` attribute can be used to denote required parameters, which at runtime will be converted to environment variables.

````markdown
## Tasks
### greet

Inputs: FORENAME, SURNAME

```
echo "Hello, $FORENAME $SURNAME."
```
````

Named Inputs can be passed via command arguments:

```sh
$ xc greet Joe Bloggs
+ echo 'Hello, Joe Bloggs.'
Hello, Joe Bloggs.
```

Or via environment variables:

```sh
$ export FORENAME=Joe
$ export SURNAME=Bloggs
$ xc greet
+ echo 'Hello, Joe Bloggs.'
Hello, Joe Bloggs.
```

xc will return an error if `Inputs` are not passed:

```sh
$ xc greet
Task has required inputs:
        xc greet [FORENAME] [SURNAME]
        FORENAME=[FORENAME] SURNAME=[SURNAME] xc greet
exit status 1
```

## Optional Inputs

Combining the `Environment` attribute and the `Inputs` attribute, you can create optional inputs to a task.

````markdown
## Tasks
### greet

Inputs: NAME
Environment: NAME=World

```
echo "Hello, $NAME."
```
````

Can be ran as:

```sh
$ xc greet Joe
+ echo 'Hello, Joe.'
Hello, Joe.
```
or

```sh
$ xc greet
+ echo 'Hello, World.'
Hello, World.
```

## Positional Syntax

As xc tasks are executed as shell scripts you can also use positional syntax of arguments.

(xc will not return an error if a script uses positional syntax and the input is not provided)

````markdown
## Tasks
### greet
```
echo "Hello, $1."
```
````

Can be ran as:

```sh
$ xc greet Joe
+ echo 'Hello, Joe.'
Hello, Joe.
```

## Variadic Syntax

As xc tasks are executed as shell scripts you can also use positional syntax of arguments.

(xc will not return an error if a script uses variadic syntax and no inputs are provided)

````markdown
## Tasks
### greet-everyone
```
echo "Hello, $@."
```
````

Can be ran as:

```sh
$ xc greet Joe Bob Steve
+ echo 'Hello, Joe Bob Steve.'
Hello, Joe Bob Steve.
```

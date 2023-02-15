---
title: "Inputs"
description:
linkTitle: "Inputs"
menu: { main: { parent: 'task-syntax', weight: 10 } }
---

## Pass an input to a task

Task definitions are run as shell scripts, therefore you can pass arguments to them as such.

`xc [task-name] [inputs...]`

### Positional Syntax

````markdown
## Tasks
### greet
```
echo "Hello, $1."
```
````

Result:

```sh
$ xc greet Joe
+ echo 'Hello, Joe.'
Hello, Joe.
```

### Variadic Syntax

````markdown
## Tasks
### greet-everyone
```
echo "Hello, $@."
```
````

Result:

```sh
$ xc greet Joe Bob Steve
+ echo 'Hello, Joe Bob Steve.'
Hello, Joe Bob Steve.
```

### Named Inputs

The `Inputs` attribute can be used to denote required parameters.

````markdown
## Tasks
### greet

Inputs: FORENAME, SURNAME

```
echo "Hello, $FORENAME $SURNAME."
```
````

Named Inputs can be passed via task arguments:

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

### Optional Inputs

Combining the `Environment` attribute and the `Inputs` attribute, you can provide optional inputs.

````markdown
## Tasks
### greet

Inputs: NAME
Environment: NAME=World

```
echo "Hello, $NAME."
```
````

Result:

```sh
$ xc greet Joe
+ echo 'Hello, Joe.'
Hello, Joe.
```

```sh
$ xc greet
+ echo 'Hello, World.'
Hello, World.
```

My Readme
==

Here is a readme.

## Subtitle
Sub text

tasks
---

### list

Lists files

```
ls
```
### list2

Lists files

Directory: ./somefolder

```
ls
```

### hello

Print a message

Requires: list, list2

Env: `somevar=val`
Inputs: FOO, BAR

```
echo "Hello, world!"
echo "Hello, world2!"
```

### all-lists

An example of a commandless task.

Requires: list, list2


## Out Of Scope


### something

Print a message

Requires: list, list2

```
echo "Hello, world!"
```

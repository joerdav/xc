---
title: "Command"
description:
linkTitle: "Command"
menu: { main: {  weight: 10 } }
---

## Help Text

```
{{< readfile file="/usage.txt" >}}
```

## Examples

`xc deploy` - runs a task named `deploy`

`xc deploy production` - runs a task named `deploy` with a single input `production`

`PLATFORM=linux xc build` - runs a task named `build` with a single input `PLATFORM` with the value `linux`

## Environment Variables

### `XC_TRACE`

Set `XC_TRACE` to "false", "no", or "0" to suppress trace output for shell commands.

```sh-session
$ xc greet
+echo "Hello, world!"
Hello, world!
$ XC_TRACE=no xc greet
Hello, world!
```

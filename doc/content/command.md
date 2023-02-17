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

---
linkTitle: IDE Support
title: IDE Support
description: Tools for using xc in IDEs
menu: main
weight: -6
---

## Visual Studio Code

Extension: <https://marketplace.visualstudio.com/items?itemName=xc-vscode.xc-vscode>

## vim

There is no vim plugin for `xc`, but [fzf.vim](https://github.com/junegunn/fzf.vim) can be used in
conjunction with [vim-run-interactive](https://github.com/christoomey/vim-run-interactive) to execute `xc` tasks.

Just use the following config:

```
:map <leader>xc :call fzf#run({'source':'xc -short', 'options': '--prompt "xc> " --preview "xc -md {}"', 'sink': 'RunInInteractiveShell xc', 'window': {'width': 0.9, 'height': 0.6}})
```

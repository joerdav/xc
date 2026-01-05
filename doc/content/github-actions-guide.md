---
linkTitle: GitHub Actions Guide
title: GitHub Actions Guide
description: Using xc for CI with GitHub Actions
menu: main
weight: -6
---

# CI Quick Start

## Create a `README` or `README` if you don't have one yet.

You can use Markdown or org-mode. This guide shows the Markdown style, but you can see an org-mode example [here](https://raw.githubusercontent.com/joerdav/xc/refs/heads/main/parser/parseorg/testdata/example.org).

## Add a heading called `Tasks` to your README with subheaders `Build`, `Test`, `Deploy`

(or adapt to your needs)  
Insert a code block under each one. Examples:

````md
## Tasks
### Build

```sh
uv lock --check
uv build --build-constraint constraints.txt --require-hashes
```

### Test

```sh
uv run pytest
```

### Deploy

```sh
rsync dist/ ci@staging:/app/dist
```
````
## Create a GitHub Actions Workflow

Choose from the [run-xc](https://github.com/joerdav/run-xc) action, which runs a single task at a time, or the [setup-xc](https://github.com/joerdav/setup-xc/) action, which installs `xc` on your runner so that you can incorporate it into scripts.

Sample flow:

`.github/workflows/push.yml`:
```yml
name: push
on: ["push", "pull_request"]

jobs:
  ci:
    name: "Run CI"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 1
      - uses: joerdav/run-xc@v1.1.0
        with:
          task: "build"
      - uses: joerdav/run-xc@v1.1.0
        with:
          task: "test"
      - uses: joerdav/run-xc@v1.1.0
        with:
          task: "deploy"
   ```
## Test on GitHub

Commit your README and workflow yaml and push to GitHub. Then make a pull request and verify that your CI runs as you expect.

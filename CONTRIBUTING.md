# Contributing to xc

Thank you for using `xc` and for your interest in contributing. Please take a moment to get familiar with our guidelines.

## Be Courteous
We expect you to read and follow our [Code of Conduct](./CODE_OF_CONDUCT.md) and to conduct yourself always in a way that creates a positive environment and clearly demonstrates your courtesy.

## Use GitHub Issues

We track and discuss work using Issues. Please [open an issue](https://github.com/joerdav/xc/issues/new/choose) to discuss changes before sending a pull request. We may close a PR without review if it's opened without discussion, especially by an unknown contributor. Whether we choose to review a PR is based on a judgment call about the maintenance burden it adds.

## Write Tests, Maintain Coverage

We expect authors of net-new code to provide tests covering the new code. Prior to merging, CI will verify that all tests pass and that coverage equals or exceeds the pre-merge state.

## Document Your Work

We expect contributors to keep documentation (in `docs/`) up to date, reflecting all substantive changes. Use `xc run-docs` (requires `hugo`) to verify that your changes build successfully & look as you expect.

# HOWTOs

## Add Support For A Markup Language
We aim to support popular, stable, lightweight plain-text markup languages. Our priorities for language support:
- New languages should work as similar to the Markdown implementation as possible
- Users should be able to read and maintain tasks without writing excessive markup or understanding implicit behavior
- New languages should be ubiquitous enough that it makes sense for the project to take on the maintenance burden to support them going forward. See for example the [list of markup languages supported by GitHub](https://github.com/github/markup/blob/master/README.md#markups) to get a rough idea.

Steps to add support:
- Create a GitHub issue to propose the addition and discuss
- Fork the repo and add a new parser module for the new language under `parser/`  
  *It may be advantageous to copy an existing parser and modify it to fit the particulars of the new addition.*
- [Write tests](#write-tests-maintain-coverage) & test fixtures in your parser module that cover all the same cases as are present in `parser/parsemd/testdata/`
- Update the `cmd/xc/` module to use your parser
- [Document your work](#document-your-work)
- Send a PR for your changes, linking the issue you created previously

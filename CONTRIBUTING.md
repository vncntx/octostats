# Contribution Guide

Please [contact the author](mailto:vincent@vincent.click) if you have any questions or want to help out.

## Style

Go code must be formatted with [gofmt](https://golang.org/cmd/gofmt/) and follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.

This project adheres to the [Conventional Commit](https://www.conventionalcommits.org) specifications.

## Getting Started

This project uses [Powershell Core](https://microsoft.com/PowerShell) to run tasks. To get started,

```ps1
Import-Module -Force ./tasks.psm1
Get-Command -Module tasks
```

The _hooks_ directory contains scripts that enforce the style guide and run tests before every commit. You can set them up by running `Install-Project`.

![](/icon.svg)

# octostats

Get stats for your pull requests on GitHub, including

- Average time to merge
- Average number of reviews
- Number of pull requests with each label

## Getting Started

To install and build the project using Powershell, run the following

```ps1
Import-Module -Force ./tasks.psm1
Install-Project
Build-Project
```

You can also use the [go](https://golang.org/doc/cmd) command.

## Authentication

Set the following environment variables for authentication: 
- `GITHUB_USER` - GitHub login username
- `GITHUB_TOKEN` - GitHub personal access token

Read more about creating a personal access token [here](https://docs.github.com/en/articles/creating-a-personal-access-token-for-the-command-line). For public repositories, you'll need the 'public_repo' scope. To access private repositories, you'll need the complete 'repo' scope. Learn more [here](https://docs.github.com/en/developers/apps/scopes-for-oauth-apps).

## Usage

```
octostats owner/repo [ date | -duration ]
```

### Examples

To inspect all your pull requests for a repo,

```sh
octostats owner/repo
```

To inspect your pull requests merged in the last 100 hours,

```sh
octostats owner/repo -100h
```

The duration _must_ be negative. See [`time.ParseDuration`](https://golang.org/pkg/time/#ParseDuration) to learn about valid duration formats.

To inspect your pull requests merged since a given date,

```sh
octostats owner/repo 2019-07-23
```

## Development

Please read the [Contribution Guide](./CONTRIBUTING.md) before you proceed.

## Copyright
Copyright 2021 [Vincent Fiestada](https://vincent.click). Released under a [BSD-3-Clause License](./license).

Icon made by [Dave Gandy](https://www.flaticon.com/authors/dave-gandy).

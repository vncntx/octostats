#!/usr/bin/env pwsh
# Copyright 2021 Vincent Fiestada

. (Join-Path 'tools' 'std' 'std.ps1')
. (Join-Path 'tools' 'std' 'build.ps1')
. (Join-Path 'tools' 'std' 'help.ps1')
. (Join-Path 'tools' 'go' 'install.ps1')
. (Join-Path 'tools' 'go' 'format.ps1')
. (Join-Path 'tools' 'go' 'check.ps1')
. (Join-Path 'tools' 'go' 'test.ps1')
. (Join-Path 'tools' 'go' 'run.ps1')
. (Join-Path 'tools' 'go' 'build.ps1')

function Invoke-Tools {
    param(
        [Parameter(ValueFromPipeline=$true)]
        [String] $Command = '',

        [Parameter(ValueFromPipeline=$true)]
        [String[]] $Arguments = @()
    )

    if ($null -eq $Command.Length) {
        Write-Error 'command required'
        exit [Errors]::NoCommand
    }

    $tools = @(
        [Tool]::new(
            'help',
            'list available commands',
            {
                Get-Toolkit $tools
            }
        ),
        [Tool]::new(
            'install',
            'check dev environment and install project',
            {
                Install-GoProject
            }
        ),
        [Tool]::new(
            'format',
            'apply style guide and tidy dependencies',
            {
                Invoke-GoFormat
            }
        ),
        [Tool]::new(
            'check',
            'detect issues using linters',
            {
                Invoke-GoChecks
            }
        ),
        [Tool]::new(
            'fix',
            'apply autofixes suggested by linters',
            {
                Invoke-GoChecks -Fix
            }
        ),
        [Tool]::new(
            'test',
            'run all tests with coverage',
            {
                Invoke-GoTests
            }
        ),
        [Tool]::new(
            'run',
            "run `e[3m[args]`e[3m",
            'compile and run',
            {
                Invoke-GoRun $Arguments
            }
        ),
        [Tool]::new(
            'build',
            'compile an executable binary',
            {
                Build-GoBinary 'bin' 'octostats_snapshot'
            }
        ),
        [Tool]::new(
            'release',
            "release `e[3mversion`e[3m",
            'build for all supported platforms',
            {
                $version = $Arguments[0]
                Build-GoBinary 'bin' "vincent.click/pkg/octostats_$version" @(
                    [BuildTarget]::new('darwin', 'amd64'),
                    [BuildTarget]::new('darwin', 'arm64'),
                    [BuildTarget]::new('linux', '386'),
                    [BuildTarget]::new('linux', 'amd64'),
                    [BuildTarget]::new('linux', 'arm'),
                    [BuildTarget]::new('linux', 'arm64'),
                    [BuildTarget]::new('windows', '386'),
                    [BuildTarget]::new('windows', 'amd64'),
                    [BuildTarget]::new('windows', 'arm')
                )
            }
        )
    )

    $target = $tools | Where-Object { $_.Command -eq $Command }
    if (-not $target) {
        Write-Error "invalid command '$Command'"
        exit [Error]::InvalidCommand
    }

    Invoke-Command -ScriptBlock $target.Script
}

enum Error {
    NoCommand = 1
}

$command = $args[0]
$arguments = $args[1..($args.Length - 1)]

Invoke-Tools $command $arguments

# Copyright 2021 Vincent Fiestada

. (Join-Path 'tools' 'std' 'std.ps1')
. (Join-Path 'tools' 'go' 'run.ps1')

function Invoke-GoBenchmarks {
    param(
        [String]$Package = './benchmarks'
    )

    Write-Info 'running benchmark tests'
    Write-Host ''

    $timeStart = (Get-Date)

    Invoke-GoPackage $Package

    $timeStop = (Get-Date)
    $durationSeconds = ($timeStop - $timeStart).TotalSeconds

    Write-Host ''
    Write-Divider
    if (Assert-ExitCode 0) {
        Write-Ok "benchmarks completed in $durationSeconds s"
    } else {
        Write-Fail "encountered an error while running benchmarks"
        exit [Error]::Fail
    }
}

enum Error {
    Fail = 1
}

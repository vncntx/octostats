# Copyright 2021 Vincent Fiestada

. (Join-Path 'tools' 'std' 'std.ps1')

class GoModule {
    [String]$Name
    [String]$Target

    GoModule([String]$Name, [String]$Target) {
        $this.Name = $Name
        $this.Target = $Target
    }
}

<#
.SYNOPSIS
Get information about the current Go module

.EXAMPLE
Get-GoModule
#>
function Get-GoModule {
    $modFile = Join-Path (Get-Location) 'go.mod'
    if (-not (Test-Path $modFile)) {
        Write-Error 'mod file not found'
        exit [Error]::NoGoModule
    }

    $contents = Get-Content $modFile
    if ($contents.Length -lt 3) {
        Write-Error 'invalid mod file'
        exit [Error]::InvalidGoMod
    }

    $name = ($contents[0] -split ' ')[-1]
    $target = ($contents[2] -split ' ')[-1]

    return [GoModule]::new($name, $target)
}

enum Error {
    NoGoModule = 1
    InvalidGoMod = 2
}

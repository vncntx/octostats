# Copyright 2021 Vincent Fiestada

. (Join-Path 'tools' 'std' 'std.ps1')

<#
.SYNOPSIS
List available toolkit commands

.DESCRIPTION
Print a list of toolkit commands

.EXAMPLE
Get-Toolkit [Command]
#>
function Get-Toolkit {
    param(
        [Parameter(Mandatory=$true, ValueFromPipeline=$true)]
        [Tool[]]$Tools
    )

    if ($Tools.Length -eq 0) {
        Write-Info 'no commands available'
        return
    }

    Write-Info 'available commands:'
    Write-Host ''
    foreach ($tool in $Tools) {
        Write-Host -NoNewLine " $($tool.Example)"
        Write-Host " - $($tool.Description)"
    }
}

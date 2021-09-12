# Copyright 2021 Vincent Fiestada

. (Join-Path 'tools' 'std' 'std.ps1')
. (Join-Path 'tools' 'std' 'build.ps1')
. (Join-Path 'tools' 'go' 'mod.ps1')

<#
.SYNOPSIS
Build an executable binary

.EXAMPLE
Build-GoBinary --help
#>
function Build-GoBinary {
    param(
        [Parameter(Mandatory=$true, ValueFromPipeline=$true)]
        [String]$Directory,

        [Parameter(Mandatory=$true, ValueFromPipeline=$true)]
        [String]$Name,

        [Parameter(ValueFromPipeline=$true)]
        [BuildTarget[]]$Targets = @()
    )
    $module = (Get-GoModule).Name

    if ($Targets.Length -lt 1) {
        $os = (go env GOOS)
        $arch = (go env GOARCH)
        
        $Targets = @([BuildTarget]::new($os, $arch))
    }

    if (-Not (Test-Path $Directory)) {
        Write-Info "creating '$Directory' directory"
        New-Item -ItemType Directory -Name $Directory > $null
    }
    $dir = Resolve-Path -Relative $Directory
    $checksums = Join-Path $dir "checksums.txt"
    
    if (Test-Path -Path $checksums -ErrorAction SilentlyContinue) {
        # remove existing checksums
        Remove-Item $checksums
    }

    $saved = @{
        'GOOS' = $env:GOOS;
        'GOARCH' = $env:GOARCH;
    }

    $c = 0
    $e = 0
    try {
        foreach ($target in $Targets) {
            $os = $target.System
            $arch = $target.Architecture

            $file = $Name + '_' + $os + '_' + $arch
            if ($os -eq 'windows') {
                $file += ".exe"
            }

            $binary = Join-Path $dir $file

            $env:GOOS = $os
            $env:GOARCH = $arch

            go build -o $binary $module

            if (Assert-ExitCode 0) {
                $c++
                $sum = (Get-FileHash -Algorithm SHA256 -Path $binary)
                Write-Ok "binary built for $os, $arch"
                Write-Output "$($sum.Algorithm) $($sum.Hash) $binary" >> $checksums
            } else {
                $e++
                Write-Error "error while building for $os, $arch"
            }
        }
    } finally {
        # restore environment variables
        $env:GOOS = $saved['GOOS']
        $env:GOARCH = $saved['GOARCH']
    }

    Write-Divider
    if ($e -lt 1) {
        Write-Ok "binaries built for $c targets"
    } else {
        Write-Error "$e errors encountered"
    }
}

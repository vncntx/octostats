# Copyright 2021 Vincent Fiestada

. (Join-Path 'tools' 'std' 'std.ps1')
. (Join-Path 'tools' 'go' 'mod.ps1')

<#
.SYNOPSIS
Publish to pkg.go.dev

.DESCRIPTION
Add the specified version to pkg.go.dev
Expects the version tag to exist in the repository

.EXAMPLE
Publish-GoModule v2.0.0
#>
function Publish-GoModule {
    param(
        [Parameter(Mandatory=$true, ValueFromPipeline=$true)]
        [String]$Version
    )

    Confirm-DeclaredVersion $Version

    $module=(Get-GoModule).Name
	Write-Info "publishing $Version of $module to pkg.go.dev"
    Invoke-WebRequest -Uri "http://proxy.golang.org/${module}/@v/${version}.info" > $null

    Write-Divider
    Write-Ok "published to https://pkg.go.dev/mod/${module}@${version}"
}

<#
.SYNOPSIS
Confirm the declared version in the doc file

.EXAMPLE
Confirm-DeclaredVersion v2.0.0
#>
function Confirm-DeclaredVersion {
    param(
        [Parameter(Mandatory=$true, ValueFromPipeline=$true)]
        [String]$Version
    )

    $doc = 'doc.go'
    if (-not (Test-Path -Path $doc -ErrorAction SilentlyContinue)) {
        Write-Info 'version not declared in code'
        return
    }

    if ($Version.StartsWith("v")) {
        $Version = $Version.Substring(1)
    }
    
    if ((Get-Content $doc | Select-String -Pattern "Version = `"$Version`"").Length -gt 0) {
        Write-Ok 'declared version is correct'
    } else {
        Write-Error "declared version in '$doc' is not $Version"
        exit [Error]::WrongVersion
    }
}

enum Error {
    WrongVersion = 1
}

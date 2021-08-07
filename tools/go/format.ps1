# Copyright 2021 Vincent Fiestada

. (Join-Path 'tools' 'std' 'std.ps1')

<#
.SYNOPSIS
Format code and optimize dependencies

.EXAMPLE
Invoke-GoFormat
#>
function Invoke-GoFormat {
    $e = @{}

    $e += (Format-GoCode)
    $e += (Format-GoMod)
    $e += (Optimize-Imports)

    foreach ($err in $e.Values) {
        $err.Write()
    }

     Write-Divider
    
    if ($e.Count -eq 0) {
        Write-Ok "go style guide applied"
    } else {
        $_errors = Pluralize -Word 'error' -Number $e.Count
        Write-Error "$($e.Count) $_errors encountered"
        exit [Error]::Fail
    }
} 

<#
.SYNOPSIS
Format all Go code in the project

.EXAMPLE
Format-GoCode
#>
function Format-GoCode {
    $errs = [CodeErrorCollection]::new()
    $fileCount = 0

    Get-ChildItem -Path . -Recurse -Filter *.go | ForEach-Object {
        $file = $_
        $fileCount++
        
        gofmt -l -w $file 2>&1 | ForEach-Object {
            $log = $_

            $pattern = '^(?<file>.+):(?<loc>[0-9]+:[0-9]+): (?<txt>.*)$'
            switch -Regex ($log) {
                $pattern {
                    $parsed = Select-String -Pattern $pattern -InputObject $log
                    $path, $loc, $txt = $parsed.Matches.Groups[1..3].Value
                    $errs.Add($path, $loc, $txt)
                }
                default {
                    Write-Ok "formatted '$file'"
                }
            }
        }
    }

    Write-Ok "examined $fileCount go files"

    return $errs.Errors
}


<#
.SYNOPSIS
Optimize the declared module dependencies

.EXAMPLE
Format-GoMod
#>
function Format-GoMod {
    $errs = [CodeErrorCollection]::new()
    $mod = Resolve-Path 'go.mod'

    if (-not (Test-Path -ErrorAction SilentlyContinue $mod)) {
        Write-Warning 'no dependencies to tidy up'
        return $err
    }

    go mod tidy 2>&1 | ForEach-Object {
        $log = $_

        $pattern = '^.+:(?<loc>[0-9]+): (?<txt>.*)$'
        switch -Regex ($log) {
            $pattern {
                $errs.Add($mod, $log, $pattern)
            }
            default {
                Write-Info $log
            }
        }
    }

    if ($errs.Count() -eq 0) {
        Write-Ok 'dependencies tidied up'
    }

    return $errs.Errors
}

<#
.SYNOPSIS
Optimize and sort package imports

.EXAMPLE
Optimize-Imports
#>
function Optimize-Imports {
    $errs = [CodeErrorCollection]::new()
    
    if (-not $env:GOBIN) {
        Write-Warning '$env:GOBIN is not set'
        return $errs.Errors
    }
    if (-not (Get-Command -Name (Join-Path $env:GOBIN 'gci') -ErrorAction SilentlyContinue)) {
        Write-Warning 'unable to optimize imports'
        return $errs.Errors
    }

    Invoke-Expression "$(Join-Path $env:GOBIN 'gci') -w ." 2>&1 | ForEach-Object {
        $log = $_

        $skip = '^skip file'
        switch -Regex ($log) {
            $skip {
                # ignore
            }
            default {
                $errs.Add((Get-Location).Path, $log)
            }
        }
    }

    if ($errs.Count() -eq 0) {
        Write-Ok 'imports optimized'
    }

    return $errs.Errors
}

enum Error {
    Fail = 1
}

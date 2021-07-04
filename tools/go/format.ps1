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

    Get-ChildItem -Path . -Recurse -Filter *.go | ForEach-Object {
        $file = $_
        
        gofmt -l -w $file 2>&1 | ForEach-Object {
            $log = $_

            $pattern = '^(?<file>.+):(?<loc>[0-9]+:[0-9]+): (?<txt>.*)$'
            switch -Regex ($log) {
                $pattern {
                    $errs.Add($log, $pattern)
                }
                default {
                    Write-Ok "formatted '$file'"
                }
            }
        }
    }

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

    if ($err.Count -eq 0) {
        Write-Ok 'dependencies tidied up'
    }
    return $errs.Errors
}

enum Error {
    Fail = 1
}

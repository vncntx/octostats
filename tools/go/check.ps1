# Copyright 2021 Vincent Fiestada

. (Join-Path 'tools' 'std' 'std.ps1')

class GoIssue : CodeIssue {
    static [GoIssue[]]FromJson([String]$Json) {
        $issues = @()
        $parsed = ConvertFrom-Json $Json

        if ($parsed.Issues) {
            foreach ($i in $parsed.Issues) {
                $exc = [CodeExcerpt]::new($i.Pos.Filename, $i.Pos.Line, $i.Pos.Column, $i.SourceLines)
                $issues += [GoIssue]::new($i.Text, $i.Fromlinter, $exc)
            }
        }

        return $issues
    }

    GoIssue([String]$Text, [String]$Linter, [CodeExcerpt]$Excerpt) : base($Text, $Linter, $Excerpt) {}
}

<#
.SYNOPSIS
Examine code for common flaws

.DESCRIPTION
Examine Go source code and report suspicious constructs

.EXAMPLE
Invoke-GoChecks
#>
function Invoke-GoChecks {
    param(
        [Switch]$Fix = $false
    )

    $JSON_PATTERN = '^{.*}$'
    $LOG_PATTERN = '^level=(?<level>[a-z]+) msg="(?<message>.+)"$'

    Write-Info 'examining packages'
    if ($Fix) {
        Write-Info 'autofixes enabled'
    }

    $e = 0
    $w = 0

    golangci-lint run --fix=$Fix --out-format=json 2>&1 | ForEach-Object {
        $log = $_

        switch -Regex ($log) {
            $JSON_PATTERN {
                $issues = [GoIssue]::FromJson($log)
                $e += $issues.Count

                foreach ($i in $issues) {
                    $i.Write()
                }
            }
            $LOG_PATTERN {
                $parsed = Select-String -Pattern $LOG_PATTERN -InputObject $log
                $level, $message = $parsed.Matches[0].Groups[1..2]

                switch ($level) {
                    'error' {
                        $e++
                        Write-Error $message
                    }
                    'warning' {
                        $w++
                        Write-Warning $message
                    }
                    default {
                        Write-Info $message
                    }
                }
            }

            default {
                Write-Info $log
            }
        }
    }

    $results = @()
    if ($e -gt 0) {
        $_errors = Pluralize -Word 'error' -Number $e
        $results += "$e $_errors"
    }
    if ($w -gt 0) {
        $_warnings = Pluralize -Word 'warning' -Number $w
        $results += "$w $_warnings"
    }
    $message = ($results -Join ', ')

    Write-Divider
    if ($e -eq 0) {
        Write-Ok 'no problems detected'
    } else {
        Write-Error $message
        exit [Error]::Fail
    }
}

enum Error {
    Fail = 1
}

# Copyright 2021 Vincent Fiestada

$env:LOG_PADDING = 5

class Tool {
    [String]$Command
    [String]$Example
    [String]$Description
    [ScriptBlock]$Script

    Tool([String]$Command, [String]$Description, [ScriptBlock]$Script) {
        $this.Command = $Command
        $this.Example = $Command
        $this.Description = $Description
        $this.Script = $Script
    }

    Tool([String]$Command, [String]$Example, [String]$Description, [ScriptBlock]$Script) {
        $this.Command = $Command
        $this.Example = $Example
        $this.Description = $Description
        $this.Script = $Script
    }
}

class CodeError {
    [String]$File
    [String[]]$Text

    CodeError([String]$File) {
        $this.File = $File 
        $this.Text = @()
    }

    CodeError([String]$File, [String[]]$Text) {
        $this.File = $File
        $this.Text = $Text
    }

    Write() {
        $padding = $env:LOG_PADDING ?? 5
        $relPath = (Resolve-Path -Path $this.File -Relative)

        Write-Host -NoNewline -ForegroundColor Red 'error'.PadLeft($padding)
        Write-Host -NoNewline ' : '
        Write-Host -NoNewline -ForegroundColor Black -BackgroundColor White " $relPath "
        Write-Host ''
        Write-Host ''
        foreach ($txt in $this.Text) {
            Write-Host ''.PadLeft($padding) "  $txt"
        }
        Write-Host ''
    }
}

class CodeErrorCollection {
    [HashTable]$Errors

    CodeErrorCollection() {
        $this.Errors = @{}
    }

    Add([String]$File, [String]$Location, [String]$Text) {
        if (-not $this.Errors.ContainsKey($File)) {
            $this.Errors[$File] = [CodeError]::new($File)
        }
        $this.Errors[$File].Text += "`e[4m$Location`e[24m : $Text"
    }

    Add([String]$Log, [String]$Pattern) {
        $parsed = Select-String -Pattern $Pattern -InputObject $Log
        $file, $loc, $txt = $parsed.Matches.Groups[1..3].Value
        $this.Add($file, $loc, $txt)        
    }

    [Int]Count() {
        return $this.Errors.Count
    }

    Write() {
        foreach ($e in $this.Errors.Values) {
            $e.Write()
        }
    }
}

class CodeExcerpt {
    [String]$File
    [Int]$Line
    [Int]$Column
    [String[]]$Code
    
    CodeExcerpt([String]$File, [Int]$Line, [Int]$Column, [String[]]$Code) {
        $this.Code = $Code
        $this.File = $File
        $this.Line = $Line
        $this.Column = $Column
    }
}

class CodeIssue {
    [String]$Text
    [String]$Linter
    [CodeExcerpt]$Excerpt

    CodeIssue([String]$Text, [String]$Linter, [CodeExcerpt]$Excerpt) {
        $this.Text = $Text
        $this.Linter = $Linter
        $this.Excerpt = $Excerpt
    }

    Write() {
        $padding = $env:LOG_PADDING ?? 5

        Write-Host -NoNewline -ForegroundColor Red 'error'.PadLeft($padding)
        Write-Host -NoNewline ' : '
        Write-Host -NoNewline -ForegroundColor Black -BackgroundColor White " $($this.Linter) "
        Write-Host " in `e[4m$($this.Excerpt.File)`e[24m @ $($this.Excerpt.Line):$($this.Excerpt.Column)".PadLeft($padding)
        Write-Host ''
        Write-Host ''.PadLeft($padding) "  $($this.Text)"
        Write-Host ''
    }
}


<#
.SYNOPSIS
Verify the last command's exit code

.DESCRIPTION
Return true if the last command's exit code is equal to the argument

.EXAMPLE
Assert-ExitCode 100
#>
function Assert-ExitCode {
    param(
        [Parameter(Mandatory=$true, ValueFromPipeline=$true)]
        [Int32]$expected
    )

    return ($LASTEXITCODE -eq $expected)
}

<#
.SYNOPSIS
Log a message

.EXAMPLE
Write-Log -Level "info" -Message "build complete" -Color Cyan
#>
function Write-Log {
    param(
        [String]$Level,
        [Object]$Message,
        [String]$Color,
        [Int32]$Padding = $env:LOG_PADDING ?? 5
    )

    Write-Host -NoNewline -ForegroundColor $Color $Level.PadLeft($Padding)
    Write-Host " : $Message"
}

<#
.SYNOPSIS
Log a successful event

.EXAMPLE
Write-Ok "build complete"
#>
function Write-Ok {
    param(
        [Parameter(ValueFromPipeline=$true)]
        [Object]$Message
    )

    Write-Log -Level 'ok' -Message $Message -Color Green
}

<#
.SYNOPSIS
Log a passing test or check

.EXAMPLE
Write-Pass "unit tests"
#>
function Write-Pass {
    param(
        [Parameter(ValueFromPipeline=$true)]
        [Object]$Message
    )

    Write-Log -Level 'pass' -Message $Message -Color Green
}

<#
.SYNOPSIS
Log a failing test or check

.EXAMPLE
Write-Fail "unit tests"
#>
function Write-Fail {
    param(
        [Parameter(ValueFromPipeline=$true)]
        [Object]$Message
    )

    Write-Log -Level 'fail' -Message $Message -Color Red
}

<#
.SYNOPSIS
Log a skipped test or check

.EXAMPLE
Write-Skip "unit tests"
#>
function Write-Skip {
    param(
        [Parameter(ValueFromPipeline=$true)]
        [Object]$Message
    )

    Write-Log -Level 'skip' -Message $Message -Color Yellow
}

<#
.SYNOPSIS
Log a warning

.EXAMPLE
Write-Warning "file missing"
#>
function Write-Warning {
    param(
        [Parameter(ValueFromPipeline=$true)]
        [Object]$Message
    )

    Write-Log -Level 'warn' -Message $Message -Color Yellow
}

<#
.SYNOPSIS
Log an error

.EXAMPLE
Write-Error "error"
#>
function Write-Error {
    param(
        [Parameter(ValueFromPipeline=$true)]
        [Object]$Message
    )

    Write-Log -Level 'error' -Message $Message -Color Red
}

<#
.SYNOPSIS
Log some information

.EXAMPLE
Write-Info "preparing environment"
#>
function Write-Info {
    param(
        [Parameter(ValueFromPipeline=$true)]
        [Object]$Message
    )

    Write-Log -Level 'info' -Message $Message -Color Blue
}

<#
.SYNOPSIS
Print a horizontal divider across the entire terminal

.EXAMPLE
Write-Divider
#>
function Write-Divider {
    Write-Host ("-" * ($env:LOG_PADDING ?? 5))
}

<#
.SYNOPSIS
Select the singular or plural version of a word based on a number

.EXAMPLE
Pluralize -Word 'apple' -Number 3
#>
function Pluralize {
    param(
        [Parameter(Mandatory=$true)]
        [String]$Word,

        [Parameter(ValueFromPipeline=$true)]
        [Int]$Number = 1
    )

    if ($Number -eq 1) {
        return $Word
    }

    return $Word + 's'
}

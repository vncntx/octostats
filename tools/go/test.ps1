# Copyright 2021 Vincent Fiestada

. (Join-Path 'tools' 'std' 'std.ps1')

enum Result {
    None = 0
    Skip = 1
    Pass = 2
    Fail = 3
}

class GoTest {
    [String]$Name
    [Result]$Result
    [String[]]$Errors

    GoTest([String]$Name) {
        $this.Name = $Name
        $this.Result = [Result]::None
        $this.Errors = @()
    }

    AddError([String]$Log) {
        $LOG_PATTERN = '^.+:\d+: (?<message>.+)$'
        $PANIC_PATTERN = '^(?<message>panic: .+)$'

        switch -Regex ($Log) {
            $LOG_PATTERN {
                $parsed = Select-String -Pattern $LOG_PATTERN -InputObject $Log
                $message = $parsed.Matches.Groups[1].Value
                $this.Errors += $message
            }
            $PANIC_PATTERN {
                $parsed = Select-String -Pattern $PANIC_PATTERN -InputObject $Log
                $message = $parsed.Matches.Groups[1].Value
                $this.Errors += $message
            }
            default {
                # ignore
            }
        }
    }

    Write() {
        $n = $this.Name
        $prefix = "Test"
        if ($n.StartsWith($prefix)) {
            $n = $n.Remove(0, $prefix.Length)
        }

        switch($this.Result) {
            Skip {
                Write-Skip $n
            }
            Pass {
                Write-Pass $n
            }
            Fail {
                Write-Fail "`e[3m$n`e[0m"
                foreach ($err in $this.Errors) {
                    Write-Host ''
                    Write-Host (Get-Pad) "  $err"
                }
                Write-Host ''
            }
        }
    }
}

class GoPackage {
    [String]$Name
    [Double]$Coverage
    [HashTable]$Tests

    GoPackage([String]$Name) {
        $this.Name = $Name
        $this.Coverage = 0
        $this.Tests = @{}
    }

    AddTestResult([String]$Name, [Result]$Result) {
        if (-not $this.Tests.ContainsKey($Name)) {
            $this.Tests[$Name] = [GoTest]::new($Name)
        }
        $this.Tests[$Name].Result = $Result
    }

    AddTestLog([String]$Name, [String]$Log) {
        if (-not $this.Tests.ContainsKey($Name)) {
            $this.Tests[$Name] = [GoTest]::new($Name)
        }
        $this.Tests[$Name].AddError($Log)
    }

    [GoTest[]]FindByResult([Result]$Result) {
        return $this.Tests.Values | Where-Object { $_.Result -eq $Result }
    }

    Write() {
        Write-Host -NoNewline -ForegroundColor Blue " $($this.Name) "
        Write-Host "[ $($this.Coverage)% coverage, $($this.Tests.Count) tests ]"
    }
}

class GoTestCollection {
    [HashTable]$Packages

    GoTestCollection() {
        $this.Packages = @{}
    }

    AddTestResult([String]$Name, [String]$Package, [Result]$Result) {
        if (-not $this.Packages.ContainsKey($Package)) {
            $this.Packages[$Package] = [GoPackage]::new($Package)
        }
        $this.Packages[$Package].AddTestResult($Name, $Result)
    }

    AddTestLog([String]$Name, [String]$Package, [String]$Log) {
        if (-not $this.Packages.ContainsKey($Package)) {
            $this.Packages[$Package] = [GoPackage]::new($Package)
        }
        $this.Packages[$Package].AddTestLog($Name, $Log)
    }

    AddCoverage([String]$Package, [Double]$Coverage) {
        if (-not $this.Packages.ContainsKey($Package)) {
            $this.Packages[$Package] = [GoPackage]::new($Package)
        }
        $this.Packages[$Package].Coverage = $Coverage
    }

    [HashTable]Count() {
        $pkgs = $this.Packages.Values
        return @{
            [Result]::Skip = $this.CountTestResultsByResult($pkgs, [Result]::Skip);
            [Result]::Pass = $this.CountTestResultsByResult($pkgs, [Result]::Pass);
            [Result]::Fail = $this.CountTestResultsByResult($pkgs, [Result]::Fail)
        }
    }

    [Int]CountTestResultsByResult([GoPackage[]]$Packages, [Result]$Result) {
        $count = 0
        foreach ($pkg in $Packages) {
            $count += $pkg.FindByResult($Result).Count
        }
        
        return $count
    }

    [Double]Coverage() {
        $cov = 0.0
        $pkgs = $this.Packages.Values
        
        foreach ($pkg in $pkgs) {
            $cov += $pkg.Coverage
        }

        return $cov / $pkgs.Count
    }

    Write() {
        foreach ($package in $this.Packages.Values) {
            $package.Write()
            foreach ($test in $package.Tests.Values | Sort-Object -Property Result,Name) {
                $test.Write()
            }
        }
    }
}

<#
.SYNOPSIS
Run tests for the entire module

.EXAMPLE
Invoke-GoTests
#>
function Invoke-GoTests {
    $t = [GoTestCollection]::new()
    $e = [CodeErrorCollection]::new()

    $timeStart = (Get-Date)
    go test './...' --cover --json 2>&1 | ForEach-Object {
        $log = $_
        try {
            $x = ConvertFrom-Json $log
        }
        catch {
            $code_error_pattern = '^(?<file>.+):(?<loc>[0-9]+:[0-9]+): (?<txt>.*)$'
            if ($log -match $code_error_pattern) {
                $parsed = Select-String -Pattern $code_error_pattern -InputObject $log
                $path, $loc, $txt = $parsed.Matches.Groups[1..3].Value
                $e.Add($path, $loc, $txt)
            }
        }

        $name = $x.Test
        $package = $x.Package
        $action = $x.Action
        $log = $x.Output

        if ($name) {
            switch ($action) {
                'skip' {
                    $t.AddTestResult($name, $package, [Result]::Skip)
                }
                'pass' {
                    $t.AddTestResult($name, $package, [Result]::Pass)
                }
                'fail' {
                    $t.AddTestResult($name, $package, [Result]::Fail)
                }
                'output' {
                    $t.AddTestLog($name, $package, $log)
                }
                default {
                    # ignore
                }
            }
        } elseif ($package) {
            $coverage_pattern = 'coverage: (?<coverage>[0-9\.]+)% of statements'
            $no_tests_pattern = '\[no test files\]'

            switch -Regex ($log) {
                $coverage_pattern {
                    $parsed = Select-String -Pattern $coverage_pattern -InputObject $log
                    $coverage = $parsed.Matches[0].Groups[1]
    
                    $t.AddCoverage($package, [Double]::Parse($coverage))
                }
                $no_tests_pattern {
                    $t.AddCoverage($package, 0)
                }
                default {
                    # ignore
                }
            }
        }
    }
    $timeStop = (Get-Date)

    $t.Write()
    $e.Write()

    $results = @()
    $count = $t.Count()

    if ($count[[Result]::Pass] -gt 0) {
        $c = $count[[Result]::Pass]
        $results += "$c passing"
    }
    if ($count[[Result]::Fail] -gt 0) {
        $c = $count[[Result]::Fail]
        $results += "$c failing"
    }
    if ($count[[Result]::Skip] -gt 0) {
        $c = $count[[Result]::Skip]
        $results += "$c skipped"
    }
    if ($e.Count() -gt 0) {
        $results += "build errors detected"
    }
    if ($results.Length -lt 1) {
        $results += "no tests ran"
    } else {
        $coverage = $t.Coverage().ToString("#.##")
        $results += "$coverage% coverage"
    }
    
    $duration = ($timeStop - $timeStart).TotalSeconds
    $message = ($results -Join ', ') + ". completed in $duration s"
    $isOk = ($count[[Result]::Fail] -eq 0) -and ($e.Count() -eq 0)

    Write-Divider
    if ($isOk) {
        Write-Ok $message
    }
    elseif ($e.Count() -eq 0) {
        Write-Fail $message
        exit [Errors]::Fail
    }
    else {
        Write-Error $message
        exit [Errors]::Build
    }
}

enum Errors {
    Fail = 1
    Build = 2
}

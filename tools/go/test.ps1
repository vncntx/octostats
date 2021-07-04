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

    GoTest([String]$Name, [Result]$Result) {
        $this.Name = $Name
        $this.Result = $Result
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
                Write-Fail $n
            }
        }
    }
}

class GoPackage {
    [String]$Name
    [Double]$Coverage
    [GoTest[]]$Tests

    GoPackage([String]$Name) {
        $this.Name = $Name
        $this.Coverage = 0
        $this.Tests = @()
    }

    [GoTest[]]FindByResult([Result]$Result) {
        return $this.Tests | Where-Object { $_.Result -eq $Result }
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

    AddTest([String]$Name, [String]$Package, [Result]$Result) {
        if (-not $this.Packages.ContainsKey($Package)) {
            $this.Packages[$Package] = [GoPackage]::new($Package)
        }
        $this.Packages[$Package].Tests += [GoTest]::new($Name, $Result)
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
            [Result]::Skip = $pkgs | foreach { $count = 0 } { $count += ($_.FindByResult([Result]::Skip)).Count } { $count };
            [Result]::Pass = $pkgs | foreach { $count = 0 } { $count += ($_.FindByResult([Result]::Pass)).Count } { $count };
            [Result]::Fail = $pkgs | foreach { $count = 0 } { $count += ($_.FindByResult([Result]::Fail)).Count } { $count }
        }
    }

    Write() {
        foreach ($package in $this.Packages.Values) {
            $package.Write()
            foreach ($test in $package.Tests | Sort-Object -Property Result,Name) {
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
                $e.Add($log, $code_error_pattern)
            }
        }

        $name = $x.Test
        $package = $x.Package
        $action = $x.Action
        $output = $x.Output

        if ($name) {
            switch ($action) {
                'skip' {
                    $t.AddTest($name, $package, [Result]::Skip)
                }
                'pass' {
                    $t.AddTest($name, $package, [Result]::Pass)
                }
                'fail' {
                    $t.AddTest($name, $package, [Result]::Fail)
                }
                default {
                    # ignore
                }
            }
        } elseif ($package) {
            $coverage_pattern = "coverage: (?<coverage>[0-9\.]+)% of statements"
            if ($output -match $coverage_pattern) {
                $parsed = Select-String -Pattern $coverage_pattern -InputObject $output
                $coverage = $parsed.Matches[0].Groups[1]

                $t.AddCoverage($package, [Double]::Parse($coverage))
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

# Copyright 2021 Vincent Fiestada

<#
.SYNOPSIS
Compile and run the project

.EXAMPLE
Invoke-GoRun --help
#>
function Invoke-GoRun {
    param(
        [Parameter(ValueFromPipeline = $true)]
        [String[]]$Arguments = @()
    )

    Invoke-GoPackage . $Arguments
}


<#
.SYNOPSIS
Compile and run a package

.EXAMPLE
Invoke-GoPackage ./pkg --help
#>
function Invoke-GoPackage {
    param(
        [Parameter(Mandatory = $true, ValueFromPipeline = $true)]
        [String]$Package,

        [Parameter(ValueFromPipeline = $true)]
        [String[]]$Arguments = @()
    )

    go run $Package ($Arguments -Join ' ')
}

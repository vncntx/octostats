#!/usr/bin/env pwsh

Import-Module .\tasks.psm1

# Format code
Format-Project

# Run tests
Invoke-Tests

# Check code quality
Invoke-Checks
# Copyright 2021 Vincent Fiestada

class BuildTarget {
    [String]$System
    [String]$Architecture

    BuildTarget([String]$System, [String]$Architecture) {
        $this.System = $System
        $this.Architecture = $Architecture
    }
}

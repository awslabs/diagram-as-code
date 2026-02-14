# Auto-positioning

Automatic link position calculation for optimal diagram layouts.

## Overview

Auto-positioning automatically determines the best connection points for links based on resource positions, eliminating the need to manually specify `SourcePosition` and `TargetPosition`.

## Usage

### Default Behavior

Links use auto-positioning by default when positions are not specified:

```yaml
Links:
  - Source: ALB
    Target: Instance1
```

### Explicit Auto-positioning

```yaml
Links:
  - Source: ALB
    SourcePosition: auto
    Target: Instance2
    TargetPosition: auto
```

### Mixed Positioning

Combine manual and automatic positioning:

```yaml
Links:
  - Source: ALB
    SourcePosition: E      # Manual
    Target: Instance3
    TargetPosition: auto   # Automatic
```

## How It Works

The auto-positioning algorithm:
1. Calculates relative positions of source and target resources
2. Determines optimal connection points based on geometry
3. Selects positions that minimize link length and overlaps

## When to Use

**Use auto-positioning when**:
- Creating new diagrams quickly
- Resource positions may change
- Optimal positions are obvious

**Use manual positioning when**:
- Precise control needed
- Multiple links from same resource
- Specific aesthetic requirements

## Related Documentation

- [Links](../links.md) - Link basics
- [UnorderedChildren](unordered-children.md) - Automatic resource reordering
- [Link Grouping](link-grouping.md) - Prevent link overlap

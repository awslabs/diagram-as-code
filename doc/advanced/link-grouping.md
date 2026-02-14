# Link Grouping Offset

Prevent link overlap when multiple links originate from or terminate at the same position.

## Overview

When multiple links share the same source or target position, they can overlap and become difficult to distinguish. Link grouping offset automatically spreads these links apart.

## Usage

Enable for specific resources:

```yaml
Resources:
  ELB:
    Type: AWS::ElasticLoadBalancingV2::LoadBalancer
    Options:
      GroupingOffset: true
```

## How It Works

- Links from the same position are offset by ±5px, ±10px, etc.
- Links are sorted by target/source position for consistent ordering
- Offset is applied perpendicular to the link direction
- Calculation: `(index - (count-1)/2.0) * 10` pixels

## Example

```yaml
Resources:
  ELB:
    Type: AWS::ElasticLoadBalancingV2::LoadBalancer
    Options:
      GroupingOffset: true
      
Links:
  - Source: ELB
    Target: Instance1
    SourcePosition: S
    TargetPosition: N
  - Source: ELB
    Target: Instance2
    SourcePosition: S  # Will be automatically offset
    TargetPosition: N
  - Source: ELB
    Target: Instance3
    SourcePosition: S  # Will be automatically offset
    TargetPosition: N
```

**Result**: Links spread horizontally to prevent overlap.

## When to Use

**Enable when**:
- Multiple links from same resource position
- Links overlap and are hard to distinguish
- Load balancer to multiple instances

**Disable when** (default):
- Single link per position
- Manual positioning preferred
- Offset causes layout issues

## Related Documentation

- [Links](../links.md) - Link basics
- [Auto-positioning](auto-positioning.md) - Automatic link positioning

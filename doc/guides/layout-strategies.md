# Layout Strategies

Optimization techniques for creating clear and organized diagrams.

## Layout Principles

### 1. Flow Direction

Choose direction based on diagram type:

**Vertical (North-South)**:
- User request flows
- API call chains
- Data pipelines

```yaml
VPC:
  Direction: "vertical"
```

**Horizontal (East-West)**:
- Redundancy across AZs
- Peer relationships
- Service mesh

```yaml
AZGroup:
  Direction: "horizontal"
```

### 2. Grouping

Use stacks to organize resources:

```yaml
# Horizontal grouping for AZs
AZs:
  Type: AWS::Diagram::HorizontalStack
  Children:
    - AZ1Resources
    - AZ2Resources

# Vertical grouping for tiers
Tiers:
  Type: AWS::Diagram::VerticalStack
  Children:
    - WebTier
    - AppTier
    - DataTier
```

### 3. Alignment

Control resource alignment within groups:

```yaml
Resources:
  MyGroup:
    Type: AWS::Diagram::HorizontalStack
    Align: "top"  # top, center, bottom
```

## Common Patterns

### Three-Tier Architecture

```yaml
Application:
  Type: AWS::Diagram::VerticalStack
  Children:
    - LoadBalancers
    - WebServers
    - Databases
```

### Multi-AZ Deployment

```yaml
VPC:
  Type: AWS::EC2::VPC
  Direction: "horizontal"
  Children:
    - AZ1
    - AZ2
    - AZ3
```

### Hub-and-Spoke

```yaml
Network:
  Type: AWS::Diagram::Resource
  Children:
    - TransitGateway  # Center
    - VPC1           # Spoke
    - VPC2           # Spoke
    - VPC3           # Spoke
```

## Optimization Tips

1. **Minimize crossing links**: Arrange resources to reduce link crossings
2. **Use UnorderedChildren**: Let the tool optimize resource order
3. **Consistent spacing**: Use same grouping patterns throughout
4. **Limit nesting depth**: Keep hierarchy shallow (max 3-4 levels)

## Related Documentation

- [Best Practices](../best-practices.md) - Design patterns
- [UnorderedChildren](../advanced/unordered-children.md) - Automatic reordering
- [Large Diagrams](large-diagrams.md) - Scaling strategies

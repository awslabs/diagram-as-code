# BorderChildren

Place resources on the borders of their parent container.

## Overview

BorderChildren allows you to position resources on the edges (North, South, East, West) of their parent container, useful for gateways, load balancers, and boundary resources.

## Usage

```yaml
Resources:
  VPC:
    Type: AWS::EC2::VPC
    BorderChildren:
      - Position: North
        Child: InternetGateway
      - Position: South
        Child: VPNGateway
    Children:
      - PublicSubnet
      - PrivateSubnet
      
  InternetGateway:
    Type: AWS::EC2::InternetGateway
    
  VPNGateway:
    Type: AWS::EC2::VPNGateway
```

## Positions

Available positions:
- `North` - Top edge
- `South` - Bottom edge
- `East` - Right edge
- `West` - Left edge

## Common Use Cases

### Internet Gateway

```yaml
VPC:
  Type: AWS::EC2::VPC
  BorderChildren:
    - Position: North
      Child: IGW
```

### VPN Gateway

```yaml
VPC:
  Type: AWS::EC2::VPC
  BorderChildren:
    - Position: South
      Child: VPNGateway
```

### Transit Gateway Attachment

```yaml
VPC:
  Type: AWS::EC2::VPC
  BorderChildren:
    - Position: West
      Child: TGWAttachment
```

## Related Documentation

- [Resource Types](../resource-types.md) - Resource basics
- [Best Practices](../best-practices.md) - Design patterns

## Links

Links are lines that show relationships between resources. Currently supports a straight line between resources.

### Link position

The start and end points of the line specify the location as the 16-wind rose of the resource (for example, NNW).

![position](static/position.png)

```
Diagrams:
  Resources:
    ALB: ...
    PublicSubnet1Instance: ...
  Links:
    - Source: ALB # (required)
      SourcePosition: NNE # (required)
      Target: PublicSubnet1Instance # (required)
      TargetPosition: S # (required)
      LineWidth: 1 # (optional)
      LineColor: 'rgba(255,255,255,255)' # (optional)
      LineStyle: `normal|dashed` (optional)
```

### Link type

#### Straight
![straight link](static/link-straight.png)
```
  Links:
    - Source: StraightLambda
      SourcePosition: N
      Target: StraightBucket
      TargetPosition: W
      TargetArrowHead:
        Type: Open
```

#### Orthogonal
![orthogonal link](static/link-orthogonal.png)
```
  Links:
      # Orthogonal (single-arm)
    - Source: Orthogonal1Lambda
      SourcePosition: N
      Target: Orthogonal1Bucket
      TargetPosition: W
      TargetArrowHead:
        Type: Open
      Type: orthogonal

      # Orthogonal (double-arm)
    - Source: Orthogonal2Lambda
      SourcePosition: E
      Target: Orthogonal2Bucket
      TargetPosition: W
      TargetArrowHead:
        Type: Open
      Type: orthogonal

      # Orthogonal (double-arm)
    - Source: Orthogonal3Lambda
      SourcePosition: E
      Target: Orthogonal3Bucket
      TargetPosition: E
      TargetArrowHead:
        Type: Open
      Type: orthogonal
```


### Arrow head

Arrows add context and meaning to a diagram by indicating the direction of flow.

![arrow head](static/arrows.png)
(generated from [static/arrows.yaml](static/arrows.yaml))

```
    - Source: ALB
      SourcePosition: NNW
      SourceArrowHead: #(optional)
        Type: Open #(required) Open/Default
        Width: Default #  (optional) Narrow/Default/Wide default="Default"
        Length: 2 # (optional) default=2
      Target: VPCPublicSubnet1Instance
      TargetPosition: SSE
      TargetArrowHead: #(optional)
        Type: Open #(required) Open/Default
        Width: Default #  (optional) Narrow/Default/Wide default="Default"
        Length: 2 # (optional) default=2
```

### Link Labels

Link Labels add labels along the link

```
  Links:
    - Source: ALB # (required)
      SourcePosition: NNE # (required)
      Target: PublicSubnet1Instance # (required)
      TargetPosition: S # (required)
      Labels: (optional)
        SourceLeft: (optional)
          Type: (optional, default: horizontal, allowed: horizontal) 
          Title: (required on SourceLeft, default: ``)
          Color: (optional, default: `` inherit from Source,Target font color)
          Font: (optional, default: `` inherit from Source,Target font name)
        SourceRight: (optional)
          Type: (optional, default: horizontal, allowed: horizontal) 
          Title: (required on SourceRight, default: ``)
          Color: (optional, default: `` inherit from Source,Target font color)
          Font: (optional, default: `` inherit from Source,Target font name)
        TargetLeft: (optional)
          Type: (optional, default: horizontal, allowed: horizontal) 
          Title: (required on TargetLeft, default: ``)
          Color: (optional, default: `` inherit from Source,Target font color)
          Font: (optional, default: `` inherit from Source,Target font name)
        TargetRight: (optional)
          Type: (optional, default: horizontal, allowed: horizontal) 
          Title: (required on TargetRight, default: ``)
          Color: (optional, default: `` inherit from Source,Target font color)
          Font: (optional, default: `` inherit from Source,Target font name)
```



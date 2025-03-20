## Resouce Types

### AWS::Diagram::Canvas

Canvas is a resource type that represents a drawable area. The Canvas resource type doesn't draw anything, but it is a special resource type, and all resources must be reached from the Canvas. There must be only one Canvas resource type in the file.

![Dependency graph](static/dependency-graph.png)

### AWS::Diagram::Cloud

A resource type that indicates that it is within the AWS cloud. It is defined internally as `AWS::Diagram::Group` resource type. It is often used mainly in contrast to users and on-premises environments but is not required.

![AWS Cloud](static/awscloud.png)

```
Diagram:
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - AWSCloud
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
```

### AWS::Diagram::Group [DEPRECATED]

~~An essential resource type that indicates that it is a group of resources. The following attributes are defined for the group, and it is possible to customize the decoration with these attributes.~~
~~`AWS::Diagram::Group` resources are often used implicitly, but it is also possible to create custom groups explicitly.~~

### AWS::Diagram::Resource

An essential resource type that represents a single resource or group. The following attributes are defined for the resource, and it is possible to customize the decoration with these attributes.
`AWS::Diagram::Resource` type are used implicitly from other predefined types, but you can specify `AWS::Diagram::Resource` explicitly and create custom resources.

| Attribute   | Type     | Default Value                              | Description                                                             |
| ----------- | -------- | ------------------------------------------ | ----------------------------------------------------------------------- |
| Icon        | string   | ` `                                        | Icon file path                                                          |
| IconFill    | IconFill | `Type: none, Color: rgba(255,255,255,255)` | Filling icon background                                                 |
| Direction   | string   | `horizontal`                               | vertical: `left`,`center`,`right` horizontal: `top`, `center`, `bottom` |
| Preset      | string   | ` `                                        | Override resource attributes from preset                                |
| Align       | string   | `center`                                   | `vertical`, `horizontal`                                                |
| FillColor   | string   | `rgba(0,0,0,0)`                            | Only group.                                                             |
| BorderColor | string   | `rgba(0,0,0,0)`                            |                                                                         |
| Title       | string   | ` `                                        |                                                                         |
| HeaderAlign | string   | `left`                                     | Only group. You can align icon and title to left/center/right.          |
| Children    | []string | `[]`                                       |                                                                         |

#### Single resource

<table>
<tr>
<td>
<pre>
    Subnet:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
</pre>
</td>
<td>
<img src="static/single-resource.png">
</td>
</tr>
</table>

#### A resource with EC2 instance children

<table>
<tr>
<td>
<pre>
    Subnet:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance
    Instance:
      Type: AWS::EC2::Instance
</pre>
</td>
<td>
<img src="static/single-resource-with-ec2-instance.png">
</td>
</tr>
</table>

#### A resource with empty resource children

<table>
<tr>
<td>
<pre>
    Subnet:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - SubnetEmptyResource
    SubnetEmptyResource:
      Type: AWS::Diagram::Resource
</pre>
</td>
<td>
<img src="static/single-resource-with-empty-resource.png">
</td>
</tr>
</table>

### AWS::Diagram::VerticalStack

A resource type that indicates a vertical stack. It is treated internally as a Group resource type but is undecorated by default.
`left` alignment, `center` alignment, or `right` alignment can be specified with `align` attribute when stacking.

![Vertical Stack](static/vertical_stack.png)

### AWS::Diagram::HorizontalStack

A resource type that indicates a horizontal stack. It is treated internally as a Group resource type but is undecorated by default.
`top` alignment, `center` alignment, or `bottom` alignment can be specified with `align` attribute when stacking.

![Horizontal Stack](static/horizontal_stack.png)

### Other predefined resource types

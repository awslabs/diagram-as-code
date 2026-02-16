# Introduction Guide

- [Part 1: Quick Start (10 minutes)](#part-1-quick-start-10-minutes)
  - [Installation](#installation)
  - [Your First Diagram](#your-first-diagram)
  - [Next Steps](#next-steps)
- [Part 2: Understanding Diagram-as-Code](#part-2-understanding-diagram-as-code)
  - [Core Concepts](#core-concepts)
  - [DAC File Structure](#dac-file-structure)
    - [DefinitionFiles section](#definitionfiles-section)
    - [Resources section](#resources-section)
    - [Links section](#links-section)
  - [\[Beta\] CloudFormation Integration](#beta-cloudformation-integration)
  - [Customization Tips](#customization-tips)

---

# Part 1: Quick Start (10 minutes)

Get your first AWS architecture diagram in under 10 minutes.

## Installation

**For macOS users:**
```bash
brew install awsdac
```

**For Go developers:**
```bash
go install github.com/awslabs/diagram-as-code/cmd/awsdac@latest
```

## Your First Diagram

Create your first diagram using an example from the repository:

```bash
awsdac https://raw.githubusercontent.com/awslabs/diagram-as-code/main/examples/alb-ec2.yaml
```

Output:
```
[Completed] AWS infrastructure diagram generated: output.png
```

Open `output.png` to see your diagram! 🎉

**Tip**: Use `-o custom-name.png` to specify a different output filename.

## Next Steps

**Encountered an error?** → See [Troubleshooting](troubleshooting.md)

**Want to create your own diagram?** → Continue to [Part 2: Understanding Diagram-as-Code](#part-2-understanding-diagram-as-code)

**Need advanced features?** → Explore [Advanced Features](advanced/)

**Using AI assistants?** → Check out [MCP Server Integration](mcp-server.md)

---

# Part 2: Understanding Diagram-as-Code

## Core Concepts

Diagram-as-code (DAC) lets you define AWS architecture diagrams using YAML files. This approach enables:

- **Version Control**: Track diagram changes with Git
- **Code Reuse**: Share and reuse diagram components
- **Automation**: Generate diagrams in CI/CD pipelines
- **Consistency**: Maintain standardized diagram styles

## DAC File Structure

To create diagrams with `awsdac`, you need to provide a "DAC (diagram as code) file" that defines the components and layout of your architecture. The DAC file uses YAML syntax and consists of 3 main sections:

```yaml
Diagram:
    DefinitionFiles:  # Specify the location of the definition file
      ...
    Resources:        # Define your AWS resources here
      ...
    Links:            # Define connections between resources
      ...
```

### DefinitionFiles section

To use the pre-defined resource definitions from the awsdac GitHub repository, specify the URL in the DefinitionFiles section:

```yaml
    DefinitionFiles:
      - Type: URL
        Url: https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml
```

If you want to customize the resource definitions locally, you can specify a local file path instead:

```yaml
    DefinitionFiles:
      - Type: LocalFile
        LocalFile: "<your definition file path (e.g. ~/Desktop/your-custom-definition.yaml)>"
```

### Resources section

`awsdac` has a unique feature where it provides resource types such as `AWS::Diagram::Canvas` and `AWS::Diagram::Cloud`, but you can define other AWS resources in a similar way to CloudFormation resource types.
And these resources are associated with each other by listing them under the "Children" property.

```yaml
    Resources:
        Canvas:
            Type: AWS::Diagram::Canvas
            Children:
                - AWSCloud
        AWSCloud:
            Type: AWS::Diagram::Cloud
            Preset: AWSCloudNoLogo
            Children:
                - VPC
        VPC:
            Type: AWS::EC2::VPC
            Children:
                - ELB            
                - Subnet1
                - Subnet2
        ELB:
            Type: AWS::ElasticLoadBalancingV2::LoadBalancer                
        Subnet1:
            Type: AWS::EC2::Subnet
            Children:
                - EC2Instance1
        Subnet2:
            Type: AWS::EC2::Subnet
            Children:
                - EC2Instance2
        EC2Instance1:
            Type: AWS::EC2::Instance
        EC2Instance2:
            Type: AWS::EC2::Instance
```

However, if you use the above file, the resulting diagram will have the ELB on the same layer as the subnets, like this:

<p align="center">
<img src="static/not-desired-image.png" width="550">
</p>

If you want to have more control over the positioning of AWS resources, such as placing the ELB in front of the subnets, you can define `AWS::Diagram::VerticalStack` or `AWS::Diagram::HorizontalStack` resources to group other resources.
For example, in this case, you could:
 * Group the 2 subnets together by creating an AWS::Diagram::HorizontalStack resource
 * Define the VPC resource to have the `AWS::Diagram::HorizontalStack` and the ELB resource as its Children
 * Set the Direction of the VPC resource to "vertical" so that the resources under the VPC are stacked vertically (note: the default value for the Direction property is "horizontal")

```diff
    Resources:
        ...
        VPC:
            Type: AWS::EC2::VPC
+           Direction: "vertical"
            Children:
                - ELB
+               - HorizontalStackGroup
-                - Subnet1
-                - Subnet2
        ...
+       HorizontalStackGroup:
+           Type: AWS::Diagram::HorizontalStack
+           Children:
+               - Subnet1
+               - Subnet2
        ...
```

<p align="center">
<img src="static/desired-image.png" width="550">
</p>

For more detailed information about "Resources" section, please refer to: [resource-types.md](resource-types.md)

### Links section

You can draw a line between the resource specified as the Source and the resource specified as the Target. At this time, you can define "where the line is drawn on the icon" by specifying SourcePosition or TargetPosition.
For example, if you make the following definition, you will get a diagram like this:

```yaml
    Links:
      - Source: ELB
        SourcePosition: S     # "S" (South) for SourcePosition means the line begins from the bottom of the icon
        Target: EC2Instance1
        TargetPosition: N     # "N" (North) for TargetPosition means the line ends at the top of the icon
        TargetArrowHead:
            Type: Open
```

<p align="center">
<img src="static/desired-image-with-link.png" width="550">
</p>

For more detailed information about Links, please refer to: [links.md](links.md)

## [Beta] CloudFormation Integration

Want to visualize existing CloudFormation templates? See the dedicated [CloudFormation Conversion Guide](cloudformation.md) for detailed instructions on both direct conversion and DAC file generation workflows.

## Customization Tips

### Changing Icons

Use the `Preset` parameter to select different icon styles. Available presets are defined in the [definition file](https://github.com/awslabs/diagram-as-code/blob/main/definitions/definition-for-aws-icons-light.yaml).

**Example**: Change S3 bucket icon to show objects

```diff
Resources:
  S3:
    Type: AWS::S3::Bucket
+   Preset: "Bucket with objects"
```

### Adding Titles

Use the `Title` parameter to add custom labels to resources.

**Example**: Add descriptive title to S3 bucket

```diff
Resources:
  S3:
    Type: AWS::S3::Bucket
    Preset: "Bucket with objects"
+   Title: "S3 (image data bucket)"
```

## Related Documentation

- **[Resource Types](resource-types.md)** - Complete list of available AWS resources and diagram elements
- **[Links](links.md)** - How to connect resources with arrows and lines
- **[Templates](template.md)** - Using Go templates for dynamic diagrams
- **[MCP Server](mcp-server.md)** - AI assistant integration
- **[CloudFormation Conversion](cloudformation.md)** [Beta] - Convert CloudFormation templates to diagrams
- **[Best Practices](best-practices.md)** - Design patterns and diagram standards
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions

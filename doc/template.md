# Templating
awsdac supports the Go standard [text/template](https://pkg.go.dev/text/template) package. To process templates in this format, enable the -t (--template) flag in awsdac command.
When this flag is enabled, awsdac first expands the strings in the file as a template, and then processes it as a YAML file. Therefore, to create a template, you need to be aware of the expanded YAML format, so it is a slightly advanced feature. However, text/template allows repetition, variable expansion, and conditional branching, which can compress lines of code and we can provide more flexibility in custom functions.
To use templating, run `awsdac -t <template>`

## Custom functions
awsdac defines the following custom functions. The custom functions are defined in [custom-func.go](https://github.com/awslabs/diagram-as-code/blob/main/internal/ctl/custom-func.go).
- `seq a` - Return an array containing the values `0` through `a-1`
- `add a b` - Returns the result of `a + b`
- `mul a b` - Returns the result of `a × b`
- `mkarr a b ...` - Returns an array of the given elements. [a, b, ...]

## Examples
### Repetition
You can execute a block repeatedly with index.
```
    {{- range $i := seq 3}}
    VPC{{$i}}:
      Type: AWS::EC2::VPC
      Title: "VPC"
      Children:
        - VPC{{$i}}AvailabilityZone0
        - VPC{{$i}}AvailabilityZone1
    {{- end}}
```

![Template repetition example](static/template-example.png)

This template will expand to the following YAML:
```yaml
    VPC0:
      Type: AWS::EC2::VPC
      Title: "VPC"
      Children:
        - VPC0AvailabilityZone0
        - VPC0AvailabilityZone1
    VPC1:
      Type: AWS::EC2::VPC
      Title: "VPC"
      Children:
        - VPC1AvailabilityZone0
        - VPC1AvailabilityZone1
    VPC2:
      Type: AWS::EC2::VPC
      Title: "VPC"
      Children:
        - VPC2AvailabilityZone0
        - VPC2AvailabilityZone1
```
if you want to use 1-indexed iteration, you can use below with add func
```
    {{- range $i := seq 3}}{{$vpc := (add $i 1)}}
    VPC{{$vpc}}:
      Type: AWS::EC2::VPC
      Title: "VPC"
      Children:
        - VPC{{$vpc}}AvailabilityZone1
        - VPC{{$vpc}}AvailabilityZone2
    {{- end}}
```

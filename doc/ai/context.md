# AI Context

## What This Project Does
`diagram-as-code` converts YAML definitions into AWS architecture diagrams.

Supported main flows:
- DAC YAML -> PNG (`awsdac input.yaml -o output.png`)
- DAC YAML -> draw.io (`awsdac input.yaml --drawio -o output.drawio`)
- CloudFormation -> PNG (`awsdac --cfn-template`)
- CloudFormation -> DAC (`awsdac --cfn-template --dac-file`)

## Core Value
The tool enables infrastructure diagrams as code: versioned, reviewable, reproducible, and automatable.

## Primary Users
- Cloud/platform engineers
- DevOps/SRE teams
- AI workflows that need diagram generation

## Key Constraints
- Layout consistency is critical.
- Backward compatibility matters for generated outputs.
- Tests should run without flaky external dependencies.

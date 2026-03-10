export const docsSectionIds = [
  'overview',
  'quick-start',
  'web-editor',
  'builder',
  'yaml-reference',
  'cli',
  'api',
  'drawio',
  'mcp',
  'local-dev',
  'troubleshooting',
  'examples',
  'architecture',
  'about',
] as const

export type DocsSectionId = (typeof docsSectionIds)[number]
export type DocsLocale = 'en' | 'pt'

export type DocsBlock =
  | { type: 'paragraph'; text: string }
  | { type: 'list'; title?: string; items: string[] }
  | { type: 'table'; title?: string; headers: string[]; rows: string[][] }
  | { type: 'code'; title?: string; lang: string; code: string }
  | { type: 'callout'; title: string; tone: 'info' | 'tip' | 'warn'; text: string }

export interface DocsSection {
  id: DocsSectionId
  category: string
  title: string
  summary: string
  keywords: string[]
  blocks: DocsBlock[]
}

export interface DocsCopy {
  ui: {
    title: string
    back: string
    contents: string
    categories: string
    searchPlaceholder: string
    searchEmpty: string
    openArticle: string
    allArticles: string
    github: string
    site: string
    repository: string
    footer: string
    previous: string
    next: string
    overview: string
  }
  hero: {
    title: string
    intro: string
    sub: string
  }
  categories: string[]
  sections: Record<DocsSectionId, DocsSection>
}

const en: DocsCopy = {
  ui: {
    title: 'Documentation',
    back: 'Back to Home',
    contents: 'Wiki Navigation',
    categories: 'Categories',
    searchPlaceholder: 'Search guides, commands, fields, or use cases…',
    searchEmpty: 'No documentation article matched the current search.',
    openArticle: 'Open article',
    allArticles: 'All articles',
    github: 'GitHub',
    site: 'Live Site',
    repository: 'Repository',
    footer: 'diagram-as-code · product wiki for the hosted editor, CLI, API, builder, and Go engine',
    previous: 'Previous',
    next: 'Next',
    overview: 'Documentation Home',
  },
  hero: {
    title: 'diagram-as-code Wiki',
    intro:
      'A complete wiki for the hosted editor, visual builder, DAC YAML format, CLI, API, PNG/PDF/draw.io export, MCP server, local development flow, and project architecture.',
    sub:
      'Use this section as the central reference for adoption, integration, troubleshooting, and maintenance. Start with Quick Start if you only need to generate a diagram now.',
  },
  categories: [
    'Getting Started',
    'Authoring',
    'Reference',
    'Integration',
    'Operations',
    'Architecture',
  ],
  sections: {
    overview: {
      id: 'overview',
      category: 'Getting Started',
      title: 'Overview',
      summary:
        'What diagram-as-code is, what it produces, and where it fits in engineering workflows.',
      keywords: ['overview', 'what is it', 'png', 'drawio', 'yaml', 'wiki', 'docs'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'diagram-as-code converts declarative YAML into AWS architecture diagrams. The same runtime model powers PNG rendering, draw.io export, the browser editor, the visual builder, and AI integrations.',
        },
        {
          type: 'paragraph',
          text:
            'The core idea is simple: architecture diagrams should be versioned, reproducible, reviewable, and generated from source instead of manually maintained in a diagramming GUI.',
        },
        {
          type: 'table',
          title: 'Primary outputs',
          headers: ['Output', 'Format', 'Main use case'],
          rows: [
            ['PNG diagram', 'image/png', 'Docs, PRs, architecture reviews, exported assets'],
            ['PDF document', 'application/pdf', 'Printable or shareable document output'],
            ['draw.io file', 'application/xml', 'Editable diagrams in diagrams.net / draw.io'],
            ['DAC YAML', 'text/yaml', 'Versioned source of the architecture definition'],
          ],
        },
        {
          type: 'list',
          title: 'Typical use cases',
          items: [
            'Document AWS architectures without manually dragging icons and arrows.',
            'Generate diagrams in CI/CD from version-controlled YAML files.',
            'Let AI assistants build diagrams through the MCP server.',
            'Use the hosted editor for fast diagram iteration with no local install.',
            'Use the builder to generate valid YAML when the team prefers forms over hand-written schema.',
          ],
        },
      ],
    },
    'quick-start': {
      id: 'quick-start',
      category: 'Getting Started',
      title: 'Quick Start',
      summary:
        'The shortest path to generate your first diagram using the website or the CLI.',
      keywords: ['quick start', 'first diagram', 'getting started', 'website', 'cli'],
      blocks: [
        {
          type: 'list',
          title: 'Fastest path',
          items: [
            'Open https://dac.moretes.com and choose Editor or Builder.',
            'Paste YAML or load one of the built-in examples.',
            'Generate PNG for visual preview, PDF for document sharing, or draw.io for editable output.',
            'If you need automation, install the CLI and run awsdac on your YAML file.',
          ],
        },
        {
          type: 'code',
          title: 'Minimal starter YAML',
          lang: 'yaml',
          code: `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - Cloud
    Cloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Children:
        - Bucket
    Bucket:
      Type: AWS::S3::Bucket`,
        },
        {
          type: 'callout',
          title: 'Recommended onboarding path',
          tone: 'tip',
          text:
            'Use the hosted editor first. Once the YAML and output are stable, move the same file into a repository and automate it with the CLI or the API.',
        },
      ],
    },
    'web-editor': {
      id: 'web-editor',
      category: 'Authoring',
      title: 'Web Editor',
      summary:
        'How the hosted editor works, when to use it, and what workflows it supports best.',
      keywords: ['editor', 'monaco', 'browser', 'preview', 'theme', 'download'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'The web editor is the fastest path from YAML to rendered diagram. It runs in the browser and calls the hosted Go backend through /api/generate.',
        },
        {
          type: 'paragraph',
          text:
            'The left panel is a Monaco editor with YAML syntax highlighting. The right panel shows the rendered output for PNG or a ready state for PDF and draw.io, depending on the selected format.',
        },
        {
          type: 'list',
          title: 'Best workflows',
          items: [
            'Fast authoring loop: write YAML, press Ctrl+Enter, inspect the output immediately.',
            'Visual tuning: adjust titles, groups, child ordering, and links until the layout stabilizes.',
            'Export loop: switch between PNG, PDF, and draw.io depending on whether the output is for preview, document delivery, or post-editing.',
          ],
        },
        {
          type: 'table',
          title: 'Editor capabilities',
          headers: ['Capability', 'Details'],
          rows: [
            ['Monaco editor', 'YAML syntax highlighting with a code-first workflow'],
            ['Examples menu', 'Built-in starter templates for common AWS architectures'],
            ['Theme and language', 'Dark/light mode and English/Portuguese interface shell'],
            ['Direct download', 'Export PNG or .drawio from the same page'],
          ],
        },
      ],
    },
    builder: {
      id: 'builder',
      category: 'Authoring',
      title: 'Visual Builder',
      summary:
        'Form-based authoring for teams that want valid DAC YAML without writing every field manually.',
      keywords: ['builder', 'forms', 'visual', 'yaml builder', 'resources', 'links'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'The Visual Builder is a form-based layer on top of the same YAML model used by the editor and CLI. It does not use a separate renderer or schema.',
        },
        {
          type: 'paragraph',
          text:
            'The builder is useful for teams that want a guided way to create resources, links, positions, labels, presets, and border children while still producing portable DAC YAML.',
        },
        {
          type: 'list',
          title: 'When the builder is a good fit',
          items: [
            'Business or platform stakeholders want to create a first draft without memorizing DAC fields.',
            'You want to teach the YAML model by showing the generated YAML next to a structured form.',
            'You need guardrails for naming, allowed resource types, or consistent shape composition.',
          ],
        },
        {
          type: 'list',
          title: 'Builder notes',
          items: [
            'Definition file can come from the official URL or a local path.',
            'Resources, links, border children, arrow heads, labels, and positions are exposed as inputs.',
            'The generated YAML can be copied or sent directly to the main editor route.',
          ],
        },
      ],
    },
    'yaml-reference': {
      id: 'yaml-reference',
      category: 'Reference',
      title: 'YAML Reference',
      summary:
        'Detailed reference for the DAC YAML structure, core fields, primitives, and best practices.',
      keywords: ['yaml', 'schema', 'reference', 'definitionfiles', 'resources', 'links'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'The DAC schema is intentionally compact. Most diagrams are defined with three core concepts: DefinitionFiles, Resources, and Links.',
        },
        {
          type: 'table',
          title: 'Core fields',
          headers: ['Field', 'Meaning'],
          rows: [
            ['Diagram.DefinitionFiles', 'Definition source for icon metadata and presets'],
            ['Resources.<name>.Type', 'AWS resource type or diagram primitive'],
            ['Resources.<name>.Children', 'Child resources placed inside a container'],
            ['Resources.<name>.Direction', 'Layout direction for children, usually vertical or horizontal'],
            ['Resources.<name>.Preset', 'Named visual preset from the definition file'],
            ['Resources.<name>.Title', 'Display label override'],
            ['Resources.<name>.BorderChildren', 'Resources attached to the border of a container'],
            ['Links[].Source / Target', 'Logical connection between resource names'],
            ['Links[].Labels.*.Title', 'Optional label placed around the edge'],
            ['Links[].TargetArrowHead.Type', 'Arrow style such as Open, Default, or none'],
          ],
        },
        {
          type: 'table',
          title: 'Common primitives',
          headers: ['Type', 'Purpose'],
          rows: [
            ['AWS::Diagram::Canvas', 'Root drawing surface'],
            ['AWS::Diagram::Cloud', 'AWS cloud grouping container'],
            ['AWS::EC2::VPC', 'VPC container with child resources'],
            ['AWS::EC2::Subnet', 'Subnet container, often with PublicSubnet or PrivateSubnet preset'],
            ['AWS::Diagram::HorizontalStack', 'Horizontal layout wrapper'],
            ['AWS::Diagram::VerticalStack', 'Vertical layout wrapper'],
            ['AWS::Diagram::Resource', 'Generic visual resource such as User or Mobile client'],
          ],
        },
        {
          type: 'list',
          title: 'Best practices',
          items: [
            'Always keep a Canvas root and attach the real top-level nodes as its children.',
            'Use stable YAML keys such as ALB, PublicSubnetA, OrdersService, and Redis.',
            'Use groups like Cloud, VPC, and Subnet to make boundaries explicit.',
            'Use Title only when the displayed label should differ from the YAML key.',
            'If a diagram has too many crossing links, split it into multiple focused diagrams.',
          ],
        },
        {
          type: 'code',
          title: 'Reference example',
          lang: 'yaml',
          code: `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"

  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children: [Cloud, User]

    Cloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Children: [VPC]

    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children: [ALB, AppTier]
      BorderChildren:
        - Position: S
          Resource: IGW

    AppTier:
      Type: AWS::Diagram::HorizontalStack
      Children: [SubnetA, SubnetB]

    SubnetA:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children: [InstanceA]

    SubnetB:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children: [InstanceB]

    InstanceA:
      Type: AWS::EC2::Instance

    InstanceB:
      Type: AWS::EC2::Instance

    ALB:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer

    IGW:
      Type: AWS::EC2::InternetGateway

    User:
      Type: AWS::Diagram::Resource
      Preset: User

  Links:
    - Source: User
      Target: IGW
      TargetArrowHead:
        Type: Open
    - Source: IGW
      Target: ALB
      TargetArrowHead:
        Type: Open
    - Source: ALB
      Target: InstanceA
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "HTTP:80"
    - Source: ALB
      Target: InstanceB
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "HTTP:80"`,
        },
      ],
    },
    cli: {
      id: 'cli',
      category: 'Reference',
      title: 'CLI Reference',
      summary:
        'Installation, commands, flags, and usage scenarios for the awsdac command line interface.',
      keywords: ['cli', 'awsdac', 'flags', 'install', 'command line', 'automation'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'The CLI is the best fit when diagrams should live with infrastructure code, run in CI/CD, or be generated in automation pipelines.',
        },
        {
          type: 'code',
          title: 'Install',
          lang: 'bash',
          code: `# macOS
brew install awsdac

# Go install
go install github.com/fernandofatech/diagram-as-code/cmd/awsdac@latest`,
        },
        {
          type: 'code',
          title: 'Basic usage',
          lang: 'bash',
          code: `# Generate PNG
awsdac examples/alb-ec2.yaml

# Custom output file
awsdac examples/alb-ec2.yaml -o my-diagram.png

# Generate draw.io XML
awsdac examples/alb-ec2.yaml --drawio -o output.drawio

# Output extension also selects format
awsdac examples/alb-ec2.yaml -o output.drawio`,
        },
        {
          type: 'table',
          title: 'Main flags',
          headers: ['Flag', 'Purpose'],
          rows: [
            ['-o, --output', 'Output file name'],
            ['--drawio', 'Generate draw.io instead of PNG'],
            ['-f, --force', 'Overwrite output without confirmation'],
            ['--width / --height', 'Resize PNG output'],
            ['-t, --template', 'Render the input as Go text/template'],
            ['-c, --cfn-template', 'Create diagram from CloudFormation'],
            ['-d, --dac-file', 'Generate DAC YAML from CloudFormation'],
            ['--allow-untrusted-definitions', 'Allow non-official definition sources'],
            ['-v, --verbose', 'Verbose logging'],
          ],
        },
        {
          type: 'list',
          title: 'CLI scenarios',
          items: [
            'Local engineering documentation stored next to IaC.',
            'CI/CD generation of diagrams for pull requests or release artifacts.',
            'Batch rendering multiple YAML files as part of docs automation.',
          ],
        },
      ],
    },
    api: {
      id: 'api',
      category: 'Integration',
      title: 'API Reference',
      summary:
        'Hosted HTTP API for browser, server, automation, and AI-agent integrations.',
      keywords: ['api', 'http', 'generate', 'fetch', 'curl', 'endpoint', 'integration'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'The hosted API powers the browser editor and can also be called directly from scripts, internal portals, CI jobs, or AI agents.',
        },
        {
          type: 'code',
          title: 'Request format',
          lang: 'http',
          code: `POST https://dac.moretes.com/api/generate
Content-Type: application/json

{
  "yaml": "Diagram:\\n  ...",
  "format": "pdf"
}`,
        },
        {
          type: 'table',
          title: 'Response types',
          headers: ['Mode', 'Content-Type', 'Body'],
          rows: [
            ['png', 'image/png', 'Raw PNG bytes'],
            ['pdf', 'application/pdf', 'Single-page PDF document'],
            ['drawio', 'application/xml', 'mxGraphModel XML'],
            ['error', 'application/json', '{"error":"message"}'],
          ],
        },
        {
          type: 'code',
          title: 'curl example',
          lang: 'bash',
          code: `curl -X POST https://dac.moretes.com/api/generate \\
  -H "Content-Type: application/json" \\
  -d '{"yaml":"Diagram:\\n  DefinitionFiles:\\n    - Type: URL\\n      Url: \\"https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml\\"\\n  Resources:\\n    Canvas:\\n      Type: AWS::Diagram::Canvas\\n      Children:\\n        - Bucket\\n    Bucket:\\n      Type: AWS::S3::Bucket","format":"pdf"}' \\
  --output diagram.pdf`,
        },
        {
          type: 'code',
          title: 'TypeScript example',
          lang: 'ts',
          code: `const yaml = \`Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children: [Bucket]
    Bucket:
      Type: AWS::S3::Bucket\`

const response = await fetch('https://dac.moretes.com/api/generate', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ yaml, format: 'drawio' }),
})

if (!response.ok) {
  const { error } = await response.json()
  throw new Error(error)
}

const xml = await response.text()`,
        },
        {
          type: 'list',
          title: 'Integration notes',
          items: [
            'The request body expects a DAC YAML string, not a file upload.',
            'Format defaults to png when omitted. Use pdf for document output or drawio for editable output.',
            'When generation fails, the endpoint returns a JSON error payload.',
            'The endpoint is suitable for browser, server-to-server, and agent workflows.',
          ],
        },
      ],
    },
    drawio: {
      id: 'drawio',
      category: 'Reference',
      title: 'draw.io Export',
      summary:
        'How draw.io export works internally and when to choose it over PNG output.',
      keywords: ['drawio', 'diagrams.net', 'mxgraphmodel', 'xml', 'editable'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'draw.io export is not a screenshot or a separate layout path. It uses the same resource model and geometry used for PNG rendering.',
        },
        {
          type: 'paragraph',
          text:
            'This keeps PNG and draw.io outputs aligned while preserving editable structure instead of flattening the result into pixels.',
        },
        {
          type: 'list',
          title: 'Export pipeline',
          items: [
            'Parse DAC YAML into runtime resources and links.',
            'Run the same layout engine used for PNG output, including Scale and ZeroAdjust.',
            'Reorder children from link topology to preserve consistent placement.',
            'Export resources as mxCell nodes and links as mxCell edges.',
            'Embed AWS SVG icons as data URIs so the file is self-contained.',
            'Write mxGraphModel XML compatible with diagrams.net / draw.io.',
          ],
        },
        {
          type: 'list',
          title: 'Best use cases',
          items: [
            'Start with code, finish with manual annotations in draw.io.',
            'Generate a baseline diagram for non-technical teams to edit later.',
            'Keep editable diagrams in repositories while still using code for the initial layout.',
          ],
        },
      ],
    },
    mcp: {
      id: 'mcp',
      category: 'Integration',
      title: 'MCP Server',
      summary:
        'Use the Model Context Protocol server to let AI assistants generate diagrams from chat.',
      keywords: ['mcp', 'claude', 'ai', 'agent', 'tools', 'server'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'The MCP server exposes diagram generation tools to AI assistants. This makes the project useful not only as a CLI or website, but also as part of agent workflows.',
        },
        {
          type: 'paragraph',
          text:
            'Typical usage is: an assistant writes DAC YAML from a human request, invokes the MCP tool, and saves the result as PNG or draw.io.',
        },
        {
          type: 'code',
          title: 'Install MCP server',
          lang: 'bash',
          code: `go install github.com/awslabs/diagram-as-code/cmd/awsdac-mcp-server@latest`,
        },
        {
          type: 'code',
          title: 'Claude Desktop configuration',
          lang: 'json',
          code: `{
  "mcpServers": {
    "diagram-as-code": {
      "command": "awsdac-mcp-server",
      "args": []
    }
  }
}`,
        },
        {
          type: 'table',
          title: 'Available tools',
          headers: ['Tool', 'Description'],
          rows: [
            ['create_diagram', 'Generate PNG from DAC YAML'],
            ['create_drawio', 'Generate draw.io XML from DAC YAML'],
          ],
        },
        {
          type: 'code',
          title: 'Example prompt',
          lang: 'text',
          code: `Create an AWS architecture diagram with:
- one VPC
- two public subnets
- an internet-facing ALB
- two EC2 instances
- an Internet Gateway

Return PNG output and save it to ~/Desktop/architecture.png`,
        },
      ],
    },
    'local-dev': {
      id: 'local-dev',
      category: 'Operations',
      title: 'Local Development',
      summary:
        'Run the Go API, Next.js frontend, tests, and local Vercel flow for development.',
      keywords: ['local dev', 'go run', 'next dev', 'vercel dev', 'tests'],
      blocks: [
        {
          type: 'table',
          title: 'Prerequisites',
          headers: ['Tool', 'Version', 'Purpose'],
          rows: [
            ['Go', '1.21+', 'Go backend, CLI, and tests'],
            ['Node.js', '18+ (Node 20+ recommended for Next 16)', 'Web frontend'],
            ['npm', '9+', 'Frontend dependencies and scripts'],
          ],
        },
        {
          type: 'code',
          title: 'Run locally',
          lang: 'bash',
          code: `git clone https://github.com/fernandofatech/diagram-as-code.git
cd diagram-as-code

# Terminal 1 - local Go API
go run ./cmd/api-dev

# Terminal 2 - Next.js app
cd web
npm install
npm run dev`,
        },
        {
          type: 'code',
          title: 'Validation commands',
          lang: 'bash',
          code: `# Go validation
go test ./...

# Frontend validation
cd web
npm run lint
npm run build`,
        },
        {
          type: 'code',
          title: 'Alternative: Vercel dev',
          lang: 'bash',
          code: `npm install --global vercel
vercel dev`,
        },
        {
          type: 'list',
          title: 'Development notes',
          items: [
            'When running the web app locally, /api/* is proxied to the Go API in development.',
            'The hosted frontend lives in web/, while the serverless endpoint lives in api/generate.go.',
            'The root vercel.json controls how Next.js and Go are deployed together.',
          ],
        },
      ],
    },
    troubleshooting: {
      id: 'troubleshooting',
      category: 'Operations',
      title: 'Troubleshooting',
      summary:
        'Common failure modes, debugging patterns, and practical steps to recover fast.',
      keywords: ['troubleshooting', 'errors', 'preview', 'icons', 'vercel', 'theme'],
      blocks: [
        {
          type: 'callout',
          title: 'Recommended debugging strategy',
          tone: 'warn',
          text:
            'Reduce the YAML to the smallest working version first. Then add groups, links, labels, and presets back incrementally until the faulty block is isolated.',
        },
        {
          type: 'list',
          title: 'Common issues',
          items: [
            'Blank preview or generation error: validate the YAML structure first. A missing Canvas, bad Type, or malformed Links block is the most common cause.',
            'Incorrect resource icon: verify that the definition file URL is valid and points to a compatible AWS definition set.',
            'Unexpected light/dark appearance: clear localStorage theme settings or use the theme switcher to resync the UI state.',
            'draw.io output opens but looks different: check zoom, previous manual edits, and whether the source YAML changed.',
            'CLI works but local web app fails: make sure the Go API is running on port 8080 before starting the frontend, or use vercel dev.',
            'Vercel deploy ignores dashboard build settings: this repository uses builds in vercel.json, so file-based config takes precedence.',
          ],
        },
      ],
    },
    examples: {
      id: 'examples',
      category: 'Reference',
      title: 'Examples and Use Cases',
      summary:
        'Representative DAC patterns for serverless, networking, application, and AI-driven workflows.',
      keywords: ['examples', 'patterns', 'serverless', 'network', 'sample yaml'],
      blocks: [
        {
          type: 'table',
          title: 'Common patterns',
          headers: ['Pattern', 'Best fit'],
          rows: [
            ['Single resource', 'S3 bucket, Lambda, or DynamoDB reference diagram'],
            ['Network topology', 'VPC, subnets, ALB, NAT, IGW, and service placement'],
            ['Application platform', 'ECS or Lambda services with data and messaging layers'],
            ['Multi-region', 'Replicated services or failover topology'],
            ['AI-generated diagram', 'Human prompt -> MCP -> DAC YAML -> PNG or draw.io'],
          ],
        },
        {
          type: 'code',
          title: 'Serverless worker pipeline',
          lang: 'yaml',
          code: `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children: [Cloud]
    Cloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Children: [App, Queue, Worker]
    App:
      Type: AWS::Lambda::Function
      Title: API Lambda
    Queue:
      Type: AWS::SQS::Queue
    Worker:
      Type: AWS::Lambda::Function
      Title: Worker Lambda
  Links:
    - Source: App
      Target: Queue
      TargetArrowHead:
        Type: Open
    - Source: Queue
      Target: Worker
      TargetArrowHead:
        Type: Open`,
        },
        {
          type: 'code',
          title: 'Frontend + API + data stack',
          lang: 'yaml',
          code: `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children: [Users, Cloud]
    Users:
      Type: AWS::Diagram::HorizontalStack
      Children: [Web, Mobile]
    Web:
      Type: AWS::Diagram::Resource
      Preset: User
      Title: Web User
    Mobile:
      Type: AWS::Diagram::Resource
      Preset: Mobile client
    Cloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Direction: vertical
      Children: [Api, Data]
    Api:
      Type: AWS::Diagram::HorizontalStack
      Children: [Gateway, Service]
    Gateway:
      Type: AWS::ApiGateway
      Title: Public API
    Service:
      Type: AWS::Lambda::Function
      Title: Orders Service
    Data:
      Type: AWS::Diagram::HorizontalStack
      Children: [Redis, Database]
    Redis:
      Type: AWS::ElastiCache::CacheCluster
    Database:
      Type: AWS::RDS::DBCluster
  Links:
    - Source: Web
      Target: Gateway
      TargetArrowHead:
        Type: Open
    - Source: Mobile
      Target: Gateway
      TargetArrowHead:
        Type: Open
    - Source: Gateway
      Target: Service
      TargetArrowHead:
        Type: Open
    - Source: Service
      Target: Redis
      TargetArrowHead:
        Type: Open
    - Source: Service
      Target: Database
      TargetArrowHead:
        Type: Open`,
        },
      ],
    },
    architecture: {
      id: 'architecture',
      category: 'Architecture',
      title: 'Project Architecture',
      summary:
        'High-level repository layers, execution flow, and deployment structure.',
      keywords: ['architecture', 'repo structure', 'internal/ctl', 'web', 'api', 'go'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'The repository is split between a Go engine and a Next.js frontend. The Go side owns parsing, layout, definitions, and export. The web side owns editing, preview interaction, theming, docs, and the builder UX.',
        },
        {
          type: 'paragraph',
          text:
            'The API route /api/generate is deployed as a Go serverless function. The frontend is deployed from web/ as a Next.js app, while vercel.json wires both runtimes together.',
        },
        {
          type: 'table',
          title: 'Repository layers',
          headers: ['Path', 'Responsibility'],
          rows: [
            ['cmd/', 'CLI entrypoints and local API dev server'],
            ['internal/ctl/', 'Main orchestration pipeline'],
            ['internal/types/', 'Runtime model for resources, links, and geometry'],
            ['internal/definition/', 'Definition loading and icon metadata'],
            ['pkg/diagram/', 'Public Go wrapper for embedding'],
            ['api/generate.go', 'Vercel serverless handler'],
            ['web/', 'Next.js frontend, builder, docs, editor, preview'],
          ],
        },
        {
          type: 'list',
          title: 'Execution flow',
          items: [
            'Input YAML is parsed.',
            'Definition files are loaded and validated.',
            'Resources and links are constructed into the runtime graph.',
            'Layout is computed deterministically where possible.',
            'Output is exported as PNG, PDF, or draw.io XML.',
          ],
        },
      ],
    },
    about: {
      id: 'about',
      category: 'Architecture',
      title: 'About the Author',
      summary:
        'Context about the maintainer of this fork and the goals behind the hosted platform.',
      keywords: ['author', 'fernando azevedo', 'maintainer', 'fork'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'Fernando Azevedo maintains this fork and expanded the original project with a hosted web frontend, draw.io export flow, public API endpoint, visual builder, and AI-friendly integrations.',
        },
        {
          type: 'paragraph',
          text:
            'The focus of this fork is practical architecture documentation: clear outputs, automation-friendly workflows, and a better bridge between code, diagrams, and AI tooling.',
        },
        {
          type: 'list',
          title: 'Profiles',
          items: [
            'Website: https://fernando.moretes.com',
            'GitHub: https://github.com/fernandofatech',
            'LinkedIn: https://www.linkedin.com/in/fernando-francisco-azevedo/',
          ],
        },
      ],
    },
  },
}

const pt: DocsCopy = {
  ui: {
    title: 'Documentação',
    back: 'Voltar ao Início',
    contents: 'Navegação da Wiki',
    categories: 'Categorias',
    searchPlaceholder: 'Buscar guias, comandos, campos ou casos de uso…',
    searchEmpty: 'Nenhum artigo de documentação corresponde à busca atual.',
    openArticle: 'Abrir artigo',
    allArticles: 'Todos os artigos',
    github: 'GitHub',
    site: 'Site',
    repository: 'Repositório',
    footer: 'diagram-as-code · wiki do produto para editor hospedado, CLI, API, builder e engine Go',
    previous: 'Anterior',
    next: 'Próximo',
    overview: 'Início da Documentação',
  },
  hero: {
    title: 'Wiki do diagram-as-code',
    intro:
      'Uma wiki completa para o editor hospedado, builder visual, formato DAC YAML, CLI, API, exportação PNG/PDF/draw.io, servidor MCP, fluxo de desenvolvimento local e arquitetura do projeto.',
    sub:
      'Use esta área como referência central para adoção, integração, troubleshooting e manutenção. Comece em Começo Rápido se você só precisa gerar um diagrama agora.',
  },
  categories: [
    'Getting Started',
    'Authoring',
    'Reference',
    'Integration',
    'Operations',
    'Architecture',
  ],
  sections: {
    overview: {
      id: 'overview',
      category: 'Getting Started',
      title: 'Visão Geral',
      summary:
        'O que é o diagram-as-code, o que ele gera e onde ele se encaixa em fluxos de engenharia.',
      keywords: ['visao geral', 'png', 'drawio', 'yaml', 'wiki', 'docs'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'diagram-as-code converte YAML declarativo em diagramas de arquitetura AWS. O mesmo modelo de runtime alimenta a renderização PNG, a exportação draw.io, o editor no navegador, o builder visual e integrações com IA.',
        },
        {
          type: 'paragraph',
          text:
            'A ideia central é simples: diagramas de arquitetura devem ser versionados, reproduzíveis, revisáveis e gerados a partir do código-fonte, não mantidos manualmente em uma GUI.',
        },
        {
          type: 'table',
          title: 'Saídas principais',
          headers: ['Saída', 'Formato', 'Caso de uso principal'],
          rows: [
            ['Diagrama PNG', 'image/png', 'Docs, PRs, revisões de arquitetura, assets exportados'],
            ['Documento PDF', 'application/pdf', 'Saída pronta para compartilhar ou imprimir'],
            ['Arquivo draw.io', 'application/xml', 'Diagramas editáveis no diagrams.net / draw.io'],
            ['YAML DAC', 'text/yaml', 'Fonte versionada da definição de arquitetura'],
          ],
        },
        {
          type: 'list',
          title: 'Casos de uso típicos',
          items: [
            'Documentar arquiteturas AWS sem arrastar ícones e conexões manualmente.',
            'Gerar diagramas em CI/CD a partir de arquivos YAML versionados.',
            'Permitir que assistentes de IA montem diagramas via servidor MCP.',
            'Usar o editor hospedado para iterar rapidamente sem instalação local.',
            'Usar o builder quando a equipe prefere formulários em vez de schema escrito à mão.',
          ],
        },
      ],
    },
    'quick-start': {
      id: 'quick-start',
      category: 'Getting Started',
      title: 'Começo Rápido',
      summary:
        'O caminho mais curto para gerar seu primeiro diagrama usando o site ou o CLI.',
      keywords: ['comeco rapido', 'primeiro diagrama', 'website', 'cli'],
      blocks: [
        {
          type: 'list',
          title: 'Caminho mais rápido',
          items: [
            'Abra https://dac.moretes.com e escolha Editor ou Builder.',
            'Cole YAML ou carregue um dos exemplos prontos.',
            'Gere PNG para preview visual, PDF para compartilhamento e draw.io para edição posterior.',
            'Se precisar de automação, instale o CLI e execute awsdac no seu arquivo YAML.',
          ],
        },
        {
          type: 'code',
          title: 'YAML inicial mínimo',
          lang: 'yaml',
          code: `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - Cloud
    Cloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Children:
        - Bucket
    Bucket:
      Type: AWS::S3::Bucket`,
        },
        {
          type: 'callout',
          title: 'Onboarding recomendado',
          tone: 'tip',
          text:
            'Use primeiro o editor hospedado. Depois que o YAML e a saída estiverem estáveis, mova o mesmo arquivo para um repositório e automatize com o CLI ou a API.',
        },
      ],
    },
    'web-editor': {
      id: 'web-editor',
      category: 'Authoring',
      title: 'Editor Web',
      summary:
        'Como o editor hospedado funciona, quando usar e quais fluxos ele atende melhor.',
      keywords: ['editor', 'monaco', 'browser', 'preview', 'tema', 'download'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'O editor web é o caminho mais rápido entre YAML e diagrama renderizado. Ele roda no navegador e chama o backend Go hospedado através de /api/generate.',
        },
        {
          type: 'paragraph',
          text:
            'O painel esquerdo é um editor Monaco com destaque de sintaxe YAML. O painel direito mostra a saída renderizada para PNG ou o estado pronto para PDF e draw.io.',
        },
        {
          type: 'list',
          title: 'Melhores fluxos',
          items: [
            'Loop rápido de autoria: escreva YAML, pressione Ctrl+Enter e veja a saída imediatamente.',
            'Ajuste visual: refine títulos, grupos, ordem dos children e links até estabilizar o layout.',
            'Loop de exportação: alterne entre PNG, PDF e draw.io dependendo do destino final.',
          ],
        },
        {
          type: 'table',
          title: 'Capacidades do editor',
          headers: ['Capacidade', 'Detalhes'],
          rows: [
            ['Editor Monaco', 'Destaque de sintaxe YAML com fluxo code-first'],
            ['Menu de exemplos', 'Templates prontos para arquiteturas AWS comuns'],
            ['Tema e idioma', 'Modo claro/escuro e shell em inglês/português'],
            ['Download direto', 'Exporta PNG ou .drawio da mesma página'],
          ],
        },
      ],
    },
    builder: {
      id: 'builder',
      category: 'Authoring',
      title: 'Builder Visual',
      summary:
        'Autoria baseada em formulários para equipes que querem YAML DAC válido sem escrever todos os campos manualmente.',
      keywords: ['builder', 'forms', 'visual', 'yaml builder', 'resources', 'links'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'O Builder Visual é uma camada baseada em formulários sobre o mesmo modelo YAML usado no editor e no CLI. Ele não usa um renderer separado.',
        },
        {
          type: 'paragraph',
          text:
            'O builder é útil para equipes que querem uma maneira guiada de criar recursos, links, posições, labels, presets e border children enquanto ainda produzem DAC YAML portável.',
        },
        {
          type: 'list',
          title: 'Quando o builder faz sentido',
          items: [
            'Stakeholders de negócio ou plataforma querem montar um primeiro rascunho sem memorizar os campos DAC.',
            'Você quer ensinar o modelo YAML mostrando o arquivo sendo gerado em paralelo.',
            'Você precisa de guardrails para nomes, tipos permitidos ou composição consistente dos diagramas.',
          ],
        },
        {
          type: 'list',
          title: 'Notas do builder',
          items: [
            'O arquivo de definição pode vir da URL oficial ou de um caminho local.',
            'Recursos, links, border children, pontas de seta, labels e posições estão expostos em inputs.',
            'O YAML gerado pode ser copiado ou enviado diretamente para o editor principal.',
          ],
        },
      ],
    },
    'yaml-reference': {
      id: 'yaml-reference',
      category: 'Reference',
      title: 'Referência YAML',
      summary:
        'Referência detalhada da estrutura DAC YAML, campos principais, primitives e boas práticas.',
      keywords: ['yaml', 'schema', 'reference', 'definitionfiles', 'resources', 'links'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'O schema DAC é intencionalmente compacto. A maior parte dos diagramas é definida com três conceitos centrais: DefinitionFiles, Resources e Links.',
        },
        {
          type: 'table',
          title: 'Campos principais',
          headers: ['Campo', 'Significado'],
          rows: [
            ['Diagram.DefinitionFiles', 'Fonte de definição para metadados de ícones e presets'],
            ['Resources.<name>.Type', 'Tipo AWS ou primitive do diagrama'],
            ['Resources.<name>.Children', 'Recursos filhos colocados dentro de um container'],
            ['Resources.<name>.Direction', 'Direção do layout dos filhos, normalmente vertical ou horizontal'],
            ['Resources.<name>.Preset', 'Preset visual nomeado vindo do arquivo de definição'],
            ['Resources.<name>.Title', 'Sobrescrita do label exibido'],
            ['Resources.<name>.BorderChildren', 'Recursos anexados à borda de um container'],
            ['Links[].Source / Target', 'Conexão lógica entre nomes de recursos'],
            ['Links[].Labels.*.Title', 'Label opcional ao redor da aresta'],
            ['Links[].TargetArrowHead.Type', 'Estilo de seta como Open, Default ou none'],
          ],
        },
        {
          type: 'table',
          title: 'Primitives comuns',
          headers: ['Tipo', 'Propósito'],
          rows: [
            ['AWS::Diagram::Canvas', 'Superfície raiz do desenho'],
            ['AWS::Diagram::Cloud', 'Container de agrupamento AWS Cloud'],
            ['AWS::EC2::VPC', 'Container de VPC com recursos filhos'],
            ['AWS::EC2::Subnet', 'Container de subnet, normalmente com preset PublicSubnet ou PrivateSubnet'],
            ['AWS::Diagram::HorizontalStack', 'Wrapper de layout horizontal'],
            ['AWS::Diagram::VerticalStack', 'Wrapper de layout vertical'],
            ['AWS::Diagram::Resource', 'Recurso visual genérico como User ou Mobile client'],
          ],
        },
        {
          type: 'list',
          title: 'Boas práticas',
          items: [
            'Sempre mantenha um Canvas raiz e conecte os nós de topo reais como children dele.',
            'Use chaves YAML estáveis como ALB, PublicSubnetA, OrdersService e Redis.',
            'Use grupos como Cloud, VPC e Subnet para tornar fronteiras explícitas.',
            'Use Title apenas quando o label exibido precisa diferir da chave YAML.',
            'Se um diagrama tiver muitos links cruzando, divida em múltiplos diagramas focados.',
          ],
        },
        {
          type: 'code',
          title: 'Exemplo de referência',
          lang: 'yaml',
          code: `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"

  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children: [Cloud, User]

    Cloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Children: [VPC]

    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children: [ALB, AppTier]
      BorderChildren:
        - Position: S
          Resource: IGW

    AppTier:
      Type: AWS::Diagram::HorizontalStack
      Children: [SubnetA, SubnetB]

    SubnetA:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children: [InstanceA]

    SubnetB:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children: [InstanceB]

    InstanceA:
      Type: AWS::EC2::Instance

    InstanceB:
      Type: AWS::EC2::Instance

    ALB:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer

    IGW:
      Type: AWS::EC2::InternetGateway

    User:
      Type: AWS::Diagram::Resource
      Preset: User

  Links:
    - Source: User
      Target: IGW
      TargetArrowHead:
        Type: Open
    - Source: IGW
      Target: ALB
      TargetArrowHead:
        Type: Open
    - Source: ALB
      Target: InstanceA
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "HTTP:80"
    - Source: ALB
      Target: InstanceB
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "HTTP:80"`,
        },
      ],
    },
    cli: {
      id: 'cli',
      category: 'Reference',
      title: 'Referência CLI',
      summary:
        'Instalação, comandos, flags e cenários de uso do comando awsdac.',
      keywords: ['cli', 'awsdac', 'flags', 'install', 'automacao'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'O CLI é a melhor opção quando diagramas precisam viver com infraestrutura como código, rodar em CI/CD ou ser gerados em pipelines automatizados.',
        },
        {
          type: 'code',
          title: 'Instalação',
          lang: 'bash',
          code: `# macOS
brew install awsdac

# Go install
go install github.com/fernandofatech/diagram-as-code/cmd/awsdac@latest`,
        },
        {
          type: 'code',
          title: 'Uso básico',
          lang: 'bash',
          code: `# Gerar PNG
awsdac examples/alb-ec2.yaml

# Arquivo de saída customizado
awsdac examples/alb-ec2.yaml -o my-diagram.png

# Gerar XML draw.io
awsdac examples/alb-ec2.yaml --drawio -o output.drawio

# A extensão da saída também seleciona o formato
awsdac examples/alb-ec2.yaml -o output.drawio`,
        },
        {
          type: 'table',
          title: 'Principais flags',
          headers: ['Flag', 'Propósito'],
          rows: [
            ['-o, --output', 'Nome do arquivo de saída'],
            ['--drawio', 'Gera draw.io em vez de PNG'],
            ['-f, --force', 'Sobrescreve a saída sem confirmação'],
            ['--width / --height', 'Redimensiona a saída PNG'],
            ['-t, --template', 'Renderiza a entrada como Go text/template'],
            ['-c, --cfn-template', 'Cria diagrama a partir de CloudFormation'],
            ['-d, --dac-file', 'Gera YAML DAC a partir de CloudFormation'],
            ['--allow-untrusted-definitions', 'Permite fontes de definição não-oficiais'],
            ['-v, --verbose', 'Log detalhado'],
          ],
        },
        {
          type: 'list',
          title: 'Cenários do CLI',
          items: [
            'Documentação local de engenharia junto com IaC.',
            'Geração em CI/CD de diagramas para pull requests ou artefatos de release.',
            'Renderização em lote de múltiplos arquivos YAML como parte de automação de documentação.',
          ],
        },
      ],
    },
    api: {
      id: 'api',
      category: 'Integration',
      title: 'Referência API',
      summary:
        'API HTTP hospedada para integrações browser, server-side, automação e agentes de IA.',
      keywords: ['api', 'http', 'generate', 'fetch', 'curl', 'endpoint'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'A API hospedada alimenta o editor no navegador e também pode ser chamada diretamente de scripts, portais internos, jobs de CI ou agentes de IA.',
        },
        {
          type: 'code',
          title: 'Formato da requisição',
          lang: 'http',
          code: `POST https://dac.moretes.com/api/generate
Content-Type: application/json

{
  "yaml": "Diagram:\\n  ...",
  "format": "pdf"
}`,
        },
        {
          type: 'table',
          title: 'Tipos de resposta',
          headers: ['Modo', 'Content-Type', 'Body'],
          rows: [
            ['png', 'image/png', 'Bytes PNG brutos'],
            ['pdf', 'application/pdf', 'Documento PDF de uma página'],
            ['drawio', 'application/xml', 'XML mxGraphModel'],
            ['error', 'application/json', '{"error":"message"}'],
          ],
        },
        {
          type: 'code',
          title: 'Exemplo com curl',
          lang: 'bash',
          code: `curl -X POST https://dac.moretes.com/api/generate \\
  -H "Content-Type: application/json" \\
  -d '{"yaml":"Diagram:\\n  DefinitionFiles:\\n    - Type: URL\\n      Url: \\"https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml\\"\\n  Resources:\\n    Canvas:\\n      Type: AWS::Diagram::Canvas\\n      Children:\\n        - Bucket\\n    Bucket:\\n      Type: AWS::S3::Bucket","format":"pdf"}' \\
  --output diagram.pdf`,
        },
        {
          type: 'code',
          title: 'Exemplo em TypeScript',
          lang: 'ts',
          code: `const yaml = \`Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children: [Bucket]
    Bucket:
      Type: AWS::S3::Bucket\`

const response = await fetch('https://dac.moretes.com/api/generate', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ yaml, format: 'drawio' }),
})

if (!response.ok) {
  const { error } = await response.json()
  throw new Error(error)
}

const xml = await response.text()`,
        },
        {
          type: 'list',
          title: 'Notas de integração',
          items: [
            'O body espera uma string YAML DAC, não upload de arquivo.',
            'Quando format é omitido, o padrão é png. Use pdf para saída em documento ou drawio para edição posterior.',
            'Quando a geração falha, o endpoint retorna um payload JSON de erro.',
            'O endpoint funciona para browser, server-to-server e fluxos com agentes.',
          ],
        },
      ],
    },
    drawio: {
      id: 'drawio',
      category: 'Reference',
      title: 'Exportação draw.io',
      summary:
        'Como a exportação draw.io funciona internamente e quando escolher esse formato em vez de PNG.',
      keywords: ['drawio', 'diagrams.net', 'mxgraphmodel', 'xml', 'editable'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'A exportação draw.io não é screenshot nem um caminho de layout separado. Ela usa o mesmo modelo de recursos e a mesma geometria usada para PNG.',
        },
        {
          type: 'paragraph',
          text:
            'Isso mantém a saída PNG e draw.io alinhadas enquanto preserva a estrutura editável em vez de achatar o resultado em pixels.',
        },
        {
          type: 'list',
          title: 'Pipeline de exportação',
          items: [
            'Parse do YAML DAC para recursos e links de runtime.',
            'Execução do mesmo layout usado no PNG, incluindo Scale e ZeroAdjust.',
            'Reordenação dos children a partir da topologia de links para manter consistência visual.',
            'Exportação dos recursos como nós mxCell e dos links como arestas mxCell.',
            'Embute SVGs AWS como data URIs para que o arquivo seja autocontido.',
            'Grava XML mxGraphModel compatível com diagrams.net / draw.io.',
          ],
        },
        {
          type: 'list',
          title: 'Melhores usos',
          items: [
            'Começar com código e finalizar com anotações manuais no draw.io.',
            'Gerar um baseline para times não técnicos editarem depois.',
            'Manter diagramas editáveis em repositórios enquanto o layout inicial continua vindo do código.',
          ],
        },
      ],
    },
    mcp: {
      id: 'mcp',
      category: 'Integration',
      title: 'Servidor MCP',
      summary:
        'Use o servidor Model Context Protocol para permitir que assistentes de IA gerem diagramas por chat.',
      keywords: ['mcp', 'claude', 'ai', 'agent', 'tools', 'server'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'O servidor MCP expõe ferramentas de geração de diagrama para assistentes de IA. Isso torna o projeto útil não só como CLI ou website, mas também em fluxos baseados em agentes.',
        },
        {
          type: 'paragraph',
          text:
            'O uso típico é: o assistente escreve DAC YAML a partir do pedido humano, invoca a ferramenta MCP e salva o resultado como PNG ou draw.io.',
        },
        {
          type: 'code',
          title: 'Instalar servidor MCP',
          lang: 'bash',
          code: `go install github.com/awslabs/diagram-as-code/cmd/awsdac-mcp-server@latest`,
        },
        {
          type: 'code',
          title: 'Configuração do Claude Desktop',
          lang: 'json',
          code: `{
  "mcpServers": {
    "diagram-as-code": {
      "command": "awsdac-mcp-server",
      "args": []
    }
  }
}`,
        },
        {
          type: 'table',
          title: 'Ferramentas disponíveis',
          headers: ['Ferramenta', 'Descrição'],
          rows: [
            ['create_diagram', 'Gera PNG a partir de DAC YAML'],
            ['create_drawio', 'Gera XML draw.io a partir de DAC YAML'],
          ],
        },
        {
          type: 'code',
          title: 'Prompt de exemplo',
          lang: 'text',
          code: `Crie um diagrama de arquitetura AWS com:
- uma VPC
- duas subnets públicas
- um ALB internet-facing
- duas instâncias EC2
- um Internet Gateway

Retorne saída PNG e salve em ~/Desktop/architecture.png`,
        },
      ],
    },
    'local-dev': {
      id: 'local-dev',
      category: 'Operations',
      title: 'Desenvolvimento Local',
      summary:
        'Rode a API Go, o frontend Next.js, testes e o fluxo local da Vercel para desenvolvimento.',
      keywords: ['local dev', 'go run', 'next dev', 'vercel dev', 'tests'],
      blocks: [
        {
          type: 'table',
          title: 'Pré-requisitos',
          headers: ['Ferramenta', 'Versão', 'Propósito'],
          rows: [
            ['Go', '1.21+', 'Backend Go, CLI e testes'],
            ['Node.js', '18+ (Node 20+ recomendado para Next 16)', 'Frontend web'],
            ['npm', '9+', 'Dependências e scripts do frontend'],
          ],
        },
        {
          type: 'code',
          title: 'Rodar localmente',
          lang: 'bash',
          code: `git clone https://github.com/fernandofatech/diagram-as-code.git
cd diagram-as-code

# Terminal 1 - API Go local
go run ./cmd/api-dev

# Terminal 2 - app Next.js
cd web
npm install
npm run dev`,
        },
        {
          type: 'code',
          title: 'Comandos de validação',
          lang: 'bash',
          code: `# Validação Go
go test ./...

# Validação frontend
cd web
npm run lint
npm run build`,
        },
        {
          type: 'code',
          title: 'Alternativa: Vercel dev',
          lang: 'bash',
          code: `npm install --global vercel
vercel dev`,
        },
        {
          type: 'list',
          title: 'Notas de desenvolvimento',
          items: [
            'Quando a app web roda localmente, /api/* é proxy para a API Go em desenvolvimento.',
            'O frontend hospedado fica em web/, enquanto o endpoint serverless fica em api/generate.go.',
            'O vercel.json na raiz conecta Next.js e Go no mesmo deploy.',
          ],
        },
      ],
    },
    troubleshooting: {
      id: 'troubleshooting',
      category: 'Operations',
      title: 'Troubleshooting',
      summary:
        'Falhas comuns, padrões de depuração e passos práticos para recuperar rápido.',
      keywords: ['troubleshooting', 'errors', 'preview', 'icons', 'vercel', 'theme'],
      blocks: [
        {
          type: 'callout',
          title: 'Estratégia recomendada de debug',
          tone: 'warn',
          text:
            'Reduza o YAML até a menor versão funcional primeiro. Depois adicione grupos, links, labels e presets de volta de forma incremental até isolar o bloco problemático.',
        },
        {
          type: 'list',
          title: 'Problemas comuns',
          items: [
            'Preview em branco ou erro de geração: valide primeiro a estrutura YAML. Falta de Canvas, Type inválido ou bloco Links malformado é a causa mais comum.',
            'Ícone de recurso incorreto: confirme se a URL do arquivo de definição é válida e aponta para um conjunto compatível.',
            'Tema claro/escuro inconsistente: limpe o localStorage do navegador ou use o switcher de tema para resincronizar a UI.',
            'Saída draw.io abre diferente do esperado: verifique zoom, edições manuais anteriores e se o YAML de origem mudou.',
            'CLI funciona mas a web local falha: confirme que a API Go está rodando na porta 8080 antes de iniciar o frontend, ou use vercel dev.',
            'Deploy da Vercel ignora Build Settings do painel: este repositório usa builds em vercel.json, então a configuração do arquivo tem precedência.',
          ],
        },
      ],
    },
    examples: {
      id: 'examples',
      category: 'Reference',
      title: 'Exemplos e Casos de Uso',
      summary:
        'Padrões DAC representativos para serverless, rede, aplicação e fluxos guiados por IA.',
      keywords: ['examples', 'patterns', 'serverless', 'network', 'sample yaml'],
      blocks: [
        {
          type: 'table',
          title: 'Padrões comuns',
          headers: ['Padrão', 'Melhor uso'],
          rows: [
            ['Recurso único', 'Diagrama simples de S3, Lambda ou DynamoDB'],
            ['Topologia de rede', 'VPC, subnets, ALB, NAT, IGW e posicionamento dos serviços'],
            ['Plataforma de aplicação', 'Serviços ECS ou Lambda com dados e mensageria'],
            ['Multi-region', 'Serviços replicados ou topologia de failover'],
            ['Diagrama gerado por IA', 'Prompt humano -> MCP -> DAC YAML -> PNG ou draw.io'],
          ],
        },
        {
          type: 'code',
          title: 'Pipeline serverless com worker',
          lang: 'yaml',
          code: `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children: [Cloud]
    Cloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Children: [App, Queue, Worker]
    App:
      Type: AWS::Lambda::Function
      Title: API Lambda
    Queue:
      Type: AWS::SQS::Queue
    Worker:
      Type: AWS::Lambda::Function
      Title: Worker Lambda
  Links:
    - Source: App
      Target: Queue
      TargetArrowHead:
        Type: Open
    - Source: Queue
      Target: Worker
      TargetArrowHead:
        Type: Open`,
        },
        {
          type: 'code',
          title: 'Frontend + API + camada de dados',
          lang: 'yaml',
          code: `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children: [Users, Cloud]
    Users:
      Type: AWS::Diagram::HorizontalStack
      Children: [Web, Mobile]
    Web:
      Type: AWS::Diagram::Resource
      Preset: User
      Title: Web User
    Mobile:
      Type: AWS::Diagram::Resource
      Preset: Mobile client
    Cloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Direction: vertical
      Children: [Api, Data]
    Api:
      Type: AWS::Diagram::HorizontalStack
      Children: [Gateway, Service]
    Gateway:
      Type: AWS::ApiGateway
      Title: Public API
    Service:
      Type: AWS::Lambda::Function
      Title: Orders Service
    Data:
      Type: AWS::Diagram::HorizontalStack
      Children: [Redis, Database]
    Redis:
      Type: AWS::ElastiCache::CacheCluster
    Database:
      Type: AWS::RDS::DBCluster
  Links:
    - Source: Web
      Target: Gateway
      TargetArrowHead:
        Type: Open
    - Source: Mobile
      Target: Gateway
      TargetArrowHead:
        Type: Open
    - Source: Gateway
      Target: Service
      TargetArrowHead:
        Type: Open
    - Source: Service
      Target: Redis
      TargetArrowHead:
        Type: Open
    - Source: Service
      Target: Database
      TargetArrowHead:
        Type: Open`,
        },
      ],
    },
    architecture: {
      id: 'architecture',
      category: 'Architecture',
      title: 'Arquitetura do Projeto',
      summary:
        'Camadas de alto nível do repositório, fluxo de execução e estrutura de deploy.',
      keywords: ['architecture', 'repo structure', 'internal/ctl', 'web', 'api', 'go'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'O repositório é dividido entre um motor Go e um frontend Next.js. O lado Go controla parse, layout, definitions e export. O lado web controla edição, preview, tema, docs e experiência do builder.',
        },
        {
          type: 'paragraph',
          text:
            'A rota /api/generate é implantada como função serverless em Go. O frontend é implantado a partir de web/ como app Next.js, enquanto o vercel.json conecta os dois runtimes.',
        },
        {
          type: 'table',
          title: 'Camadas do repositório',
          headers: ['Path', 'Responsabilidade'],
          rows: [
            ['cmd/', 'Entrypoints CLI e servidor local da API'],
            ['internal/ctl/', 'Pipeline principal de orquestração'],
            ['internal/types/', 'Modelo de runtime para recursos, links e geometria'],
            ['internal/definition/', 'Carregamento de definitions e metadados de ícones'],
            ['pkg/diagram/', 'Wrapper público Go para embedding'],
            ['api/generate.go', 'Handler serverless da Vercel'],
            ['web/', 'Frontend Next.js, builder, docs, editor e preview'],
          ],
        },
        {
          type: 'list',
          title: 'Fluxo de execução',
          items: [
            'O YAML de entrada é parseado.',
            'Os arquivos de definição são carregados e validados.',
            'Recursos e links formam o grafo de runtime.',
            'O layout é calculado de forma determinística quando possível.',
            'A saída é exportada como PNG, PDF ou draw.io XML.',
          ],
        },
      ],
    },
    about: {
      id: 'about',
      category: 'Architecture',
      title: 'Sobre o Autor',
      summary:
        'Contexto sobre o mantenedor deste fork e os objetivos da plataforma hospedada.',
      keywords: ['author', 'fernando azevedo', 'maintainer', 'fork'],
      blocks: [
        {
          type: 'paragraph',
          text:
            'Fernando Azevedo mantém este fork e expandiu o projeto original com frontend hospedado, fluxo de exportação draw.io, endpoint público de API, builder visual e integrações amigáveis para IA.',
        },
        {
          type: 'paragraph',
          text:
            'O foco deste fork é documentação arquitetural prática: saídas claras, fluxos amigáveis para automação e uma ponte melhor entre código, diagramas e tooling de IA.',
        },
        {
          type: 'list',
          title: 'Perfis',
          items: [
            'Website: https://fernando.moretes.com',
            'GitHub: https://github.com/fernandofatech',
            'LinkedIn: https://www.linkedin.com/in/fernando-francisco-azevedo/',
          ],
        },
      ],
    },
  },
}

export const docsContent: Record<DocsLocale, DocsCopy> = { en, pt }

export function getDocsCopy(locale: DocsLocale) {
  return docsContent[locale]
}

export function getDocsSections(locale: DocsLocale): DocsSection[] {
  const copy = getDocsCopy(locale)
  return docsSectionIds.map((id) => copy.sections[id])
}

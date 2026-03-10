'use client'

import { useMemo, useState } from 'react'
import Link from 'next/link'
import {
  ArrowLeft,
  Terminal,
  Globe,
  Code2,
  Cpu,
  BookOpen,
  Layers,
  User,
  ExternalLink,
  ChevronRight,
  Github,
  Linkedin,
  Wrench,
  FileCode2,
  Server,
  ShieldCheck,
  Workflow,
  TriangleAlert,
  Sparkles,
} from 'lucide-react'
import { useLanguage } from '@/lib/i18n'
import LanguageSwitcher from '@/components/LanguageSwitcher'
import ThemeSwitcher from '@/components/ThemeSwitcher'

type SectionId =
  | 'overview'
  | 'quick-start'
  | 'web-editor'
  | 'builder'
  | 'yaml-reference'
  | 'cli'
  | 'api'
  | 'drawio'
  | 'mcp'
  | 'local-dev'
  | 'troubleshooting'
  | 'examples'
  | 'architecture'
  | 'about'

function H2({ id, children }: { id: string; children: React.ReactNode }) {
  return (
    <h2
      id={id}
      className="text-2xl font-semibold text-[var(--text)] mt-14 mb-4 flex items-center gap-3 scroll-mt-20"
    >
      <span className="w-1.5 h-6 bg-[#FF9900] rounded-full inline-block flex-shrink-0" />
      {children}
    </h2>
  )
}

function H3({ children }: { children: React.ReactNode }) {
  return <h3 className="text-base font-semibold text-[var(--text-2)] mt-8 mb-2">{children}</h3>
}

function P({ children }: { children: React.ReactNode }) {
  return <p className="text-[var(--text-3)] text-sm leading-7 mb-3">{children}</p>
}

function Code({ children }: { children: React.ReactNode }) {
  return (
    <code className="bg-[var(--surface)] border border-[var(--border)] text-[#FF9900] text-xs px-1.5 py-0.5 rounded font-mono">
      {children}
    </code>
  )
}

function Pre({ children, lang }: { children: string; lang?: string }) {
  return (
    <div className="relative my-4">
      {lang && (
        <span className="absolute top-2 right-3 text-[10px] text-[var(--text-5)] font-mono uppercase">
          {lang}
        </span>
      )}
      <pre className="bg-[var(--code-bg)] border border-[var(--border)] rounded-xl p-4 text-xs text-[var(--text-2)] font-mono overflow-x-auto leading-6 whitespace-pre">
        {children}
      </pre>
    </div>
  )
}

function Badge({ children, color = 'orange' }: { children: React.ReactNode; color?: 'orange' | 'blue' | 'green' | 'slate' }) {
  const colors = {
    orange: 'bg-[#FF9900]/10 text-[#FF9900] border-[#FF9900]/20',
    blue: 'bg-sky-500/10 text-sky-400 border-sky-500/20',
    green: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
    slate: 'bg-[var(--surface)] text-[var(--text-3)] border-[var(--border)]',
  }
  return <span className={`inline-block border text-[10px] font-mono px-2 py-0.5 rounded ${colors[color]}`}>{children}</span>
}

function Table({ headers, rows }: { headers: string[]; rows: string[][] }) {
  return (
    <div className="overflow-x-auto my-4 rounded-xl border border-[var(--border)]">
      <table className="w-full text-xs text-left border-collapse">
        <thead className="bg-[var(--surface)]">
          <tr className="border-b border-[var(--border)]">
            {headers.map((h) => (
              <th key={h} className="px-4 py-3 text-[var(--text-4)] font-medium uppercase tracking-wider text-[10px]">
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-b last:border-b-0 border-[var(--border)]">
              {row.map((cell, j) => (
                <td key={j} className="px-4 py-3 text-[var(--text-2)] align-top font-mono">
                  {cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function Callout({ title, children, tone = 'info' }: { title: string; children: React.ReactNode; tone?: 'info' | 'tip' | 'warn' }) {
  const styles = {
    info: 'bg-sky-500/8 border-sky-500/20 text-sky-100',
    tip: 'bg-[#FF9900]/5 border-[#FF9900]/20 text-[#ffd7a3]',
    warn: 'bg-amber-500/8 border-amber-500/20 text-amber-100',
  }
  return (
    <div className={`border rounded-xl px-4 py-3 text-xs my-4 ${styles[tone]}`}>
      <p className="font-semibold mb-1">{title}</p>
      <div className="leading-6">{children}</div>
    </div>
  )
}

function BulletList({ items }: { items: React.ReactNode[] }) {
  return (
    <ul className="space-y-2 mb-4">
      {items.map((item, index) => (
        <li key={index} className="text-sm text-[var(--text-3)] leading-7 flex gap-3">
          <span className="mt-2 h-1.5 w-1.5 rounded-full bg-[#FF9900] flex-shrink-0" />
          <span>{item}</span>
        </li>
      ))}
    </ul>
  )
}

const CONTENT = {
  en: {
    labels: {
      title: 'Documentation',
      contents: 'Wiki',
      back: 'Back to Home',
      github: 'GitHub',
      site: 'Live Site',
      repo: 'Repository',
      fork: 'diagram-as-code · Comprehensive fork by Fernando Azevedo',
      viewGithub: 'Open repository',
    },
    nav: {
      overview: 'Overview',
      quickStart: 'Quick Start',
      webEditor: 'Web Editor',
      builder: 'Visual Builder',
      yamlReference: 'YAML Reference',
      cli: 'CLI Reference',
      api: 'API Reference',
      drawio: 'draw.io Export',
      mcp: 'MCP Server',
      localDev: 'Local Development',
      troubleshooting: 'Troubleshooting',
      examples: 'Examples & Use Cases',
      architecture: 'Architecture',
      about: 'About the Author',
    },
    hero: {
      intro:
        'This page is the complete product wiki for the hosted editor, CLI, API, builder, draw.io export, local development flow, and architecture model behind diagram-as-code.',
      sub:
        'If you only need to generate a diagram quickly, start with Quick Start. If you need to integrate the tool into an application, jump to API, CLI, or MCP. If you want to extend or maintain the project, use Local Development and Architecture.',
    },
    overview: {
      paragraphs: [
        'diagram-as-code converts declarative YAML into AWS architecture diagrams. The same runtime model powers PNG rendering, draw.io export, the browser editor, the visual builder, and AI integrations.',
        'The core idea is simple: infrastructure diagrams should be versioned, reviewable, reproducible, and generated from source, not manually maintained in a GUI.',
      ],
      outputs: [
        ['PNG diagram', 'image/png', 'Documentation, PRs, architecture reviews, exports for presentations'],
        ['draw.io file', 'application/xml', 'Editable diagrams in diagrams.net / draw.io'],
        ['YAML source', 'text/yaml', 'Reusable infrastructure diagram definitions'],
      ],
      useCases: [
        'Document AWS architectures without manually dragging shapes.',
        'Generate diagrams in CI/CD from versioned YAML.',
        'Create diagrams from chat through an MCP server.',
        'Build diagrams in the browser and export them as PNG or draw.io.',
        'Use the builder when a team wants structure without writing YAML by hand.',
      ],
    },
    quickStart: {
      steps: [
        'Open https://dac.moretes.com and choose Editor or Builder.',
        'Paste YAML or load one of the built-in examples.',
        'Generate a PNG for documentation or a draw.io file for later editing.',
        'If you prefer local automation, install the CLI and run awsdac against your YAML file.',
      ],
      yaml: `Diagram:
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
    webEditor: {
      paragraphs: [
        'The web editor is the fastest path from YAML to rendered diagram. It runs in the browser and calls the hosted Go backend through /api/generate.',
        'The left panel is a Monaco editor with YAML syntax highlighting. The right panel shows the generated output or the draw.io ready state.',
      ],
      workflows: [
        'Fast authoring: write YAML, press Ctrl+Enter, inspect the output immediately.',
        'Review loop: tweak layout, labels, or children ordering, then regenerate until the diagram is stable.',
        'Export loop: switch between PNG and draw.io depending on whether the target is a document or a diagram editor.',
      ],
      features: [
        ['Monaco editor', 'Rich editing with syntax highlighting and fast keyboard-driven workflow'],
        ['Examples menu', 'Load realistic starter templates without writing from scratch'],
        ['Theme + language', 'Supports dark/light UI and English/Portuguese shell text'],
        ['Direct export', 'Download PNG or .drawio output from the same UI'],
      ],
    },
    builder: {
      paragraphs: [
        'The Visual Builder is a form-based layer on top of the same YAML model. It is useful for teams that want valid DAC YAML without remembering every field name.',
        'The builder is not a separate rendering engine. It generates YAML, and that YAML is rendered by the same backend used everywhere else.',
      ],
      cases: [
        'Business analysts or architects can assemble a first draft without typing YAML.',
        'Platform teams can standardize naming and structure before switching to hand-edited YAML.',
        'The builder can be used as a teaching tool to learn the DAC schema by seeing YAML update live.',
      ],
      notes: [
        'Definition file can come from the official URL or a local file path.',
        'Resources, links, border children, arrow heads, positions, and labels are exposed as form controls.',
        'The generated YAML can be copied or sent directly to the editor.',
      ],
    },
    yamlReference: {
      intro:
        'The DAC schema is intentionally small. Most diagrams are defined with three top-level concepts: DefinitionFiles, Resources, and Links.',
      concepts: [
        ['Diagram.DefinitionFiles', 'Definition source for icon metadata and presets'],
        ['Resources.<name>.Type', 'Resource type or diagram primitive'],
        ['Resources.<name>.Children', 'Child resources inside a container'],
        ['Resources.<name>.Direction', 'Layout direction for children, usually vertical or horizontal'],
        ['Resources.<name>.Preset', 'Named visual preset from the definition file'],
        ['Resources.<name>.Title', 'Display label override'],
        ['Resources.<name>.BorderChildren', 'Resources attached to the border of a container'],
        ['Resources.<name>.Align', 'Alignment for child layout when supported'],
        ['Links[].Source / Target', 'Logical connection between resource names'],
        ['Links[].Labels.*.Title', 'Optional labels placed around the edge'],
        ['Links[].TargetArrowHead.Type', 'Arrow style such as Open, Default, or none'],
      ],
      primitives: [
        ['AWS::Diagram::Canvas', 'Root drawing surface'],
        ['AWS::Diagram::Cloud', 'AWS cloud grouping container'],
        ['AWS::EC2::VPC', 'VPC container with child resources'],
        ['AWS::EC2::Subnet', 'Subnet container, often with PublicSubnet or PrivateSubnet preset'],
        ['AWS::Diagram::HorizontalStack', 'Horizontal layout wrapper'],
        ['AWS::Diagram::VerticalStack', 'Vertical layout wrapper'],
        ['AWS::Diagram::Resource', 'Generic visual resource such as User or Mobile client'],
      ],
      bestPractices: [
        'Always keep a Canvas root and attach the real top-level nodes as its children.',
        'Use stable YAML keys like ALB, PublicSubnetA, OrdersService, and Redis instead of random names.',
        'Prefer containers such as Cloud, VPC, and Subnet to express architectural boundaries clearly.',
        'Use Title only when the displayed label should differ from the YAML key.',
        'Keep links intentional. Too many crossing links usually means the diagram should be decomposed.',
      ],
      example: `Diagram:
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
    cli: {
      install: `# macOS
brew install awsdac

# Go install
go install github.com/fernandofatech/diagram-as-code/cmd/awsdac@latest`,
      usage: `# Generate PNG
awsdac examples/alb-ec2.yaml

# Custom output file
awsdac examples/alb-ec2.yaml -o my-diagram.png

# Generate draw.io XML
awsdac examples/alb-ec2.yaml --drawio -o output.drawio

# Output extension also selects format
awsdac examples/alb-ec2.yaml -o output.drawio`,
      flags: [
        ['-o, --output', 'Output file name'],
        ['--drawio', 'Generate draw.io instead of PNG'],
        ['-f, --force', 'Overwrite output without confirmation'],
        ['--width / --height', 'Resize PNG output'],
        ['-t, --template', 'Render the input as Go text/template'],
        ['-c, --cfn-template', 'Build diagram from CloudFormation template'],
        ['-d, --dac-file', 'Generate DAC YAML from CloudFormation template'],
        ['--allow-untrusted-definitions', 'Allow non-official definition sources'],
        ['-v, --verbose', 'Verbose logging'],
      ],
      scenarios: [
        'Local authoring: ideal for engineers who keep architecture diagrams in the same repo as infrastructure code.',
        'CI/CD generation: generate PNGs during pull requests or release pipelines.',
        'Batch rendering: render multiple YAML diagrams as part of internal docs automation.',
      ],
    },
    api: {
      intro:
        'The hosted API powers the browser editor and can also be called directly from scripts, internal portals, CI jobs, or AI agents.',
      request: `POST https://dac.moretes.com/api/generate
Content-Type: application/json

{
  "yaml": "Diagram:\\n  ...",
  "format": "png"
}`,
      responseRows: [
        ['png', 'image/png', 'Raw PNG bytes'],
        ['drawio', 'application/xml', 'mxGraphModel XML'],
        ['error', 'application/json', '{"error":"message"}'],
      ],
      curl: `curl -X POST https://dac.moretes.com/api/generate \\
  -H "Content-Type: application/json" \\
  -d '{"yaml":"Diagram:\\n  DefinitionFiles:\\n    - Type: URL\\n      Url: \\"https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml\\"\\n  Resources:\\n    Canvas:\\n      Type: AWS::Diagram::Canvas\\n      Children:\\n        - Bucket\\n    Bucket:\\n      Type: AWS::S3::Bucket","format":"png"}' \\
  --output diagram.png`,
      ts: `const yaml = \`Diagram:
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
      notes: [
        'The request body expects a DAC YAML string, not a file upload.',
        'Format defaults to png when omitted.',
        'When generation fails, the API returns a JSON error payload.',
        'The endpoint is suitable for server-to-server usage as well as browser usage.',
      ],
    },
    drawio: {
      paragraphs: [
        'draw.io export is not a screenshot or a separate layout path. It uses the same resource model and geometry used by PNG rendering.',
        'This is important because it keeps PNG and draw.io outputs aligned. The editable file preserves the diagram structure rather than flattening it.',
      ],
      pipeline: [
        'Parse DAC YAML into runtime resources and links.',
        'Run the same layout engine used for PNG output, including Scale and ZeroAdjust.',
        'Reorder children from link topology to preserve consistent placement.',
        'Export resources as mxCell nodes and connections as mxCell edges.',
        'Embed AWS SVG icons as data URIs so the file is self-contained.',
        'Write mxGraphModel XML compatible with diagrams.net / draw.io.',
      ],
      useCases: [
        'Start with code, finish with manual annotations in draw.io.',
        'Generate a baseline AWS diagram and let non-technical users adjust labels or framing later.',
        'Store editable diagrams in repositories without maintaining the entire layout manually.',
      ],
    },
    mcp: {
      paragraphs: [
        'The MCP server exposes diagram generation tools to AI assistants. This makes the repository useful not only as a CLI or website, but also as part of agent workflows.',
        'Typical usage is: the assistant writes DAC YAML from a human request, calls the MCP tool, and saves the output as PNG or draw.io.',
      ],
      install: `go install github.com/awslabs/diagram-as-code/cmd/awsdac-mcp-server@latest`,
      config: `{
  "mcpServers": {
    "diagram-as-code": {
      "command": "awsdac-mcp-server",
      "args": []
    }
  }
}`,
      tools: [
        ['create_diagram', 'Generate PNG from DAC YAML'],
        ['create_drawio', 'Generate draw.io XML from DAC YAML'],
      ],
      prompt: `Create an AWS architecture diagram with:
- one VPC
- two public subnets
- an internet-facing ALB
- two EC2 instances
- an Internet Gateway

Return PNG output and save it to ~/Desktop/architecture.png`,
    },
    localDev: {
      prerequisites: [
        ['Go', '1.21+', 'Go backend, CLI, and tests'],
        ['Node.js', '18+ (Node 20+ recommended for Next 16)', 'Web frontend'],
        ['npm', '9+', 'Frontend dependencies and scripts'],
      ],
      run: `git clone https://github.com/fernandofatech/diagram-as-code.git
cd diagram-as-code

# Terminal 1 - local Go API
go run ./cmd/api-dev

# Terminal 2 - Next.js app
cd web
npm install
npm run dev`,
      test: `# Go unit and package tests
go test ./...

# Frontend quality
cd web
npm run lint
npm run build`,
      vercel: `npm install --global vercel
vercel dev`,
      notes: [
        'When running the web app locally, /api/* is proxied to the Go API in development.',
        'The hosted frontend lives in web/, while the serverless endpoint lives in api/generate.go.',
        'The root vercel.json controls how Next.js and Go are deployed together.',
      ],
    },
    troubleshooting: {
      items: [
        'Blank preview or error state: validate the YAML structure first. A missing Canvas, bad Type, or malformed Links block is the most common cause.',
        'Incorrect resource icon: verify that the definition file URL is valid and points to a compatible AWS definition file.',
        'Unexpected light/dark appearance: clear localStorage theme settings or use the header theme switcher to resync the UI state.',
        'draw.io output opens but looks different: check whether the file was opened in diagrams.net with its default zoom or if manual edits were already applied.',
        'CLI works but the web app fails locally: make sure the Go API is running on port 8080 before starting the frontend, or use vercel dev.',
        'Vercel deploy ignores dashboard build settings: this repo uses builds in vercel.json, so the file-based config takes precedence.',
      ],
    },
    examples: {
      useCases: [
        ['Single resource', 'S3 bucket, Lambda, or DynamoDB table reference diagram'],
        ['Network topology', 'VPC, subnets, ALB, NAT, IGW, and service placement'],
        ['App platform', 'ECS or Lambda services with storage, messaging, and identity layers'],
        ['Multi-region', 'Replicated services or failover topology'],
        ['AI-generated diagram', 'Human prompt -> MCP -> DAC YAML -> PNG or draw.io output'],
      ],
      sample1: `Diagram:
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
      sample2: `Diagram:
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
    architecture: {
      paragraphs: [
        'The repository is split into a Go engine and a Next.js frontend. The Go side owns parsing, layout, definitions, and export. The web side owns editing, preview interaction, theming, docs, and the builder UX.',
        'The API route /api/generate is deployed as a Go serverless function. The frontend is deployed from web/ as a Next.js app, while vercel.json wires both runtimes together.',
      ],
      layers: [
        ['cmd/', 'CLI entrypoints and local API dev server'],
        ['internal/ctl/', 'Main orchestration pipeline'],
        ['internal/types/', 'Runtime model for resources, links, and geometry'],
        ['internal/definition/', 'Definition loading and icon metadata'],
        ['pkg/diagram/', 'Public Go wrapper for embedding'],
        ['api/generate.go', 'Vercel serverless handler'],
        ['web/', 'Next.js frontend, builder, docs, editor, preview'],
      ],
      flow: [
        'Input YAML is parsed.',
        'Definition files are loaded and validated.',
        'Resources and links are constructed into the runtime graph.',
        'Layout is computed deterministically where possible.',
        'Output is exported either as PNG or draw.io XML.',
      ],
    },
    about: {
      subtitle: 'Senior Solutions Architect · 16+ years of global experience',
      bio: [
        'Fernando Azevedo maintains this fork and expanded the original project with a hosted web frontend, draw.io export flow, API endpoint, builder, and AI-friendly integrations.',
        'The focus of this fork is practical architecture documentation: clear outputs, automation-friendly workflows, and a better bridge between code, diagrams, and AI tooling.',
      ],
    },
  },
  pt: {
    labels: {
      title: 'Documentação',
      contents: 'Wiki',
      back: 'Voltar ao Início',
      github: 'GitHub',
      site: 'Site',
      repo: 'Repositório',
      fork: 'diagram-as-code · Fork completo por Fernando Azevedo',
      viewGithub: 'Abrir repositório',
    },
    nav: {
      overview: 'Visão Geral',
      quickStart: 'Começo Rápido',
      webEditor: 'Editor Web',
      builder: 'Builder Visual',
      yamlReference: 'Referência YAML',
      cli: 'Referência CLI',
      api: 'Referência API',
      drawio: 'Exportação draw.io',
      mcp: 'Servidor MCP',
      localDev: 'Desenvolvimento Local',
      troubleshooting: 'Troubleshooting',
      examples: 'Exemplos e Casos de Uso',
      architecture: 'Arquitetura',
      about: 'Sobre o Autor',
    },
    hero: {
      intro:
        'Esta página é a wiki completa do produto para o editor hospedado, CLI, API, builder, exportação draw.io, fluxo de desenvolvimento local e arquitetura do projeto diagram-as-code.',
      sub:
        'Se você só quer gerar um diagrama rápido, comece em Começo Rápido. Se quiser integrar em outra aplicação, vá para API, CLI ou MCP. Se quiser manter ou evoluir o projeto, use Desenvolvimento Local e Arquitetura.',
    },
    overview: {
      paragraphs: [
        'diagram-as-code converte YAML declarativo em diagramas de arquitetura AWS. O mesmo modelo de runtime alimenta a renderização PNG, a exportação draw.io, o editor no navegador, o builder visual e integrações com IA.',
        'A ideia central é simples: diagramas de infraestrutura devem ser versionados, revisáveis, reproduzíveis e gerados a partir de código-fonte, não mantidos manualmente em GUI.',
      ],
      outputs: [
        ['Diagrama PNG', 'image/png', 'Documentação, PRs, revisões de arquitetura, apresentações'],
        ['Arquivo draw.io', 'application/xml', 'Diagramas editáveis no diagrams.net / draw.io'],
        ['YAML fonte', 'text/yaml', 'Definições reutilizáveis de diagramas'],
      ],
      useCases: [
        'Documentar arquiteturas AWS sem arrastar shapes manualmente.',
        'Gerar diagramas em CI/CD a partir de YAML versionado.',
        'Criar diagramas por chat usando servidor MCP.',
        'Montar diagramas no navegador e exportar como PNG ou draw.io.',
        'Usar o builder quando o time quer estrutura sem escrever YAML manualmente.',
      ],
    },
    quickStart: {
      steps: [
        'Abra https://dac.moretes.com e escolha Editor ou Builder.',
        'Cole seu YAML ou carregue um dos exemplos prontos.',
        'Gere um PNG para documentação ou um arquivo draw.io para edição posterior.',
        'Se preferir automação local, instale o CLI e execute awsdac no seu arquivo YAML.',
      ],
      yaml: `Diagram:
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
    webEditor: {
      paragraphs: [
        'O editor web é o caminho mais rápido entre YAML e diagrama renderizado. Ele roda no navegador e chama o backend Go hospedado através de /api/generate.',
        'O painel esquerdo é um editor Monaco com destaque de sintaxe YAML. O painel direito mostra a saída gerada ou o estado pronto para draw.io.',
      ],
      workflows: [
        'Autoria rápida: escreva YAML, pressione Ctrl+Enter e inspecione a saída imediatamente.',
        'Ciclo de revisão: ajuste layout, labels ou ordem dos children e gere novamente até estabilizar o diagrama.',
        'Ciclo de exportação: alterne entre PNG e draw.io conforme o destino for documento ou editor de diagramas.',
      ],
      features: [
        ['Editor Monaco', 'Edição rica com sintaxe destacada e fluxo orientado a teclado'],
        ['Menu de exemplos', 'Carregue templates realistas sem começar do zero'],
        ['Tema + idioma', 'Suporte a UI claro/escuro e shell textual em inglês/português'],
        ['Exportação direta', 'Baixe PNG ou .drawio da mesma interface'],
      ],
    },
    builder: {
      paragraphs: [
        'O Visual Builder é uma camada baseada em formulários sobre o mesmo modelo YAML. Ele é útil para equipes que querem gerar YAML DAC válido sem memorizar todos os campos.',
        'O builder não é um motor de renderização separado. Ele gera YAML, e esse YAML é renderizado pelo mesmo backend usado em todo o projeto.',
      ],
      cases: [
        'Analistas de negócio ou arquitetos podem montar um rascunho inicial sem digitar YAML.',
        'Times de plataforma podem padronizar nomes e estrutura antes de migrar para YAML manual.',
        'O builder funciona como ferramenta de aprendizado ao mostrar o YAML sendo montado em tempo real.',
      ],
      notes: [
        'O arquivo de definição pode vir da URL oficial ou de um caminho local.',
        'Recursos, links, border children, pontas de seta, posições e labels estão disponíveis via formulário.',
        'O YAML gerado pode ser copiado ou enviado direto ao editor.',
      ],
    },
    yamlReference: {
      intro:
        'O schema DAC é intencionalmente pequeno. A maior parte dos diagramas é definida com três conceitos principais: DefinitionFiles, Resources e Links.',
      concepts: [
        ['Diagram.DefinitionFiles', 'Fonte de definição para metadados de ícones e presets'],
        ['Resources.<name>.Type', 'Tipo do recurso ou primitive do diagrama'],
        ['Resources.<name>.Children', 'Recursos filhos dentro de um container'],
        ['Resources.<name>.Direction', 'Direção do layout dos filhos, normalmente vertical ou horizontal'],
        ['Resources.<name>.Preset', 'Preset visual nomeado vindo do arquivo de definição'],
        ['Resources.<name>.Title', 'Sobrescrita do label exibido'],
        ['Resources.<name>.BorderChildren', 'Recursos anexados à borda do container'],
        ['Resources.<name>.Align', 'Alinhamento do layout quando suportado'],
        ['Links[].Source / Target', 'Conexão lógica entre nomes de recursos'],
        ['Links[].Labels.*.Title', 'Labels opcionais ao redor da aresta'],
        ['Links[].TargetArrowHead.Type', 'Estilo da seta como Open, Default ou none'],
      ],
      primitives: [
        ['AWS::Diagram::Canvas', 'Superfície raiz do desenho'],
        ['AWS::Diagram::Cloud', 'Container de agrupamento AWS Cloud'],
        ['AWS::EC2::VPC', 'Container de VPC com recursos filhos'],
        ['AWS::EC2::Subnet', 'Container de subnet, normalmente com PublicSubnet ou PrivateSubnet'],
        ['AWS::Diagram::HorizontalStack', 'Wrapper de layout horizontal'],
        ['AWS::Diagram::VerticalStack', 'Wrapper de layout vertical'],
        ['AWS::Diagram::Resource', 'Recurso visual genérico como User ou Mobile client'],
      ],
      bestPractices: [
        'Sempre mantenha um Canvas raiz e conecte os nós de topo reais como children dele.',
        'Use chaves YAML estáveis como ALB, PublicSubnetA, OrdersService e Redis.',
        'Prefira containers como Cloud, VPC e Subnet para expressar fronteiras arquiteturais.',
        'Use Title apenas quando o label exibido precisa diferir da chave YAML.',
        'Mantenha os links intencionais. Muitos cruzamentos normalmente indicam que o diagrama deve ser quebrado.',
      ],
      example: `Diagram:
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
    cli: {
      install: `# macOS
brew install awsdac

# Go install
go install github.com/fernandofatech/diagram-as-code/cmd/awsdac@latest`,
      usage: `# Gerar PNG
awsdac examples/alb-ec2.yaml

# Arquivo de saída customizado
awsdac examples/alb-ec2.yaml -o meu-diagrama.png

# Gerar XML draw.io
awsdac examples/alb-ec2.yaml --drawio -o output.drawio

# A extensão da saída também seleciona o formato
awsdac examples/alb-ec2.yaml -o output.drawio`,
      flags: [
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
      scenarios: [
        'Autoria local: ideal para engenheiros que mantêm diagramas no mesmo repositório do código.',
        'Geração em CI/CD: gere PNGs durante pull requests ou pipelines de release.',
        'Renderização em lote: gere múltiplos diagramas YAML como parte da automação de documentação.',
      ],
    },
    api: {
      intro:
        'A API hospedada alimenta o editor do navegador e também pode ser chamada diretamente de scripts, portais internos, jobs de CI ou agentes de IA.',
      request: `POST https://dac.moretes.com/api/generate
Content-Type: application/json

{
  "yaml": "Diagram:\\n  ...",
  "format": "png"
}`,
      responseRows: [
        ['png', 'image/png', 'Bytes PNG brutos'],
        ['drawio', 'application/xml', 'XML mxGraphModel'],
        ['error', 'application/json', '{"error":"message"}'],
      ],
      curl: `curl -X POST https://dac.moretes.com/api/generate \\
  -H "Content-Type: application/json" \\
  -d '{"yaml":"Diagram:\\n  DefinitionFiles:\\n    - Type: URL\\n      Url: \\"https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml\\"\\n  Resources:\\n    Canvas:\\n      Type: AWS::Diagram::Canvas\\n      Children:\\n        - Bucket\\n    Bucket:\\n      Type: AWS::S3::Bucket","format":"png"}' \\
  --output diagram.png`,
      ts: `const yaml = \`Diagram:
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
      notes: [
        'O corpo espera uma string YAML DAC, não upload de arquivo.',
        'Quando format é omitido, o padrão é png.',
        'Quando a geração falha, a API retorna um payload JSON com erro.',
        'O endpoint funciona tanto para uso browser quanto server-to-server.',
      ],
    },
    drawio: {
      paragraphs: [
        'A exportação draw.io não é screenshot nem um caminho de layout separado. Ela usa o mesmo modelo de recursos e a mesma geometria usada no PNG.',
        'Isso é importante porque mantém PNG e draw.io alinhados. O arquivo editável preserva a estrutura do diagrama em vez de achatar o resultado.',
      ],
      pipeline: [
        'Parse do YAML DAC para recursos e links de runtime.',
        'Execução do mesmo layout usado no PNG, incluindo Scale e ZeroAdjust.',
        'Reordenação dos children a partir da topologia de links para manter consistência visual.',
        'Exportação dos recursos como nós mxCell e das conexões como arestas mxCell.',
        'Embute ícones SVG da AWS como data URIs para que o arquivo seja autocontido.',
        'Gravação do XML mxGraphModel compatível com diagrams.net / draw.io.',
      ],
      useCases: [
        'Começar com código e finalizar com anotações manuais no draw.io.',
        'Gerar um baseline AWS e permitir que usuários não técnicos ajustem labels depois.',
        'Armazenar diagramas editáveis em repositórios sem manter layout 100% manual.',
      ],
    },
    mcp: {
      paragraphs: [
        'O servidor MCP expõe ferramentas de geração de diagrama para assistentes de IA. Isso torna o projeto útil não só como CLI ou website, mas também dentro de fluxos baseados em agentes.',
        'O uso típico é: o assistente escreve DAC YAML a partir do pedido humano, chama a ferramenta MCP e salva a saída como PNG ou draw.io.',
      ],
      install: `go install github.com/awslabs/diagram-as-code/cmd/awsdac-mcp-server@latest`,
      config: `{
  "mcpServers": {
    "diagram-as-code": {
      "command": "awsdac-mcp-server",
      "args": []
    }
  }
}`,
      tools: [
        ['create_diagram', 'Gera PNG a partir de DAC YAML'],
        ['create_drawio', 'Gera XML draw.io a partir de DAC YAML'],
      ],
      prompt: `Crie um diagrama de arquitetura AWS com:
- uma VPC
- duas subnets públicas
- um ALB internet-facing
- duas instâncias EC2
- um Internet Gateway

Retorne saída PNG e salve em ~/Desktop/architecture.png`,
    },
    localDev: {
      prerequisites: [
        ['Go', '1.21+', 'Backend Go, CLI e testes'],
        ['Node.js', '18+ (Node 20+ recomendado para Next 16)', 'Frontend web'],
        ['npm', '9+', 'Dependências e scripts do frontend'],
      ],
      run: `git clone https://github.com/fernandofatech/diagram-as-code.git
cd diagram-as-code

# Terminal 1 - API Go local
go run ./cmd/api-dev

# Terminal 2 - app Next.js
cd web
npm install
npm run dev`,
      test: `# Testes Go
go test ./...

# Qualidade do frontend
cd web
npm run lint
npm run build`,
      vercel: `npm install --global vercel
vercel dev`,
      notes: [
        'Quando a aplicação web roda localmente, /api/* é proxy para a API Go em desenvolvimento.',
        'O frontend hospedado fica em web/, enquanto o endpoint serverless fica em api/generate.go.',
        'O arquivo vercel.json na raiz conecta Next.js e Go no mesmo deploy.',
      ],
    },
    troubleshooting: {
      items: [
        'Preview em branco ou erro: valide primeiro a estrutura do YAML. Falta de Canvas, Type inválido ou bloco Links malformado é a causa mais comum.',
        'Ícone incorreto: confirme se a DefinitionFiles URL é válida e compatível com o conjunto AWS esperado.',
        'Tema claro/escuro inconsistente: limpe o localStorage do navegador ou use o switcher do cabeçalho para ressincronizar o estado visual.',
        'draw.io abre diferente do esperado: verifique se o arquivo foi aberto no diagrams.net com zoom padrão ou se já sofreu edição manual.',
        'CLI funciona mas a web local falha: confirme que a API Go está rodando na porta 8080 antes de subir o frontend, ou use vercel dev.',
        'Deploy da Vercel ignora Build Settings do painel: este repositório usa builds em vercel.json, então a configuração do arquivo tem precedência.',
      ],
    },
    examples: {
      useCases: [
        ['Recurso único', 'Diagrama simples de S3, Lambda ou DynamoDB'],
        ['Topologia de rede', 'VPC, subnets, ALB, NAT, IGW e posicionamento dos serviços'],
        ['Plataforma de aplicação', 'Serviços ECS ou Lambda com storage, mensageria e identidade'],
        ['Multi-region', 'Serviços replicados ou topologia de failover'],
        ['Diagrama gerado por IA', 'Prompt humano -> MCP -> DAC YAML -> PNG ou draw.io'],
      ],
      sample1: `Diagram:
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
      sample2: `Diagram:
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
    architecture: {
      paragraphs: [
        'O repositório é dividido entre um motor Go e um frontend Next.js. O lado Go controla parse, layout, definitions e export. O lado web controla edição, preview, tema, docs e experiência do builder.',
        'A rota /api/generate é implantada como função serverless em Go. O frontend sai de web/ como app Next.js, e o vercel.json conecta ambos os runtimes.',
      ],
      layers: [
        ['cmd/', 'Entrypoints CLI e servidor local da API'],
        ['internal/ctl/', 'Pipeline principal de orquestração'],
        ['internal/types/', 'Modelo de runtime para recursos, links e geometria'],
        ['internal/definition/', 'Carregamento de definitions e metadados de ícones'],
        ['pkg/diagram/', 'Wrapper público Go para embedding'],
        ['api/generate.go', 'Handler serverless da Vercel'],
        ['web/', 'Frontend Next.js, builder, docs, editor e preview'],
      ],
      flow: [
        'O YAML de entrada é parseado.',
        'Os arquivos de definição são carregados e validados.',
        'Recursos e links formam o grafo de runtime.',
        'O layout é calculado de forma determinística quando possível.',
        'A saída é exportada como PNG ou draw.io XML.',
      ],
    },
    about: {
      subtitle: 'Arquiteto de Soluções Sênior · 16+ anos de experiência global',
      bio: [
        'Fernando Azevedo mantém este fork e expandiu o projeto original com frontend hospedado, fluxo draw.io, endpoint de API, builder e integrações amigáveis para IA.',
        'O foco deste fork é documentação arquitetural prática: saídas claras, fluxos automatizáveis e uma ponte melhor entre código, diagramas e tooling de IA.',
      ],
    },
  },
}

export default function DocsPage() {
  const { lang } = useLanguage()
  const copy = CONTENT[lang] ?? CONTENT.en
  const [active, setActive] = useState<SectionId>('overview')

  const nav = useMemo(
    () => [
      { id: 'overview' as SectionId, label: copy.nav.overview, icon: <BookOpen size={14} /> },
      { id: 'quick-start' as SectionId, label: copy.nav.quickStart, icon: <Sparkles size={14} /> },
      { id: 'web-editor' as SectionId, label: copy.nav.webEditor, icon: <Globe size={14} /> },
      { id: 'builder' as SectionId, label: copy.nav.builder, icon: <Wrench size={14} /> },
      { id: 'yaml-reference' as SectionId, label: copy.nav.yamlReference, icon: <FileCode2 size={14} /> },
      { id: 'cli' as SectionId, label: copy.nav.cli, icon: <Terminal size={14} /> },
      { id: 'api' as SectionId, label: copy.nav.api, icon: <Code2 size={14} /> },
      { id: 'drawio' as SectionId, label: copy.nav.drawio, icon: <Layers size={14} /> },
      { id: 'mcp' as SectionId, label: copy.nav.mcp, icon: <Cpu size={14} /> },
      { id: 'local-dev' as SectionId, label: copy.nav.localDev, icon: <Server size={14} /> },
      { id: 'troubleshooting' as SectionId, label: copy.nav.troubleshooting, icon: <TriangleAlert size={14} /> },
      { id: 'examples' as SectionId, label: copy.nav.examples, icon: <Workflow size={14} /> },
      { id: 'architecture' as SectionId, label: copy.nav.architecture, icon: <ShieldCheck size={14} /> },
      { id: 'about' as SectionId, label: copy.nav.about, icon: <User size={14} /> },
    ],
    [copy],
  )

  function scrollTo(id: SectionId) {
    setActive(id)
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
  }

  return (
    <div className="flex flex-col h-screen bg-[var(--bg)] overflow-hidden">
      <header className="flex items-center justify-between px-4 h-12 border-b border-[var(--border)] flex-shrink-0">
        <div className="flex items-center gap-3 min-w-0">
          <Link
            href="/"
            className="flex items-center gap-1.5 text-xs text-[var(--text-4)] hover:text-[var(--text-2)] transition-colors"
          >
            <ArrowLeft size={13} />
            {copy.labels.back}
          </Link>
          <div className="w-px h-4 bg-[var(--border)]" />
          <div className="flex items-center gap-2 min-w-0">
            <div className="w-5 h-5 bg-[#FF9900] rounded flex items-center justify-center flex-shrink-0">
              <svg viewBox="0 0 16 16" fill="white" className="w-3 h-3">
                <rect x="1" y="1" width="6" height="6" rx="1" />
                <rect x="9" y="1" width="6" height="6" rx="1" />
                <rect x="1" y="9" width="6" height="6" rx="1" />
                <rect x="9" y="9" width="6" height="6" rx="1" />
              </svg>
            </div>
            <span className="text-sm font-semibold text-[var(--text)] truncate">{copy.labels.title}</span>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <LanguageSwitcher />
          <ThemeSwitcher />
          <a
            href="https://github.com/fernandofatech/diagram-as-code"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-1.5 text-xs text-[var(--text-4)] hover:text-[var(--text-2)] transition-colors"
          >
            <Github size={14} />
            {copy.labels.github}
          </a>
          <a
            href="https://dac.moretes.com"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-1.5 text-xs text-[#FF9900] hover:text-[#ffb340] transition-colors"
          >
            <ExternalLink size={12} />
            {copy.labels.site}
          </a>
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        <aside className="w-64 flex-shrink-0 border-r border-[var(--border)] overflow-y-auto py-4">
          <p className="px-4 text-[10px] text-[var(--text-6)] uppercase tracking-widest font-medium mb-2">
            {copy.labels.contents}
          </p>
          <nav className="space-y-0.5 px-2">
            {nav.map(({ id, label, icon }) => (
              <button
                key={id}
                type="button"
                onClick={() => scrollTo(id)}
                className={`w-full flex items-center gap-2 px-3 py-2 rounded text-xs text-left transition-colors ${
                  active === id
                    ? 'bg-[#FF9900]/10 text-[#FF9900]'
                    : 'text-[var(--text-4)] hover:text-[var(--text-2)] hover:bg-[var(--surface)]'
                }`}
              >
                {icon}
                {label}
                {active === id && <ChevronRight size={10} className="ml-auto" />}
              </button>
            ))}
          </nav>
        </aside>

        <main className="flex-1 overflow-y-auto px-8 py-8">
          <div className="max-w-5xl">
            <div className="mb-8 rounded-2xl border border-[var(--border)] bg-[var(--code-bg)] p-6">
              <div className="flex flex-wrap gap-2 mb-4">
                <Badge>https://dac.moretes.com</Badge>
                <Badge color="blue">PNG</Badge>
                <Badge color="green">draw.io</Badge>
                <Badge color="slate">CLI</Badge>
                <Badge color="slate">API</Badge>
                <Badge color="slate">MCP</Badge>
              </div>
              <p className="text-[var(--text)] text-lg leading-8 mb-2">{copy.hero.intro}</p>
              <p className="text-sm text-[var(--text-3)] leading-7">{copy.hero.sub}</p>
            </div>

            <H2 id="overview">{copy.nav.overview}</H2>
            {copy.overview.paragraphs.map((paragraph) => (
              <P key={paragraph}>{paragraph}</P>
            ))}
            <H3>Outputs</H3>
            <Table headers={['Output', 'Format', 'When to use']} rows={copy.overview.outputs} />
            <H3>Typical use cases</H3>
            <BulletList items={copy.overview.useCases} />

            <H2 id="quick-start">{copy.nav.quickStart}</H2>
            <BulletList items={copy.quickStart.steps} />
            <H3>Minimal starter YAML</H3>
            <Pre lang="yaml">{copy.quickStart.yaml}</Pre>
            <Callout title="Recommended first path" tone="tip">
              Start in the hosted editor, generate a PNG once, then switch to the CLI only if you need automation or repository-driven workflows.
            </Callout>

            <H2 id="web-editor">{copy.nav.webEditor}</H2>
            {copy.webEditor.paragraphs.map((paragraph) => (
              <P key={paragraph}>{paragraph}</P>
            ))}
            <H3>Main workflows</H3>
            <BulletList items={copy.webEditor.workflows} />
            <H3>Editor capabilities</H3>
            <Table headers={['Feature', 'Details']} rows={copy.webEditor.features} />

            <H2 id="builder">{copy.nav.builder}</H2>
            {copy.builder.paragraphs.map((paragraph) => (
              <P key={paragraph}>{paragraph}</P>
            ))}
            <H3>When to use the builder</H3>
            <BulletList items={copy.builder.cases} />
            <H3>Builder notes</H3>
            <BulletList items={copy.builder.notes} />

            <H2 id="yaml-reference">{copy.nav.yamlReference}</H2>
            <P>{copy.yamlReference.intro}</P>
            <H3>Core fields</H3>
            <Table headers={['Field', 'Meaning']} rows={copy.yamlReference.concepts} />
            <H3>Common primitives</H3>
            <Table headers={['Type', 'Purpose']} rows={copy.yamlReference.primitives} />
            <H3>YAML best practices</H3>
            <BulletList items={copy.yamlReference.bestPractices} />
            <H3>Full reference example</H3>
            <Pre lang="yaml">{copy.yamlReference.example}</Pre>

            <H2 id="cli">{copy.nav.cli}</H2>
            <P>
              The CLI is the best option when diagrams should live with infrastructure code, be generated in automation, or be versioned as part of engineering workflows.
            </P>
            <H3>Install</H3>
            <Pre lang="bash">{copy.cli.install}</Pre>
            <H3>Basic usage</H3>
            <Pre lang="bash">{copy.cli.usage}</Pre>
            <H3>Main flags</H3>
            <Table headers={['Flag', 'Purpose']} rows={copy.cli.flags} />
            <H3>CLI scenarios</H3>
            <BulletList items={copy.cli.scenarios} />

            <H2 id="api">{copy.nav.api}</H2>
            <P>{copy.api.intro}</P>
            <H3>Request format</H3>
            <Pre lang="http">{copy.api.request}</Pre>
            <H3>Response types</H3>
            <Table headers={['Mode', 'Content-Type', 'Body']} rows={copy.api.responseRows} />
            <H3>curl example</H3>
            <Pre lang="bash">{copy.api.curl}</Pre>
            <H3>TypeScript example</H3>
            <Pre lang="ts">{copy.api.ts}</Pre>
            <H3>API notes</H3>
            <BulletList items={copy.api.notes} />

            <H2 id="drawio">{copy.nav.drawio}</H2>
            {copy.drawio.paragraphs.map((paragraph) => (
              <P key={paragraph}>{paragraph}</P>
            ))}
            <H3>Export pipeline</H3>
            <BulletList items={copy.drawio.pipeline} />
            <H3>Best draw.io use cases</H3>
            <BulletList items={copy.drawio.useCases} />

            <H2 id="mcp">{copy.nav.mcp}</H2>
            {copy.mcp.paragraphs.map((paragraph) => (
              <P key={paragraph}>{paragraph}</P>
            ))}
            <H3>Install MCP server</H3>
            <Pre lang="bash">{copy.mcp.install}</Pre>
            <H3>Claude Desktop configuration</H3>
            <Pre lang="json">{copy.mcp.config}</Pre>
            <H3>Available tools</H3>
            <Table headers={['Tool', 'Description']} rows={copy.mcp.tools} />
            <H3>Example prompt</H3>
            <Pre lang="text">{copy.mcp.prompt}</Pre>

            <H2 id="local-dev">{copy.nav.localDev}</H2>
            <H3>Prerequisites</H3>
            <Table headers={['Tool', 'Version', 'Purpose']} rows={copy.localDev.prerequisites} />
            <H3>Run locally</H3>
            <Pre lang="bash">{copy.localDev.run}</Pre>
            <H3>Validation commands</H3>
            <Pre lang="bash">{copy.localDev.test}</Pre>
            <H3>Alternative: Vercel dev</H3>
            <Pre lang="bash">{copy.localDev.vercel}</Pre>
            <H3>Development notes</H3>
            <BulletList items={copy.localDev.notes} />

            <H2 id="troubleshooting">{copy.nav.troubleshooting}</H2>
            <Callout title="Troubleshooting strategy" tone="warn">
              When a diagram does not render correctly, reduce the YAML to the smallest working version first, then add resources, groups, and links back incrementally.
            </Callout>
            <BulletList items={copy.troubleshooting.items} />

            <H2 id="examples">{copy.nav.examples}</H2>
            <H3>Common use case patterns</H3>
            <Table headers={['Pattern', 'Best fit']} rows={copy.examples.useCases} />
            <H3>Serverless worker pipeline</H3>
            <Pre lang="yaml">{copy.examples.sample1}</Pre>
            <H3>Frontend + API + data stack</H3>
            <Pre lang="yaml">{copy.examples.sample2}</Pre>

            <H2 id="architecture">{copy.nav.architecture}</H2>
            {copy.architecture.paragraphs.map((paragraph) => (
              <P key={paragraph}>{paragraph}</P>
            ))}
            <H3>Repository layers</H3>
            <Table headers={['Path', 'Responsibility']} rows={copy.architecture.layers} />
            <H3>Execution flow</H3>
            <BulletList items={copy.architecture.flow} />

            <H2 id="about">{copy.nav.about}</H2>
            <div className="border border-[var(--border)] rounded-2xl p-6 bg-[var(--code-bg)] mt-4">
              <div className="flex items-start gap-4">
                <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#FF9900] to-[#ff6600] flex items-center justify-center flex-shrink-0 text-white font-bold text-lg">
                  FA
                </div>
                <div className="flex-1 min-w-0">
                  <h3 className="text-[var(--text)] font-semibold text-base">Fernando Azevedo</h3>
                  <p className="text-[var(--text-4)] text-xs mb-3">{copy.about.subtitle}</p>
                  <div className="flex items-center gap-3 mb-4 flex-wrap">
                    <a
                      href="https://fernando.moretes.com"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1.5 text-xs text-[#FF9900] hover:text-[#ffb340] transition-colors"
                    >
                      <ExternalLink size={11} />
                      fernando.moretes.com
                    </a>
                    <a
                      href="https://github.com/fernandofatech"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1.5 text-xs text-[var(--text-4)] hover:text-[var(--text-2)] transition-colors"
                    >
                      <Github size={11} />
                      fernandofatech
                    </a>
                    <a
                      href="https://www.linkedin.com/in/fernando-francisco-azevedo/"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1.5 text-xs text-[var(--text-4)] hover:text-[var(--text-2)] transition-colors"
                    >
                      <Linkedin size={11} />
                      LinkedIn
                    </a>
                  </div>
                  <div className="space-y-3 text-xs text-[var(--text-3)] leading-7">
                    {copy.about.bio.map((paragraph) => (
                      <p key={paragraph}>{paragraph}</p>
                    ))}
                  </div>
                </div>
              </div>
            </div>

            <div className="mt-12 pt-6 border-t border-[var(--border)] flex items-center justify-between">
              <span className="text-[10px] text-[var(--text-6)]">{copy.labels.fork}</span>
              <a
                href="https://github.com/fernandofatech/diagram-as-code"
                target="_blank"
                rel="noopener noreferrer"
                className="text-[10px] text-[var(--text-6)] hover:text-[var(--text-4)] transition-colors flex items-center gap-1"
              >
                <Github size={10} />
                {copy.labels.viewGithub}
              </a>
            </div>
          </div>
        </main>
      </div>
    </div>
  )
}

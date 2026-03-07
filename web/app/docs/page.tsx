'use client'

import { useState } from 'react'
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
} from 'lucide-react'

// ── Section IDs ───────────────────────────────────────────────────────────────

type SectionId =
  | 'overview'
  | 'web-editor'
  | 'cli'
  | 'drawio'
  | 'api'
  | 'mcp'
  | 'local-dev'
  | 'examples'
  | 'about'

const NAV: { id: SectionId; label: string; icon: React.ReactNode }[] = [
  { id: 'overview',   label: 'Overview',          icon: <BookOpen size={14} /> },
  { id: 'web-editor', label: 'Web Editor',         icon: <Globe size={14} /> },
  { id: 'cli',        label: 'CLI Usage',          icon: <Terminal size={14} /> },
  { id: 'drawio',     label: 'Draw.io Export',     icon: <Layers size={14} /> },
  { id: 'api',        label: 'API Reference',      icon: <Code2 size={14} /> },
  { id: 'mcp',        label: 'MCP Server',         icon: <Cpu size={14} /> },
  { id: 'local-dev',  label: 'Local Development',  icon: <Terminal size={14} /> },
  { id: 'examples',   label: 'YAML Examples',      icon: <Layers size={14} /> },
  { id: 'about',      label: 'About the Author',   icon: <User size={14} /> },
]

// ── Reusable primitives ───────────────────────────────────────────────────────

function H2({ id, children }: { id: string; children: React.ReactNode }) {
  return (
    <h2
      id={id}
      className="text-xl font-semibold text-[#e5e5e5] mt-12 mb-4 flex items-center gap-2 scroll-mt-20"
    >
      <span className="w-1 h-5 bg-[#FF9900] rounded-full inline-block flex-shrink-0" />
      {children}
    </h2>
  )
}

function H3({ children }: { children: React.ReactNode }) {
  return <h3 className="text-base font-semibold text-[#ccc] mt-6 mb-2">{children}</h3>
}

function P({ children }: { children: React.ReactNode }) {
  return <p className="text-[#999] text-sm leading-relaxed mb-3">{children}</p>
}

function Code({ children }: { children: React.ReactNode }) {
  return (
    <code className="bg-[#1a1a1a] border border-[#2a2a2a] text-[#FF9900] text-xs px-1.5 py-0.5 rounded font-mono">
      {children}
    </code>
  )
}

function Pre({ children, lang }: { children: string; lang?: string }) {
  return (
    <div className="relative my-4">
      {lang && (
        <span className="absolute top-2 right-3 text-[10px] text-[#555] font-mono uppercase">
          {lang}
        </span>
      )}
      <pre className="bg-[#111] border border-[#2a2a2a] rounded-lg p-4 text-xs text-[#ccc] font-mono overflow-x-auto leading-relaxed whitespace-pre">
        {children}
      </pre>
    </div>
  )
}

function Badge({ children, color = 'orange' }: { children: React.ReactNode; color?: 'orange' | 'blue' | 'green' | 'purple' }) {
  const colors = {
    orange: 'bg-[#FF9900]/10 text-[#FF9900] border-[#FF9900]/20',
    blue:   'bg-blue-900/20 text-blue-400 border-blue-800/30',
    green:  'bg-green-900/20 text-green-400 border-green-800/30',
    purple: 'bg-purple-900/20 text-purple-400 border-purple-800/30',
  }
  return (
    <span className={`inline-block border text-[10px] font-mono px-2 py-0.5 rounded ${colors[color]}`}>
      {children}
    </span>
  )
}

function Table({ headers, rows }: { headers: string[]; rows: string[][] }) {
  return (
    <div className="overflow-x-auto my-4">
      <table className="w-full text-xs text-left border-collapse">
        <thead>
          <tr className="border-b border-[#2a2a2a]">
            {headers.map((h) => (
              <th key={h} className="pb-2 pr-6 text-[#666] font-medium uppercase tracking-wider text-[10px]">
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-b border-[#1a1a1a]">
              {row.map((cell, j) => (
                <td key={j} className="py-2.5 pr-6 text-[#aaa] align-top font-mono">
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

function Callout({ type, children }: { type: 'info' | 'tip'; children: React.ReactNode }) {
  const styles = {
    info: 'bg-blue-950/30 border-blue-800/40 text-blue-300',
    tip:  'bg-[#FF9900]/5 border-[#FF9900]/20 text-[#FF9900]/80',
  }
  return (
    <div className={`border rounded-lg px-4 py-3 text-xs my-4 ${styles[type]}`}>
      {children}
    </div>
  )
}

// ── Page ──────────────────────────────────────────────────────────────────────

export default function DocsPage() {
  const [active, setActive] = useState<SectionId>('overview')

  function scrollTo(id: SectionId) {
    setActive(id)
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
  }

  return (
    <div className="flex flex-col h-screen bg-[#0f0f0f] overflow-hidden">

      {/* Header */}
      <header className="flex items-center justify-between px-4 h-12 border-b border-[#2a2a2a] flex-shrink-0">
        <div className="flex items-center gap-3">
          <Link
            href="/"
            className="flex items-center gap-1.5 text-xs text-[#666] hover:text-[#ccc] transition-colors"
          >
            <ArrowLeft size={13} />
            Back to Editor
          </Link>
          <div className="w-px h-4 bg-[#2a2a2a]" />
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 bg-[#FF9900] rounded flex items-center justify-center">
              <svg viewBox="0 0 16 16" fill="white" className="w-3 h-3">
                <rect x="1" y="1" width="6" height="6" rx="1" />
                <rect x="9" y="1" width="6" height="6" rx="1" />
                <rect x="1" y="9" width="6" height="6" rx="1" />
                <rect x="9" y="9" width="6" height="6" rx="1" />
              </svg>
            </div>
            <span className="text-sm font-semibold text-[#e5e5e5]">Documentation</span>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <a
            href="https://github.com/fernandofatech/diagram-as-code"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-1.5 text-xs text-[#666] hover:text-[#ccc] transition-colors"
          >
            <Github size={14} />
            GitHub
          </a>
          <a
            href="https://fernando.moretes.com"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-1.5 text-xs text-[#FF9900] hover:text-[#ffb340] transition-colors"
          >
            <ExternalLink size={12} />
            fernando.moretes.com
          </a>
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">

        {/* Sidebar */}
        <aside className="w-52 flex-shrink-0 border-r border-[#2a2a2a] overflow-y-auto py-4">
          <p className="px-4 text-[10px] text-[#444] uppercase tracking-widest font-medium mb-2">
            Contents
          </p>
          <nav className="space-y-0.5 px-2">
            {NAV.map(({ id, label, icon }) => (
              <button
                key={id}
                onClick={() => scrollTo(id)}
                className={`w-full flex items-center gap-2 px-3 py-2 rounded text-xs text-left transition-colors ${
                  active === id
                    ? 'bg-[#FF9900]/10 text-[#FF9900]'
                    : 'text-[#666] hover:text-[#ccc] hover:bg-[#1a1a1a]'
                }`}
              >
                {icon}
                {label}
                {active === id && <ChevronRight size={10} className="ml-auto" />}
              </button>
            ))}
          </nav>
        </aside>

        {/* Main content */}
        <main className="flex-1 overflow-y-auto px-10 py-8 max-w-3xl">

          {/* ── Overview ──────────────────────────────────────────────────── */}
          <H2 id="overview">Overview</H2>
          <P>
            <strong className="text-[#ccc]">diagram-as-code</strong> is an open-source CLI tool that
            generates AWS architecture diagrams from YAML. Write human-readable YAML, get
            pixel-perfect PNG diagrams or draw.io files — all without touching any GUI diagramming tool.
          </P>
          <P>
            This fork, maintained by{' '}
            <a href="https://fernando.moretes.com" className="text-[#FF9900] hover:underline" target="_blank" rel="noopener noreferrer">
              Fernando Azevedo
            </a>
            , extends the original{' '}
            <a href="https://github.com/awslabs/diagram-as-code" className="text-[#FF9900] hover:underline" target="_blank" rel="noopener noreferrer">
              awslabs/diagram-as-code
            </a>{' '}
            with a browser-based editor hosted on Vercel and a native draw.io export pipeline.
          </P>
          <Callout type="tip">
            Fork repository:{' '}
            <a href="https://github.com/fernandofatech/diagram-as-code" className="underline" target="_blank" rel="noopener noreferrer">
              github.com/fernandofatech/diagram-as-code
            </a>
          </Callout>

          <H3>What this tool produces</H3>
          <Table
            headers={['Output', 'Format', 'Use case']}
            rows={[
              ['PNG diagram', 'image/png', 'Docs, wikis, PRs, CI artifacts'],
              ['draw.io file', 'application/xml', 'Editable diagrams in diagrams.net'],
            ]}
          />

          {/* ── Web Editor ────────────────────────────────────────────────── */}
          <H2 id="web-editor">Web Editor</H2>
          <P>
            The web editor lets you write YAML and generate diagrams directly in the browser.
            No installation needed. Powered by Vercel serverless functions running the same Go
            engine as the CLI.
          </P>
          <H3>How to use it</H3>
          <ol className="list-decimal list-inside space-y-2 text-sm text-[#999] mb-4 ml-2">
            <li>Open the editor at <Link href="/" className="text-[#FF9900] hover:underline">diagram-as-code-ruddy.vercel.app</Link></li>
            <li>Write or paste your YAML in the left panel, or pick an example from the <strong className="text-[#ccc]">Examples</strong> dropdown</li>
            <li>Choose output format: <Badge>PNG</Badge> or <Badge color="blue">draw.io</Badge></li>
            <li>Click <Badge color="orange">⚡ Generate</Badge> or press <Code>Ctrl+Enter</Code></li>
            <li>Download the result with the button in the preview panel</li>
          </ol>
          <Callout type="info">
            The editor uses Monaco (the VS Code engine) with YAML syntax highlighting. draw.io files
            can be opened in <a href="https://app.diagrams.net" className="underline" target="_blank" rel="noopener noreferrer">diagrams.net</a> for further editing.
          </Callout>

          {/* ── CLI Usage ─────────────────────────────────────────────────── */}
          <H2 id="cli">CLI Usage</H2>
          <H3>Install</H3>
          <Pre lang="bash">{`# Go 1.21+
go install github.com/awslabs/diagram-as-code/cmd/awsdac@latest

# macOS (Homebrew)
brew install awsdac`}</Pre>

          <H3>Basic usage</H3>
          <Pre lang="bash">{`# Generate PNG
awsdac examples/alb-ec2.yaml

# Custom output file
awsdac examples/alb-ec2.yaml -o my-diagram.png

# Generate draw.io file
awsdac examples/alb-ec2.yaml --drawio -o output.drawio

# Shorthand: output extension auto-selects format
awsdac examples/alb-ec2.yaml -o output.drawio`}</Pre>

          <H3>All flags</H3>
          <Table
            headers={['Flag', 'Default', 'Description']}
            rows={[
              ['-o, --output', 'output.png', 'Output file name'],
              ['--drawio', 'false', 'Generate draw.io file instead of PNG'],
              ['-f, --force', 'false', 'Overwrite output without confirmation'],
              ['-v, --verbose', 'false', 'Enable verbose logging'],
              ['--width', '0 (no resize)', 'Resize output image width (PNG only)'],
              ['--height', '0 (no resize)', 'Resize output image height (PNG only)'],
              ['-t, --template', 'false', 'Process input as Go text/template'],
              ['-c, --cfn-template', 'false', '[Beta] Create diagram from CloudFormation template'],
              ['--allow-untrusted-definitions', 'false', 'Allow definition files from non-official URLs'],
            ]}
          />

          {/* ── Draw.io Export ────────────────────────────────────────────── */}
          <H2 id="drawio">Draw.io Export</H2>
          <P>
            The draw.io export uses the same resource/link model and layout engine as PNG
            rendering. Resources become <Code>mxCell</Code> nodes; links become edges. Official
            AWS SVG icons are embedded as data URIs so the file is fully self-contained.
          </P>
          <H3>Export pipeline</H3>
          <ol className="list-decimal list-inside space-y-2 text-sm text-[#999] mb-4 ml-2">
            <li>YAML is parsed into the resource graph</li>
            <li>The same layout pass runs (<Code>Scale</Code> + <Code>ZeroAdjust</Code>)</li>
            <li>Children are reordered by link topology (matches PNG ordering)</li>
            <li>Leaf resources get AWS icons embedded as base64 data URIs</li>
            <li>Groups get AWS4 group styles with correct fill/stroke colors</li>
            <li>Links become <Code>mxCell</Code> edges with optional source/target labels</li>
            <li>Output is written as <Code>mxGraphModel</Code> XML</li>
          </ol>
          <Callout type="info">
            Open the generated <Code>.drawio</Code> file at{' '}
            <a href="https://app.diagrams.net" className="underline" target="_blank" rel="noopener noreferrer">app.diagrams.net</a>{' '}
            — all icons, labels and group borders render without internet access.
          </Callout>

          {/* ── API Reference ─────────────────────────────────────────────── */}
          <H2 id="api">API Reference</H2>
          <P>
            The serverless function at <Code>/api/generate</Code> powers the web editor and can
            be called directly from any HTTP client or AI agent.
          </P>

          <H3>POST /api/generate</H3>
          <Table
            headers={['Property', 'Value']}
            rows={[
              ['Method', 'POST'],
              ['Content-Type', 'application/json'],
              ['URL (production)', 'https://diagram-as-code-ruddy.vercel.app/api/generate'],
            ]}
          />

          <H3>Request body</H3>
          <Pre lang="json">{`{
  "yaml": "Diagram:\\n  ...",   // Required. DAC YAML string.
  "format": "png"              // Optional. "png" (default) or "drawio"
}`}</Pre>

          <H3>Response</H3>
          <Table
            headers={['format', 'Content-Type', 'Body']}
            rows={[
              ['png', 'image/png', 'Raw PNG bytes'],
              ['drawio', 'application/xml', 'draw.io XML (mxGraphModel)'],
              ['error', 'application/json', '{"error": "message"}'],
            ]}
          />

          <H3>curl example</H3>
          <Pre lang="bash">{`# Generate PNG
curl -X POST https://diagram-as-code-ruddy.vercel.app/api/generate \\
  -H "Content-Type: application/json" \\
  -d '{"yaml":"Diagram:\\n  DefinitionFiles:\\n    - Type: URL\\n      Url: \\"https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml\\"\\n  Resources:\\n    Canvas:\\n      Type: AWS::Diagram::Canvas\\n      Children:\\n        - MyBucket\\n    MyBucket:\\n      Type: AWS::S3::Bucket"}' \\
  --output diagram.png

# Generate draw.io
curl -X POST https://diagram-as-code-ruddy.vercel.app/api/generate \\
  -H "Content-Type: application/json" \\
  -d '{"yaml":"...","format":"drawio"}' \\
  --output diagram.drawio`}</Pre>

          <H3>JavaScript / TypeScript example</H3>
          <Pre lang="typescript">{`const yaml = \`Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - Bucket
    Bucket:
      Type: AWS::S3::Bucket\`

const res = await fetch('https://diagram-as-code-ruddy.vercel.app/api/generate', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ yaml, format: 'png' }),
})

if (!res.ok) {
  const { error } = await res.json()
  throw new Error(error)
}

const blob = await res.blob()
const url = URL.createObjectURL(blob)
// use url in <img src={url} />`}</Pre>

          {/* ── MCP Server ────────────────────────────────────────────────── */}
          <H2 id="mcp">MCP Server</H2>
          <P>
            The MCP (Model Context Protocol) server lets AI assistants like Claude generate AWS
            architecture diagrams directly from chat — no CLI required.
          </P>

          <H3>Install the MCP server</H3>
          <Pre lang="bash">{`go install github.com/awslabs/diagram-as-code/cmd/awsdac-mcp-server@latest`}</Pre>

          <H3>Configure in Claude Desktop</H3>
          <P>
            Add the following to your <Code>claude_desktop_config.json</Code>:
          </P>
          <Pre lang="json">{`{
  "mcpServers": {
    "diagram-as-code": {
      "command": "awsdac-mcp-server",
      "args": []
    }
  }
}`}</Pre>
          <P>
            On macOS the config file is at{' '}
            <Code>~/Library/Application Support/Claude/claude_desktop_config.json</Code>.
          </P>

          <H3>Available MCP tools</H3>
          <Table
            headers={['Tool', 'Description']}
            rows={[
              ['create_diagram', 'Generate a PNG diagram from a DAC YAML string'],
              ['create_drawio', 'Generate a draw.io XML file from a DAC YAML string'],
            ]}
          />

          <H3>Example prompt for Claude</H3>
          <Pre lang="text">{`Create an AWS architecture diagram showing:
- A VPC with two public subnets
- An Application Load Balancer
- Two EC2 instances behind the ALB
- An Internet Gateway

Export as PNG and save to ~/Desktop/architecture.png`}</Pre>

          <Callout type="info">
            Claude will use the MCP server to invoke the Go diagram engine, applying the same
            layout and icon rules as the CLI tool.
          </Callout>

          {/* ── Local Dev ─────────────────────────────────────────────────── */}
          <H2 id="local-dev">Local Development</H2>
          <H3>Prerequisites</H3>
          <Table
            headers={['Tool', 'Version', 'Install']}
            rows={[
              ['Go', '1.21+', 'go.dev/dl'],
              ['Node.js', '18+', 'nodejs.org'],
              ['npm', '9+', 'included with Node.js'],
            ]}
          />

          <H3>Clone and run the full stack locally</H3>
          <Pre lang="bash">{`git clone https://github.com/fernandofatech/diagram-as-code.git
cd diagram-as-code

# Terminal 1 — Go API server on :8080
go run ./cmd/api-dev

# Terminal 2 — Next.js frontend on :3000 (proxies /api/* → :8080)
cd web && npm install && npm run dev`}</Pre>
          <P>
            Open <Code>http://localhost:3000</Code> — the editor connects to your local Go server
            automatically. No Vercel CLI needed.
          </P>

          <H3>Run tests</H3>
          <Pre lang="bash">{`# Unit tests
go test ./internal/...

# Golden-file integration tests (generates PNGs and compares)
go test ./test/...

# Update golden files after intentional rendering changes
go test ./test/... -update`}</Pre>

          <H3>Build the CLI</H3>
          <Pre lang="bash">{`go build -o awsdac ./cmd/awsdac
./awsdac examples/alb-ec2.yaml`}</Pre>

          <H3>Run with Vercel CLI (full stack)</H3>
          <Pre lang="bash">{`npm i -g vercel
vercel dev   # runs Next.js + Go serverless function together`}</Pre>

          {/* ── Examples ──────────────────────────────────────────────────── */}
          <H2 id="examples">YAML Examples</H2>
          <H3>Minimal — S3 bucket</H3>
          <Pre lang="yaml">{`Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - AWSCloud
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Children:
        - Bucket
    Bucket:
      Type: AWS::S3::Bucket`}</Pre>

          <H3>ALB + EC2 in a VPC</H3>
          <Pre lang="yaml">{`Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - AWSCloud
        - User
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Children:
        - VPC
    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - ALB
        - SubnetStack
      BorderChildren:
        - Position: S
          Resource: IGW
    SubnetStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Subnet1
        - Subnet2
    Subnet1:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance1
    Instance1:
      Type: AWS::EC2::Instance
    Subnet2:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance2
    Instance2:
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
    - Source: ALB
      Target: Instance1
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "HTTP:80"
    - Source: ALB
      Target: Instance2
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "HTTP:80"
    - Source: IGW
      Target: ALB
      TargetArrowHead:
        Type: Open
    - Source: User
      Target: IGW
      TargetArrowHead:
        Type: Open`}</Pre>

          <H3>Key YAML concepts</H3>
          <Table
            headers={['Field', 'Description']}
            rows={[
              ['Diagram.DefinitionFiles', 'URLs or local files defining AWS resource types and icons'],
              ['Resources.<name>.Type', 'AWS resource type (e.g. AWS::EC2::Instance) or diagram element'],
              ['Resources.<name>.Children', 'List of resource names placed inside this container'],
              ['Resources.<name>.Direction', '"vertical" or "horizontal" layout for children'],
              ['Resources.<name>.Preset', 'Named preset from the definition file (e.g. "PublicSubnet")'],
              ['Resources.<name>.Title', 'Override the display label'],
              ['Resources.<name>.BorderChildren', 'Place resources on the border of a container'],
              ['Links[].Source / Target', 'Resource names to connect'],
              ['Links[].Labels.SourceLeft', 'Label on the source side of the edge'],
              ['Links[].TargetArrowHead.Type', '"Open" = open arrow, "none" = no arrow'],
            ]}
          />

          {/* ── About ─────────────────────────────────────────────────────── */}
          <H2 id="about">About the Author</H2>
          <div className="border border-[#2a2a2a] rounded-xl p-6 bg-[#111] mt-4">
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#FF9900] to-[#ff6600] flex items-center justify-center flex-shrink-0 text-white font-bold text-lg">
                FA
              </div>
              <div className="flex-1 min-w-0">
                <h3 className="text-[#e5e5e5] font-semibold text-base">Fernando Azevedo</h3>
                <p className="text-[#666] text-xs mb-3">Senior Solutions Architect · 16+ years global experience</p>
                <div className="flex items-center gap-3 mb-4">
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
                    className="flex items-center gap-1.5 text-xs text-[#666] hover:text-[#ccc] transition-colors"
                  >
                    <Github size={11} />
                    fernandofatech
                  </a>
                  <a
                    href="https://www.linkedin.com/in/fernando-francisco-azevedo/"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-1.5 text-xs text-[#666] hover:text-[#ccc] transition-colors"
                  >
                    <Linkedin size={11} />
                    LinkedIn
                  </a>
                </div>
                <div className="space-y-3 text-xs text-[#888] leading-relaxed">
                  <p>
                    Senior Solutions Architect with 16+ years of global experience delivering
                    impactful, secure, and scalable digital solutions. Specialized in designing
                    scalable, secure, and cost-efficient cloud architectures. Currently at{' '}
                    <strong className="text-[#aaa]">Banco Itaú</strong>, bridging business goals
                    and technology through innovation and best practices.
                  </p>
                  <p>
                    Deep expertise in <strong className="text-[#aaa]">Clean Architecture</strong>,{' '}
                    <strong className="text-[#aaa]">DDD</strong>, <strong className="text-[#aaa]">CQRS</strong>,
                    Event-Driven systems, Microservices patterns, and high-scale architectures.
                    Experienced with Data Mesh, Kafka, EKS, and AWS Well-Architected Framework.
                  </p>
                  <p>
                    Security-first development approach with extensive experience in CI/CD,
                    Infrastructure as Code (Terraform, CDK), GitOps, and Zero Trust Architecture.
                    Compliance with PCI-DSS, ISO 27001, and GDPR.
                  </p>
                </div>
                <div className="flex flex-wrap gap-1.5 mt-4">
                  {[
                    'AWS', 'Clean Architecture', 'DDD', 'CQRS', 'Kafka', 'EKS',
                    'Terraform', 'CDK', 'GitOps', 'Zero Trust', 'PCI-DSS', 'Go',
                  ].map((tag) => (
                    <span
                      key={tag}
                      className="text-[10px] bg-[#1a1a1a] border border-[#2a2a2a] text-[#666] px-2 py-0.5 rounded-full"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </div>

          <div className="mt-12 pt-6 border-t border-[#1a1a1a] flex items-center justify-between">
            <span className="text-[10px] text-[#333]">
              diagram-as-code · Fork by Fernando Azevedo
            </span>
            <a
              href="https://github.com/fernandofatech/diagram-as-code"
              target="_blank"
              rel="noopener noreferrer"
              className="text-[10px] text-[#444] hover:text-[#666] transition-colors flex items-center gap-1"
            >
              <Github size={10} />
              View on GitHub
            </a>
          </div>

        </main>
      </div>
    </div>
  )
}

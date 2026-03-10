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
import { useLanguage } from '@/lib/i18n'
import LanguageSwitcher from '@/components/LanguageSwitcher'
import ThemeSwitcher from '@/components/ThemeSwitcher'

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

// ── Reusable primitives ───────────────────────────────────────────────────────

function H2({ id, children }: { id: string; children: React.ReactNode }) {
  return (
    <h2
      id={id}
      className="text-xl font-semibold text-[var(--text)] mt-12 mb-4 flex items-center gap-2 scroll-mt-20"
    >
      <span className="w-1 h-5 bg-[#FF9900] rounded-full inline-block flex-shrink-0" />
      {children}
    </h2>
  )
}

function H3({ children }: { children: React.ReactNode }) {
  return <h3 className="text-base font-semibold text-[var(--text-2)] mt-6 mb-2">{children}</h3>
}

function P({ children }: { children: React.ReactNode }) {
  return <p className="text-[var(--text-3)] text-sm leading-relaxed mb-3">{children}</p>
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
      <pre className="bg-[var(--code-bg)] border border-[var(--border)] rounded-lg p-4 text-xs text-[var(--text-2)] font-mono overflow-x-auto leading-relaxed whitespace-pre">
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
          <tr className="border-b border-[var(--border)]">
            {headers.map((h) => (
              <th key={h} className="pb-2 pr-6 text-[var(--text-4)] font-medium uppercase tracking-wider text-[10px]">
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-b border-[var(--border)]">
              {row.map((cell, j) => (
                <td key={j} className="py-2.5 pr-6 text-[var(--text-2)] align-top font-mono">
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
  const { t, lang } = useLanguage()
  const [active, setActive] = useState<SectionId>('overview')

  const NAV: { id: SectionId; label: string; icon: React.ReactNode }[] = [
    { id: 'overview',   label: t.navOverview,   icon: <BookOpen size={14} /> },
    { id: 'web-editor', label: t.navWebEditor,  icon: <Globe size={14} /> },
    { id: 'cli',        label: t.navCli,        icon: <Terminal size={14} /> },
    { id: 'drawio',     label: t.navDrawio,     icon: <Layers size={14} /> },
    { id: 'api',        label: t.navApi,        icon: <Code2 size={14} /> },
    { id: 'mcp',        label: t.navMcp,        icon: <Cpu size={14} /> },
    { id: 'local-dev',  label: t.navLocalDev,   icon: <Terminal size={14} /> },
    { id: 'examples',   label: t.navExamples,   icon: <Layers size={14} /> },
    { id: 'about',      label: t.navAbout,      icon: <User size={14} /> },
  ]

  function scrollTo(id: SectionId) {
    setActive(id)
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
  }

  return (
    <div className="flex flex-col h-screen bg-[var(--bg)] overflow-hidden">

      {/* Header */}
      <header className="flex items-center justify-between px-4 h-12 border-b border-[var(--border)] flex-shrink-0">
        <div className="flex items-center gap-3">
          <Link
            href="/"
            className="flex items-center gap-1.5 text-xs text-[var(--text-4)] hover:text-[var(--text-2)] transition-colors"
          >
            <ArrowLeft size={13} />
            {t.backToEditor}
          </Link>
          <div className="w-px h-4 bg-[var(--border)]" />
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 bg-[#FF9900] rounded flex items-center justify-center">
              <svg viewBox="0 0 16 16" fill="white" className="w-3 h-3">
                <rect x="1" y="1" width="6" height="6" rx="1" />
                <rect x="9" y="1" width="6" height="6" rx="1" />
                <rect x="1" y="9" width="6" height="6" rx="1" />
                <rect x="9" y="9" width="6" height="6" rx="1" />
              </svg>
            </div>
            <span className="text-sm font-semibold text-[var(--text)]">{t.documentation}</span>
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
        <aside className="w-52 flex-shrink-0 border-r border-[var(--border)] overflow-y-auto py-4">
          <p className="px-4 text-[10px] text-[var(--text-6)] uppercase tracking-widest font-medium mb-2">
            {t.contents}
          </p>
          <nav className="space-y-0.5 px-2">
            {NAV.map(({ id, label, icon }) => (
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

        {/* Main content */}
        <main className="flex-1 overflow-y-auto px-10 py-8 max-w-3xl">

          {/* ── Overview ──────────────────────────────────────────────────── */}
          <H2 id="overview">{t.navOverview}</H2>
          <P>
            <strong className="text-[var(--text-2)]">diagram-as-code</strong>{' '}
            {t.overviewP1.replace('diagram-as-code is', '').replace('diagram-as-code é', '')}
          </P>
          <P>
            {t.overviewP2Fork}{' '}
            <a href="https://fernando.moretes.com" className="text-[#FF9900] hover:underline" target="_blank" rel="noopener noreferrer">
              Fernando Azevedo
            </a>
            {t.overviewP2Extends}{' '}
            <a href="https://github.com/awslabs/diagram-as-code" className="text-[#FF9900] hover:underline" target="_blank" rel="noopener noreferrer">
              awslabs/diagram-as-code
            </a>{' '}
            {t.overviewP2With}
          </P>
          <Callout type="tip">
            {t.overviewForkRepo}{' '}
            <a href="https://github.com/fernandofatech/diagram-as-code" className="underline" target="_blank" rel="noopener noreferrer">
              github.com/fernandofatech/diagram-as-code
            </a>
          </Callout>

          <H3>{t.whatToolProduces}</H3>
          <Table
            headers={[t.outputCol, t.formatCol, t.useCaseCol]}
            rows={[
              [t.pngDiagram, 'image/png', t.pngUseCase],
              [t.drawioFile, 'application/xml', t.drawioUseCase],
            ]}
          />

          {/* ── Web Editor ────────────────────────────────────────────────── */}
          <H2 id="web-editor">{t.navWebEditor}</H2>
          <P>{t.webEditorP1}</P>
          <H3>{t.howToUse}</H3>
          <ol className="list-decimal list-inside space-y-2 text-sm text-[var(--text-3)] mb-4 ml-2">
            <li>{t.webEditorStep1} <Link href="/" className="text-[#FF9900] hover:underline">diagram-as-code-ruddy.vercel.app</Link></li>
            <li>{t.webEditorStep2} <strong className="text-[var(--text-2)]">{t.examples}</strong> {t.webEditorStep2b}</li>
            <li>{t.webEditorStep3} <Badge>PNG</Badge> {lang === 'pt' ? 'ou' : 'or'} <Badge color="blue">draw.io</Badge></li>
            <li>{t.webEditorStep4} <Badge color="orange">⚡ {t.generate}</Badge> {t.webEditorStep4b} <Code>Ctrl+Enter</Code></li>
            <li>{t.webEditorStep5}</li>
          </ol>
          <Callout type="info">
            {t.webEditorCallout}{' '}
            <a href="https://app.diagrams.net" className="underline" target="_blank" rel="noopener noreferrer">diagrams.net</a>{' '}
            {t.webEditorCalloutSuffix}
          </Callout>

          {/* ── CLI Usage ─────────────────────────────────────────────────── */}
          <H2 id="cli">{t.navCli}</H2>
          <H3>{t.install}</H3>
          <Pre lang="bash">{`# Go 1.21+
go install github.com/fernandofatech/diagram-as-code/cmd/awsdac@latest

# macOS (Homebrew)
brew install awsdac`}</Pre>

          <H3>{t.basicUsage}</H3>
          <Pre lang="bash">{`# Generate PNG
awsdac examples/alb-ec2.yaml

# Custom output file
awsdac examples/alb-ec2.yaml -o my-diagram.png

# Generate draw.io file
awsdac examples/alb-ec2.yaml --drawio -o output.drawio

# Shorthand: output extension auto-selects format
awsdac examples/alb-ec2.yaml -o output.drawio`}</Pre>

          <H3>{t.allFlags}</H3>
          <Table
            headers={[t.flagCol, t.defaultCol, t.descriptionCol]}
            rows={[
              ['-o, --output', 'output.png', t.flagOutput],
              ['--drawio', 'false', t.flagDrawio],
              ['-f, --force', 'false', t.flagForce],
              ['-v, --verbose', 'false', t.flagVerbose],
              ['--width', '0 (no resize)', t.flagWidth],
              ['--height', '0 (no resize)', t.flagHeight],
              ['-t, --template', 'false', t.flagTemplate],
              ['-c, --cfn-template', 'false', t.flagCfn],
              ['--allow-untrusted-definitions', 'false', t.flagUntrusted],
            ]}
          />

          {/* ── Draw.io Export ────────────────────────────────────────────── */}
          <H2 id="drawio">{t.navDrawio}</H2>
          <P>
            {t.drawioP1} <Code>mxCell</Code> {t.drawioP1b}
          </P>
          <H3>{t.exportPipeline}</H3>
          <ol className="list-decimal list-inside space-y-2 text-sm text-[var(--text-3)] mb-4 ml-2">
            <li>{t.drawioStep1}</li>
            <li>{t.drawioStep2} (<Code>Scale</Code> + <Code>ZeroAdjust</Code>)</li>
            <li>{t.drawioStep3}</li>
            <li>{t.drawioStep4}</li>
            <li>{t.drawioStep5}</li>
            <li>{t.drawioStep6} <Code>mxCell</Code> {t.drawioStep6b}</li>
            <li>{t.drawioStep7} <Code>mxGraphModel</Code> {t.drawioStep7b}</li>
          </ol>
          <Callout type="info">
            {t.drawioCallout} <Code>.drawio</Code> {t.drawioCalloutAt}{' '}
            <a href="https://app.diagrams.net" className="underline" target="_blank" rel="noopener noreferrer">app.diagrams.net</a>{' '}
            {t.drawioCalloutSuffix}
          </Callout>

          {/* ── API Reference ─────────────────────────────────────────────── */}
          <H2 id="api">{t.navApi}</H2>
          <P>
            {t.apiP1} <Code>/api/generate</Code> {t.apiP1b}
          </P>

          <H3>POST /api/generate</H3>
          <Table
            headers={[t.propertyCol, t.valueCol]}
            rows={[
              ['Method', 'POST'],
              ['Content-Type', 'application/json'],
              ['URL (production)', 'https://diagram-as-code-ruddy.vercel.app/api/generate'],
            ]}
          />

          <H3>{t.requestBody}</H3>
          <Pre lang="json">{`{
  "yaml": "Diagram:\\n  ...",   // Required. DAC YAML string.
  "format": "png"              // Optional. "png" (default) or "drawio"
}`}</Pre>

          <H3>{t.response}</H3>
          <Table
            headers={['format', 'Content-Type', t.bodyCol]}
            rows={[
              ['png', 'image/png', 'Raw PNG bytes'],
              ['drawio', 'application/xml', 'draw.io XML (mxGraphModel)'],
              ['error', 'application/json', '{"error": "message"}'],
            ]}
          />

          <H3>{t.curlExample}</H3>
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

          <H3>{t.jsExample}</H3>
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
          <H2 id="mcp">{t.navMcp}</H2>
          <P>{t.mcpP1}</P>

          <H3>{t.installMcp}</H3>
          <Pre lang="bash">{`go install github.com/awslabs/diagram-as-code/cmd/awsdac-mcp-server@latest`}</Pre>

          <H3>{t.configureClaude}</H3>
          <P>
            {t.configureClaudeP} <Code>claude_desktop_config.json</Code>:
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
            {t.configureClaudeP2}{' '}
            <Code>~/Library/Application Support/Claude/claude_desktop_config.json</Code>.
          </P>

          <H3>{t.availableMcpTools}</H3>
          <Table
            headers={[t.toolCol, t.descriptionCol]}
            rows={[
              ['create_diagram', t.mcpTool1Desc],
              ['create_drawio', t.mcpTool2Desc],
            ]}
          />

          <H3>{t.examplePrompt}</H3>
          <Pre lang="text">{`Create an AWS architecture diagram showing:
- A VPC with two public subnets
- An Application Load Balancer
- Two EC2 instances behind the ALB
- An Internet Gateway

Export as PNG and save to ~/Desktop/architecture.png`}</Pre>

          <Callout type="info">{t.mcpCallout}</Callout>

          {/* ── Local Dev ─────────────────────────────────────────────────── */}
          <H2 id="local-dev">{t.navLocalDev}</H2>
          <H3>{t.prerequisites}</H3>
          <Table
            headers={[t.toolCol, t.versionCol, t.install]}
            rows={[
              ['Go', '1.21+', 'go.dev/dl'],
              ['Node.js', '18+', 'nodejs.org'],
              ['npm', '9+', 'included with Node.js'],
            ]}
          />

          <H3>{t.cloneAndRun}</H3>
          <Pre lang="bash">{`git clone https://github.com/fernandofatech/diagram-as-code.git
cd diagram-as-code

# Terminal 1 — Go API server on :8080
go run ./cmd/api-dev

# Terminal 2 — Next.js frontend on :3001 (proxies /api/* → :8080)
cd web && npm install && npm run dev`}</Pre>
          <P>
            {t.localDevP} <Code>http://localhost:3001</Code> {t.localDevPSuffix}
          </P>

          <H3>{t.runTests}</H3>
          <Pre lang="bash">{`# Unit tests
go test ./internal/...

# Golden-file integration tests (generates PNGs and compares)
go test ./test/...

# Update golden files after intentional rendering changes
go test ./test/... -update`}</Pre>

          <H3>{t.buildCli}</H3>
          <Pre lang="bash">{`go build -o awsdac ./cmd/awsdac
./awsdac examples/alb-ec2.yaml`}</Pre>

          <H3>{t.runWithVercel}</H3>
          <Pre lang="bash">{`npm i -g vercel
vercel dev   # runs Next.js + Go serverless function together`}</Pre>

          {/* ── Examples ──────────────────────────────────────────────────── */}
          <H2 id="examples">{t.navExamples}</H2>
          <H3>{t.minimalS3}</H3>
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

          <H3>{t.albEc2}</H3>
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

          <H3>{t.keyYamlConcepts}</H3>
          <Table
            headers={[t.fieldCol, t.descriptionCol]}
            rows={[
              ['Diagram.DefinitionFiles', t.yamlConcept1],
              ['Resources.<name>.Type', t.yamlConcept2],
              ['Resources.<name>.Children', t.yamlConcept3],
              ['Resources.<name>.Direction', t.yamlConcept4],
              ['Resources.<name>.Preset', t.yamlConcept5],
              ['Resources.<name>.Title', t.yamlConcept6],
              ['Resources.<name>.BorderChildren', t.yamlConcept7],
              ['Links[].Source / Target', t.yamlConcept8],
              ['Links[].Labels.SourceLeft', t.yamlConcept9],
              ['Links[].TargetArrowHead.Type', t.yamlConcept10],
            ]}
          />

          {/* ── About ─────────────────────────────────────────────────────── */}
          <H2 id="about">{t.navAbout}</H2>
          <div className="border border-[var(--border)] rounded-xl p-6 bg-[var(--code-bg)] mt-4">
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#FF9900] to-[#ff6600] flex items-center justify-center flex-shrink-0 text-white font-bold text-lg">
                FA
              </div>
              <div className="flex-1 min-w-0">
                <h3 className="text-[var(--text)] font-semibold text-base">Fernando Azevedo</h3>
                <p className="text-[var(--text-4)] text-xs mb-3">{t.authorSubtitle}</p>
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
                <div className="space-y-3 text-xs text-[var(--text-3)] leading-relaxed">
                  <p>
                    {t.authorBio1}{' '}
                    <strong className="text-[var(--text-2)]">Banco Itaú</strong>
                    {t.authorBio1b}
                  </p>
                  <p>
                    {t.authorBio2}{' '}
                    <strong className="text-[var(--text-2)]">Clean Architecture</strong>,{' '}
                    <strong className="text-[var(--text-2)]">DDD</strong>,{' '}
                    <strong className="text-[var(--text-2)]">CQRS</strong>
                    {t.authorBio2b}
                  </p>
                  <p>{t.authorBio3}</p>
                </div>
                <div className="flex flex-wrap gap-1.5 mt-4">
                  {[
                    'AWS', 'Clean Architecture', 'DDD', 'CQRS', 'Kafka', 'EKS',
                    'Terraform', 'CDK', 'GitOps', 'Zero Trust', 'PCI-DSS', 'Go',
                  ].map((tag) => (
                    <span
                      key={tag}
                      className="text-[10px] bg-[var(--surface)] border border-[var(--border)] text-[var(--text-4)] px-2 py-0.5 rounded-full"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </div>

          <div className="mt-12 pt-6 border-t border-[var(--border)] flex items-center justify-between">
            <span className="text-[10px] text-[var(--text-6)]">{t.footerFork}</span>
            <a
              href="https://github.com/fernandofatech/diagram-as-code"
              target="_blank"
              rel="noopener noreferrer"
              className="text-[10px] text-[var(--text-6)] hover:text-[var(--text-4)] transition-colors flex items-center gap-1"
            >
              <Github size={10} />
              {t.viewOnGithub}
            </a>
          </div>

        </main>
      </div>
    </div>
  )
}

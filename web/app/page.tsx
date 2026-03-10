'use client'

import { useState, useCallback, useEffect, useRef } from 'react'
import dynamic from 'next/dynamic'
import Link from 'next/link'
import { ChevronDown, Zap, Github, BookOpen, Wrench } from 'lucide-react'
import DiagramPreview from '@/components/DiagramPreview'
import LanguageSwitcher from '@/components/LanguageSwitcher'
import { useLanguage } from '@/lib/i18n'

const YamlEditor = dynamic(() => import('@/components/YamlEditor'), {
  ssr: false,
  loading: () => (
    <div className="flex items-center justify-center h-full text-[#444] text-sm">
      Loading editor…
    </div>
  ),
})

// ── Example templates ────────────────────────────────────────────────────────

const EXAMPLES: Record<string, string> = {
  'ALB + EC2': `Diagram:
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
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - VPC
    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - VPCPublicStack
        - ALB
      BorderChildren:
        - Position: S
          Resource: IGW
    VPCPublicStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - VPCPublicSubnet1
        - VPCPublicSubnet2
    VPCPublicSubnet1:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - VPCPublicSubnet1Instance
    VPCPublicSubnet1Instance:
      Type: AWS::EC2::Instance
    VPCPublicSubnet2:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - VPCPublicSubnet2Instance
    VPCPublicSubnet2Instance:
      Type: AWS::EC2::Instance
    ALB:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer
    IGW:
      Type: AWS::EC2::InternetGateway
      IconFill:
        Type: rect
    User:
      Type: AWS::Diagram::Resource
      Preset: User
  Links:
    - Source: ALB
      SourcePosition: NNW
      Target: VPCPublicSubnet1Instance
      TargetPosition: SSE
      TargetArrowHead:
        Type: Open
    - Source: ALB
      SourcePosition: NNE
      Target: VPCPublicSubnet2Instance
      TargetPosition: SSW
      TargetArrowHead:
        Type: Open
    - Source: IGW
      SourcePosition: N
      Target: ALB
      TargetPosition: S
      TargetArrowHead:
        Type: Open
    - Source: User
      SourcePosition: N
      Target: IGW
      TargetPosition: S
      TargetArrowHead:
        Type: Open
`,

  'VPC + NAT Gateway': `Diagram:
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
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - VPC
    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - VPCPublicSubnetStack
        - VPCPrivateSubnetStack
    VPCPublicSubnetStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - VPCPublicSubnet1
        - VPCPublicSubnet2
    VPCPublicSubnet1:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - VPCPublicSubnet1NatGateway
    VPCPublicSubnet1NatGateway:
      Type: AWS::EC2::NatGateway
    VPCPublicSubnet2:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - VPCPublicSubnet2NatGateway
    VPCPublicSubnet2NatGateway:
      Type: AWS::EC2::NatGateway
    VPCPrivateSubnetStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - VPCPrivateSubnet1
        - VPCPrivateSubnet2
    VPCPrivateSubnet1:
      Type: AWS::EC2::Subnet
      Preset: PrivateSubnet
      Children:
        - VPCPrivateSubnet1Instance
    VPCPrivateSubnet1Instance:
      Type: AWS::EC2::Instance
    VPCPrivateSubnet2:
      Type: AWS::EC2::Subnet
      Preset: PrivateSubnet
      Children:
        - VPCPrivateSubnet2Instance
    VPCPrivateSubnet2Instance:
      Type: AWS::EC2::Instance
    User:
      Type: AWS::Diagram::Resource
      Preset: User
`,

  'ALB + Auto Scaling': `Diagram:
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
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - VPC
    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - AutoScalingGroup
        - ALB
      BorderChildren:
        - Position: S
          Resource: IGW
    AutoScalingGroup:
      Type: AWS::AutoScaling::AutoScalingGroup
      Children:
        - Instance1
        - Instance2
    Instance1:
      Type: AWS::EC2::Instance
    Instance2:
      Type: AWS::EC2::Instance
    ALB:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer
    IGW:
      Type: AWS::EC2::InternetGateway
      IconFill:
        Type: rect
    User:
      Type: AWS::Diagram::Resource
      Preset: User
  Links:
    - Source: ALB
      SourcePosition: NNW
      Target: Instance1
      TargetPosition: S
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "HTTP:80"
    - Source: ALB
      SourcePosition: NNE
      Target: Instance2
      TargetPosition: S
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "HTTP:80"
    - Source: IGW
      SourcePosition: N
      Target: ALB
      TargetPosition: S
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "HTTP:80"
        SourceRight:
          Title: "HTTPS:443"
    - Source: User
      SourcePosition: N
      Target: IGW
      TargetPosition: S
      TargetArrowHead:
        Type: Open
`,

  'Multi-Region': `Diagram:
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
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - Regions
    Regions:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - UsEast1
        - UsWest2
    UsEast1:
      Type: AWS::Region
      Title: us-east-1
      Direction: vertical
      Children:
        - VpcEast
    VpcEast:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - SubnetStackEast
        - ALBEast
      BorderChildren:
        - Position: S
          Resource: IGWEast
    SubnetStackEast:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Subnet1East
        - Subnet2East
    Subnet1East:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance1East
    Instance1East:
      Type: AWS::EC2::Instance
    Subnet2East:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance2East
    Instance2East:
      Type: AWS::EC2::Instance
    ALBEast:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer
    IGWEast:
      Type: AWS::EC2::InternetGateway
      IconFill:
        Type: rect
    UsWest2:
      Type: AWS::Region
      Title: us-west-2
      Direction: vertical
      Children:
        - VpcWest
    VpcWest:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - SubnetStackWest
        - ALBWest
      BorderChildren:
        - Position: S
          Resource: IGWWest
    SubnetStackWest:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Subnet1West
        - Subnet2West
    Subnet1West:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance1West
    Instance1West:
      Type: AWS::EC2::Instance
    Subnet2West:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance2West
    Instance2West:
      Type: AWS::EC2::Instance
    ALBWest:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer
    IGWWest:
      Type: AWS::EC2::InternetGateway
      IconFill:
        Type: rect
    User:
      Type: AWS::Diagram::Resource
      Preset: User
  Links:
    - Source: ALBEast
      Target: Instance1East
      TargetArrowHead:
        Type: Open
    - Source: ALBEast
      Target: Instance2East
      TargetArrowHead:
        Type: Open
    - Source: IGWEast
      Target: ALBEast
      TargetArrowHead:
        Type: Open
    - Source: User
      Target: IGWEast
      TargetArrowHead:
        Type: Open
    - Source: ALBWest
      Target: Instance1West
      TargetArrowHead:
        Type: Open
    - Source: ALBWest
      Target: Instance2West
      TargetArrowHead:
        Type: Open
    - Source: IGWWest
      Target: ALBWest
      TargetArrowHead:
        Type: Open
    - Source: User
      Target: IGWWest
      TargetArrowHead:
        Type: Open
`,
}

const EXAMPLE_NAMES = Object.keys(EXAMPLES)

// ── Page component ───────────────────────────────────────────────────────────

export default function Home() {
  const { t } = useLanguage()
  const [yaml, setYaml] = useState(EXAMPLES['ALB + EC2'])
  const [format, setFormat] = useState<'png' | 'drawio'>('png')
  const [imageUrl, setImageUrl] = useState<string | null>(null)
  const [drawioContent, setDrawioContent] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [examplesOpen, setExamplesOpen] = useState(false)
  const prevImageUrl = useRef<string | null>(null)

  // Pick up YAML from builder
  useEffect(() => {
    const saved = localStorage.getItem('builder_yaml')
    if (saved) {
      setYaml(saved)
      localStorage.removeItem('builder_yaml')
    }
  }, [])

  // Ctrl+Enter shortcut
  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault()
        generate()
      }
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [yaml, format])

  const generate = useCallback(async () => {
    if (!yaml.trim()) return
    setLoading(true)
    setError(null)

    try {
      const res = await fetch('/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ yaml, format }),
      })

      if (!res.ok) {
        const body = await res.json().catch(() => ({ error: 'Unknown error' }))
        throw new Error(body.error ?? 'Generation failed')
      }

      if (format === 'png') {
        const blob = await res.blob()
        const url = URL.createObjectURL(blob)
        if (prevImageUrl.current) URL.revokeObjectURL(prevImageUrl.current)
        prevImageUrl.current = url
        setImageUrl(url)
        setDrawioContent(null)
      } else {
        const text = await res.text()
        setDrawioContent(text)
        setImageUrl(null)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
    } finally {
      setLoading(false)
    }
  }, [yaml, format])

  function loadExample(name: string) {
    setYaml(EXAMPLES[name])
    setExamplesOpen(false)
    setImageUrl(null)
    setDrawioContent(null)
    setError(null)
  }

  return (
    <div className="flex flex-col h-screen bg-[#0f0f0f] overflow-hidden">

      {/* ── Header ─────────────────────────────────────────────────────────── */}
      <header className="flex items-center justify-between px-4 h-12 border-b border-[#2a2a2a] flex-shrink-0">
        <div className="flex items-center gap-3">
          {/* Logo */}
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 bg-[#FF9900] rounded flex items-center justify-center">
              <svg viewBox="0 0 16 16" fill="white" className="w-4 h-4">
                <rect x="1" y="1" width="6" height="6" rx="1" />
                <rect x="9" y="1" width="6" height="6" rx="1" />
                <rect x="1" y="9" width="6" height="6" rx="1" />
                <rect x="9" y="9" width="6" height="6" rx="1" />
              </svg>
            </div>
            <span className="text-sm font-semibold text-[#e5e5e5] tracking-tight">
              diagram-as-code
            </span>
          </div>

          {/* Divider */}
          <div className="w-px h-4 bg-[#2a2a2a]" />

          {/* Examples dropdown */}
          <div className="relative">
            <button
              onClick={() => setExamplesOpen((o) => !o)}
              className="flex items-center gap-1.5 text-xs text-[#999] hover:text-[#e5e5e5] transition-colors px-2 py-1 rounded hover:bg-[#1a1a1a]"
            >
              {t.examples}
              <ChevronDown size={12} className={`transition-transform ${examplesOpen ? 'rotate-180' : ''}`} />
            </button>
            {examplesOpen && (
              <>
                <div
                  className="fixed inset-0 z-10"
                  onClick={() => setExamplesOpen(false)}
                />
                <div className="absolute top-full left-0 mt-1 bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg shadow-xl z-20 min-w-[180px] py-1 overflow-hidden">
                  {EXAMPLE_NAMES.map((name) => (
                    <button
                      key={name}
                      onClick={() => loadExample(name)}
                      className="w-full text-left px-3 py-2 text-xs text-[#ccc] hover:bg-[#252525] hover:text-[#e5e5e5] transition-colors"
                    >
                      {name}
                    </button>
                  ))}
                </div>
              </>
            )}
          </div>
        </div>

        {/* Right controls */}
        <div className="flex items-center gap-2">
          {/* Language switcher */}
          <LanguageSwitcher />

          {/* Format toggle */}
          <div className="flex items-center bg-[#1a1a1a] border border-[#2a2a2a] rounded-md p-0.5 text-xs">
            {(['png', 'drawio'] as const).map((f) => (
              <button
                key={f}
                onClick={() => setFormat(f)}
                className={`px-3 py-1 rounded transition-all ${
                  format === f
                    ? 'bg-[#FF9900] text-[#0f0f0f] font-semibold'
                    : 'text-[#888] hover:text-[#ccc]'
                }`}
              >
                {f === 'png' ? 'PNG' : 'draw.io'}
              </button>
            ))}
          </div>

          {/* Generate button */}
          <button
            onClick={generate}
            disabled={loading || !yaml.trim()}
            className="flex items-center gap-1.5 px-4 py-1.5 bg-[#FF9900] hover:bg-[#ffb340] disabled:opacity-40 disabled:cursor-not-allowed text-[#0f0f0f] text-xs font-semibold rounded-md transition-colors"
          >
            {loading ? (
              <div className="w-3.5 h-3.5 border-2 border-[#0f0f0f] border-t-transparent rounded-full animate-spin" />
            ) : (
              <Zap size={13} />
            )}
            {t.generate}
          </button>

          {/* Builder link */}
          <Link
            href="/builder"
            className="flex items-center gap-1.5 text-xs text-[#555] hover:text-[#999] transition-colors px-2 py-1 rounded hover:bg-[#1a1a1a]"
          >
            <Wrench size={13} />
            {t.builderLink}
          </Link>

          {/* Docs link */}
          <Link
            href="/docs"
            className="flex items-center gap-1.5 text-xs text-[#555] hover:text-[#999] transition-colors px-2 py-1 rounded hover:bg-[#1a1a1a]"
          >
            <BookOpen size={13} />
            {t.docs}
          </Link>

          {/* GitHub link */}
          <a
            href="https://github.com/fernandofatech/diagram-as-code"
            target="_blank"
            rel="noopener noreferrer"
            className="text-[#555] hover:text-[#999] transition-colors p-1"
            aria-label="GitHub"
          >
            <Github size={16} />
          </a>
        </div>
      </header>

      {/* ── Main content ────────────────────────────────────────────────────── */}
      <main className="flex flex-1 overflow-hidden">
        {/* Editor panel */}
        <div className="w-1/2 flex flex-col overflow-hidden border-r border-[#2a2a2a]">
          <div className="flex items-center justify-between px-4 py-2 border-b border-[#2a2a2a] flex-shrink-0">
            <span className="text-xs text-[#555] font-medium uppercase tracking-wider">
              {t.yamlEditor}
            </span>
            <span className="text-xs text-[#444]">
              {t.lines(yaml.split('\n').length)}
            </span>
          </div>
          <div className="flex-1 overflow-hidden">
            <YamlEditor value={yaml} onChange={setYaml} />
          </div>
        </div>

        {/* Preview panel */}
        <div className="w-1/2 flex flex-col overflow-hidden">
          <div className="flex items-center px-4 py-2 border-b border-[#2a2a2a] flex-shrink-0">
            <span className="text-xs text-[#555] font-medium uppercase tracking-wider">
              {t.preview}
            </span>
          </div>
          <div className="flex-1 overflow-hidden">
            <DiagramPreview
              imageUrl={imageUrl}
              drawioContent={drawioContent}
              loading={loading}
              error={error}
            />
          </div>
        </div>
      </main>

      {/* ── Status bar ─────────────────────────────────────────────────────── */}
      <footer className="flex items-center justify-between px-4 h-6 border-t border-[#2a2a2a] flex-shrink-0">
        <span className="text-[10px] text-[#444]">
          diagram-as-code · AWS Architecture Diagrams from YAML ·{' '}
          <a
            href="https://fernando.moretes.com"
            target="_blank"
            rel="noopener noreferrer"
            className="text-[#555] hover:text-[#888] transition-colors"
          >
            {t.footerBy}
          </a>
        </span>
        <span className="text-[10px] text-[#444]">
          {t.mode(format)} · <kbd className="font-mono">Ctrl+Enter</kbd>
        </span>
      </footer>
    </div>
  )
}

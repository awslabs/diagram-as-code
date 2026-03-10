'use client'

import { useState, useCallback, useEffect, useRef } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import {
  ArrowLeft, Copy, Check, Zap, ChevronDown, ChevronRight,
  Trash2, Plus, Github, BookOpen
} from 'lucide-react'
import { useLanguage } from '@/lib/i18n'
import LanguageSwitcher from '@/components/LanguageSwitcher'

// ── AWS Resource Catalog ──────────────────────────────────────────────────────

const AWS_TYPES: Record<string, string[]> = {
  'Diagram Elements': [
    'AWS::Diagram::Canvas',
    'AWS::Diagram::Cloud',
    'AWS::Diagram::VerticalStack',
    'AWS::Diagram::HorizontalStack',
    'AWS::Diagram::Resource',
    'AWS::Region',
  ],
  'Compute': [
    'AWS::EC2::Instance',
    'AWS::AutoScaling::AutoScalingGroup',
    'AWS::Lambda::Function',
    'AWS::ECS::Cluster',
    'AWS::ECS::Service',
    'AWS::ECS::TaskDefinition',
    'AWS::Batch::ComputeEnvironment',
    'AWS::AppRunner::Service',
    'AWS::Lightsail::Instance',
    'AWS::Outposts::Outpost',
  ],
  'Networking': [
    'AWS::EC2::VPC',
    'AWS::EC2::Subnet',
    'AWS::EC2::InternetGateway',
    'AWS::EC2::NatGateway',
    'AWS::EC2::SecurityGroup',
    'AWS::EC2::VPNGateway',
    'AWS::EC2::TransitGateway',
    'AWS::EC2::VPCPeeringConnection',
    'AWS::ElasticLoadBalancingV2::LoadBalancer',
    'AWS::ElasticLoadBalancingV2::TargetGroup',
    'AWS::Route53::HostedZone',
    'AWS::Route53::RecordSet',
    'AWS::CloudFront::Distribution',
    'AWS::ApiGateway::RestApi',
    'AWS::ApiGatewayV2::Api',
    'AWS::NetworkFirewall::Firewall',
    'AWS::GlobalAccelerator::Accelerator',
    'AWS::DirectConnect::Connection',
    'AWS::PrivateLink::VpcEndpoint',
  ],
  'Storage': [
    'AWS::S3::Bucket',
    'AWS::EFS::FileSystem',
    'AWS::FSx::FileSystem',
    'AWS::Backup::BackupVault',
    'AWS::StorageGateway::Gateway',
    'AWS::S3Outposts::Bucket',
  ],
  'Database': [
    'AWS::RDS::DBInstance',
    'AWS::RDS::DBCluster',
    'AWS::DynamoDB::Table',
    'AWS::ElastiCache::ReplicationGroup',
    'AWS::ElastiCache::CacheCluster',
    'AWS::Redshift::Cluster',
    'AWS::Neptune::DBCluster',
    'AWS::QLDB::Ledger',
    'AWS::DocumentDB::DBCluster',
    'AWS::Keyspaces::Table',
    'AWS::Timestream::Database',
    'AWS::MemoryDB::Cluster',
  ],
  'Messaging & Integration': [
    'AWS::SQS::Queue',
    'AWS::SNS::Topic',
    'AWS::Kinesis::Stream',
    'AWS::KinesisFirehose::DeliveryStream',
    'AWS::MSK::Cluster',
    'AWS::Events::EventBus',
    'AWS::EventSchemas::Registry',
    'AWS::StepFunctions::StateMachine',
    'AWS::AppSync::GraphQLApi',
    'AWS::MQ::Broker',
  ],
  'Security & Identity': [
    'AWS::IAM::Role',
    'AWS::IAM::Policy',
    'AWS::KMS::Key',
    'AWS::SecretsManager::Secret',
    'AWS::WAFv2::WebACL',
    'AWS::Cognito::UserPool',
    'AWS::Cognito::IdentityPool',
    'AWS::ACM::Certificate',
    'AWS::Shield::Protection',
    'AWS::SSO::PermissionSet',
    'AWS::Inspector::AssessmentTarget',
    'AWS::GuardDuty::Detector',
    'AWS::SecurityHub::Hub',
    'AWS::Macie::Session',
  ],
  'Developer Tools': [
    'AWS::CodePipeline::Pipeline',
    'AWS::CodeBuild::Project',
    'AWS::CodeDeploy::Application',
    'AWS::CodeCommit::Repository',
    'AWS::ECR::Repository',
    'AWS::Cloud9::EnvironmentEC2',
    'AWS::Amplify::App',
    'AWS::XRay::Group',
  ],
  'Monitoring & Management': [
    'AWS::CloudWatch::Alarm',
    'AWS::CloudWatch::Dashboard',
    'AWS::CloudTrail::Trail',
    'AWS::Config::ConfigRule',
    'AWS::SSM::Parameter',
    'AWS::Organizations::Account',
    'AWS::ServiceCatalog::Portfolio',
    'AWS::OpsWorks::Stack',
    'AWS::SystemsManager::Document',
  ],
  'Analytics': [
    'AWS::Athena::WorkGroup',
    'AWS::Glue::Database',
    'AWS::Glue::Job',
    'AWS::EMR::Cluster',
    'AWS::QuickSight::Dashboard',
    'AWS::DataSync::Task',
    'AWS::LakeFormation::DataLake',
    'AWS::OpenSearchService::Domain',
  ],
  'AI & Machine Learning': [
    'AWS::SageMaker::NotebookInstance',
    'AWS::SageMaker::Endpoint',
    'AWS::SageMaker::Pipeline',
    'AWS::Rekognition::Collection',
    'AWS::Comprehend::DocumentClassifier',
    'AWS::Textract::DocumentAnalysis',
    'AWS::Bedrock::Agent',
    'AWS::Lex::Bot',
    'AWS::Polly::Lexicon',
    'AWS::Translate::Terminology',
  ],
  'IoT': [
    'AWS::IoT::TopicRule',
    'AWS::IoT::Thing',
    'AWS::IoTEvents::DetectorModel',
    'AWS::Greengrass::Group',
  ],
}

const PRESETS: Record<string, string[]> = {
  'AWS::ElasticLoadBalancingV2::LoadBalancer': [
    'Application Load Balancer',
    'Network Load Balancer',
    'Gateway Load Balancer',
  ],
  'AWS::EC2::Subnet': ['PublicSubnet', 'PrivateSubnet'],
  'AWS::Diagram::Cloud': ['AWSCloudNoLogo', 'AWSCloud'],
  'AWS::Diagram::Resource': [
    'User', 'Users', 'Mobile client', 'Internet',
    'Corporate data center', 'Traditional server',
  ],
  'AWS::EC2::Instance': ['EC2 Instance Contents'],
  'AWS::S3::Bucket': ['Bucket with objects'],
}

const POSITIONS = [
  'N', 'NNE', 'NE', 'ENE', 'E', 'ESE', 'SE', 'SSE',
  'S', 'SSW', 'SW', 'WSW', 'W', 'WNW', 'NW', 'NNW', 'auto',
] as const

const ARROW_HEADS = ['Open', 'Default', 'none'] as const
const ICON_FILLS = ['rect', 'circle', 'none'] as const
const ALIGNS = ['left', 'center', 'right'] as const

// ── Data types ────────────────────────────────────────────────────────────────

interface BorderChild {
  id: string
  position: string
  resource: string
}

interface ResourceEntry {
  id: string
  name: string
  type: string
  title: string
  preset: string
  direction: string
  align: string
  iconFill: string
  children: string[]
  borderChildren: BorderChild[]
}

interface EdgeLabel {
  SourceLeft: string
  SourceRight: string
  TargetLeft: string
  TargetRight: string
  AutoLeft: string
  AutoRight: string
}

interface LinkEntry {
  id: string
  source: string
  sourcePosition: string
  target: string
  targetPosition: string
  arrowHead: string
  labels: EdgeLabel
}

interface FormState {
  defType: 'URL' | 'LocalFile'
  defUrl: string
  defFile: string
  canvasDirection: string
  resources: ResourceEntry[]
  links: LinkEntry[]
}

// ── YAML Generator ────────────────────────────────────────────────────────────

function generateYaml(form: FormState): string {
  const lines: string[] = ['Diagram:']

  // Definition files
  lines.push('  DefinitionFiles:')
  lines.push(`    - Type: ${form.defType}`)
  if (form.defType === 'URL') {
    lines.push(`      Url: "${form.defUrl}"`)
  } else {
    lines.push(`      LocalFile: "${form.defFile}"`)
  }

  // Resources
  lines.push('  Resources:')

  // Canvas (auto-generated — children = resources not nested in others)
  const allChildren = new Set<string>()
  form.resources.forEach(r => {
    r.children.forEach(c => allChildren.add(c))
    r.borderChildren.forEach(bc => allChildren.add(bc.resource))
  })
  const topLevel = form.resources.filter(r => !allChildren.has(r.name))

  lines.push('    Canvas:')
  lines.push('      Type: AWS::Diagram::Canvas')
  if (form.canvasDirection) lines.push(`      Direction: ${form.canvasDirection}`)
  if (topLevel.length > 0) {
    lines.push('      Children:')
    topLevel.forEach(r => lines.push(`        - ${r.name}`))
  }

  // User resources
  form.resources.forEach(r => {
    if (!r.name.trim()) return
    lines.push(`    ${r.name}:`)
    lines.push(`      Type: ${r.type}`)
    if (r.title.trim()) lines.push(`      Title: "${r.title}"`)
    if (r.preset) lines.push(`      Preset: ${r.preset}`)
    if (r.direction) lines.push(`      Direction: ${r.direction}`)
    if (r.align) lines.push(`      Align: ${r.align}`)
    if (r.iconFill) {
      lines.push('      IconFill:')
      lines.push(`        Type: ${r.iconFill}`)
    }
    const validChildren = r.children.filter(c => c && form.resources.some(x => x.name === c))
    if (validChildren.length > 0) {
      lines.push('      Children:')
      validChildren.forEach(c => lines.push(`        - ${c}`))
    }
    if (r.borderChildren.length > 0) {
      lines.push('      BorderChildren:')
      r.borderChildren.forEach(bc => {
        if (!bc.resource) return
        lines.push(`        - Position: ${bc.position}`)
        lines.push(`          Resource: ${bc.resource}`)
      })
    }
  })

  // Links
  const validLinks = form.links.filter(l => l.source && l.target)
  if (validLinks.length > 0) {
    lines.push('  Links:')
    validLinks.forEach(l => {
      lines.push(`    - Source: ${l.source}`)
      if (l.sourcePosition) lines.push(`      SourcePosition: ${l.sourcePosition}`)
      lines.push(`      Target: ${l.target}`)
      if (l.targetPosition) lines.push(`      TargetPosition: ${l.targetPosition}`)
      if (l.arrowHead && l.arrowHead !== 'none') {
        lines.push('      TargetArrowHead:')
        lines.push(`        Type: ${l.arrowHead}`)
      }
      const hasLabel = Object.values(l.labels).some(v => v.trim())
      if (hasLabel) {
        lines.push('      Labels:')
        const keys = ['SourceLeft', 'SourceRight', 'TargetLeft', 'TargetRight', 'AutoLeft', 'AutoRight'] as const
        keys.forEach(k => {
          const v = l.labels[k]
          if (v.trim()) {
            lines.push(`        ${k}:`)
            lines.push(`          Title: "${v}"`)
          }
        })
      }
    })
  }

  return lines.join('\n')
}

// ── Helpers ───────────────────────────────────────────────────────────────────

let idCounter = 0
function newId() { return `id_${++idCounter}` }

function newResource(): ResourceEntry {
  return {
    id: newId(), name: '', type: 'AWS::EC2::Instance',
    title: '', preset: '', direction: '', align: '',
    iconFill: '', children: [], borderChildren: [],
  }
}

function newLink(): LinkEntry {
  return {
    id: newId(), source: '', sourcePosition: '',
    target: '', targetPosition: '', arrowHead: 'Open',
    labels: { SourceLeft: '', SourceRight: '', TargetLeft: '', TargetRight: '', AutoLeft: '', AutoRight: '' },
  }
}

const DEFAULT_URL = 'https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml'

// ── Sub-components ────────────────────────────────────────────────────────────

function SectionHeader({ title, open, onToggle, action }: {
  title: string; open: boolean; onToggle: () => void; action?: React.ReactNode
}) {
  return (
    <div className="flex items-center justify-between px-4 py-2.5 bg-[#141414] border-b border-[#2a2a2a] cursor-pointer select-none"
      onClick={onToggle}>
      <div className="flex items-center gap-2">
        {open ? <ChevronDown size={13} className="text-[#555]" /> : <ChevronRight size={13} className="text-[#555]" />}
        <span className="text-xs font-semibold text-[#ccc] uppercase tracking-wider">{title}</span>
      </div>
      {action && <div onClick={e => e.stopPropagation()}>{action}</div>}
    </div>
  )
}

function Label({ children }: { children: React.ReactNode }) {
  return <span className="text-[10px] text-[#666] uppercase tracking-wider font-medium">{children}</span>
}

function Input({ value, onChange, placeholder, className = '' }: {
  value: string; onChange: (v: string) => void; placeholder?: string; className?: string
}) {
  return (
    <input
      type="text"
      value={value}
      onChange={e => onChange(e.target.value)}
      placeholder={placeholder}
      className={`w-full bg-[#111] border border-[#2a2a2a] rounded text-xs text-[#ccc] px-2.5 py-1.5 focus:outline-none focus:border-[#FF9900]/60 placeholder-[#444] ${className}`}
    />
  )
}

function Select({ value, onChange, options, placeholder }: {
  value: string; onChange: (v: string) => void
  options: { value: string; label: string }[]; placeholder?: string
}) {
  return (
    <select
      value={value}
      onChange={e => onChange(e.target.value)}
      className="w-full bg-[#111] border border-[#2a2a2a] rounded text-xs text-[#ccc] px-2.5 py-1.5 focus:outline-none focus:border-[#FF9900]/60"
    >
      {placeholder && <option value="">{placeholder}</option>}
      {options.map(o => (
        <option key={o.value} value={o.value}>{o.label}</option>
      ))}
    </select>
  )
}

function TypeSearchSelect({ value, onChange, placeholder }: {
  value: string; onChange: (v: string) => void; placeholder: string
}) {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handler(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  const filtered = Object.entries(AWS_TYPES).flatMap(([cat, types]) =>
    types.filter(t => t.toLowerCase().includes(query.toLowerCase())).map(t => ({ cat, type: t }))
  )

  const catColors: Record<string, string> = {
    'Diagram Elements': 'text-purple-400',
    'Compute': 'text-orange-400',
    'Networking': 'text-blue-400',
    'Storage': 'text-green-400',
    'Database': 'text-teal-400',
    'Messaging & Integration': 'text-yellow-400',
    'Security & Identity': 'text-red-400',
    'Developer Tools': 'text-cyan-400',
    'Monitoring & Management': 'text-gray-400',
    'Analytics': 'text-indigo-400',
    'AI & Machine Learning': 'text-pink-400',
    'IoT': 'text-lime-400',
  }

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen(o => !o)}
        className="w-full bg-[#111] border border-[#2a2a2a] rounded text-xs text-left px-2.5 py-1.5 focus:outline-none focus:border-[#FF9900]/60 flex items-center justify-between gap-1"
      >
        <span className={value ? 'text-[#FF9900] font-mono text-[11px]' : 'text-[#444]'}>
          {value || placeholder}
        </span>
        <ChevronDown size={11} className="text-[#555] flex-shrink-0" />
      </button>
      {open && (
        <div className="absolute z-50 top-full left-0 right-0 mt-1 bg-[#181818] border border-[#2a2a2a] rounded-lg shadow-2xl overflow-hidden">
          <div className="p-2 border-b border-[#2a2a2a]">
            <input
              autoFocus
              type="text"
              value={query}
              onChange={e => setQuery(e.target.value)}
              placeholder={placeholder}
              className="w-full bg-[#111] border border-[#333] rounded text-xs text-[#ccc] px-2 py-1.5 focus:outline-none focus:border-[#FF9900]/60 placeholder-[#444]"
            />
          </div>
          <div className="max-h-64 overflow-y-auto">
            {filtered.length === 0 ? (
              <p className="text-xs text-[#555] px-3 py-3 text-center">No results</p>
            ) : (
              filtered.map(({ cat, type }) => (
                <button
                  key={type}
                  type="button"
                  onClick={() => { onChange(type); setOpen(false); setQuery('') }}
                  className="w-full text-left px-3 py-1.5 hover:bg-[#252525] group flex items-center gap-2"
                >
                  <span className={`text-[9px] font-mono shrink-0 ${catColors[cat] ?? 'text-[#666]'}`}>
                    {cat.split(' ')[0].slice(0, 4).toUpperCase()}
                  </span>
                  <span className="text-[11px] text-[#ccc] font-mono truncate group-hover:text-[#e5e5e5]">
                    {type}
                  </span>
                </button>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}

function ChildrenSelect({ value, onChange, allResources, selfName }: {
  value: string[]; onChange: (v: string[]) => void
  allResources: ResourceEntry[]; selfName: string
}) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handler(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  const available = allResources.filter(r => r.name && r.name !== selfName)

  function toggle(name: string) {
    onChange(value.includes(name) ? value.filter(v => v !== name) : [...value, name])
  }

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen(o => !o)}
        className="w-full bg-[#111] border border-[#2a2a2a] rounded text-xs text-left px-2.5 py-1.5 focus:outline-none focus:border-[#FF9900]/60 flex items-center justify-between gap-1 min-h-[30px]"
      >
        <span className="text-[#ccc] truncate">
          {value.length === 0 ? <span className="text-[#444]">—</span> : value.join(', ')}
        </span>
        <ChevronDown size={11} className="text-[#555] flex-shrink-0" />
      </button>
      {open && available.length > 0 && (
        <div className="absolute z-40 top-full left-0 right-0 mt-1 bg-[#181818] border border-[#2a2a2a] rounded-lg shadow-xl max-h-40 overflow-y-auto py-1">
          {available.map(r => (
            <label key={r.id} className="flex items-center gap-2 px-3 py-1.5 hover:bg-[#252525] cursor-pointer">
              <input
                type="checkbox"
                checked={value.includes(r.name)}
                onChange={() => toggle(r.name)}
                className="accent-[#FF9900]"
              />
              <span className="text-xs text-[#ccc] font-mono">{r.name}</span>
            </label>
          ))}
        </div>
      )}
    </div>
  )
}

// ── Resource card ─────────────────────────────────────────────────────────────

function ResourceCard({ res, index, resources, onChange, onRemove, t }: {
  res: ResourceEntry; index: number; resources: ResourceEntry[]
  onChange: (updated: ResourceEntry) => void
  onRemove: () => void
  t: ReturnType<typeof useLanguage>['t']
}) {
  const [open, setOpen] = useState(true)
  const presets = PRESETS[res.type] ?? []
  const allNames = resources.map(r => r.name).filter(Boolean)

  function update(patch: Partial<ResourceEntry>) {
    onChange({ ...res, ...patch })
  }

  function addBorderChild() {
    update({ borderChildren: [...res.borderChildren, { id: newId(), position: 'S', resource: '' }] })
  }

  function updateBorderChild(id: string, patch: Partial<BorderChild>) {
    update({
      borderChildren: res.borderChildren.map(bc => bc.id === id ? { ...bc, ...patch } : bc)
    })
  }

  function removeBorderChild(id: string) {
    update({ borderChildren: res.borderChildren.filter(bc => bc.id !== id) })
  }

  const catColor = Object.entries(AWS_TYPES).find(([, types]) => types.includes(res.type))?.[0]

  return (
    <div className="border border-[#2a2a2a] rounded-lg overflow-hidden">
      {/* Card header */}
      <div
        className="flex items-center justify-between px-3 py-2 bg-[#161616] cursor-pointer"
        onClick={() => setOpen(o => !o)}
      >
        <div className="flex items-center gap-2 min-w-0">
          {open ? <ChevronDown size={12} className="text-[#555] shrink-0" /> : <ChevronRight size={12} className="text-[#555] shrink-0" />}
          <span className="text-xs font-mono text-[#FF9900] truncate">{res.name || `resource_${index + 1}`}</span>
          <span className="text-[10px] text-[#555] truncate hidden sm:block">{res.type.split('::').slice(1).join('::')}</span>
        </div>
        <button
          type="button"
          onClick={e => { e.stopPropagation(); onRemove() }}
          className="text-[#555] hover:text-red-400 transition-colors p-1 shrink-0"
          aria-label={t.removeLabel}
        >
          <Trash2 size={12} />
        </button>
      </div>

      {open && (
        <div className="p-3 space-y-3 bg-[#0f0f0f]">
          {/* Name + Type */}
          <div className="grid grid-cols-2 gap-2">
            <div className="space-y-1">
              <Label>{t.resourceName}</Label>
              <Input
                value={res.name}
                onChange={v => update({ name: v.replace(/\s/g, '') })}
                placeholder={t.resourceNamePlaceholder}
              />
            </div>
            <div className="space-y-1">
              <Label>{t.resourceType}</Label>
              <TypeSearchSelect
                value={res.type}
                onChange={v => update({ type: v, preset: '' })}
                placeholder={t.searchTypes}
              />
            </div>
          </div>

          {/* Title + Preset */}
          <div className="grid grid-cols-2 gap-2">
            <div className="space-y-1">
              <Label>{t.resourceTitle}</Label>
              <Input value={res.title} onChange={v => update({ title: v })} placeholder="—" />
            </div>
            <div className="space-y-1">
              <Label>{t.resourcePreset}</Label>
              {presets.length > 0 ? (
                <Select
                  value={res.preset}
                  onChange={v => update({ preset: v })}
                  options={presets.map(p => ({ value: p, label: p }))}
                  placeholder={t.selectOptional}
                />
              ) : (
                <Input value={res.preset} onChange={v => update({ preset: v })} placeholder={t.selectOptional} />
              )}
            </div>
          </div>

          {/* Direction + Align + IconFill */}
          <div className="grid grid-cols-3 gap-2">
            <div className="space-y-1">
              <Label>{t.resourceDirection}</Label>
              <Select
                value={res.direction}
                onChange={v => update({ direction: v })}
                options={[{ value: 'vertical', label: 'vertical' }, { value: 'horizontal', label: 'horizontal' }]}
                placeholder={t.selectOptional}
              />
            </div>
            <div className="space-y-1">
              <Label>{t.resourceAlign}</Label>
              <Select
                value={res.align}
                onChange={v => update({ align: v })}
                options={ALIGNS.map(a => ({ value: a, label: a }))}
                placeholder={t.selectOptional}
              />
            </div>
            <div className="space-y-1">
              <Label>{t.iconFillLabel}</Label>
              <Select
                value={res.iconFill}
                onChange={v => update({ iconFill: v })}
                options={ICON_FILLS.map(f => ({ value: f, label: f }))}
                placeholder={t.selectOptional}
              />
            </div>
          </div>

          {/* Children */}
          <div className="space-y-1">
            <Label>{t.resourceChildren}</Label>
            <ChildrenSelect
              value={res.children}
              onChange={v => update({ children: v })}
              allResources={resources}
              selfName={res.name}
            />
          </div>

          {/* Border children */}
          <div className="space-y-1">
            <div className="flex items-center justify-between">
              <Label>{t.borderChildrenLabel}</Label>
              <button
                type="button"
                onClick={addBorderChild}
                className="text-[10px] text-[#FF9900] hover:text-[#ffb340] transition-colors"
              >
                {t.addBorderChild}
              </button>
            </div>
            {res.borderChildren.map(bc => (
              <div key={bc.id} className="flex items-center gap-2">
                <Select
                  value={bc.position}
                  onChange={v => updateBorderChild(bc.id, { position: v })}
                  options={POSITIONS.map(p => ({ value: p, label: p }))}
                />
                <Select
                  value={bc.resource}
                  onChange={v => updateBorderChild(bc.id, { resource: v })}
                  options={allNames.filter(n => n !== res.name).map(n => ({ value: n, label: n }))}
                  placeholder="—"
                />
                <button
                  type="button"
                  onClick={() => removeBorderChild(bc.id)}
                  className="text-[#555] hover:text-red-400 shrink-0"
                >
                  <Trash2 size={12} />
                </button>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

// ── Link card ─────────────────────────────────────────────────────────────────

function LinkCard({ link, resources, onChange, onRemove, t }: {
  link: LinkEntry; resources: ResourceEntry[]
  onChange: (updated: LinkEntry) => void; onRemove: () => void
  t: ReturnType<typeof useLanguage>['t']
}) {
  const [open, setOpen] = useState(true)
  const [labelsOpen, setLabelsOpen] = useState(false)
  const names = resources.map(r => r.name).filter(Boolean)

  function update(patch: Partial<LinkEntry>) { onChange({ ...link, ...patch }) }
  function updateLabel(k: keyof EdgeLabel, v: string) {
    update({ labels: { ...link.labels, [k]: v } })
  }

  const labelKeys: (keyof EdgeLabel)[] = ['SourceLeft', 'SourceRight', 'TargetLeft', 'TargetRight', 'AutoLeft', 'AutoRight']

  return (
    <div className="border border-[#2a2a2a] rounded-lg overflow-hidden">
      <div
        className="flex items-center justify-between px-3 py-2 bg-[#161616] cursor-pointer"
        onClick={() => setOpen(o => !o)}
      >
        <div className="flex items-center gap-2 min-w-0">
          {open ? <ChevronDown size={12} className="text-[#555] shrink-0" /> : <ChevronRight size={12} className="text-[#555] shrink-0" />}
          <span className="text-xs font-mono text-[#ccc] truncate">
            {link.source || '?'} <span className="text-[#FF9900]">→</span> {link.target || '?'}
          </span>
        </div>
        <button
          type="button"
          onClick={e => { e.stopPropagation(); onRemove() }}
          className="text-[#555] hover:text-red-400 transition-colors p-1 shrink-0"
        >
          <Trash2 size={12} />
        </button>
      </div>

      {open && (
        <div className="p-3 space-y-3 bg-[#0f0f0f]">
          {/* Source + Target */}
          <div className="grid grid-cols-2 gap-2">
            <div className="space-y-1">
              <Label>{t.linkSource}</Label>
              <Select
                value={link.source}
                onChange={v => update({ source: v })}
                options={names.map(n => ({ value: n, label: n }))}
                placeholder="—"
              />
            </div>
            <div className="space-y-1">
              <Label>{t.linkTarget}</Label>
              <Select
                value={link.target}
                onChange={v => update({ target: v })}
                options={names.map(n => ({ value: n, label: n }))}
                placeholder="—"
              />
            </div>
          </div>

          {/* Positions + ArrowHead */}
          <div className="grid grid-cols-3 gap-2">
            <div className="space-y-1">
              <Label>{t.sourcePosLabel}</Label>
              <Select
                value={link.sourcePosition}
                onChange={v => update({ sourcePosition: v })}
                options={POSITIONS.map(p => ({ value: p, label: p }))}
                placeholder={t.selectOptional}
              />
            </div>
            <div className="space-y-1">
              <Label>{t.targetPosLabel}</Label>
              <Select
                value={link.targetPosition}
                onChange={v => update({ targetPosition: v })}
                options={POSITIONS.map(p => ({ value: p, label: p }))}
                placeholder={t.selectOptional}
              />
            </div>
            <div className="space-y-1">
              <Label>{t.arrowHeadLabel}</Label>
              <Select
                value={link.arrowHead}
                onChange={v => update({ arrowHead: v })}
                options={ARROW_HEADS.map(a => ({ value: a, label: a }))}
              />
            </div>
          </div>

          {/* Edge labels (collapsible) */}
          <div>
            <button
              type="button"
              onClick={() => setLabelsOpen(o => !o)}
              className="flex items-center gap-1 text-[10px] text-[#555] hover:text-[#999] transition-colors uppercase tracking-wider"
            >
              {labelsOpen ? <ChevronDown size={10} /> : <ChevronRight size={10} />}
              {t.edgeLabels}
            </button>
            {labelsOpen && (
              <div className="grid grid-cols-2 gap-2 mt-2">
                {labelKeys.map(k => (
                  <div key={k} className="space-y-0.5">
                    <Label>{k}</Label>
                    <Input
                      value={link.labels[k]}
                      onChange={v => updateLabel(k, v)}
                      placeholder='e.g. "HTTP:80"'
                    />
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

// ── Main page ─────────────────────────────────────────────────────────────────

export default function BuilderPage() {
  const { t } = useLanguage()
  const router = useRouter()

  const [form, setForm] = useState<FormState>({
    defType: 'URL',
    defUrl: DEFAULT_URL,
    defFile: '',
    canvasDirection: 'vertical',
    resources: [],
    links: [],
  })

  const [sections, setSections] = useState({
    def: false,
    canvas: false,
    resources: true,
    links: true,
  })

  const [copied, setCopied] = useState(false)

  const yaml = generateYaml(form)

  function toggleSection(key: keyof typeof sections) {
    setSections(s => ({ ...s, [key]: !s[key] }))
  }

  function updateForm(patch: Partial<FormState>) {
    setForm(f => ({ ...f, ...patch }))
  }

  function addResource() {
    updateForm({ resources: [...form.resources, newResource()] })
  }

  function updateResource(id: string, updated: ResourceEntry) {
    updateForm({ resources: form.resources.map(r => r.id === id ? updated : r) })
  }

  function removeResource(id: string) {
    const name = form.resources.find(r => r.id === id)?.name ?? ''
    updateForm({
      resources: form.resources.filter(r => r.id !== id).map(r => ({
        ...r,
        children: r.children.filter(c => c !== name),
        borderChildren: r.borderChildren.filter(bc => bc.resource !== name),
      })),
      links: form.links.map(l => ({
        ...l,
        source: l.source === name ? '' : l.source,
        target: l.target === name ? '' : l.target,
      })),
    })
  }

  function addLink() {
    updateForm({ links: [...form.links, newLink()] })
  }

  function updateLink(id: string, updated: LinkEntry) {
    updateForm({ links: form.links.map(l => l.id === id ? updated : l) })
  }

  function removeLink(id: string) {
    updateForm({ links: form.links.filter(l => l.id !== id) })
  }

  async function copyYaml() {
    await navigator.clipboard.writeText(yaml)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  function useInEditor() {
    localStorage.setItem('builder_yaml', yaml)
    router.push('/')
  }

  // Pick up builder_yaml when returning to editor
  useEffect(() => {
    const saved = localStorage.getItem('builder_yaml')
    if (saved) localStorage.removeItem('builder_yaml')
  }, [])

  const lineCount = yaml.split('\n').length

  return (
    <div className="flex flex-col h-screen bg-[#0f0f0f] overflow-hidden">

      {/* ── Header ── */}
      <header className="flex items-center justify-between px-4 h-12 border-b border-[#2a2a2a] flex-shrink-0">
        <div className="flex items-center gap-3">
          <Link href="/" className="flex items-center gap-1.5 text-xs text-[#666] hover:text-[#ccc] transition-colors">
            <ArrowLeft size={13} />
            {t.backToEditor}
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
            <span className="text-sm font-semibold text-[#e5e5e5]">{t.builderTitle}</span>
            <span className="hidden sm:block text-xs text-[#444]">— {t.builderDesc}</span>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <LanguageSwitcher />
          <Link href="/docs" className="flex items-center gap-1.5 text-xs text-[#555] hover:text-[#999] transition-colors px-2 py-1 rounded hover:bg-[#1a1a1a]">
            <BookOpen size={13} />
            {t.docs}
          </Link>
          <a href="https://github.com/fernandofatech/diagram-as-code" target="_blank" rel="noopener noreferrer"
            className="text-[#555] hover:text-[#999] transition-colors p-1" aria-label="GitHub">
            <Github size={16} />
          </a>
        </div>
      </header>

      {/* ── Body ── */}
      <div className="flex flex-1 overflow-hidden">

        {/* Left: Form */}
        <div className="w-[55%] flex flex-col overflow-hidden border-r border-[#2a2a2a]">
          <div className="flex-1 overflow-y-auto">

            {/* Definition File */}
            <div className="border-b border-[#2a2a2a]">
              <SectionHeader title={t.defFileSection} open={sections.def} onToggle={() => toggleSection('def')} />
              {sections.def && (
                <div className="p-4 space-y-3">
                  <div className="flex gap-3">
                    {(['URL', 'LocalFile'] as const).map(type => (
                      <label key={type} className="flex items-center gap-1.5 cursor-pointer">
                        <input type="radio" name="defType" value={type}
                          checked={form.defType === type}
                          onChange={() => updateForm({ defType: type })}
                          className="accent-[#FF9900]"
                        />
                        <span className="text-xs text-[#999]">
                          {type === 'URL' ? t.officialUrl : t.localFileLabel}
                        </span>
                      </label>
                    ))}
                  </div>
                  {form.defType === 'URL' ? (
                    <div className="space-y-1">
                      <Label>URL</Label>
                      <Input value={form.defUrl} onChange={v => updateForm({ defUrl: v })} placeholder={t.customUrlPlaceholder} />
                    </div>
                  ) : (
                    <div className="space-y-1">
                      <Label>{t.filePathLabel}</Label>
                      <Input value={form.defFile} onChange={v => updateForm({ defFile: v })} placeholder={t.localFilePlaceholder} />
                    </div>
                  )}
                </div>
              )}
            </div>

            {/* Canvas */}
            <div className="border-b border-[#2a2a2a]">
              <SectionHeader title={t.canvasSection} open={sections.canvas} onToggle={() => toggleSection('canvas')} />
              {sections.canvas && (
                <div className="p-4">
                  <div className="space-y-1">
                    <Label>{t.canvasDir}</Label>
                    <div className="flex gap-2">
                      {(['vertical', 'horizontal'] as const).map(d => (
                        <button
                          key={d}
                          type="button"
                          onClick={() => updateForm({ canvasDirection: d })}
                          className={`px-4 py-1.5 rounded text-xs font-medium transition-all border ${
                            form.canvasDirection === d
                              ? 'bg-[#FF9900] border-[#FF9900] text-[#0f0f0f]'
                              : 'border-[#2a2a2a] text-[#666] hover:text-[#ccc] hover:border-[#444]'
                          }`}
                        >
                          {d}
                        </button>
                      ))}
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* Resources */}
            <div className="border-b border-[#2a2a2a]">
              <SectionHeader
                title={`${t.resourcesSection} (${form.resources.length})`}
                open={sections.resources}
                onToggle={() => toggleSection('resources')}
                action={
                  <button
                    type="button"
                    onClick={addResource}
                    className="flex items-center gap-1 text-[10px] text-[#FF9900] hover:text-[#ffb340] transition-colors font-medium"
                  >
                    <Plus size={11} />
                    {t.addResource}
                  </button>
                }
              />
              {sections.resources && (
                <div className="p-3 space-y-2">
                  {form.resources.length === 0 ? (
                    <p className="text-xs text-[#444] text-center py-6">{t.noResources}</p>
                  ) : (
                    form.resources.map((res, i) => (
                      <ResourceCard
                        key={res.id}
                        res={res}
                        index={i}
                        resources={form.resources}
                        onChange={updated => updateResource(res.id, updated)}
                        onRemove={() => removeResource(res.id)}
                        t={t}
                      />
                    ))
                  )}
                  {form.resources.length > 0 && (
                    <button
                      type="button"
                      onClick={addResource}
                      className="w-full py-2 border border-dashed border-[#333] rounded-lg text-xs text-[#555] hover:border-[#FF9900]/40 hover:text-[#FF9900] transition-colors"
                    >
                      {t.addResource}
                    </button>
                  )}
                </div>
              )}
            </div>

            {/* Links */}
            <div>
              <SectionHeader
                title={`${t.linksSection} (${form.links.length})`}
                open={sections.links}
                onToggle={() => toggleSection('links')}
                action={
                  <button
                    type="button"
                    onClick={addLink}
                    className="flex items-center gap-1 text-[10px] text-[#FF9900] hover:text-[#ffb340] transition-colors font-medium"
                  >
                    <Plus size={11} />
                    {t.addLink}
                  </button>
                }
              />
              {sections.links && (
                <div className="p-3 space-y-2">
                  {form.links.length === 0 ? (
                    <p className="text-xs text-[#444] text-center py-6">{t.noLinks}</p>
                  ) : (
                    form.links.map(link => (
                      <LinkCard
                        key={link.id}
                        link={link}
                        resources={form.resources}
                        onChange={updated => updateLink(link.id, updated)}
                        onRemove={() => removeLink(link.id)}
                        t={t}
                      />
                    ))
                  )}
                  {form.links.length > 0 && (
                    <button
                      type="button"
                      onClick={addLink}
                      className="w-full py-2 border border-dashed border-[#333] rounded-lg text-xs text-[#555] hover:border-[#FF9900]/40 hover:text-[#FF9900] transition-colors"
                    >
                      {t.addLink}
                    </button>
                  )}
                </div>
              )}
            </div>

          </div>
        </div>

        {/* Right: YAML Preview */}
        <div className="w-[45%] flex flex-col overflow-hidden">
          <div className="flex items-center justify-between px-4 py-2 border-b border-[#2a2a2a] flex-shrink-0">
            <div className="flex items-center gap-3">
              <span className="text-xs text-[#555] font-medium uppercase tracking-wider">{t.yamlPreview}</span>
              <span className="text-[10px] text-[#444]">{lineCount} lines</span>
            </div>
            <div className="flex items-center gap-2">
              <button
                type="button"
                onClick={copyYaml}
                className="flex items-center gap-1.5 text-xs text-[#555] hover:text-[#ccc] transition-colors px-2 py-1 rounded hover:bg-[#1a1a1a] border border-[#2a2a2a]"
              >
                {copied ? <Check size={12} className="text-green-400" /> : <Copy size={12} />}
                {copied ? t.copied : t.copyYaml}
              </button>
              <button
                type="button"
                onClick={useInEditor}
                className="flex items-center gap-1.5 px-3 py-1.5 bg-[#FF9900] hover:bg-[#ffb340] text-[#0f0f0f] text-xs font-semibold rounded transition-colors"
              >
                <Zap size={12} />
                {t.useInEditor}
              </button>
            </div>
          </div>
          <div className="flex-1 overflow-auto">
            <pre className="text-xs font-mono text-[#ccc] leading-relaxed p-4 whitespace-pre min-h-full">
              <code>
                {yaml.split('\n').map((line, i) => (
                  <span key={i} className="block">
                    <span className="text-[#333] select-none mr-3 text-[10px]">{String(i + 1).padStart(3, ' ')}</span>
                    {line
                      .replace(/^(\s*)([\w:]+:)(\s|$)/, (_, indent, key, rest) =>
                        `${indent}<k>${key}</k>${rest}`)
                      .split(/<k>|<\/k>/)
                      .map((part, j) =>
                        j % 2 === 1
                          ? <span key={j} className="text-[#7eb8f7]">{part}</span>
                          : <span key={j}>{part
                              .replace(/(AWS::[A-Za-z:]+)/g, '')
                              .split(/(AWS::[A-Za-z:]+)/)
                              .map((p, k) =>
                                /^AWS::/.test(p)
                                  ? <span key={k} className="text-[#FF9900]">{p}</span>
                                  : <span key={k}>{p
                                      .split(/(- |"|')/)
                                      .map((s, m) =>
                                        s === '"' || s === "'"
                                          ? <span key={m} className="text-green-400">{s}</span>
                                          : s === '- '
                                            ? <span key={m} className="text-[#888]">{s}</span>
                                            : s
                                      )
                                    }</span>
                              )
                            }</span>
                      )
                    }
                  </span>
                ))}
              </code>
            </pre>
          </div>
        </div>
      </div>
    </div>
  )
}

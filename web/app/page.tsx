'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Github, Zap, Terminal, Globe, Wrench, ArrowRight, Code2, Layers, BookOpen, Menu, X } from 'lucide-react'
import LanguageSwitcher from '@/components/LanguageSwitcher'
import ThemeSwitcher from '@/components/ThemeSwitcher'
import { useLanguage } from '@/lib/i18n'

const YAML_DEMO = `Diagram:
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children: [AWSCloud]
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Children: [VPC]
    VPC:
      Type: AWS::EC2::VPC
      Children: [Subnet]
    Subnet:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children: [Instance]
    Instance:
      Type: AWS::EC2::Instance`

const CONTENT = {
  en: {
    nav: { editor: 'Editor', builder: 'Builder', docs: 'Docs' },
    hero: {
      badge: 'Open Source · Free · No Install Required',
      heading1: 'AWS Architecture Diagrams',
      heading2: 'from YAML',
      sub: 'Write YAML, get beautiful AWS architecture diagrams. Generate PNG or draw.io files instantly — in the browser or via CLI.',
      cta: 'Open Editor',
      ctaBuilder: 'Try the Builder',
    },
    how: {
      title: 'How it works',
      steps: [
        { icon: '✍️', title: 'Write YAML', desc: 'Describe your AWS architecture using a simple, readable YAML syntax.' },
        { icon: '⚡', title: 'Generate', desc: 'Click Generate or press Ctrl+Enter. Our Go backend renders your diagram instantly.' },
        { icon: '📥', title: 'Download', desc: 'Export as PNG for presentations or draw.io for further editing.' },
      ],
    },
    features: {
      title: 'Everything you need',
      items: [
        { icon: <Code2 size={20} />, title: 'YAML-first', desc: 'No drag-and-drop needed. Define your entire architecture in code.' },
        { icon: <Layers size={20} />, title: '100+ AWS Icons', desc: 'Full AWS icon set with proper colors, borders, and grouping.' },
        { icon: <Zap size={20} />, title: 'Instant Preview', desc: 'Real-time diagram preview as you type.' },
        { icon: <Wrench size={20} />, title: 'Visual Builder', desc: 'Prefer forms? Use the Builder to generate YAML without writing a line.' },
        { icon: <Terminal size={20} />, title: 'CLI Tool', desc: 'Use the CLI in your CI/CD pipeline to automate diagram generation.' },
        { icon: <Globe size={20} />, title: 'Open Source', desc: 'MIT licensed. Contribute, self-host, or fork it freely.' },
      ],
    },
    cta: {
      title: 'Start building your diagrams',
      sub: 'No sign-up. No install. Just open the editor and start writing.',
      btn: 'Open Editor →',
    },
    footer: 'Made with ❤️ by Fernando Azevedo · Open Source',
  },
  pt: {
    nav: { editor: 'Editor', builder: 'Builder', docs: 'Docs' },
    hero: {
      badge: 'Open Source · Gratuito · Sem instalação',
      heading1: 'Diagramas de Arquitetura AWS',
      heading2: 'a partir de YAML',
      sub: 'Escreva YAML, obtenha diagramas AWS bonitos. Gere arquivos PNG ou draw.io instantaneamente — no navegador ou via CLI.',
      cta: 'Abrir Editor',
      ctaBuilder: 'Usar o Builder',
    },
    how: {
      title: 'Como funciona',
      steps: [
        { icon: '✍️', title: 'Escreva YAML', desc: 'Descreva sua arquitetura AWS usando uma sintaxe YAML simples e legível.' },
        { icon: '⚡', title: 'Gere', desc: 'Clique em Gerar ou pressione Ctrl+Enter. Nosso backend Go renderiza seu diagrama.' },
        { icon: '📥', title: 'Baixe', desc: 'Exporte como PNG para apresentações ou draw.io para edição adicional.' },
      ],
    },
    features: {
      title: 'Tudo que você precisa',
      items: [
        { icon: <Code2 size={20} />, title: 'YAML primeiro', desc: 'Sem arrastar e soltar. Defina toda sua arquitetura em código.' },
        { icon: <Layers size={20} />, title: '100+ ícones AWS', desc: 'Conjunto completo de ícones AWS com cores, bordas e agrupamentos corretos.' },
        { icon: <Zap size={20} />, title: 'Preview instantâneo', desc: 'Visualização em tempo real enquanto você digita.' },
        { icon: <Wrench size={20} />, title: 'Builder Visual', desc: 'Prefere formulários? Use o Builder para gerar YAML sem escrever uma linha.' },
        { icon: <Terminal size={20} />, title: 'CLI Tool', desc: 'Use o CLI no seu pipeline CI/CD para automatizar a geração de diagramas.' },
        { icon: <Globe size={20} />, title: 'Open Source', desc: 'Licença MIT. Contribua, faça self-host ou fork livremente.' },
      ],
    },
    cta: {
      title: 'Comece a criar seus diagramas',
      sub: 'Sem cadastro. Sem instalação. Só abrir o editor e começar a escrever.',
      btn: 'Abrir Editor →',
    },
    footer: 'Feito com ❤️ por Fernando Azevedo · Open Source',
  },
}

export default function LandingPage() {
  const { lang } = useLanguage()
  const c = CONTENT[lang] ?? CONTENT.en
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)

  return (
    <div className="min-h-screen bg-[var(--bg-elevated)] text-[var(--text)] flex flex-col">

      {/* ── Sticky Header ─────────────────────────────────────────────────── */}
      <header className="sticky top-0 z-50 flex items-center justify-between px-6 h-14 border-b border-[var(--border)] bg-[var(--bg-elevated)]/90 backdrop-blur">
        <div className="flex items-center gap-2">
          <div className="w-7 h-7 bg-[#FF9900] rounded-md flex items-center justify-center">
            <svg viewBox="0 0 16 16" fill="white" className="w-4 h-4">
              <rect x="1" y="1" width="6" height="6" rx="1" />
              <rect x="9" y="1" width="6" height="6" rx="1" />
              <rect x="1" y="9" width="6" height="6" rx="1" />
              <rect x="9" y="9" width="6" height="6" rx="1" />
            </svg>
          </div>
          <span className="font-semibold text-sm tracking-tight text-[var(--text)]">diagram-as-code</span>
        </div>

        {/* Desktop nav */}
        <nav className="hidden sm:flex items-center gap-1">
          <Link href="/editor" className="text-xs text-[var(--text-3)] hover:text-[var(--text)] transition-colors px-3 py-1.5 rounded hover:bg-[var(--surface)]">{c.nav.editor}</Link>
          <Link href="/builder" className="text-xs text-[var(--text-3)] hover:text-[var(--text)] transition-colors px-3 py-1.5 rounded hover:bg-[var(--surface)]">{c.nav.builder}</Link>
          <Link href="/docs" className="text-xs text-[var(--text-3)] hover:text-[var(--text)] transition-colors px-3 py-1.5 rounded hover:bg-[var(--surface)]">{c.nav.docs}</Link>
          <div className="w-px h-4 bg-[var(--border)] mx-1" />
          <LanguageSwitcher />
          <ThemeSwitcher />
          <a
            href="https://github.com/fernandofatech/diagram-as-code"
            target="_blank"
            rel="noopener noreferrer"
            className="text-[var(--text-5)] hover:text-[var(--text)] transition-colors p-1.5 rounded hover:bg-[var(--surface)]"
            aria-label="GitHub"
          >
            <Github size={16} />
          </a>
        </nav>

        {/* Mobile controls */}
        <div className="flex sm:hidden items-center gap-1">
          <LanguageSwitcher />
          <ThemeSwitcher />
          <button
            onClick={() => setMobileMenuOpen(o => !o)}
            className="p-1.5 rounded hover:bg-[var(--surface)] text-[var(--text-3)] transition-colors"
            aria-label="Menu"
          >
            {mobileMenuOpen ? <X size={18} /> : <Menu size={18} />}
          </button>
        </div>
      </header>

      {/* Mobile menu */}
      {mobileMenuOpen && (
        <>
          <div className="fixed inset-0 z-40 sm:hidden" onClick={() => setMobileMenuOpen(false)} />
          <div className="fixed top-14 left-0 right-0 z-50 sm:hidden bg-[var(--bg-elevated)] border-b border-[var(--border)] py-2 px-4 flex flex-col gap-0.5 shadow-lg">
            <Link href="/editor" onClick={() => setMobileMenuOpen(false)} className="flex items-center gap-2.5 px-3 py-2.5 rounded-lg hover:bg-[var(--surface)] text-sm text-[var(--text-2)] transition-colors"><Zap size={15} className="text-[#FF9900]" />{c.nav.editor}</Link>
            <Link href="/builder" onClick={() => setMobileMenuOpen(false)} className="flex items-center gap-2.5 px-3 py-2.5 rounded-lg hover:bg-[var(--surface)] text-sm text-[var(--text-2)] transition-colors"><Wrench size={15} className="text-[#FF9900]" />{c.nav.builder}</Link>
            <Link href="/docs" onClick={() => setMobileMenuOpen(false)} className="flex items-center gap-2.5 px-3 py-2.5 rounded-lg hover:bg-[var(--surface)] text-sm text-[var(--text-2)] transition-colors"><BookOpen size={15} className="text-[#FF9900]" />{c.nav.docs}</Link>
            <a href="https://github.com/fernandofatech/diagram-as-code" target="_blank" rel="noopener noreferrer" className="flex items-center gap-2.5 px-3 py-2.5 rounded-lg hover:bg-[var(--surface)] text-sm text-[var(--text-2)] transition-colors"><Github size={15} /> GitHub</a>
          </div>
        </>
      )}

      {/* ── Hero ──────────────────────────────────────────────────────────── */}
      <section className="flex flex-col items-center justify-center text-center px-6 pt-24 pb-20 gap-6">
        <span className="text-xs text-[#FF9900] border border-[#FF9900]/30 bg-[#FF9900]/10 px-3 py-1 rounded-full">
          {c.hero.badge}
        </span>

        <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold leading-tight max-w-3xl">
          <span className="text-[var(--text)]">{c.hero.heading1}</span>
          <br />
          <span className="text-[#FF9900]">{c.hero.heading2}</span>
        </h1>

        <p className="text-base text-[var(--text-3)] max-w-xl leading-relaxed">
          {c.hero.sub}
        </p>

        <div className="flex items-center gap-3 flex-wrap justify-center">
          <Link
            href="/editor"
            className="flex items-center gap-2 px-6 py-2.5 bg-[#FF9900] hover:bg-[#ffb340] text-[var(--accent-contrast)] font-semibold text-sm rounded-lg transition-colors"
          >
            <Zap size={15} />
            {c.hero.cta}
          </Link>
          <Link
            href="/builder"
            className="flex items-center gap-2 px-6 py-2.5 bg-[var(--surface)] hover:bg-[var(--surface-hover)] border border-[var(--border)] text-[var(--text-2)] font-medium text-sm rounded-lg transition-colors"
          >
            <Wrench size={15} />
            {c.hero.ctaBuilder}
          </Link>
        </div>

        {/* Demo block */}
        <div className="w-full max-w-4xl mt-8 grid grid-cols-1 md:grid-cols-2 gap-4 text-left">
          {/* YAML code */}
          <div className="bg-[var(--code-bg)] border border-[var(--border)] rounded-xl overflow-hidden">
            <div className="flex items-center gap-2 px-4 py-2.5 border-b border-[var(--border)]">
              <div className="w-2.5 h-2.5 rounded-full bg-[#ff5f57]" />
              <div className="w-2.5 h-2.5 rounded-full bg-[#febc2e]" />
              <div className="w-2.5 h-2.5 rounded-full bg-[#28c840]" />
              <span className="ml-2 text-[10px] text-[var(--text-6)] font-mono">architecture.yaml</span>
            </div>
            <pre className="px-4 py-4 text-[11px] text-[var(--text-2)] font-mono leading-relaxed overflow-auto">
              <code>{YAML_DEMO}</code>
            </pre>
          </div>

          {/* Diagram illustration */}
          <div className="bg-[var(--code-bg)] border border-[var(--border)] rounded-xl overflow-hidden flex items-center justify-center p-8">
            <svg viewBox="0 0 220 200" className="w-full max-w-[220px]" fill="none" xmlns="http://www.w3.org/2000/svg">
              {/* Cloud outline */}
              <rect x="4" y="4" width="212" height="192" rx="12" stroke="#FF9900" strokeWidth="1.5" strokeDasharray="4 3" fill="none" />
              {/* VPC */}
              <rect x="20" y="20" width="180" height="160" rx="8" stroke="#8b5cf6" strokeWidth="1.5" fill="#8b5cf620" />
              <text x="26" y="34" fontSize="8" fill="#8b5cf6" fontFamily="monospace">VPC</text>
              {/* Subnet */}
              <rect x="36" y="44" width="148" height="100" rx="6" stroke="#10b981" strokeWidth="1.5" fill="#10b98115" />
              <text x="42" y="57" fontSize="7" fill="#10b981" fontFamily="monospace">PublicSubnet</text>
              {/* EC2 icon */}
              <rect x="90" y="68" width="40" height="40" rx="6" fill="#FF9900" opacity="0.9" />
              <text x="95" y="92" fontSize="9" fill="white" fontFamily="monospace" fontWeight="bold">EC2</text>
              {/* Arrow down */}
              <line x1="110" y1="148" x2="110" y2="168" stroke="var(--text-5)" strokeWidth="1.5" />
              <polygon points="105,165 110,173 115,165" fill="var(--text-5)" />
              {/* IGW */}
              <rect x="85" y="173" width="50" height="18" rx="4" fill="var(--surface)" stroke="var(--text-6)" strokeWidth="1" />
              <text x="98" y="185" fontSize="7" fill="var(--text-3)" fontFamily="monospace">IGW</text>
            </svg>
          </div>
        </div>
      </section>

      {/* ── How it works ──────────────────────────────────────────────────── */}
      <section className="px-6 py-20 border-t border-[var(--border)]">
        <div className="max-w-4xl mx-auto">
          <h2 className="text-2xl font-bold text-center mb-12 text-[var(--text)]">{c.how.title}</h2>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-8">
            {c.how.steps.map((step, i) => (
              <div key={i} className="flex flex-col items-center text-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-[var(--surface)] border border-[var(--border)] flex items-center justify-center text-2xl">
                  {step.icon}
                </div>
                <div className="w-6 h-6 rounded-full bg-[#FF9900]/20 border border-[#FF9900]/40 flex items-center justify-center text-xs font-bold text-[#FF9900]">
                  {i + 1}
                </div>
                <h3 className="font-semibold text-[var(--text)]">{step.title}</h3>
                <p className="text-sm text-[var(--text-4)] leading-relaxed">{step.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Features ──────────────────────────────────────────────────────── */}
      <section className="px-6 py-20 border-t border-[var(--border)]">
        <div className="max-w-4xl mx-auto">
          <h2 className="text-2xl font-bold text-center mb-12 text-[var(--text)]">{c.features.title}</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {c.features.items.map((item, i) => (
              <div key={i} className="bg-[var(--code-bg)] border border-[var(--border)] rounded-xl p-5 hover:border-[#FF9900]/30 transition-colors">
                <div className="text-[#FF9900] mb-3">{item.icon}</div>
                <h3 className="font-semibold text-sm text-[var(--text)] mb-1">{item.title}</h3>
                <p className="text-xs text-[var(--text-4)] leading-relaxed">{item.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── CTA ───────────────────────────────────────────────────────────── */}
      <section className="px-6 py-20 border-t border-[var(--border)]">
        <div className="max-w-2xl mx-auto text-center flex flex-col items-center gap-6">
          <h2 className="text-3xl font-bold text-[var(--text)]">{c.cta.title}</h2>
          <p className="text-[var(--text-4)] text-sm leading-relaxed">{c.cta.sub}</p>
          <Link
            href="/editor"
            className="flex items-center gap-2 px-8 py-3 bg-[#FF9900] hover:bg-[#ffb340] text-[var(--accent-contrast)] font-semibold rounded-lg transition-colors text-sm"
          >
            {c.cta.btn}
            <ArrowRight size={15} />
          </Link>
        </div>
      </section>

      {/* ── Footer ────────────────────────────────────────────────────────── */}
      <footer className="border-t border-[var(--border)] px-6 py-6 flex flex-col sm:flex-row items-center sm:justify-between gap-3">
        <div className="flex items-center gap-2">
          <div className="w-5 h-5 bg-[#FF9900] rounded flex items-center justify-center">
            <svg viewBox="0 0 16 16" fill="white" className="w-3 h-3">
              <rect x="1" y="1" width="6" height="6" rx="1" />
              <rect x="9" y="1" width="6" height="6" rx="1" />
              <rect x="1" y="9" width="6" height="6" rx="1" />
              <rect x="9" y="9" width="6" height="6" rx="1" />
            </svg>
          </div>
          <span className="text-xs text-[var(--text-6)]">{c.footer}</span>
        </div>
        <div className="flex items-center gap-4 flex-wrap justify-center">
          <Link href="/editor" className="text-xs text-[var(--text-6)] hover:text-[var(--text-3)] transition-colors flex items-center gap-1">
            <Zap size={11} /> Editor
          </Link>
          <Link href="/builder" className="text-xs text-[var(--text-6)] hover:text-[var(--text-3)] transition-colors flex items-center gap-1">
            <Wrench size={11} /> Builder
          </Link>
          <Link href="/docs" className="text-xs text-[var(--text-6)] hover:text-[var(--text-3)] transition-colors flex items-center gap-1">
            <BookOpen size={11} /> Docs
          </Link>
          <a
            href="https://github.com/fernandofatech/diagram-as-code"
            target="_blank"
            rel="noopener noreferrer"
            className="text-[var(--text-6)] hover:text-[var(--text-3)] transition-colors"
          >
            <Github size={13} />
          </a>
        </div>
      </footer>
    </div>
  )
}

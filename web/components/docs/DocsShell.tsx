'use client'

import { useState } from 'react'
import Link from 'next/link'
import {
  ArrowLeft,
  BookOpen,
  ChevronRight,
  Code2,
  Cpu,
  ExternalLink,
  FileCode2,
  Github,
  Globe,
  Layers,
  Linkedin,
  Menu,
  Server,
  ShieldCheck,
  Sparkles,
  Terminal,
  TriangleAlert,
  User,
  Workflow,
  Wrench,
  X,
} from 'lucide-react'
import LanguageSwitcher from '@/components/LanguageSwitcher'
import ThemeSwitcher from '@/components/ThemeSwitcher'
import type { DocsCopy, DocsSection, DocsSectionId } from '@/lib/docs'

const iconMap: Record<DocsSectionId, React.ReactNode> = {
  overview: <BookOpen size={14} />,
  'quick-start': <Sparkles size={14} />,
  'web-editor': <Globe size={14} />,
  builder: <Wrench size={14} />,
  'yaml-reference': <FileCode2 size={14} />,
  cli: <Terminal size={14} />,
  api: <Code2 size={14} />,
  drawio: <Layers size={14} />,
  mcp: <Cpu size={14} />,
  'local-dev': <Server size={14} />,
  troubleshooting: <TriangleAlert size={14} />,
  examples: <Workflow size={14} />,
  architecture: <ShieldCheck size={14} />,
  about: <User size={14} />,
}

interface DocsShellProps {
  copy: DocsCopy
  sections: DocsSection[]
  activeId?: DocsSectionId
  children: React.ReactNode
}

export default function DocsShell({ copy, sections, activeId, children }: DocsShellProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const byCategory = copy.categories.map((category) => ({
    category,
    sections: sections.filter((section) => section.category === category),
  }))

  return (
    <div className="flex flex-col h-screen bg-[var(--bg)] overflow-hidden">
      <header className="flex items-center justify-between px-4 h-12 border-b border-[var(--border)] flex-shrink-0">
        <div className="flex items-center gap-3 min-w-0">
          <Link
            href="/"
            className="flex items-center gap-1.5 text-xs text-[var(--text-4)] hover:text-[var(--text-2)] transition-colors"
          >
            <ArrowLeft size={13} />
            {copy.ui.back}
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
            <span className="text-sm font-semibold text-[var(--text)] truncate">{copy.ui.title}</span>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <LanguageSwitcher />
          <ThemeSwitcher />
          <a
            href="https://github.com/fernandofatech/diagram-as-code"
            target="_blank"
            rel="noopener noreferrer"
            className="hidden sm:flex items-center gap-1.5 text-xs text-[var(--text-4)] hover:text-[var(--text-2)] transition-colors"
          >
            <Github size={14} />
            {copy.ui.github}
          </a>
          <a
            href="https://dac.moretes.com"
            target="_blank"
            rel="noopener noreferrer"
            className="hidden sm:flex items-center gap-1.5 text-xs text-[#FF9900] hover:text-[#ffb340] transition-colors"
          >
            <ExternalLink size={12} />
            {copy.ui.site}
          </a>
          {/* Mobile sidebar toggle */}
          <button
            onClick={() => setSidebarOpen(o => !o)}
            className="lg:hidden p-1.5 rounded hover:bg-[var(--surface)] text-[var(--text-3)] transition-colors"
            aria-label="Toggle navigation"
          >
            {sidebarOpen ? <X size={18} /> : <Menu size={18} />}
          </button>
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        {/* Mobile overlay */}
        {sidebarOpen && (
          <div
            className="fixed inset-0 z-30 bg-black/50 lg:hidden"
            onClick={() => setSidebarOpen(false)}
          />
        )}
        <aside className={`
          fixed lg:relative inset-y-0 left-0 z-40
          w-72 flex-shrink-0 border-r border-[var(--border)] overflow-y-auto py-4
          bg-[var(--bg)] lg:bg-transparent
          transition-transform duration-200
          top-12 lg:top-auto h-[calc(100vh-3rem)] lg:h-auto
          ${sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
        `}>
          <p className="px-4 text-[10px] text-[var(--text-6)] uppercase tracking-widest font-medium mb-2">
            {copy.ui.contents}
          </p>
          <nav className="space-y-4 px-2">
            <Link
              href="/docs"
              onClick={() => setSidebarOpen(false)}
              className={`w-full flex items-center gap-2 px-3 py-2 rounded text-xs text-left transition-colors ${
                !activeId ? 'bg-[#FF9900]/10 text-[#FF9900]' : 'text-[var(--text-4)] hover:text-[var(--text-2)] hover:bg-[var(--surface)]'
              }`}
            >
              <BookOpen size={14} />
              {copy.ui.overview}
            </Link>

            {byCategory.map(({ category, sections: grouped }) => (
              <div key={category}>
                <p className="px-2 pb-1 text-[10px] text-[var(--text-6)] uppercase tracking-widest">
                  {category}
                </p>
                <div className="space-y-0.5">
                  {grouped.map((section) => (
                    <Link
                      key={section.id}
                      href={`/docs/${section.id}`}
                      onClick={() => setSidebarOpen(false)}
                      className={`w-full flex items-center gap-2 px-3 py-2 rounded text-xs text-left transition-colors ${
                        activeId === section.id
                          ? 'bg-[#FF9900]/10 text-[#FF9900]'
                          : 'text-[var(--text-4)] hover:text-[var(--text-2)] hover:bg-[var(--surface)]'
                      }`}
                    >
                      {iconMap[section.id]}
                      <span>{section.title}</span>
                      {activeId === section.id ? <ChevronRight size={10} className="ml-auto" /> : null}
                    </Link>
                  ))}
                </div>
              </div>
            ))}
          </nav>

          <div className="mt-6 px-4 pt-4 border-t border-[var(--border)] flex items-center gap-3">
            <a
              href="https://fernando.moretes.com"
              target="_blank"
              rel="noopener noreferrer"
              className="text-[10px] text-[var(--text-6)] hover:text-[var(--text-3)] transition-colors"
            >
              fernando.moretes.com
            </a>
            <a
              href="https://www.linkedin.com/in/fernando-francisco-azevedo/"
              target="_blank"
              rel="noopener noreferrer"
              className="text-[var(--text-6)] hover:text-[var(--text-3)] transition-colors"
              aria-label="LinkedIn"
            >
              <Linkedin size={12} />
            </a>
          </div>
        </aside>

        <main className="flex-1 overflow-y-auto px-4 sm:px-8 py-6 sm:py-8">
          <div className="max-w-5xl">{children}</div>
        </main>
      </div>
    </div>
  )
}

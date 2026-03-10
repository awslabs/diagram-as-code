'use client'

import type { DocsBlock } from '@/lib/docs'

export function DocsH2({ children }: { children: React.ReactNode }) {
  return (
    <h2 className="text-2xl font-semibold text-[var(--text)] mb-4 flex items-center gap-3">
      <span className="w-1.5 h-6 bg-[#FF9900] rounded-full inline-block flex-shrink-0" />
      {children}
    </h2>
  )
}

export function DocsH3({ children }: { children: React.ReactNode }) {
  return <h3 className="text-base font-semibold text-[var(--text-2)] mt-8 mb-2">{children}</h3>
}

export function DocsP({ children }: { children: React.ReactNode }) {
  return <p className="text-[var(--text-3)] text-sm leading-7 mb-3">{children}</p>
}

export function DocsPre({ children, lang }: { children: string; lang?: string }) {
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

export function DocsTable({ headers, rows }: { headers: string[]; rows: string[][] }) {
  return (
    <div className="overflow-x-auto my-4 rounded-xl border border-[var(--border)]">
      <table className="w-full text-xs text-left border-collapse">
        <thead className="bg-[var(--surface)]">
          <tr className="border-b border-[var(--border)]">
            {headers.map((header) => (
              <th key={header} className="px-4 py-3 text-[var(--text-4)] font-medium uppercase tracking-wider text-[10px]">
                {header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, rowIndex) => (
            <tr key={rowIndex} className="border-b last:border-b-0 border-[var(--border)]">
              {row.map((cell, cellIndex) => (
                <td key={cellIndex} className="px-4 py-3 text-[var(--text-2)] align-top font-mono">
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

export function DocsList({ items }: { items: string[] }) {
  return (
    <ul className="space-y-2 mb-4">
      {items.map((item) => (
        <li key={item} className="text-sm text-[var(--text-3)] leading-7 flex gap-3">
          <span className="mt-2 h-1.5 w-1.5 rounded-full bg-[#FF9900] flex-shrink-0" />
          <span>{item}</span>
        </li>
      ))}
    </ul>
  )
}

export function DocsCallout({
  title,
  tone,
  text,
}: {
  title: string
  tone: 'info' | 'tip' | 'warn'
  text: string
}) {
  const styles = {
    info: 'bg-sky-500/8 border-sky-500/20 text-sky-100',
    tip: 'bg-[#FF9900]/5 border-[#FF9900]/20 text-[#ffd7a3]',
    warn: 'bg-amber-500/8 border-amber-500/20 text-amber-100',
  }

  return (
    <div className={`border rounded-xl px-4 py-3 text-xs my-4 ${styles[tone]}`}>
      <p className="font-semibold mb-1">{title}</p>
      <p className="leading-6">{text}</p>
    </div>
  )
}

export function DocsBlocks({ blocks }: { blocks: DocsBlock[] }) {
  return (
    <>
      {blocks.map((block, index) => {
        if (block.type === 'paragraph') return <DocsP key={index}>{block.text}</DocsP>
        if (block.type === 'list') {
          return (
            <div key={index}>
              {block.title ? <DocsH3>{block.title}</DocsH3> : null}
              <DocsList items={block.items} />
            </div>
          )
        }
        if (block.type === 'table') {
          return (
            <div key={index}>
              {block.title ? <DocsH3>{block.title}</DocsH3> : null}
              <DocsTable headers={block.headers} rows={block.rows} />
            </div>
          )
        }
        if (block.type === 'code') {
          return (
            <div key={index}>
              {block.title ? <DocsH3>{block.title}</DocsH3> : null}
              <DocsPre lang={block.lang}>{block.code}</DocsPre>
            </div>
          )
        }

        return <DocsCallout key={index} title={block.title} tone={block.tone} text={block.text} />
      })}
    </>
  )
}

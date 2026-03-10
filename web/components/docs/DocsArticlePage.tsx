'use client'

import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import DocsShell from '@/components/docs/DocsShell'
import { DocsBlocks, DocsH2, DocsP } from '@/components/docs/DocsPrimitives'
import { docsSectionIds, getDocsCopy, getDocsSections, type DocsSectionId } from '@/lib/docs'
import { useLanguage } from '@/lib/i18n'

export default function DocsArticlePage({ slug }: { slug: DocsSectionId }) {
  const { lang } = useLanguage()
  const copy = getDocsCopy(lang)
  const sections = getDocsSections(lang)
  const section = copy.sections[slug]
  const index = docsSectionIds.indexOf(slug)
  const previous = index > 0 ? copy.sections[docsSectionIds[index - 1]] : null
  const next = index < docsSectionIds.length - 1 ? copy.sections[docsSectionIds[index + 1]] : null

  return (
    <DocsShell copy={copy} sections={sections} activeId={slug}>
      <div className="mb-8 rounded-2xl border border-[var(--border)] bg-[var(--code-bg)] p-6">
        <p className="text-[10px] uppercase tracking-[0.18em] text-[var(--text-5)] mb-2">{section.category}</p>
        <DocsH2>{section.title}</DocsH2>
        <DocsP>{section.summary}</DocsP>
      </div>

      <DocsBlocks blocks={section.blocks} />

      <div className="mt-10 pt-6 border-t border-[var(--border)] flex items-center justify-between gap-4">
        {previous ? (
          <Link
            href={`/docs/${previous.id}`}
            className="inline-flex items-center gap-2 text-sm text-[var(--text-3)] hover:text-[var(--text)] transition-colors"
          >
            <ArrowLeft size={14} />
            <span>
              {copy.ui.previous}: {previous.title}
            </span>
          </Link>
        ) : (
          <span />
        )}

        {next ? (
          <Link
            href={`/docs/${next.id}`}
            className="inline-flex items-center gap-2 text-sm text-[var(--text-3)] hover:text-[var(--text)] transition-colors"
          >
            <span>
              {copy.ui.next}: {next.title}
            </span>
            <ArrowRight size={14} />
          </Link>
        ) : null}
      </div>
    </DocsShell>
  )
}

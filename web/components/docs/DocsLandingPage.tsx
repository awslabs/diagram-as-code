'use client'

import Link from 'next/link'
import { useMemo, useState } from 'react'
import { ArrowRight, BookOpen, Search } from 'lucide-react'
import DocsShell from '@/components/docs/DocsShell'
import { DocsH2, DocsP } from '@/components/docs/DocsPrimitives'
import { getDocsCopy, getDocsSections } from '@/lib/docs'
import { useLanguage } from '@/lib/i18n'

export default function DocsLandingPage() {
  const { lang } = useLanguage()
  const copy = getDocsCopy(lang)
  const sections = getDocsSections(lang)
  const [query, setQuery] = useState('')
  const [category, setCategory] = useState<string>(copy.ui.allArticles)

  const filtered = useMemo(() => {
    const lowered = query.trim().toLowerCase()
    return sections.filter((section) => {
      const matchesCategory = category === copy.ui.allArticles || section.category === category
      if (!matchesCategory) return false
      if (!lowered) return true
      const haystack = [section.title, section.summary, ...section.keywords].join(' ').toLowerCase()
      return haystack.includes(lowered)
    })
  }, [category, copy.ui.allArticles, query, sections])

  return (
    <DocsShell copy={copy} sections={sections}>
      <div className="mb-8 rounded-2xl border border-[var(--border)] bg-[var(--code-bg)] p-6">
        <div className="flex items-center gap-2 mb-3 text-[#FF9900]">
          <BookOpen size={16} />
          <span className="text-xs uppercase tracking-[0.18em] font-semibold">{copy.hero.title}</span>
        </div>
        <p className="text-[var(--text)] text-lg leading-8 mb-2">{copy.hero.intro}</p>
        <p className="text-sm text-[var(--text-3)] leading-7">{copy.hero.sub}</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-[1fr_220px] gap-4 mb-10">
        <div className="relative">
          <Search size={15} className="absolute left-4 top-3.5 text-[var(--text-5)]" />
          <input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder={copy.ui.searchPlaceholder}
            className="w-full rounded-xl border border-[var(--border)] bg-[var(--surface)] pl-11 pr-4 py-3 text-sm text-[var(--text-2)] placeholder:text-[var(--text-5)] focus:outline-none focus:border-[#FF9900]/60"
          />
        </div>
        <select
          value={category}
          onChange={(event) => setCategory(event.target.value)}
          className="rounded-xl border border-[var(--border)] bg-[var(--surface)] px-4 py-3 text-sm text-[var(--text-2)] focus:outline-none focus:border-[#FF9900]/60"
        >
          <option value={copy.ui.allArticles}>{copy.ui.allArticles}</option>
          {copy.categories.map((item) => (
            <option key={item} value={item}>
              {item}
            </option>
          ))}
        </select>
      </div>

      <DocsH2>{copy.ui.contents}</DocsH2>
      <DocsP>
        {filtered.length} {filtered.length === 1 ? 'article' : 'articles'} available in the wiki.
      </DocsP>

      {filtered.length === 0 ? (
        <div className="rounded-2xl border border-[var(--border)] bg-[var(--surface)] p-6 text-sm text-[var(--text-3)]">
          {copy.ui.searchEmpty}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {filtered.map((section) => (
            <Link
              key={section.id}
              href={`/docs/${section.id}`}
              className="rounded-2xl border border-[var(--border)] bg-[var(--surface)] p-5 hover:border-[#FF9900]/30 hover:bg-[var(--code-bg)] transition-colors"
            >
              <p className="text-[10px] uppercase tracking-[0.18em] text-[var(--text-5)] mb-2">{section.category}</p>
              <h3 className="text-[var(--text)] font-semibold text-base mb-2">{section.title}</h3>
              <p className="text-sm text-[var(--text-3)] leading-7 mb-4">{section.summary}</p>
              <div className="flex flex-wrap gap-1.5 mb-4">
                {section.keywords.slice(0, 4).map((keyword) => (
                  <span
                    key={keyword}
                    className="text-[10px] bg-[var(--code-bg)] border border-[var(--border)] text-[var(--text-4)] px-2 py-0.5 rounded-full"
                  >
                    {keyword}
                  </span>
                ))}
              </div>
              <span className="inline-flex items-center gap-1 text-xs text-[#FF9900]">
                {copy.ui.openArticle}
                <ArrowRight size={12} />
              </span>
            </Link>
          ))}
        </div>
      )}
    </DocsShell>
  )
}

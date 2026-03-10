'use client'

import { useLanguage } from '@/lib/i18n'

export default function LanguageSwitcher() {
  const { lang, toggle } = useLanguage()

  return (
    <button
      onClick={toggle}
      title={lang === 'en' ? 'Mudar para Português' : 'Switch to English'}
      className="flex items-center gap-1 text-xs text-[var(--text-5)] hover:text-[var(--text)] transition-colors px-2 py-1 rounded hover:bg-[var(--surface)] border border-[var(--border)] font-mono"
    >
      <span className={lang === 'en' ? 'text-[#FF9900]' : 'text-[var(--text-5)]'}>EN</span>
      <span className="text-[var(--text-6)]">/</span>
      <span className={lang === 'pt' ? 'text-[#FF9900]' : 'text-[var(--text-5)]'}>PT</span>
    </button>
  )
}

'use client'

import { useLanguage } from '@/lib/i18n'

export default function LanguageSwitcher() {
  const { lang, toggle } = useLanguage()

  return (
    <button
      onClick={toggle}
      title={lang === 'en' ? 'Mudar para Português' : 'Switch to English'}
      className="flex items-center gap-1 text-xs text-[#555] hover:text-[#e5e5e5] transition-colors px-2 py-1 rounded hover:bg-[#1a1a1a] border border-[#2a2a2a] font-mono"
    >
      <span className={lang === 'en' ? 'text-[#FF9900]' : 'text-[#555]'}>EN</span>
      <span className="text-[#333]">/</span>
      <span className={lang === 'pt' ? 'text-[#FF9900]' : 'text-[#555]'}>PT</span>
    </button>
  )
}

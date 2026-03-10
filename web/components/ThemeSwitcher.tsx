'use client'

import { Sun, Moon } from 'lucide-react'
import { useTheme } from '@/lib/theme'

export default function ThemeSwitcher() {
  const { theme, toggleTheme } = useTheme()
  return (
    <button
      type="button"
      onClick={toggleTheme}
      title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
      className="flex items-center justify-center w-7 h-7 rounded text-[#555] hover:text-[#e5e5e5] hover:bg-[#1a1a1a] transition-colors border border-[#2a2a2a]"
    >
      {theme === 'dark' ? <Sun size={13} /> : <Moon size={13} />}
    </button>
  )
}

'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'

export type Theme = 'dark' | 'light'

type ThemeContextType = { theme: Theme; toggleTheme: () => void }

const ThemeContext = createContext<ThemeContextType>({ theme: 'dark', toggleTheme: () => {} })

function resolveTheme(): Theme {
  if (typeof window === 'undefined') return 'dark'
  const stored = localStorage.getItem('theme')
  if (stored === 'light' || stored === 'dark') return stored
  return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark'
}

function applyTheme(theme: Theme) {
  document.documentElement.setAttribute('data-theme', theme)
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setTheme] = useState<Theme>('dark')

  useEffect(() => {
    const nextTheme = resolveTheme()
    setTheme(nextTheme)
    applyTheme(nextTheme)

    const media = window.matchMedia('(prefers-color-scheme: light)')
    const onChange = () => {
      if (localStorage.getItem('theme')) return
      const systemTheme = media.matches ? 'light' : 'dark'
      setTheme(systemTheme)
      applyTheme(systemTheme)
    }

    const onStorage = (event: StorageEvent) => {
      if (event.key !== 'theme') return
      const storedTheme = resolveTheme()
      setTheme(storedTheme)
      applyTheme(storedTheme)
    }

    media.addEventListener('change', onChange)
    window.addEventListener('storage', onStorage)

    return () => {
      media.removeEventListener('change', onChange)
      window.removeEventListener('storage', onStorage)
    }
  }, [])

  function toggleTheme() {
    setTheme(prev => {
      const next: Theme = prev === 'dark' ? 'light' : 'dark'
      localStorage.setItem('theme', next)
      applyTheme(next)
      return next
    })
  }

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  )
}

export function useTheme() { return useContext(ThemeContext) }

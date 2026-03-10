'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'

export type Theme = 'dark' | 'light'

type ThemeContextType = { theme: Theme; toggleTheme: () => void }

const ThemeContext = createContext<ThemeContextType>({ theme: 'dark', toggleTheme: () => {} })

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setTheme] = useState<Theme>('dark')

  useEffect(() => {
    const stored = localStorage.getItem('theme') as Theme | null
    const t: Theme = stored === 'light' ? 'light' : 'dark'
    setTheme(t)
    document.documentElement.setAttribute('data-theme', t)
  }, [])

  function toggleTheme() {
    setTheme(prev => {
      const next: Theme = prev === 'dark' ? 'light' : 'dark'
      localStorage.setItem('theme', next)
      document.documentElement.setAttribute('data-theme', next)
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

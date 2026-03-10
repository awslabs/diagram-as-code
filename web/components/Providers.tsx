'use client'

import { LanguageProvider } from '@/lib/i18n'
import { ThemeProvider } from '@/lib/theme'

export default function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider>
      <LanguageProvider>{children}</LanguageProvider>
    </ThemeProvider>
  )
}

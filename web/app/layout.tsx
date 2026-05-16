import type { Metadata } from 'next'
import { Analytics } from '@vercel/analytics/react'
import { SpeedInsights } from '@vercel/speed-insights/next'
import './globals.css'
import Providers from '@/components/Providers'

export const metadata: Metadata = {
  title: {
    default: 'Diagram as Code | Fernando Moretes',
    template: '%s | Diagram as Code',
  },
  description:
    'Generate AWS architecture diagrams from YAML. Public diagram-as-code tool connected to Fernando Moretes solution architecture portfolio.',
  metadataBase: new URL('https://dac.moretes.com'),
  authors: [{ name: 'Fernando Francisco Azevedo', url: 'https://fernando.moretes.com' }],
  creator: 'Fernando Francisco Azevedo',
  publisher: 'Fernando Francisco Azevedo',
  keywords: [
    'diagram as code',
    'AWS architecture',
    'YAML diagrams',
    'solution architecture',
    'cloud architecture',
    'Fernando Moretes',
  ],
  alternates: {
    canonical: 'https://dac.moretes.com',
    languages: {
      en: 'https://dac.moretes.com',
      'pt-BR': 'https://dac.moretes.com',
    },
  },
  openGraph: {
    type: 'website',
    title: 'Diagram as Code | Fernando Moretes',
    description:
      'Generate AWS architecture diagrams from YAML and document cloud designs as code.',
    url: 'https://dac.moretes.com',
    siteName: 'Diagram as Code',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Diagram as Code | Fernando Moretes',
    description:
      'Generate AWS architecture diagrams from YAML and document cloud designs as code.',
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
  icons: { icon: '/favicon.ico' },
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const structuredData = {
    '@context': 'https://schema.org',
    '@type': 'WebApplication',
    name: 'Diagram as Code',
    url: 'https://dac.moretes.com',
    description:
      'Generate AWS architecture diagrams from YAML and document cloud architecture as code.',
    applicationCategory: 'DeveloperApplication',
    operatingSystem: 'Web Browser',
    author: {
      '@type': 'Person',
      name: 'Fernando Francisco Azevedo',
      url: 'https://fernando.moretes.com',
      sameAs: [
        'https://www.linkedin.com/in/fernando-francisco-azevedo/',
        'https://github.com/fernandofatech',
        'https://fernando.moretes.com',
      ],
    },
  }

  return (
    <html lang="en" className="h-full">
      <head>
        {/* Prevent flash of wrong theme */}
        <script
          dangerouslySetInnerHTML={{
            __html: `(function(){try{var t=localStorage.getItem('theme');var r=t==='light'||t==='dark'?t:(window.matchMedia('(prefers-color-scheme: light)').matches?'light':'dark');document.documentElement.setAttribute('data-theme',r)}catch(e){document.documentElement.setAttribute('data-theme','dark')}})()`,
          }}
        />
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(structuredData) }}
        />
      </head>
      <body className="h-full">
        <Providers>{children}</Providers>
        <Analytics />
        <SpeedInsights />
      </body>
    </html>
  )
}

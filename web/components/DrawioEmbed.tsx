'use client'

import { useEffect, useRef, useState } from 'react'

const DRAWIO_EMBED_URL = 'https://embed.diagrams.net/?embed=1&proto=json&spin=1'
const DRAWIO_ORIGIN = 'https://embed.diagrams.net'

interface DrawioEmbedProps {
  xml: string
  onChange?: (xml: string) => void
}

type DrawioMessage =
  | { event: 'init' | 'ready' }
  | { event: 'load'; xml?: string }
  | { event: 'save' | 'autosave'; xml: string }
  | { event: 'exit' }
  | { event: 'configure' }
  | { event: string; xml?: string }

export default function DrawioEmbed({ xml, onChange }: DrawioEmbedProps) {
  const iframeRef = useRef<HTMLIFrameElement | null>(null)
  const loadedXmlRef = useRef('')
  const [isReady, setIsReady] = useState(false)
  const [loadError, setLoadError] = useState<string | null>(null)

  function postMessage(message: Record<string, unknown>) {
    iframeRef.current?.contentWindow?.postMessage(JSON.stringify(message), DRAWIO_ORIGIN)
  }

  function loadDiagram(nextXml: string) {
    loadedXmlRef.current = nextXml
    postMessage({
      action: 'load',
      xml: nextXml,
      autosave: 1,
      title: 'diagram.drawio',
      modified: 'unsavedChanges',
    })
  }

  useEffect(() => {
    function handleMessage(event: MessageEvent) {
      if (event.origin !== DRAWIO_ORIGIN) {
        return
      }

      let payload: DrawioMessage | null = null

      if (typeof event.data === 'string') {
        try {
          payload = JSON.parse(event.data) as DrawioMessage
        } catch {
          return
        }
      } else if (typeof event.data === 'object' && event.data !== null) {
        payload = event.data as DrawioMessage
      }

      if (!payload) {
        return
      }

      if (payload.event === 'configure') {
        postMessage({
          action: 'configure',
          config: {
            defaultFonts: ['Helvetica', 'Arial'],
          },
        })
        return
      }

      if (payload.event === 'init' || payload.event === 'ready') {
        setIsReady(true)
        setLoadError(null)
        loadDiagram(xml)
        return
      }

      if ((payload.event === 'save' || payload.event === 'autosave') && typeof payload.xml === 'string') {
        loadedXmlRef.current = payload.xml
        onChange?.(payload.xml)
        return
      }

      if (payload.event === 'exit') {
        setIsReady(false)
      }
    }

    window.addEventListener('message', handleMessage)
    return () => window.removeEventListener('message', handleMessage)
  }, [xml, onChange])

  useEffect(() => {
    if (!isReady || xml === loadedXmlRef.current) {
      return
    }
    loadDiagram(xml)
  }, [isReady, xml])

  return (
    <div className="flex h-full flex-col">
      <div className="border-b border-[var(--border)] px-4 py-2 text-xs text-[var(--text-5)]">
        Edit directly in diagrams.net inside the app. Changes are kept in the current session.
      </div>
      <div className="relative flex-1 bg-[var(--surface)]">
        <iframe
          ref={iframeRef}
          src={DRAWIO_EMBED_URL}
          title="Embedded draw.io editor"
          className="h-full w-full border-0"
          onError={() => setLoadError('Failed to load diagrams.net editor')}
        />
        {!isReady && !loadError ? (
          <div className="pointer-events-none absolute inset-0 flex items-center justify-center bg-[var(--bg)]/70 text-sm text-[var(--text-4)]">
            Loading diagrams.net editor…
          </div>
        ) : null}
        {loadError ? (
          <div className="absolute inset-0 flex items-center justify-center px-6 text-center text-sm text-red-400">
            {loadError}
          </div>
        ) : null}
      </div>
    </div>
  )
}

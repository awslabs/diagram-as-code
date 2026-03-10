'use client'

import { Download, FileCode2, FileText, ImageIcon } from 'lucide-react'
import { useLanguage } from '@/lib/i18n'
import DrawioEmbed from '@/components/DrawioEmbed'

interface DiagramPreviewProps {
  imageUrl: string | null
  drawioContent: string | null
  pdfUrl: string | null
  loading: boolean
  error: string | null
  onDrawioChange?: (xml: string) => void
}

export default function DiagramPreview({
  imageUrl,
  drawioContent,
  pdfUrl,
  loading,
  error,
  onDrawioChange,
}: DiagramPreviewProps) {
  const { t } = useLanguage()

  function downloadDrawio() {
    if (!drawioContent) return
    const blob = new Blob([drawioContent], { type: 'application/xml' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'diagram.drawio'
    a.click()
    URL.revokeObjectURL(url)
  }

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-3 text-[var(--text-4)]">
        <div className="w-8 h-8 border-2 border-[#FF9900] border-t-transparent rounded-full animate-spin" />
        <span className="text-sm">{t.generating}</span>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-3 px-8">
        <div className="w-full max-w-lg bg-red-950/40 border border-red-800/60 rounded-lg p-4">
          <p className="text-red-400 text-sm font-medium mb-1">{t.generationFailed}</p>
          <p className="text-red-300/80 text-xs font-mono whitespace-pre-wrap break-words">{error}</p>
        </div>
      </div>
    )
  }

  if (imageUrl) {
    return (
      <div className="flex flex-col h-full">
        <div className="flex-1 overflow-auto flex items-center justify-center p-4">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src={imageUrl}
            alt={t.generatedDiagram}
            className="max-w-full max-h-full object-contain rounded shadow-lg"
          />
        </div>
        <div className="border-t border-[var(--border)] px-4 py-2 flex items-center justify-end">
          <a
            href={imageUrl}
            download="diagram.png"
            className="flex items-center gap-2 text-xs text-[#FF9900] hover:text-[#ffb340] transition-colors px-3 py-1.5 rounded border border-[#FF9900]/30 hover:border-[#FF9900]/60"
          >
            <Download size={13} />
            {t.downloadPng}
          </a>
        </div>
      </div>
    )
  }

  if (drawioContent) {
    return (
      <div className="flex flex-col h-full">
        <div className="flex-1 overflow-hidden">
          <DrawioEmbed xml={drawioContent} onChange={onDrawioChange} />
        </div>
        <div className="border-t border-[var(--border)] px-4 py-2 flex items-center justify-between">
          <div className="flex items-center gap-2 text-xs text-[var(--text-4)]">
            <FileCode2 size={14} className="text-[#4a9eff]" />
            <span>{t.drawioEmbedded}</span>
          </div>
          <button
            type="button"
            onClick={downloadDrawio}
            className="flex items-center gap-2 text-xs text-[#FF9900] hover:text-[#ffb340] transition-colors px-3 py-1.5 rounded border border-[#FF9900]/30 hover:border-[#FF9900]/60"
          >
            <Download size={13} />
            {t.downloadDrawio}
          </button>
        </div>
      </div>
    )
  }

  if (pdfUrl) {
    return (
      <div className="flex flex-col h-full">
        <div className="flex-1 flex flex-col items-center justify-center gap-4">
          <div className="w-14 h-14 bg-[#5b3a1e] rounded-xl flex items-center justify-center">
            <FileText size={28} className="text-[#ffb340]" />
          </div>
          <div className="text-center">
            <p className="text-sm text-[var(--text)] font-medium">{t.pdfReady}</p>
            <p className="text-xs text-[var(--text-4)] mt-1">{t.pdfHelper}</p>
          </div>
          <a
            href={pdfUrl}
            download="diagram.pdf"
            className="flex items-center gap-2 text-sm text-[#FF9900] hover:text-[#ffb340] transition-colors px-4 py-2 rounded border border-[#FF9900]/30 hover:border-[#FF9900]/60"
          >
            <Download size={14} />
            {t.downloadPdf}
          </a>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col items-center justify-center h-full gap-3 text-[var(--text-6)]">
      <ImageIcon size={40} strokeWidth={1} />
      <p className="text-sm">
        {t.emptyStatePre} <span className="text-[#FF9900]">{t.generate}</span> {t.emptyStatePost}
      </p>
    </div>
  )
}

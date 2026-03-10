'use client'

import { Download, FileCode2, ImageIcon } from 'lucide-react'
import { useLanguage } from '@/lib/i18n'

interface DiagramPreviewProps {
  imageUrl: string | null
  drawioContent: string | null
  loading: boolean
  error: string | null
}

export default function DiagramPreview({
  imageUrl,
  drawioContent,
  loading,
  error,
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
      <div className="flex flex-col items-center justify-center h-full gap-3 text-[#666]">
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
        <div className="border-t border-[#2a2a2a] px-4 py-2 flex items-center justify-end">
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
        <div className="flex-1 flex flex-col items-center justify-center gap-4">
          <div className="w-14 h-14 bg-[#1e3a5f] rounded-xl flex items-center justify-center">
            <FileCode2 size={28} className="text-[#4a9eff]" />
          </div>
          <div className="text-center">
            <p className="text-sm text-[#e5e5e5] font-medium">{t.drawioReady}</p>
            <p className="text-xs text-[#666] mt-1">{t.drawioHelper}</p>
          </div>
          <button
            type="button"
            onClick={downloadDrawio}
            className="flex items-center gap-2 text-sm text-[#FF9900] hover:text-[#ffb340] transition-colors px-4 py-2 rounded border border-[#FF9900]/30 hover:border-[#FF9900]/60"
          >
            <Download size={14} />
            {t.downloadDrawio}
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col items-center justify-center h-full gap-3 text-[#444]">
      <ImageIcon size={40} strokeWidth={1} />
      <p className="text-sm">
        {t.emptyStatePre} <span className="text-[#FF9900]">{t.generate}</span> {t.emptyStatePost}
      </p>
    </div>
  )
}

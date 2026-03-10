'use client'

import { useState, useEffect, useCallback } from 'react'
import { X, ChevronRight, ChevronLeft } from 'lucide-react'

export interface TourStep {
  target: string        // data-tour="xxx" attribute value
  title: string
  description: string
  position?: 'top' | 'bottom' | 'left' | 'right'
}

interface Props {
  id: string            // unique tour ID stored in localStorage
  steps: TourStep[]
  onDone?: () => void
}

interface Rect { top: number; left: number; width: number; height: number }

export default function Tour({ id, steps, onDone }: Props) {
  const [active, setActive] = useState(false)
  const [step, setStep] = useState(0)
  const [rect, setRect] = useState<Rect | null>(null)
  const storageKey = `tour_done_${id}`

  useEffect(() => {
    if (!localStorage.getItem(storageKey)) {
      const t = setTimeout(() => setActive(true), 600)
      return () => clearTimeout(t)
    }
  }, [storageKey])

  const updateRect = useCallback(() => {
    const target = steps[step]?.target
    if (!target) return
    const el = document.querySelector(`[data-tour="${target}"]`) as HTMLElement | null
    if (!el) { setRect(null); return }
    el.classList.add('tour-highlight')
    const r = el.getBoundingClientRect()
    setRect({ top: r.top, left: r.left, width: r.width, height: r.height })
  }, [step, steps])

  useEffect(() => {
    if (!active) return
    // Remove highlight from previous
    document.querySelectorAll('.tour-highlight').forEach(el => el.classList.remove('tour-highlight'))
    updateRect()
    window.addEventListener('resize', updateRect)
    return () => window.removeEventListener('resize', updateRect)
  }, [active, step, updateRect])

  function done() {
    document.querySelectorAll('.tour-highlight').forEach(el => el.classList.remove('tour-highlight'))
    localStorage.setItem(storageKey, '1')
    setActive(false)
    onDone?.()
  }

  function next() { if (step < steps.length - 1) setStep(s => s + 1); else done() }
  function prev() { if (step > 0) setStep(s => s - 1) }

  if (!active || steps.length === 0) return null

  const current = steps[step]
  const GAP = 12

  // Tooltip position relative to highlighted element
  let tipStyle: React.CSSProperties = { position: 'fixed', zIndex: 1100 }
  if (rect) {
    const pos = current.position ?? 'bottom'
    const tipW = 300
    if (pos === 'bottom') {
      tipStyle = { ...tipStyle, top: rect.top + rect.height + GAP, left: Math.max(8, Math.min(rect.left, window.innerWidth - tipW - 8)) }
    } else if (pos === 'top') {
      tipStyle = { ...tipStyle, bottom: window.innerHeight - rect.top + GAP, left: Math.max(8, Math.min(rect.left, window.innerWidth - tipW - 8)) }
    } else if (pos === 'right') {
      tipStyle = { ...tipStyle, top: rect.top, left: rect.left + rect.width + GAP }
    } else {
      tipStyle = { ...tipStyle, top: rect.top, right: window.innerWidth - rect.left + GAP }
    }
  } else {
    // Centered fallback
    tipStyle = { ...tipStyle, top: '50%', left: '50%', transform: 'translate(-50%,-50%)' }
  }

  return (
    <>
      {/* Overlay */}
      <div
        className="fixed inset-0 bg-black/50 z-[1000] pointer-events-none"
        style={rect ? {
          background: `radial-gradient(ellipse at ${rect.left + rect.width / 2}px ${rect.top + rect.height / 2}px, transparent ${Math.max(rect.width, rect.height) * 0.6}px, rgba(0,0,0,0.55) ${Math.max(rect.width, rect.height) * 0.65}px)`,
        } : {}}
      />

      {/* Tooltip card */}
      <div style={{ ...tipStyle, width: 300 }}
        className="bg-[var(--surface)] border border-[#FF9900]/40 rounded-xl p-4 shadow-2xl"
      >
        <div className="flex items-start justify-between gap-2 mb-2">
          <h3 className="text-sm font-semibold text-[var(--text)] leading-snug">{current.title}</h3>
          <button type="button" onClick={done} className="text-[var(--text-5)] hover:text-[var(--text-3)] transition-colors shrink-0 mt-0.5">
            <X size={14} />
          </button>
        </div>
        <p className="text-xs text-[var(--text-3)] leading-relaxed mb-4">{current.description}</p>

        <div className="flex items-center justify-between">
          {/* Step dots */}
          <div className="flex items-center gap-1.5">
            {steps.map((_, i) => (
              <div key={i} className={`w-1.5 h-1.5 rounded-full transition-colors ${i === step ? 'bg-[#FF9900]' : 'bg-[var(--border-strong)]'}`} />
            ))}
          </div>

          <div className="flex items-center gap-1">
            {step > 0 && (
              <button type="button" onClick={prev}
                className="flex items-center gap-1 text-xs text-[var(--text-4)] hover:text-[var(--text-2)] transition-colors px-2 py-1 rounded hover:bg-[var(--surface-hover)]">
                <ChevronLeft size={12} /> Prev
              </button>
            )}
            <button type="button" onClick={next}
              className="flex items-center gap-1 text-xs bg-[#FF9900] hover:bg-[#ffb340] text-[var(--accent-contrast)] font-semibold px-3 py-1 rounded transition-colors">
              {step === steps.length - 1 ? 'Done' : 'Next'} {step < steps.length - 1 && <ChevronRight size={12} />}
            </button>
          </div>
        </div>
      </div>
    </>
  )
}

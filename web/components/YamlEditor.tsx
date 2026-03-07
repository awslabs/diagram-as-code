'use client'

import Editor from '@monaco-editor/react'

interface YamlEditorProps {
  value: string
  onChange: (value: string) => void
}

export default function YamlEditor({ value, onChange }: YamlEditorProps) {
  return (
    <Editor
      language="yaml"
      value={value}
      onChange={(v) => onChange(v ?? '')}
      theme="vs-dark"
      options={{
        fontSize: 13,
        lineHeight: 20,
        fontFamily: "'JetBrains Mono', 'Fira Code', 'Cascadia Code', Consolas, monospace",
        minimap: { enabled: false },
        scrollBeyondLastLine: false,
        wordWrap: 'on',
        renderWhitespace: 'boundary',
        tabSize: 2,
        padding: { top: 16, bottom: 16 },
        scrollbar: { verticalScrollbarSize: 6, horizontalScrollbarSize: 6 },
      }}
    />
  )
}

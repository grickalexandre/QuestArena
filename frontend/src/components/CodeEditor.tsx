import { useEffect, useLayoutEffect, useMemo, useRef, type KeyboardEvent } from 'react'
import { tokenizeLines } from '../lib/highlight'

const INDENT = '  '

type Props = {
  value: string
  onChange: (value: string) => void
  language?: string
  placeholder?: string
  minRows?: number
}

export default function CodeEditor({
  value,
  onChange,
  language = 'plain',
  placeholder,
  minRows = 8,
}: Props) {
  const taRef = useRef<HTMLTextAreaElement>(null)
  const preRef = useRef<HTMLPreElement>(null)
  const pendingSel = useRef<[number, number] | null>(null)
  const lines = useMemo(() => tokenizeLines(value, language), [value, language])

  useLayoutEffect(() => {
    const sel = pendingSel.current
    if (!sel) return
    pendingSel.current = null
    taRef.current?.setSelectionRange(sel[0], sel[1])
  })

  useEffect(() => {
    syncScroll()
  }, [value])

  function syncScroll() {
    const ta = taRef.current
    const pre = preRef.current
    if (!ta || !pre) return
    pre.scrollTop = ta.scrollTop
    pre.scrollLeft = ta.scrollLeft
  }

  function apply(next: string, selStart: number, selEnd = selStart) {
    pendingSel.current = [selStart, selEnd]
    onChange(next)
  }

  function onKeyDown(e: KeyboardEvent<HTMLTextAreaElement>) {
    const ta = e.currentTarget
    const start = ta.selectionStart
    const end = ta.selectionEnd

    if (e.key === 'Tab') {
      e.preventDefault()
      if (start === end && !e.shiftKey) {
        apply(value.slice(0, start) + INDENT + value.slice(end), start + INDENT.length)
        return
      }
      const blockStart = value.lastIndexOf('\n', start - 1) + 1
      const nextBreak = value.indexOf('\n', end)
      const blockEnd = nextBreak === -1 ? value.length : nextBreak
      const block = value.slice(blockStart, blockEnd)
      const shifted = block
        .split('\n')
        .map((line) =>
          e.shiftKey ? line.replace(/^ {1,2}|^\t/, '') : line.length ? INDENT + line : line,
        )
        .join('\n')
      apply(
        value.slice(0, blockStart) + shifted + value.slice(blockEnd),
        blockStart,
        blockStart + shifted.length,
      )
      return
    }

    if (e.key === 'Enter') {
      const lineStart = value.lastIndexOf('\n', start - 1) + 1
      const currentLine = value.slice(lineStart, start)
      const indent = currentLine.match(/^[ \t]*/)?.[0] ?? ''
      const extra = /[{[(:]\s*$/.test(currentLine) ? INDENT : ''
      if (!indent && !extra) return
      e.preventDefault()
      const insert = `\n${indent}${extra}`
      apply(value.slice(0, start) + insert + value.slice(end), start + insert.length)
    }
  }

  return (
    <div className="code-editor">
      <pre className="code-editor-highlight" ref={preRef} aria-hidden="true">
        <code>
          {lines.map((tokens, i) => (
            <span className="code-line" key={i}>
              <span className="code-ln">{i + 1}</span>
              <span className="code-src">
                {tokens.map((token, j) => (
                  <span key={j} className={`tk-${token.k}`}>
                    {token.t}
                  </span>
                ))}
                {'\n'}
              </span>
            </span>
          ))}
        </code>
      </pre>
      <textarea
        ref={taRef}
        className="code-editor-input"
        value={value}
        rows={Math.max(minRows, lines.length + 1)}
        wrap="off"
        spellCheck={false}
        autoCapitalize="off"
        autoCorrect="off"
        placeholder={placeholder}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={onKeyDown}
        onScroll={syncScroll}
      />
    </div>
  )
}

import { useMemo, useState } from 'react'
import { languageLabel, tokenizeLines } from '../lib/highlight'

type Props = {
  code: string
  language?: string
  showLineNumbers?: boolean
  copyable?: boolean
  compact?: boolean
}

export default function CodeBlock({
  code,
  language = 'plain',
  showLineNumbers = true,
  copyable = true,
  compact = false,
}: Props) {
  const lines = useMemo(() => tokenizeLines(code, language), [code, language])
  const [copied, setCopied] = useState(false)

  async function copy() {
    try {
      await navigator.clipboard.writeText(code)
      setCopied(true)
      setTimeout(() => setCopied(false), 1600)
    } catch {
      setCopied(false)
    }
  }

  return (
    <div className={`code-block ${compact ? 'compact' : ''}`}>
      <div className="code-block-head">
        <span className="code-lang">{languageLabel(language)}</span>
        {copyable && (
          <button type="button" className="code-copy" onClick={copy}>
            {copied ? 'Copiado' : 'Copiar'}
          </button>
        )}
      </div>
      <pre className="code-block-body">
        <code>
          {lines.map((tokens, i) => (
            <span className="code-line" key={i}>
              {showLineNumbers && (
                <span className="code-ln" aria-hidden="true">
                  {i + 1}
                </span>
              )}
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
    </div>
  )
}

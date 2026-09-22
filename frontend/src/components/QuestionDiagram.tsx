const BANNED = /<\s*(script|foreignobject|iframe|object|embed|use|animate|set|handler|style)\b|javascript:|data:/i
const EVENT_ATTR = /\s+on[a-z]+\s*=\s*("[^"]*"|'[^']*'|[^\s>]+)/gi

function safeDiagram(svg: string) {
  const raw = svg.trim()
  if (!raw.toLowerCase().startsWith('<svg') || !raw.toLowerCase().includes('</svg>')) {
    return ''
  }
  if (BANNED.test(raw)) {
    return ''
  }
  return raw.replace(EVENT_ATTR, '')
}

export default function QuestionDiagram({ svg, caption }: { svg?: string; caption?: string }) {
  const clean = svg ? safeDiagram(svg) : ''
  if (!clean) return null
  return (
    <figure className="question-der">
      <div className="question-der-inner" dangerouslySetInnerHTML={{ __html: clean }} />
      {caption ? <figcaption>{caption}</figcaption> : null}
    </figure>
  )
}

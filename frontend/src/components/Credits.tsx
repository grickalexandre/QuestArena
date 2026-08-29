export const AUTHOR_NAME = 'Professor Luis Alexandre de Oliveira'

export default function Credits({ compact = false }: { compact?: boolean }) {
  return (
    <p className={`credits ${compact ? 'compact' : ''}`}>
      <span className="credits-label">Créditos</span>
      <strong>{AUTHOR_NAME}</strong>
      <span>Análise e Desenvolvimento de Sistemas</span>
    </p>
  )
}

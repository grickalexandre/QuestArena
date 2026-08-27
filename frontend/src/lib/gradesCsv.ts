export type GradeCsvRow = {
  ra?: string
  nickname: string
  correctCount?: number
  total?: number
  score: number
  maxScore?: number
  grade?: number
}

export function formatGrade(n?: number) {
  if (n == null || Number.isNaN(n)) return '—'
  return n.toFixed(1).replace('.', ',')
}

function csvCell(value: string | number) {
  const s = String(value)
  if (/[;"\n]/.test(s)) return `"${s.replace(/"/g, '""')}"`
  return s
}

export function downloadGradesCsv(quizTitle: string, pin: string, grades: GradeCsvRow[]) {
  const header = ['RA', 'Apelido', 'Acertos', 'Total', 'XP', 'XP_max', 'Nota']
  const lines = [
    '# Nota = 70% acertos + 30% XP (resposta mais rapida vale mais). Some as planilhas pelo RA.',
    `# Quiz: ${quizTitle || 'QuestArena'} | PIN: ${pin}`,
    header.join(';'),
    ...grades.map((g) =>
      [
        csvCell(g.ra || ''),
        csvCell(g.nickname),
        g.correctCount ?? '',
        g.total ?? '',
        g.score,
        g.maxScore ?? '',
        formatGrade(g.grade),
      ].join(';'),
    ),
  ]
  const blob = new Blob(['\uFEFF' + lines.join('\n')], { type: 'text/csv;charset=utf-8' })
  const slug = (quizTitle || 'quiz')
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-zA-Z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
    .slice(0, 40)
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `notas-${slug || 'quiz'}-PIN${pin}.csv`
  a.click()
  URL.revokeObjectURL(a.href)
}

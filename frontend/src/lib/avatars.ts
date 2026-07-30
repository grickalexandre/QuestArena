export const AVATARS = [
  { id: 0, emoji: '🧙', label: 'Mago' },
  { id: 1, emoji: '🥷', label: 'Ninja' },
  { id: 2, emoji: '🐉', label: 'Dragão' },
  { id: 3, emoji: '🦊', label: 'Raposa' },
  { id: 4, emoji: '🤖', label: 'Robô' },
  { id: 5, emoji: '🦸', label: 'Herói' },
  { id: 6, emoji: '👻', label: 'Fantasma' },
  { id: 7, emoji: '🐺', label: 'Lobo' },
  { id: 8, emoji: '🐱', label: 'Gato' },
  { id: 9, emoji: '🦉', label: 'Coruja' },
  { id: 10, emoji: '🦁', label: 'Leão' },
  { id: 11, emoji: '🐸', label: 'Sapo' },
] as const

export const AVATAR_COUNT = AVATARS.length

export function avatarEmoji(id: number | undefined | null): string {
  if (id == null || id < 0) return '🎮'
  return AVATARS[id % AVATARS.length]?.emoji ?? '🎮'
}

export function avatarLabel(id: number | undefined | null): string {
  if (id == null || id < 0) return 'Avatar'
  return AVATARS[id % AVATARS.length]?.label ?? 'Avatar'
}

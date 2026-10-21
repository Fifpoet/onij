const EMOJIS = [
  '🎧', '🎹', '🎸', '🥁', '🎤', '📻', '🌙', '⭐',
  '🔥', '💧', '🌸', '🍀', '🦊', '🐱', '🐻', '🐼',
  '🐸', '🐧', '🦋', '🐝', '🍉', '🍋', '🍑', '🍓',
  '🍪', '☕', '🍵', '🧋', '🎪', '🎨', '🌈', '⚡',
  '❄️', '🌊', '🍁', '🍄', '🧸', '🎯', '🎲', '💿',
]

const BACKGROUNDS = [
  '#FDE68A', '#FECACA', '#BFDBFE', '#A7F3D0',
  '#DDD6FE', '#FBCFE8', '#FED7AA', '#C7D2FE',
  '#99F6E4', '#FDE68A', '#E9D5FF', '#BAE6FD',
]

function hashSeed(seed: string): number {
  let h = 2166136261
  for (let i = 0; i < seed.length; i++) {
    h ^= seed.charCodeAt(i)
    h = Math.imul(h, 16777619)
  }
  return h >>> 0
}

export type PresenceAvatarVisual = {
  emoji: string
  background: string
}

/** 由 client_id 决定，所有人看到的同一用户头像一致 */
export function presenceAvatar(seed: string): PresenceAvatarVisual {
  const h = hashSeed(seed || 'onij')
  return {
    emoji: EMOJIS[h % EMOJIS.length]!,
    background: BACKGROUNDS[Math.floor(h / EMOJIS.length) % BACKGROUNDS.length]!,
  }
}

import type { PitchNote } from '@/api/uvr'

const NOTE_NAMES = ['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B'] as const

/** MIDI 60 = C4（科学音高记谱） */
export function midiToNoteName(midi: number): string {
  const m = Math.round(midi)
  const name = NOTE_NAMES[((m % 12) + 12) % 12]
  const octave = Math.floor(m / 12) - 1
  return `${name}${octave}`
}

export function isNaturalC(midi: number): boolean {
  return ((Math.round(midi) % 12) + 12) % 12 === 0
}

/** 后端已切段；前端只去掉非法段，不再 sticky，避免半音台阶被抹平。 */
export function mergePitchNotes(notes: PitchNote[]): PitchNote[] {
  if (!notes.length) return notes
  return notes
    .filter((n) => n.t1 > n.t0)
    .map((n) => ({
      t0: n.t0,
      t1: n.t1,
      midi: Math.round(n.midi),
    }))
    .sort((a, b) => a.t0 - b.t0)
}
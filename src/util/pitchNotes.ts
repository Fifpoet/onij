import type { PitchNote } from '@/api/uvr'

const NOTE_NAMES = ['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B'] as const
const NATURAL_PC = new Set([0, 2, 4, 5, 7, 9, 11])

/** MIDI 60 = C4（科学音高记谱） */
export function midiToNoteName(midi: number): string {
  const m = Math.round(midi)
  const name = NOTE_NAMES[((m % 12) + 12) % 12]
  const octave = Math.floor(m / 12) - 1
  return `${name}${octave}`
}

export function pitchClass(midi: number): number {
  return ((Math.round(midi) % 12) + 12) % 12
}

export function isNaturalNote(midi: number): boolean {
  return NATURAL_PC.has(pitchClass(midi))
}

export function isNaturalC(midi: number): boolean {
  return pitchClass(midi) === 0
}

/** 后端 v5 已切段并滤 Live；前端只 round，避免再滤掉轻开口。 */
export function mergePitchNotes(notes: PitchNote[]): PitchNote[] {
  if (!notes.length) return notes
  return notes
    .filter((n) => n.t1 > n.t0)
    .map((n) => ({ t0: n.t0, t1: n.t1, midi: Math.round(n.midi) }))
    .sort((a, b) => a.t0 - b.t0)
}

export type PitchPhrase = {
  t0: number
  t1: number
  notes: PitchNote[]
}

/** 无歌词时：相邻 F0 横条间隔超过该值则另起一句 */
const PHRASE_GAP_SEC = 0.8

export function groupPitchPhrases(notes: PitchNote[], lyricStarts?: number[]): PitchPhrase[] {
  const list = mergePitchNotes(notes)
  if (!list.length) return []
  const starts = [...new Set((lyricStarts ?? []).filter((t) => Number.isFinite(t)))].sort((a, b) => a - b)
  if (starts.length >= 2) return groupByLyricStarts(list, starts)

  const phrases: PitchPhrase[] = []
  let cur: PitchNote[] = [list[0]!]
  for (let i = 1; i < list.length; i++) {
    const n = list[i]!
    const prev = cur[cur.length - 1]!
    if (n.t0 - prev.t1 > PHRASE_GAP_SEC) {
      phrases.push({ t0: cur[0]!.t0, t1: prev.t1, notes: cur })
      cur = [n]
    } else {
      cur.push(n)
    }
  }
  phrases.push({ t0: cur[0]!.t0, t1: cur[cur.length - 1]!.t1, notes: cur })
  return phrases
}

function groupByLyricStarts(list: PitchNote[], starts: number[]): PitchPhrase[] {
  const phrases: PitchPhrase[] = []
  for (let i = 0; i < starts.length; i++) {
    const a = starts[i]!
    const b = starts[i + 1] ?? Number.POSITIVE_INFINITY
    const chunk = list.filter((n) => {
      const mid = (n.t0 + n.t1) / 2
      return mid >= a && mid < b
    })
    if (!chunk.length) continue
    phrases.push({ t0: chunk[0]!.t0, t1: chunk[chunk.length - 1]!.t1, notes: chunk })
  }
  return phrases
}

export function pitchNoteKey(n: PitchNote): string {
  return `${n.t0.toFixed(3)}:${n.t1.toFixed(3)}:${n.midi}`
}

const RANGE_PAD = 0.7
const MIN_NOTE_DUR = 0.08

/** 一屏铺满本曲实际最高/最低横条，只留一点边 */
export function songPitchExtent(notes: PitchNote[]): { lo: number; hi: number } {
  const solid = notes.filter((n) => n.t1 - n.t0 >= MIN_NOTE_DUR)
  const src = solid.length ? solid : notes
  if (!src.length) return { lo: 55, hi: 72 }
  let lo = src[0]!.midi
  let hi = src[0]!.midi
  for (const n of src) {
    if (n.midi < lo) lo = n.midi
    if (n.midi > hi) hi = n.midi
  }
  lo -= RANGE_PAD
  hi += RANGE_PAD
  if (hi - lo < 7) {
    const mid = (lo + hi) / 2
    lo = mid - 3.5
    hi = mid + 3.5
  }
  return { lo, hi }
}

/** 左侧标尺：只标自然音（C D E F G A B），E–F / B–C 仍是半音间距 */
export function pitchScaleLabels(lo: number, hi: number): number[] {
  const a = Math.ceil(lo)
  const b = Math.floor(hi)
  const out: number[] = []
  for (let m = a; m <= b; m++) {
    if (isNaturalNote(m)) out.push(m)
  }
  return out.length ? out : [Math.round((lo + hi) / 2)]
}
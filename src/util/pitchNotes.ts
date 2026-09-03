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

/** Live：丢掉观众尖叫（短、离主体音域远）和讲话（密、跳、时值短） */
function cleanLiveNotes(notes: PitchNote[]): PitchNote[] {
  const list = notes
    .filter((n) => n.t1 > n.t0)
    .map((n) => ({ t0: n.t0, t1: n.t1, midi: Math.round(n.midi) }))
    .sort((a, b) => a.t0 - b.t0)
  if (list.length < 4) return list

  const long = list.filter((n) => n.t1 - n.t0 >= 0.16)
  const coreSrc = long.length >= 3 ? long : list
  const items = coreSrc.map((n) => ({ midi: n.midi, w: n.t1 - n.t0 }))
  const lo = weightedPercentile(items, 0.08)
  const hi = weightedPercentile(items, 0.92)

  const kept: PitchNote[] = []
  for (const n of list) {
    const dur = n.t1 - n.t0
    if ((n.midi < lo - 4 || n.midi > hi + 4) && dur < 0.22) continue
    if (n.midi < lo - 7 || n.midi > hi + 7) continue
    kept.push(n)
  }

  const drop = new Set<number>()
  for (let i = 0; i < kept.length; i++) {
    const win0 = kept[i]!.t0 - 0.12
    const win1 = kept[i]!.t1 + 0.12
    const local: number[] = []
    for (let j = 0; j < kept.length; j++) {
      const n = kept[j]!
      if (n.t1 < win0 || n.t0 > win1) continue
      local.push(j)
    }
    if (local.length < 4) continue
    const durs = local.map((j) => kept[j]!.t1 - kept[j]!.t0)
    const meanDur = durs.reduce((a, b) => a + b, 0) / durs.length
    let jumps = 0
    for (let k = 1; k < local.length; k++) {
      if (Math.abs(kept[local[k]!]!.midi - kept[local[k - 1]!]!.midi) >= 3) jumps += 1
    }
    if (meanDur < 0.14 && jumps >= 3) {
      for (const j of local) {
        if (kept[j]!.t1 - kept[j]!.t0 < 0.2) drop.add(j)
      }
    }
  }
  return kept.filter((_, i) => !drop.has(i))
}

/** 后端已切段；前端清 Live 杂质后只 round。 */
export function mergePitchNotes(notes: PitchNote[]): PitchNote[] {
  if (!notes.length) return notes
  return cleanLiveNotes(notes)
}

export type PitchPhrase = {
  t0: number
  t1: number
  notes: PitchNote[]
}

/** 相邻 F0 横条间隔超过该值则另起一句 */
const PHRASE_GAP_SEC = 0.45

export function groupPitchPhrases(notes: PitchNote[]): PitchPhrase[] {
  const list = mergePitchNotes(notes)
  if (!list.length) return []
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

export function pitchNoteKey(n: PitchNote): string {
  return `${n.t0.toFixed(3)}:${n.t1.toFixed(3)}:${n.midi}`
}

const RANGE_PAD = 2
const MIN_NOTE_DUR = 0.08

function weightedPercentile(items: { midi: number; w: number }[], p: number): number {
  const sorted = [...items].sort((a, b) => a.midi - b.midi)
  const total = sorted.reduce((s, i) => s + i.w, 0)
  if (total <= 0) return sorted[0]?.midi ?? 60
  const target = total * p
  let acc = 0
  for (const i of sorted) {
    acc += i.w
    if (acc >= target) return i.midi
  }
  return sorted[sorted.length - 1]!.midi
}

/** 按 F0 横条时长加权取主体音域，忽略短暂飞点，只留一点边距 */
export function songPitchExtent(notes: PitchNote[]): { lo: number; hi: number } {
  const solid = notes.filter((n) => n.t1 - n.t0 >= MIN_NOTE_DUR)
  const src = solid.length ? solid : notes
  if (!src.length) return { lo: 55, hi: 72 }
  const items = src.map((n) => ({ midi: n.midi, w: Math.max(n.t1 - n.t0, 0.02) }))
  let lo = weightedPercentile(items, 0.04)
  let hi = weightedPercentile(items, 0.96)
  if (hi < lo) {
    const t = lo
    lo = hi
    hi = t
  }
  lo -= RANGE_PAD
  hi += RANGE_PAD
  if (hi - lo < 10) {
    const mid = (lo + hi) / 2
    lo = mid - 5
    hi = mid + 5
  }
  return { lo, hi }
}

/** 左侧标尺：范围内每个半音 */
export function pitchScaleLabels(lo: number, hi: number): number[] {
  const a = Math.ceil(lo)
  const b = Math.floor(hi)
  const out: number[] = []
  for (let m = a; m <= b; m++) out.push(m)
  return out.length ? out : [Math.round((lo + hi) / 2)]
}
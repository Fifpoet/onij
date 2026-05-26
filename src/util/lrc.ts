export interface LrcLine {
  time: number
  text: string
}

const TAG_RE = /^\[(ar|ti|al|by|offset|length|id|hash|sign|qq|total):/i

/** 解析 LRC 文本为带时间戳的行（多时间轴一行会拆成多行） */
export function parseLrc(raw: string | null | undefined): LrcLine[] {
  if (!raw || !raw.trim()) return []
  const out: LrcLine[] = []
  for (const line of raw.split(/\r?\n/)) {
    const t = line.trim()
    if (!t || TAG_RE.test(t)) continue
    const timeRe = /\[(\d{2}):(\d{2})(?:\.(\d{1,3}))?\]/g
    const times: number[] = []
    let m: RegExpExecArray | null
    while ((m = timeRe.exec(t)) !== null) {
      const min = parseInt(m[1], 10)
      const sec = parseInt(m[2], 10)
      const ms = m[3] ? parseInt(m[3].padEnd(3, '0').slice(0, 3), 10) / 1000 : 0
      times.push(min * 60 + sec + ms)
    }
    if (!times.length) continue
    const last = t.lastIndexOf(']')
    const text = last >= 0 ? t.slice(last + 1).trim() : t
    for (const time of times) out.push({ time, text })
  }
  out.sort((a, b) => a.time - b.time || a.text.localeCompare(b.text))
  return out
}

/** KTV 双行：先播完上行再播下行，两行都结束后随 activeIdx 切到下一对 */
export interface KtvDualLineState {
  baseIdx: number
  line1: string
  line2: string
  line1Progress: number
  line2Progress: number
}

function lineEndTime(lines: LrcLine[], idx: number, currentSec: number): number {
  const start = lines[idx].time
  const next = lines[idx + 1]?.time
  if (next != null && next > start) return next
  return Number.isFinite(currentSec) && currentSec > start ? currentSec + 2 : start + 4
}

export function ktvDualLineState(lines: LrcLine[], currentSec: number): KtvDualLineState {
  const empty: KtvDualLineState = {
    baseIdx: -1,
    line1: '',
    line2: '',
    line1Progress: 0,
    line2Progress: 0,
  }
  if (!lines.length) return empty

  const t = Number(currentSec)
  let baseIdx = activeLrcIndex(lines, currentSec)
  if (baseIdx < 0) baseIdx = 0

  const line1 = lines[baseIdx]?.text ?? ''
  const line2 = lines[baseIdx + 1]?.text ?? ''
  const start1 = lines[baseIdx].time
  const end1 = lineEndTime(lines, baseIdx, t)

  let line1Progress = end1 <= start1 ? 1 : Math.min(1, Math.max(0, (t - start1) / (end1 - start1)))

  if (!line2) {
    return { baseIdx, line1, line2: '', line1Progress, line2Progress: 0 }
  }

  const start2 = lines[baseIdx + 1].time
  const end2 = lineEndTime(lines, baseIdx + 1, t)
  let line2Progress = 0
  if (line1Progress >= 1) {
    line2Progress = end2 <= start2 ? 1 : Math.min(1, Math.max(0, (t - start2) / (end2 - start2)))
  }

  return { baseIdx, line1, line2, line1Progress, line2Progress }
}

/** 当前行内播放进度 0–1（用于卡拉 OK 逐行高亮） */
export function activeLineProgress(lines: LrcLine[], currentSec: number): number {
  const idx = activeLrcIndex(lines, currentSec)
  if (idx < 0) return 0
  const start = lines[idx].time
  const end =
    lines[idx + 1]?.time ??
    (Number.isFinite(currentSec) && currentSec > start ? currentSec + 2 : start + 4)
  if (end <= start) return 1
  return Math.min(1, Math.max(0, (currentSec - start) / (end - start)))
}

/** 当前时间对应最后一行歌词下标 */
export function activeLrcIndex(lines: LrcLine[], currentSec: number): number {
  let idx = -1
  const t = Number(currentSec)
  if (!Number.isFinite(t)) return -1
  for (let i = 0; i < lines.length; i++) {
    if (lines[i].time <= t + 0.2) idx = i
    else break
  }
  return idx
}

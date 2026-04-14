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

/** 人声基频：16k 降采样 + YIN，并优先取基频而非谐波。 */

export type PitchDetectResult = {
  hz: number
  midi: number
  prob: number
  rms: number
}

const MIN_HZ = 80
const MAX_HZ = 520
const TARGET_SR = 16000
const YIN_THRESHOLD = 0.12

export function hzToMidi(hz: number): number {
  return 69 + 12 * Math.log2(hz / 440)
}

export function midiCents(a: number, b: number): number {
  return (a - b) * 100
}

/** 折到最近八度再比音分，低/高八度唱同一音级分数相同 */
export function octaveAwareCents(user: number, ref: number): number {
  const semi = user - ref
  return (semi - 12 * Math.round(semi / 12)) * 100
}

/** 0~100：±30 音分满分，±200 音分零分 */
export function scoreCents(cents: number): number {
  const a = Math.abs(cents)
  if (a <= 30) return 100
  if (a >= 200) return 0
  return Math.round(100 * (1 - (a - 30) / 170))
}

export function bufferRms(buf: Float32Array): number {
  let s = 0
  for (let i = 0; i < buf.length; i++) {
    const v = buf[i] ?? 0
    s += v * v
  }
  return Math.sqrt(s / Math.max(buf.length, 1))
}

/** 有参考音时一律折到最近八度，避免噪声把音高钉在顶/底 */
export function snapOctave(midi: number, ref: number | null): number {
  if (ref == null) return midi
  return midi - 12 * Math.round((midi - ref) / 12)
}

function downsample(src: Float32Array, sampleRate: number, dest: Float32Array): number {
  if (sampleRate <= TARGET_SR * 1.05) {
    const n = Math.min(src.length, dest.length)
    dest.set(src.subarray(0, n))
    return n
  }
  const step = sampleRate / TARGET_SR
  const n = Math.min(dest.length, Math.floor(src.length / step))
  for (let i = 0; i < n; i++) {
    const a = Math.floor(i * step)
    const b = Math.min(src.length, Math.floor((i + 1) * step))
    let s = 0
    for (let j = a; j < b; j++) s += src[j] ?? 0
    dest[i] = s / Math.max(b - a, 1)
  }
  return n
}

export class VoicePitchTracker {
  private work = new Float32Array(1024)
  private d = new Float32Array(512)
  private cmnd = new Float32Array(512)
  private hist: number[] = []

  detect(buf: Float32Array, sampleRate: number): PitchDetectResult | null {
    const rms = bufferRms(buf)
    if (rms < 0.012) return null

    const n = downsample(buf, sampleRate, this.work)
    const sr = sampleRate <= TARGET_SR * 1.05 ? sampleRate : TARGET_SR
    const raw = this.yin(this.work, n, sr)
    if (!raw) return null

    let midi = raw.midi
    this.hist.push(midi)
    if (this.hist.length > 5) this.hist.shift()
    if (this.hist.length >= 3) {
      const sorted = [...this.hist].sort((a, b) => a - b)
      midi = sorted[Math.floor(sorted.length / 2)]!
    }

    return { hz: raw.hz, midi, prob: raw.prob, rms }
  }

  private yin(buf: Float32Array, n: number, sampleRate: number): PitchDetectResult | null {
    const tauMin = Math.max(2, Math.floor(sampleRate / MAX_HZ))
    const tauMax = Math.min(n - 2, Math.floor(sampleRate / MIN_HZ))
    if (tauMax <= tauMin + 2) return null
    if (this.d.length < tauMax + 1) {
      this.d = new Float32Array(tauMax + 1)
      this.cmnd = new Float32Array(tauMax + 1)
    }

    const d = this.d
    const cmnd = this.cmnd
    for (let tau = 1; tau <= tauMax; tau++) {
      let sum = 0
      const limit = n - tau
      for (let i = 0; i < limit; i++) {
        const diff = (buf[i] ?? 0) - (buf[i + tau] ?? 0)
        sum += diff * diff
      }
      d[tau] = sum
    }

    cmnd[0] = 1
    let running = 0
    for (let tau = 1; tau <= tauMax; tau++) {
      running += d[tau] ?? 0
      cmnd[tau] = (d[tau] ?? 0) * tau / Math.max(running, 1e-12)
    }

    let tauEst = 0
    for (let tau = tauMin; tau < tauMax; tau++) {
      const v = cmnd[tau] ?? 1
      if (v < YIN_THRESHOLD) {
        while (tau + 1 < tauMax && (cmnd[tau + 1] ?? 1) < v) tau += 1
        tauEst = tau
        break
      }
    }
    if (!tauEst) return null

    const doubled = tauEst * 2
    if (doubled < tauMax && (cmnd[doubled] ?? 1) <= (cmnd[tauEst] ?? 1) + 0.06) {
      tauEst = doubled
    }

    const y1 = cmnd[tauEst - 1] ?? 0
    const y2 = cmnd[tauEst] ?? 0
    const y3 = cmnd[tauEst + 1] ?? 0
    const denom = 2 * (2 * y2 - y1 - y3)
    const shift = Math.abs(denom) > 1e-6 ? (y3 - y1) / denom : 0
    const period = tauEst + Math.max(-1, Math.min(1, shift))
    if (period <= 0) return null

    const hz = sampleRate / period
    if (hz < MIN_HZ || hz > MAX_HZ) return null
    const prob = Math.max(0, Math.min(1, 1 - (cmnd[tauEst] ?? 1)))
    if (prob < 0.62) return null

    return { hz, midi: hzToMidi(hz), prob, rms: 0 }
  }
}

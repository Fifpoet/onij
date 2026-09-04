import { onBeforeUnmount, ref, watch, type Ref } from 'vue'
import type { PitchNote } from '@/api/uvr'
import { VoicePitchTracker, midiCents, octaveAwareCents, scoreCents } from '@/util/pitchDetect'
import { groupPitchPhrases, pitchNoteKey, type PitchPhrase } from '@/util/pitchNotes'

const RMS_GATE = 0.012
const HOLD_MS = 450
const CONFIRM_FRAMES = 2
const STICKY_CENTS = 80
const JUMP_CENTS = 180
const IN_TUNE_CENTS = 70
const HIT_SCORE = 60
const PHRASE_GAP_RESET = 1.2

export function useSingPitch(opts: {
  enabled: Ref<boolean>
  currentSec: () => number
  notes: () => PitchNote[]
  songKey: () => string | number
}) {
  const userMidi = ref<number | null>(null)
  const voiced = ref(false)
  const inTune = ref(false)
  const phraseScore = ref<number | null>(null)
  const avgScore = ref<number | null>(null)
  const hitKeys = ref<string[]>([])
  const micError = ref('')
  const listening = ref(false)

  let stream: MediaStream | null = null
  let audioCtx: AudioContext | null = null
  let analyser: AnalyserNode | null = null
  let source: MediaStreamAudioSourceNode | null = null
  let buf: Float32Array | null = null
  let hp: BiquadFilterNode | null = null
  let lp: BiquadFilterNode | null = null
  let raf = 0
  let lastSong = ''
  const tracker = new VoicePitchTracker()

  let stableMidi: number | null = null
  let pendingMidi: number | null = null
  let pendingCount = 0
  let holdUntil = 0
  let lastSec = 0
  let noiseRms = 0.008

  let phrases: PitchPhrase[] = []
  let phraseIdx = -1
  let barIdx = -1
  let barSum = 0
  let barFrames = 0
  let barVoiced = 0
  let barLocked: number[] = []
  let finishedPhrases: number[] = []
  let finishedMask: boolean[] = []
  let hits = new Set<string>()

  function refNoteAt(sec: number): PitchNote | null {
    for (const p of phrases) {
      for (const n of p.notes) {
        if (sec >= n.t0 && sec <= n.t1) return n
      }
    }
    return null
  }

  function phraseAt(sec: number): number {
    for (let i = 0; i < phrases.length; i++) {
      const p = phrases[i]!
      if (sec >= p.t0 && sec <= p.t1) return i
    }
    return -1
  }

  function barAt(phrase: PitchPhrase, sec: number): number {
    for (let i = 0; i < phrase.notes.length; i++) {
      const n = phrase.notes[i]!
      if (sec >= n.t0 && sec <= n.t1) return i
    }
    return -1
  }

  function publishHits() {
    hitKeys.value = [...hits]
  }

  function finalizeBar(phrase: PitchPhrase, idx: number) {
    const note = phrase.notes[idx]
    if (!note) return
    const dur = Math.max(note.t1 - note.t0, 0.04)
    const quality = barFrames > 0 ? barSum / barFrames : 0
    const cover = Math.min(1, barVoiced / (dur * 0.5))
    const locked = Math.round(quality * (0.65 + 0.35 * cover))
    barLocked[idx] = locked
    if (locked >= HIT_SCORE) {
      hits.add(pitchNoteKey(note))
      publishHits()
    }
    barSum = 0
    barFrames = 0
    barVoiced = 0
  }

  function updateAvg() {
    if (!finishedPhrases.length) {
      avgScore.value = null
      return
    }
    const sum = finishedPhrases.reduce((a, b) => a + b, 0)
    avgScore.value = Math.round(sum / finishedPhrases.length)
  }

  function closePhrase() {
    if (phraseIdx < 0) return
    const phrase = phrases[phraseIdx]
    if (!phrase) return
    if (barIdx >= 0) finalizeBar(phrase, barIdx)
    for (let i = 0; i < phrase.notes.length; i++) {
      if (barLocked[i] == null) barLocked[i] = 0
    }
    const vals = phrase.notes.map((_, i) => barLocked[i] ?? 0)
    const locked = Math.round(vals.reduce((a, b) => a + b, 0) / Math.max(vals.length, 1))
    finishedPhrases.push(locked)
    finishedMask[phraseIdx] = true
    phraseScore.value = locked
    updateAvg()
    phraseIdx = -1
    barIdx = -1
    barLocked = []
    barSum = 0
    barFrames = 0
    barVoiced = 0
  }

  function missPhrase(i: number) {
    if (finishedMask[i]) return
    finishedMask[i] = true
    finishedPhrases.push(0)
    phraseScore.value = 0
    updateAvg()
  }

  function enterPhrase(idx: number) {
    if (phraseIdx === idx || finishedMask[idx]) return
    if (phraseIdx >= 0) closePhrase()
    phraseIdx = idx
    barIdx = -1
    barLocked = []
    barSum = 0
    barFrames = 0
    barVoiced = 0
  }

  function advanceByTime(sec: number) {
    for (let i = 0; i < phrases.length; i++) {
      if (finishedMask[i]) continue
      if (sec > (phrases[i]?.t1 ?? Infinity)) {
        if (phraseIdx === i) closePhrase()
        else missPhrase(i)
      }
    }
    const idx = phraseAt(sec)
    if (idx >= 0 && !finishedMask[idx] && phraseIdx !== idx) enterPhrase(idx)
  }

  function syncPhrases() {
    phrases = groupPitchPhrases(opts.notes())
    finishedMask = phrases.map(() => false)
  }

  function resetScoresOnly() {
    phraseIdx = -1
    barIdx = -1
    barLocked = []
    barSum = 0
    barFrames = 0
    barVoiced = 0
    finishedPhrases = []
    phraseScore.value = null
    avgScore.value = null
    hits = new Set()
    hitKeys.value = []
    lastSec = 0
    syncPhrases()
  }

  function resetScore() {
    stableMidi = null
    pendingMidi = null
    pendingCount = 0
    holdUntil = 0
    lastSec = 0
    inTune.value = false
    userMidi.value = null
    voiced.value = false
    resetScoresOnly()
  }

  function acceptPitch(raw: number, nowMs: number): number {
    holdUntil = nowMs + HOLD_MS
    if (stableMidi == null) {
      stableMidi = raw
      pendingMidi = null
      pendingCount = 0
      return stableMidi
    }

    const delta = Math.abs(midiCents(raw, stableMidi))
    if (delta <= STICKY_CENTS) {
      stableMidi = stableMidi * 0.7 + raw * 0.3
      pendingMidi = null
      pendingCount = 0
      return stableMidi
    }
    if (delta <= JUMP_CENTS) {
      stableMidi = stableMidi * 0.5 + raw * 0.5
      pendingMidi = null
      pendingCount = 0
      return stableMidi
    }
    if (pendingMidi != null && Math.abs(midiCents(raw, pendingMidi)) <= STICKY_CENTS) {
      pendingCount += 1
      pendingMidi = pendingMidi * 0.5 + raw * 0.5
    } else {
      pendingMidi = raw
      pendingCount = 1
    }
    if (pendingCount >= CONFIRM_FRAMES) {
      stableMidi = pendingMidi
      pendingMidi = null
      pendingCount = 0
    }
    return stableMidi
  }

  function scoreFrame(sec: number, midi: number, dt: number) {
    if (phraseIdx < 0) return
    const phrase = phrases[phraseIdx]
    if (!phrase) return
    const b = barAt(phrase, sec)
    if (b < 0) return
    if (barIdx >= 0 && b !== barIdx) finalizeBar(phrase, barIdx)
    barIdx = b
    const note = phrase.notes[b]
    if (!note) return
    const frame = scoreCents(octaveAwareCents(midi, note.midi))
    barSum += frame
    barFrames += 1
    barVoiced += dt
  }

  function tick() {
    if (!analyser || !buf) return
    analyser.getFloatTimeDomainData(buf)
    const now = opts.currentSec()
    const nowMs = performance.now()
    const moved = lastSec > 0 ? now - lastSec : 0
    if (lastSec > 0 && Math.abs(moved) > PHRASE_GAP_RESET) resetScoresOnly()
    const dt = moved > 0 ? Math.min(0.08, moved) : 0
    lastSec = now
    advanceByTime(now)

    const note = refNoteAt(now)
    const ref = note?.midi ?? null
    const found = tracker.detect(buf, audioCtx?.sampleRate || 48000)
    const loud = !!(found && found.rms >= RMS_GATE)
    if (found && !loud) {
      noiseRms = noiseRms * 0.96 + found.rms * 0.04
    }
    const aboveNoise = !found || found.rms > noiseRms * 1.35

    let shown: number | null = null
    if (loud && aboveNoise && found) {
      shown = acceptPitch(found.midi, nowMs)
    } else if (nowMs < holdUntil && stableMidi != null) {
      shown = stableMidi
    } else {
      pendingMidi = null
      pendingCount = 0
    }

    if (shown != null) {
      userMidi.value = shown
      voiced.value = true
      if (ref != null) {
        const cents = Math.abs(octaveAwareCents(shown, ref))
        inTune.value = cents <= IN_TUNE_CENTS
        if (dt > 0) scoreFrame(now, shown, dt)
      } else {
        inTune.value = false
      }
    } else {
      voiced.value = false
      inTune.value = false
    }
    raf = requestAnimationFrame(tick)
  }

  async function start() {
    if (listening.value) return
    micError.value = ''
    const gum = navigator.mediaDevices?.getUserMedia?.bind(navigator.mediaDevices)
    if (!gum) {
      micError.value = window.isSecureContext
        ? '当前浏览器不支持麦克风'
        : '麦克风需要 HTTPS，当前 HTTP 页面无法打开（请用 https://onij.fun 或本机 localhost）'
      return
    }
    try {
      stream = await gum({
        audio: {
          echoCancellation: false,
          noiseSuppression: false,
          autoGainControl: false,
          channelCount: 1,
        },
      })
      audioCtx = new AudioContext()
      if (audioCtx.state === 'suspended') void audioCtx.resume()
      hp = audioCtx.createBiquadFilter()
      hp.type = 'highpass'
      hp.frequency.value = 80
      hp.Q.value = 0.7
      lp = audioCtx.createBiquadFilter()
      lp.type = 'lowpass'
      lp.frequency.value = 520
      lp.Q.value = 0.7
      analyser = audioCtx.createAnalyser()
      analyser.fftSize = 2048
      analyser.smoothingTimeConstant = 0
      buf = new Float32Array(analyser.fftSize)
      source = audioCtx.createMediaStreamSource(stream)
      source.connect(hp)
      hp.connect(lp)
      lp.connect(analyser)
      listening.value = true
      syncPhrases()
      raf = requestAnimationFrame(tick)
    } catch (e) {
      micError.value = e instanceof Error ? e.message : '无法打开麦克风'
      listening.value = false
    }
  }

  function stop() {
    cancelAnimationFrame(raf)
    raf = 0
    source?.disconnect()
    hp?.disconnect()
    lp?.disconnect()
    void audioCtx?.close()
    stream?.getTracks().forEach((t) => t.stop())
    source = null
    hp = null
    lp = null
    analyser = null
    audioCtx = null
    stream = null
    buf = null
    listening.value = false
    inTune.value = false
    userMidi.value = null
    voiced.value = false
  }

  watch(
    () => opts.enabled.value,
    (on) => {
      if (on) void start()
      else stop()
    },
    { immediate: true },
  )

  watch(
    () => String(opts.songKey()),
    (key) => {
      if (key !== lastSong) {
        lastSong = key
        resetScore()
      }
    },
    { immediate: true },
  )

  watch(
    () => opts.notes().length,
    () => syncPhrases(),
  )

  onBeforeUnmount(stop)

  return {
    userMidi,
    voiced,
    inTune,
    phraseScore,
    avgScore,
    hitKeys,
    micError,
    listening,
  }
}

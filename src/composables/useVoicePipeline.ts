import { onMounted, onBeforeUnmount, ref } from 'vue'
import * as OpenCC from 'opencc-js'
import { transcribeAudio } from '@/api/voice'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'
import { usePlayQueueStore } from '@/store/playQueue'
import { useAiChatStore } from '@/store/aiChat'

const toSimplified = OpenCC.Converter({ from: 'tw', to: 'cn' })

/** 仅用于判断一句是否说完 / 是否有声，不再作为「要不要进 LLM」的开关 */
const SILENCE_RMS = 0.01
const SPEECH_END_MS = 800
const MIN_RECORD_MS = 400
const MAX_RECORD_MS = 8000
const DUPLICATE_TTL_MS = 2500

type LocalAction = 'player_next' | 'player_toggle'

const WAKE_RE = /(小雪|晓雪|小薛|小靴|xiaoxue)/i

function pickMimeType(): string {
  const candidates = ['audio/webm;codecs=opus', 'audio/webm', 'audio/mp4']
  for (const t of candidates) {
    if (typeof MediaRecorder !== 'undefined' && MediaRecorder.isTypeSupported(t)) return t
  }
  return ''
}

function rmsFromTimeDomain(data: Float32Array): number {
  let sum = 0
  for (let i = 0; i < data.length; i++) {
    const v = data[i] ?? 0
    sum += v * v
  }
  return Math.sqrt(sum / Math.max(data.length, 1))
}

function extractWakeCommand(text: string): string | null {
  const t = text.replace(/\s+/g, '')
  if (!t) return null
  const m = t.match(WAKE_RE)
  if (!m || m.index == null) return null
  if (m.index > 6) return null
  return t.slice(m.index + m[0].length).replace(/^[,，。.!！、？?]+/, '').trim()
}

function matchLocalAction(text: string): LocalAction | null {
  const t = text.replace(/\s+/g, '')
  if (!t) return null
  if (/下一曲|下一首|切歌/.test(t)) return 'player_next'
  if (t === '暂停' || t === '播放' || t === '继续播放' || t === '继续') return 'player_toggle'
  return null
}

export function useVoicePipeline() {
  const transcript = ref('')
  const status = ref('idle')
  const listening = ref(false)
  const error = ref('')

  const playQueue = usePlayQueueStore()
  const aiChat = useAiChatStore()

  let stream: MediaStream | null = null
  let audioCtx: AudioContext | null = null
  let analyser: AnalyserNode | null = null
  let source: MediaStreamAudioSourceNode | null = null
  let rafId = 0
  let recorder: MediaRecorder | null = null
  let chunks: BlobPart[] = []
  let belowSince = 0
  let recordingStartedAt = 0
  let peaked = false
  let busy = false
  let maxTimer: ReturnType<typeof setTimeout> | null = null
  let lastCommand = ''
  let lastCommandAt = 0

  async function handleUtterance(text: string) {
    const command = extractWakeCommand(text)
    if (command == null) return
    if (!command) {
      status.value = 'listening'
      return
    }

    const now = Date.now()
    if (command === lastCommand && now - lastCommandAt < DUPLICATE_TTL_MS) return
    lastCommand = command
    lastCommandAt = now

    const action = matchLocalAction(command)
    if (action === 'player_next') {
      playQueue.skipToNext()
      return
    }
    if (action === 'player_toggle') {
      neteasePlayerControl.togglePlay()
      return
    }

    if (aiChat.sending || aiChat.pendingConfirm) return
    status.value = 'intent'
    await aiChat.sendVoice(command)
  }

  async function processBlob(blob: Blob) {
    if (busy) return
    busy = true
    status.value = 'transcribing'
    try {
      const mime = blob.type || 'audio/webm'
      const ext = mime.includes('mp4') ? 'm4a' : 'webm'
      const result = await transcribeAudio(blob, `voice.${ext}`)
      const text = toSimplified((result.text || '').trim())
      if (!text) {
        status.value = 'listening'
        return
      }
      transcript.value = text
      await handleUtterance(text)
    } catch (e) {
      error.value = e instanceof Error ? e.message : String(e)
      console.warn('voice pipeline failed:', error.value)
    } finally {
      busy = false
      status.value = listening.value ? 'listening' : 'idle'
    }
  }

  function stopRecorder(send: boolean) {
    if (maxTimer) {
      clearTimeout(maxTimer)
      maxTimer = null
    }
    const rec = recorder
    recorder = null
    const hadPeak = peaked
    peaked = false
    if (!rec || rec.state === 'inactive') return
    rec.onstop = () => {
      const blob = new Blob(chunks, { type: rec.mimeType || 'audio/webm' })
      chunks = []
      const dur = Date.now() - recordingStartedAt
      if (send && hadPeak && dur >= MIN_RECORD_MS && blob.size > 800) {
        void processBlob(blob)
      }
    }
    try {
      rec.stop()
    } catch {
      /* ignore */
    }
  }

  function startRecorder() {
    if (!stream || recorder || busy) return
    const mime = pickMimeType()
    try {
      recorder = mime
        ? new MediaRecorder(stream, { mimeType: mime })
        : new MediaRecorder(stream)
    } catch (e) {
      error.value = e instanceof Error ? e.message : String(e)
      return
    }
    chunks = []
    peaked = false
    belowSince = 0
    recordingStartedAt = Date.now()
    status.value = 'recording'
    recorder.ondataavailable = (ev) => {
      if (ev.data.size > 0) chunks.push(ev.data)
    }
    recorder.start(200)
    maxTimer = setTimeout(() => stopRecorder(true), MAX_RECORD_MS)
  }

  function tick() {
    if (!analyser) return
    const buf = new Float32Array(analyser.fftSize)
    analyser.getFloatTimeDomainData(buf)
    const rms = rmsFromTimeDomain(buf)
    const now = performance.now()

    if (!recorder && !busy) {
      startRecorder()
    }

    if (recorder) {
      if (rms >= SILENCE_RMS) {
        peaked = true
        belowSince = 0
      } else if (peaked) {
        if (!belowSince) belowSince = now
        if (now - belowSince >= SPEECH_END_MS) {
          stopRecorder(true)
          belowSince = 0
        }
      }
    }
    rafId = requestAnimationFrame(tick)
  }

  async function startListening() {
    if (listening.value) return
    error.value = ''
    try {
      stream = await navigator.mediaDevices.getUserMedia({
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          channelCount: 1,
        },
      })
      audioCtx = new AudioContext()
      analyser = audioCtx.createAnalyser()
      analyser.fftSize = 2048
      source = audioCtx.createMediaStreamSource(stream)
      source.connect(analyser)
      listening.value = true
      status.value = 'listening'
      rafId = requestAnimationFrame(tick)
    } catch (e) {
      error.value = e instanceof Error ? e.message : String(e)
      console.warn('microphone unavailable:', error.value)
    }
  }

  function stopListening() {
    cancelAnimationFrame(rafId)
    rafId = 0
    stopRecorder(false)
    source?.disconnect()
    void audioCtx?.close()
    stream?.getTracks().forEach((t) => t.stop())
    source = null
    analyser = null
    audioCtx = null
    stream = null
    listening.value = false
    status.value = 'idle'
  }

  onMounted(() => {
    void startListening()
  })

  onBeforeUnmount(() => {
    stopListening()
  })

  return {
    transcript,
    status,
    listening,
    error,
    startListening,
    stopListening,
  }
}

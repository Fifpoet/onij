import { onMounted, onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { parseVoiceIntent, transcribeAudio } from '@/api/voice'
import { getSongDetail } from '@/api/netease/search'
import { usePlayQueueStore } from '@/store/playQueue'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'
import { useSearch } from '@/composables/searchMusic'
import type { ViewMusicListItem } from '@/api/view/music'

const LOUDNESS_THRESHOLD = 0.045
const SPEECH_START_MS = 180
const SPEECH_END_MS = 900
const MIN_RECORD_MS = 450
const MAX_RECORD_MS = 12000

type LocalAction = 'player_next' | 'player_toggle' | 'search'

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

function matchLocalAction(text: string): LocalAction | null {
  const t = text.replace(/\s+/g, '')
  if (!t) return null
  if (/下一曲|下一首|切歌/.test(t)) return 'player_next'
  if (t === '暂停' || t === '播放' || t === '继续播放') return 'player_toggle'
  if (t.includes('搜索')) return 'search'
  return null
}

function buildSearchQuery(args: Record<string, unknown>): string {
  const query = String(args.query ?? '').trim()
  if (query) return query
  const artist = String(args.artist ?? '').trim()
  const title = String(args.title ?? '').trim()
  return [artist, title].filter(Boolean).join(' ').trim()
}

export function useVoicePipeline() {
  const transcript = ref('')
  const status = ref('idle')
  const listening = ref(false)
  const error = ref('')

  const router = useRouter()
  const playQueue = usePlayQueueStore()
  const { searchValue, handleSearch } = useSearch()

  let stream: MediaStream | null = null
  let audioCtx: AudioContext | null = null
  let analyser: AnalyserNode | null = null
  let source: MediaStreamAudioSourceNode | null = null
  let rafId = 0
  let recorder: MediaRecorder | null = null
  let chunks: BlobPart[] = []
  let aboveSince = 0
  let belowSince = 0
  let recordingStartedAt = 0
  let busy = false
  let maxTimer: ReturnType<typeof setTimeout> | null = null

  async function searchAndEnqueueFirst(query: string) {
    searchValue.value = query
    await router.push({ path: '/search', query: { q: query } })
    const results = await handleSearch([1], 0, 10)
    const first = results?.[0]?.result?.songs?.[0]
    if (!first) return
    const detailResp = await getSongDetail([first.id])
    const detail = detailResp.songs?.[0]
    if (!detail) return
    const item: ViewMusicListItem = {
      id: detail.id,
      name: detail.name,
      time_long: detail.dt / 1000,
      album_id: detail.al.id,
      album_name: detail.al.name,
      cover_file_url: detail.al.picUrl,
      artists: detail.ar.map((artist) => ({
        artist_id: artist.id,
        artist_name: artist.name,
      })),
    }
    playQueue.enqueue(item)
  }

  async function handleUtterance(text: string) {
    const action = matchLocalAction(text)
    if (!action) return

    if (action === 'player_next') {
      playQueue.skipToNext()
      return
    }
    if (action === 'player_toggle') {
      neteasePlayerControl.togglePlay()
      return
    }

    status.value = 'intent'
    const intent = await parseVoiceIntent(text)
    const searchAction = intent.actions.find((a) => a.name === 'search_music')
    const query = buildSearchQuery(searchAction?.args ?? {})
    if (!query) return
    await searchAndEnqueueFirst(query)
  }

  async function processBlob(blob: Blob) {
    if (busy) return
    busy = true
    status.value = 'transcribing'
    try {
      const mime = blob.type || 'audio/webm'
      const ext = mime.includes('mp4') ? 'm4a' : 'webm'
      const result = await transcribeAudio(blob, `voice.${ext}`)
      const text = (result.text || '').trim()
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
    if (!rec || rec.state === 'inactive') return
    rec.onstop = () => {
      const blob = new Blob(chunks, { type: rec.mimeType || 'audio/webm' })
      chunks = []
      const dur = Date.now() - recordingStartedAt
      if (send && dur >= MIN_RECORD_MS && blob.size > 800) {
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

    if (rms >= LOUDNESS_THRESHOLD) {
      belowSince = 0
      if (!aboveSince) aboveSince = now
      if (!recorder && !busy && now - aboveSince >= SPEECH_START_MS) {
        startRecorder()
      }
    } else {
      aboveSince = 0
      if (recorder) {
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

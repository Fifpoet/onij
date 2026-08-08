import { post, uploadFile } from '@/util/http'

export type VoiceAction = {
  name: string
  args?: Record<string, unknown>
}

export type VoiceIntentResult = {
  reply: string
  actions: VoiceAction[]
}

export type WhisperTranscribeResult = {
  text: string
  language?: string | null
  duration?: number | null
  elapsed_sec?: number
}

export async function transcribeAudio(blob: Blob, filename = 'voice.webm'): Promise<WhisperTranscribeResult> {
  const form = new FormData()
  form.append('file', blob, filename)
  form.append('language', 'zh')
  form.append('vad_filter', 'true')
  return uploadFile<WhisperTranscribeResult>('/uvr/whisper/transcribe', form, {
    timeout: 120000,
  })
}

export async function parseVoiceIntent(text: string): Promise<VoiceIntentResult> {
  return post('/voice/intent', { text }) as Promise<VoiceIntentResult>
}

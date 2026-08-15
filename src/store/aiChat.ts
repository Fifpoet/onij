import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  AI_SESSION_KIND_VOICE,
  AiChatStream,
  AiConfirmStream,
  ListAiSessions,
  type AiChatStreamHandlers,
  type AiConfirmDecision,
  type AiConfirmDTO,
  type AiSessionDTO,
} from '@/api/ai'
import { useAiClientTools } from '@/composables/useAiClientTools'

export type LiveTool = { id: string; name: string; status: string }
export type AiChatOrigin = 'page' | 'voice'

export const useAiChatStore = defineStore('aiChat', () => {
  const router = useRouter()
  const { runClientAction } = useAiClientTools()

  const sending = ref(false)
  const pendingConfirm = ref<AiConfirmDTO | null>(null)
  const liveTools = ref<LiveTool[]>([])
  const streamSessionId = ref('')
  const focusSessionId = ref('')
  const lastReply = ref('')
  const voiceSessionId = ref('')
  const origin = ref<AiChatOrigin | ''>('')

  const listeners = new Set<AiChatStreamHandlers>()

  function subscribe(h: AiChatStreamHandlers) {
    listeners.add(h)
    return () => {
      listeners.delete(h)
    }
  }

  function emit<K extends keyof AiChatStreamHandlers>(
    key: K,
    ...args: Parameters<NonNullable<AiChatStreamHandlers[K]>>
  ) {
    for (const h of listeners) {
      const fn = h[key]
      if (typeof fn === 'function') {
        void (fn as any)(...args)
      }
    }
  }

  async function ensureVoiceSession(): Promise<string> {
    if (voiceSessionId.value) return voiceSessionId.value
    const res = await ListAiSessions()
    const voice = (res.items || []).find((s) => s.kind === AI_SESSION_KIND_VOICE)
    if (!voice?.id) throw new Error('语音助手会话未就绪')
    voiceSessionId.value = voice.id
    return voice.id
  }

  function rememberVoice(session: AiSessionDTO) {
    if (session.kind === AI_SESSION_KIND_VOICE) {
      voiceSessionId.value = session.id
    }
  }

  function wrapHandlers(fromVoice: boolean): AiChatStreamHandlers {
    return {
      onSession(session) {
        streamSessionId.value = session.id
        rememberVoice(session)
        emit('onSession', session)
      },
      onUser(msg) {
        emit('onUser', msg)
      },
      onDelta(content) {
        emit('onDelta', content)
      },
      onToolCall(payload) {
        liveTools.value = [
          ...liveTools.value.filter((t) => t.id !== payload.id),
          { id: payload.id, name: payload.name, status: '调用中…' },
        ]
        emit('onToolCall', payload)
      },
      onToolResult(payload) {
        liveTools.value = liveTools.value.map((t) =>
          t.id === payload.id ? { ...t, status: payload.ok ? '完成' : '失败' } : t,
        )
        emit('onToolResult', payload)
      },
      async onClientAction(payload) {
        await runClientAction(payload)
      },
      onNeedsConfirm(payload) {
        pendingConfirm.value = payload
        liveTools.value = []
        if (fromVoice) {
          focusSessionId.value = payload.session_id
          if (router.currentRoute.value.path !== '/ai') {
            void router.push({ path: '/ai' })
          }
        }
        emit('onNeedsConfirm', payload)
      },
      onPaused(payload) {
        liveTools.value = []
        emit('onPaused', payload)
      },
      onDone(payload) {
        liveTools.value = []
        pendingConfirm.value = null
        lastReply.value = payload.reply || payload.message?.content || ''
        emit('onDone', payload)
      },
      onError(message) {
        emit('onError', message)
      },
    }
  }

  async function runChat(sessionId: string | null, content: string, from: AiChatOrigin) {
    if (sending.value) throw new Error('正在处理上一条')
    sending.value = true
    origin.value = from
    pendingConfirm.value = from === 'page' ? null : pendingConfirm.value
    if (from === 'page') pendingConfirm.value = null
    liveTools.value = []
    lastReply.value = ''
    let errMsg = ''
    try {
      const handlers = wrapHandlers(from === 'voice')
      await AiChatStream(sessionId, content, {
        ...handlers,
        onError(message) {
          errMsg = message
          handlers.onError?.(message)
        },
      })
      if (errMsg) throw new Error(errMsg)
    } finally {
      sending.value = false
      origin.value = ''
    }
  }

  async function sendVoice(content: string) {
    if (pendingConfirm.value) {
      throw new Error('请先在 AI 页完成确认')
    }
    const id = await ensureVoiceSession()
    await runChat(id, content, 'voice')
  }

  async function runConfirm(confirmId: string, decision: AiConfirmDecision) {
    if (sending.value) throw new Error('正在处理上一条')
    sending.value = true
    origin.value = 'page'
    const confirmSessionId = pendingConfirm.value?.session_id || ''
    pendingConfirm.value = null
    liveTools.value = []
    lastReply.value = ''
    let errMsg = ''
    try {
      const handlers = wrapHandlers(false)
      await AiConfirmStream(confirmId, decision, {
        ...handlers,
        onSession(session) {
          streamSessionId.value = session.id || confirmSessionId
          rememberVoice(session)
          emit('onSession', session)
        },
        onError(message) {
          errMsg = message
          handlers.onError?.(message)
        },
      })
      if (errMsg) throw new Error(errMsg)
    } finally {
      sending.value = false
      origin.value = ''
    }
  }

  return {
    sending,
    pendingConfirm,
    liveTools,
    streamSessionId,
    focusSessionId,
    lastReply,
    voiceSessionId,
    origin,
    subscribe,
    ensureVoiceSession,
    rememberVoice,
    runChat,
    sendVoice,
    runConfirm,
  }
})

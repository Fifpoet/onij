import { getServerBaseUrl, post } from '../util/http'

export type AiSessionDTO = {
  id: string
  title: string
  kind?: number
  token_total: number
  updated_at: number
  created_at: number
}

export const AI_SESSION_KIND_VOICE = 1

export type AiMessageDTO = {
  id: string
  session_id: string
  role: string
  content: string
  tool_name?: string
  tool_call_id?: string
  token_prompt: number
  token_completion: number
  token_total: number
  created_at: number
}

export type AiConfirmOption = {
  id: string
  label: string
  meta?: any
}

export type AiConfirmUI = {
  type: 'yes_no' | 'single_select' | 'multi_select' | string
  title: string
  yes_label?: string
  no_label?: string
  options?: AiConfirmOption[]
  min?: number
  max?: number
}

export type AiConfirmDTO = {
  confirm_id: string
  session_id: string
  tool_call_id: string
  tool_name: string
  draft_args: any
  ui: AiConfirmUI
  expires_at: number
}

export type AiConfirmDecision = {
  type: 'yes_no' | 'single_select' | 'multi_select' | 'cancel'
  value?: boolean
  option_id?: string
  option_ids?: string[]
}

type SessionItemResp = {
  code: number
  message: string
  item: AiSessionDTO
}

type SessionListResp = {
  code: number
  message: string
  items: AiSessionDTO[]
}

type MessageListResp = {
  code: number
  message: string
  items: AiMessageDTO[]
  token_total: number
}

type ChatResp = {
  code: number
  message: string
  session: AiSessionDTO
  messages: AiMessageDTO[]
  reply: string
}

export type AiChatStreamHandlers = {
  onSession?: (session: AiSessionDTO) => void
  onUser?: (msg: AiMessageDTO) => void
  onDelta?: (content: string) => void
  onToolCall?: (payload: { id: string; name: string; arguments: string }) => void
  onToolResult?: (payload: { id: string; name: string; ok: boolean; content: string }) => void
  onClientAction?: (payload: { name: string; args?: Record<string, any> }) => void | Promise<void>
  onNeedsConfirm?: (payload: AiConfirmDTO) => void
  onPaused?: (payload: { session: AiSessionDTO; confirm_id: string }) => void
  onDone?: (payload: { session: AiSessionDTO; message: AiMessageDTO; reply: string }) => void
  onError?: (message: string) => void
}

export const CreateAiSession = async (): Promise<SessionItemResp> => {
  return post('/ai/session/create', {}) as Promise<SessionItemResp>
}

export const ListAiSessions = async (): Promise<SessionListResp> => {
  return post('/ai/session/list', {}) as Promise<SessionListResp>
}

export const ListAiMessages = async (sessionId: string): Promise<MessageListResp> => {
  return post('/ai/session/messages', { session_id: sessionId }) as Promise<MessageListResp>
}

export const DeleteAiSession = async (sessionId: string): Promise<{ code: number; message: string }> => {
  return post('/ai/session/delete', { session_id: sessionId }) as Promise<{ code: number; message: string }>
}

export const AiChat = async (sessionId: string, content: string): Promise<ChatResp> => {
  return post('/ai/chat', { session_id: sessionId, content }, { timeout: 120000 }) as Promise<ChatResp>
}

function parseSseChunk(buffer: string): { events: { event: string; data: string }[]; rest: string } {
  const events: { event: string; data: string }[] = []
  let rest = buffer
  while (true) {
    const sep = rest.indexOf('\n\n')
    if (sep < 0) break
    const block = rest.slice(0, sep)
    rest = rest.slice(sep + 2)
    let event = 'message'
    const dataLines: string[] = []
    for (const line of block.split('\n')) {
      if (line.startsWith('event:')) event = line.slice(6).trim()
      else if (line.startsWith('data:')) dataLines.push(line.slice(5).trimStart())
    }
    if (dataLines.length) events.push({ event, data: dataLines.join('\n') })
  }
  return { events, rest }
}

async function readAiSse(resp: Response, handlers: AiChatStreamHandlers): Promise<void> {
  if (!resp.ok) {
    const text = await resp.text().catch(() => '')
    throw new Error(text || `HTTP ${resp.status}`)
  }
  if (!resp.body) throw new Error('浏览器不支持流式响应')

  const reader = resp.body.getReader()
  const decoder = new TextDecoder('utf-8')
  let buf = ''
  let finished = false

  while (true) {
    const { done, value } = await reader.read()
    if (done) break
    buf += decoder.decode(value, { stream: true })
    const parsed = parseSseChunk(buf)
    buf = parsed.rest
    for (const ev of parsed.events) {
      let payload: any = null
      try {
        payload = JSON.parse(ev.data)
      } catch {
        continue
      }
      if (ev.event === 'session') handlers.onSession?.(payload as AiSessionDTO)
      else if (ev.event === 'user') handlers.onUser?.(payload as AiMessageDTO)
      else if (ev.event === 'delta') handlers.onDelta?.(String(payload?.content ?? ''))
      else if (ev.event === 'tool_call') handlers.onToolCall?.(payload)
      else if (ev.event === 'tool_result') handlers.onToolResult?.(payload)
      else if (ev.event === 'client_action') await handlers.onClientAction?.(payload)
      else if (ev.event === 'needs_confirm') handlers.onNeedsConfirm?.(payload as AiConfirmDTO)
      else if (ev.event === 'paused') {
        finished = true
        handlers.onPaused?.(payload)
      } else if (ev.event === 'done') {
        finished = true
        handlers.onDone?.(payload)
      } else if (ev.event === 'error') {
        finished = true
        handlers.onError?.(String(payload?.message || '流式对话失败'))
      }
    }
  }

  if (!finished) {
    handlers.onError?.('流式连接中断')
  }
}

function authHeaders(): HeadersInit {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    Accept: 'text/event-stream',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  }
}

export async function AiChatStream(
  sessionId: string | null | undefined,
  content: string,
  handlers: AiChatStreamHandlers,
  signal?: AbortSignal,
): Promise<void> {
  const resp = await fetch(`${getServerBaseUrl()}/ai/chat/stream`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({
      session_id: sessionId && sessionId !== '0' ? sessionId : '0',
      content,
    }),
    signal,
  })
  await readAiSse(resp, handlers)
}

export async function AiConfirmStream(
  confirmId: string,
  decision: AiConfirmDecision,
  handlers: AiChatStreamHandlers,
  signal?: AbortSignal,
): Promise<void> {
  const resp = await fetch(`${getServerBaseUrl()}/ai/confirm/stream`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({
      confirm_id: confirmId,
      decision,
    }),
    signal,
  })
  await readAiSse(resp, handlers)
}

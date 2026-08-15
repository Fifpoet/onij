<template>
  <PageContent>
    <div class="ai-page">
      <aside class="ai-sidebar">
        <div class="ai-sidebar__head">
          <h1 class="browse-section-title border-l-4 border-blue-500 pl-2">AI</h1>
          <button type="button" class="ai-btn" :disabled="busy" @click="startNewChat">新对话</button>
        </div>
        <p v-if="sessionsLoading" class="browse-muted px-1 text-sm">加载会话…</p>
        <ul v-else class="ai-session-list">
          <li v-for="s in sessions" :key="s.id" class="ai-session-row">
            <button
              type="button"
              class="ai-session-item"
              :class="{
                'ai-session-item--active': s.id === currentSessionId,
                'ai-session-item--voice': isVoiceSession(s),
              }"
              @click="selectSession(s.id)"
            >
              <span class="ai-session-item__title">{{ s.title || '新对话' }}</span>
              <span class="ai-session-item__meta">{{ isVoiceSession(s) ? '语音 · ' : '' }}{{ formatTokens(s.token_total) }}</span>
            </button>
            <button
              v-if="!isVoiceSession(s)"
              type="button"
              class="ai-session-del"
              title="删除会话"
              :disabled="busy || deletingId === s.id"
              @click.stop="deleteSession(s.id)"
            >
              ×
            </button>
          </li>
          <li v-if="!sessions.length" class="browse-muted px-1 py-4 text-sm">暂无会话</li>
        </ul>
      </aside>

      <section class="ai-main">
        <div class="ai-main__bar">
          <span class="ai-main__title">{{ currentTitle }}</span>
          <span class="ai-main__tokens">本会话 {{ formatTokens(tokenTotal) }}</span>
        </div>

        <div ref="scrollRef" class="ai-messages select-text-all">
          <p v-if="isDraftChat && !messagesLoading" class="browse-muted py-10 text-center text-sm">
            直接输入消息开始新对话；会话会在发送第一条时创建。
          </p>
          <p v-else-if="messagesLoading" class="browse-muted py-10 text-center text-sm">加载消息…</p>
          <template v-else>
            <div
              v-for="m in displayMessages"
              :key="m.id"
              class="ai-bubble"
              :class="`ai-bubble--${m.role}`"
            >
              <div class="ai-bubble__role">{{ roleLabel(m.role, m.tool_name) }}</div>
              <div class="ai-bubble__content" v-html="formatMessageHtml(m)" />
              <div v-if="m.token_total > 0" class="ai-bubble__tokens">
                {{ m.token_prompt }}+{{ m.token_completion }} = {{ m.token_total }} tokens
              </div>
            </div>
            <div v-if="liveTools.length" class="ai-live-tools">
              <div v-for="t in liveTools" :key="t.id" class="ai-live-tool">
                <span class="ai-live-tool__name">{{ t.name }}</span>
                <span class="ai-live-tool__status">{{ t.status }}</span>
              </div>
            </div>
            <div v-if="pendingConfirm" class="ai-confirm">
              <p class="ai-confirm__title">{{ pendingConfirm.ui.title }}</p>
              <template v-if="pendingConfirm.ui.type === 'yes_no'">
                <div class="ai-confirm__actions">
                  <button type="button" class="ai-btn ai-btn--primary" :disabled="sending" @click="confirmYes">
                    {{ pendingConfirm.ui.yes_label || '确认' }}
                  </button>
                  <button type="button" class="ai-btn" :disabled="sending" @click="confirmNo">
                    {{ pendingConfirm.ui.no_label || '取消' }}
                  </button>
                </div>
              </template>
              <template v-else-if="pendingConfirm.ui.type === 'multi_select' || pendingConfirm.ui.type === 'single_select'">
                <ul class="ai-confirm__options">
                  <li v-for="opt in pendingConfirm.ui.options || []" :key="opt.id">
                    <label class="ai-confirm__opt">
                      <input
                        v-if="pendingConfirm.ui.type === 'multi_select'"
                        v-model="selectedOptionIds"
                        type="checkbox"
                        :value="opt.id"
                      >
                      <input
                        v-else
                        v-model="selectedOptionId"
                        type="radio"
                        name="ai-confirm-opt"
                        :value="opt.id"
                      >
                      <span>{{ opt.label }}</span>
                    </label>
                  </li>
                </ul>
                <div class="ai-confirm__actions">
                  <button type="button" class="ai-btn ai-btn--primary" :disabled="sending" @click="confirmSelect">
                    确认选择
                  </button>
                  <button type="button" class="ai-btn" :disabled="sending" @click="confirmNo">取消</button>
                </div>
              </template>
            </div>
            <p v-if="!displayMessages.length && !liveTools.length && !pendingConfirm" class="browse-muted py-10 text-center text-sm">
              发一条消息试试。可试：搜索林俊杰、新建合集、把歌加入合集（需确认）。
            </p>
          </template>
        </div>

        <form class="ai-composer" @submit.prevent="send">
          <textarea
            v-model="draft"
            class="ai-composer__input"
            rows="3"
            placeholder="输入消息…"
            :disabled="sending"
            @keydown.enter.exact.prevent="send"
          />
          <div class="ai-composer__actions">
            <p v-if="errorText" class="ai-error">{{ errorText }}</p>
            <button type="submit" class="ai-btn ai-btn--primary" :disabled="!canSend">
              {{ sending ? '生成中…' : '发送' }}
            </button>
          </div>
        </form>
      </section>
    </div>
  </PageContent>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import PageContent from '@/components/layout/PageContent.vue'
import {
  AI_SESSION_KIND_VOICE,
  DeleteAiSession,
  ListAiMessages,
  ListAiSessions,
  type AiConfirmDecision,
  type AiConfirmDTO,
  type AiMessageDTO,
  type AiSessionDTO,
} from '@/api/ai'
import { useAiChatStore, type LiveTool } from '@/store/aiChat'

const aiChat = useAiChatStore()
const {
  sending,
  pendingConfirm: storePending,
  liveTools: storeLiveTools,
  streamSessionId,
  focusSessionId,
} = storeToRefs(aiChat)

const sessions = ref<AiSessionDTO[]>([])
const messages = ref<AiMessageDTO[]>([])
const selectedOptionIds = ref<string[]>([])
const selectedOptionId = ref('')
const currentSessionId = ref('')
const tokenTotal = ref(0)
const draft = ref('')
const sessionsLoading = ref(false)
const messagesLoading = ref(false)
const deletingId = ref('')
const errorText = ref('')
const scrollRef = ref<HTMLElement | null>(null)

const busy = computed(() => sessionsLoading.value || sending.value || !!deletingId.value)
const isDraftChat = computed(() => !currentSessionId.value)
const pendingConfirm = computed(() => {
  const p = storePending.value
  if (!p || p.session_id !== currentSessionId.value) return null
  return p
})
const liveTools = computed<LiveTool[]>(() => {
  if (streamSessionId.value !== currentSessionId.value) return []
  return storeLiveTools.value
})
const canSend = computed(() => !!draft.value.trim() && !sending.value && !pendingConfirm.value)

function isVoiceSession(s: AiSessionDTO) {
  return s.kind === AI_SESSION_KIND_VOICE
}

function pinSessions(list: AiSessionDTO[]) {
  const voice = list.find((s) => s.kind === AI_SESSION_KIND_VOICE)
  const rest = list.filter((s) => s.kind !== AI_SESSION_KIND_VOICE)
  return voice ? [voice, ...rest] : rest
}

const currentTitle = computed(() => {
  if (!currentSessionId.value) return '新对话'
  const s = sessions.value.find((x) => x.id === currentSessionId.value)
  return s?.title || '新对话'
})
const displayMessages = computed(() =>
  messages.value.filter((m) => {
    if (m.role === 'assistant' && m.tool_name === '__tool_calls__') return false
    return m.role === 'user' || m.role === 'assistant' || m.role === 'tool'
  }),
)

function formatTokens(n: number) {
  if (!n) return '0 tok'
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k tok`
  return `${n} tok`
}

function roleLabel(role: string, toolName?: string) {
  if (role === 'user') return '你'
  if (role === 'assistant') return 'AI'
  if (role === 'tool') return `tool · ${toolName || ''}`
  return role
}

function formatMessageContent(m: AiMessageDTO) {
  if (m.role !== 'tool') return m.content
  try {
    const parsed = JSON.parse(m.content)
    if (parsed?.ok === false) return `失败：${parsed.error || ''} ${parsed.data?.detail || ''}`.trim()
    if (parsed?.data?.side === 'frontend') return `已派发前端动作 ${parsed.data.action || m.tool_name}`
    if (Array.isArray(parsed?.data?.songs)) return `返回 ${parsed.data.songs.length} 首`
    if (parsed?.data?.name) return `合集：${parsed.data.name}`
    if (Array.isArray(parsed?.data)) return `返回 ${parsed.data.length} 条`
    return m.content.length > 180 ? `${m.content.slice(0, 180)}…` : m.content
  } catch {
    return m.content.length > 180 ? `${m.content.slice(0, 180)}…` : m.content
  }
}

function escapeHtml(s: string) {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

function formatMessageHtml(m: AiMessageDTO) {
  let html = escapeHtml(formatMessageContent(m) || '')
  if (m.role === 'assistant' || m.role === 'user') {
    html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
  }
  return html
}

function upsertSession(session: AiSessionDTO) {
  sessions.value = pinSessions([session, ...sessions.value.filter((s) => s.id !== session.id)])
}

async function scrollToBottom() {
  await nextTick()
  const el = scrollRef.value
  if (el) el.scrollTop = el.scrollHeight
}

async function refreshSessions() {
  sessionsLoading.value = true
  errorText.value = ''
  try {
    const res = await ListAiSessions()
    sessions.value = pinSessions(res.items ?? [])
    const voice = sessions.value.find((s) => s.kind === AI_SESSION_KIND_VOICE)
    if (voice) aiChat.rememberVoice(voice)
  } catch (e: any) {
    errorText.value = e?.response?.data || e?.message || '加载会话失败'
  } finally {
    sessionsLoading.value = false
  }
}

function startNewChat() {
  if (sending.value) return
  errorText.value = ''
  currentSessionId.value = ''
  messages.value = []
  selectedOptionIds.value = []
  selectedOptionId.value = ''
  tokenTotal.value = 0
}

async function deleteSession(id: string) {
  if (!id || sending.value || deletingId.value) return
  const target = sessions.value.find((s) => s.id === id)
  if (target && isVoiceSession(target)) return
  if (!window.confirm('确定删除该会话？消息将一并删除。')) return
  deletingId.value = id
  errorText.value = ''
  try {
    await DeleteAiSession(id)
    sessions.value = sessions.value.filter((s) => s.id !== id)
    if (currentSessionId.value === id) startNewChat()
  } catch (e: any) {
    errorText.value = e?.response?.data || e?.message || '删除失败'
  } finally {
    deletingId.value = ''
  }
}

async function selectSession(id: string) {
  if (!id) return
  if (sending.value && id !== focusSessionId.value && id !== streamSessionId.value) return
  currentSessionId.value = id
  messagesLoading.value = true
  errorText.value = ''
  try {
    const res = await ListAiMessages(id)
    messages.value = res.items ?? []
    tokenTotal.value = res.token_total ?? 0
    await scrollToBottom()
  } catch (e: any) {
    errorText.value = e?.response?.data || e?.message || '加载消息失败'
  } finally {
    messagesLoading.value = false
  }
}

function forCurrentSession(sessionId?: string) {
  return !!sessionId && sessionId === currentSessionId.value
}

function bindStreamHandlers() {
  const state = { assistantIdx: -1, streamingId: '' }
  return {
    onSession(session: AiSessionDTO) {
      upsertSession(session)
      if (aiChat.origin === 'page' && !currentSessionId.value) {
        currentSessionId.value = session.id
      }
      if (focusSessionId.value && session.id === focusSessionId.value) {
        currentSessionId.value = session.id
      }
    },
    onUser(msg: AiMessageDTO) {
      if (!forCurrentSession(msg.session_id)) return
      state.assistantIdx = -1
      state.streamingId = `stream-${msg.id}`
      messages.value.push(msg)
      void scrollToBottom()
    },
    onToolCall() {
      if (streamSessionId.value !== currentSessionId.value) return
      void scrollToBottom()
    },
    onToolResult(payload: { id: string; name: string; ok: boolean; content: string }) {
      if (streamSessionId.value !== currentSessionId.value) return
      messages.value.push({
        id: `tool-${payload.id}-${Date.now()}`,
        session_id: currentSessionId.value,
        role: 'tool',
        content: payload.content,
        tool_name: payload.name,
        tool_call_id: payload.id,
        token_prompt: 0,
        token_completion: 0,
        token_total: 0,
        created_at: Date.now(),
      })
      void scrollToBottom()
    },
    onNeedsConfirm(payload: AiConfirmDTO) {
      if (payload.session_id === currentSessionId.value) {
        selectedOptionIds.value = (payload.ui.options || []).map((o) => o.id)
        selectedOptionId.value = payload.ui.options?.[0]?.id || ''
        void scrollToBottom()
      }
    },
    onPaused(payload: { session: AiSessionDTO }) {
      if (payload.session) {
        upsertSession(payload.session)
        if (forCurrentSession(payload.session.id) || payload.session.id === focusSessionId.value) {
          currentSessionId.value = payload.session.id
          tokenTotal.value = payload.session.token_total
        }
      }
    },
    onDelta(chunk: string) {
      if (!chunk || streamSessionId.value !== currentSessionId.value) return
      if (state.assistantIdx < 0) {
        messages.value.push({
          id: state.streamingId || `stream-${Date.now()}`,
          session_id: currentSessionId.value,
          role: 'assistant',
          content: chunk,
          token_prompt: 0,
          token_completion: 0,
          token_total: 0,
          created_at: Date.now(),
        })
        state.assistantIdx = messages.value.length - 1
      } else {
        messages.value[state.assistantIdx] = {
          ...messages.value[state.assistantIdx],
          content: messages.value[state.assistantIdx].content + chunk,
        }
      }
      void scrollToBottom()
    },
    onDone(payload: { session: AiSessionDTO; message: AiMessageDTO }) {
      if (payload.session) {
        upsertSession(payload.session)
        if (forCurrentSession(payload.session.id) || aiChat.origin === 'page') {
          currentSessionId.value = payload.session.id
          tokenTotal.value = payload.session.token_total
        }
      }
      if (payload.message && forCurrentSession(payload.message.session_id || payload.session?.id)) {
        if (state.assistantIdx >= 0) messages.value[state.assistantIdx] = payload.message
        else messages.value.push(payload.message)
      }
      state.assistantIdx = -1
      void scrollToBottom()
    },
    onError(message: string) {
      errorText.value = message
    },
  }
}

async function send() {
  const text = draft.value.trim()
  if (!canSend.value || !text) return
  errorText.value = ''
  draft.value = ''
  const requestSessionId = currentSessionId.value
  try {
    await aiChat.runChat(requestSessionId || null, text, 'page')
  } catch (e: any) {
    draft.value = text
    errorText.value = e?.message || '发送失败'
    if (!requestSessionId) {
      currentSessionId.value = ''
      messages.value = []
    }
  }
}

async function submitConfirm(decision: AiConfirmDecision) {
  if (!pendingConfirm.value || sending.value) return
  const confirmId = pendingConfirm.value.confirm_id
  errorText.value = ''
  try {
    await aiChat.runConfirm(confirmId, decision)
  } catch (e: any) {
    errorText.value = e?.message || '确认失败'
  }
}

function confirmYes() {
  void submitConfirm({ type: 'yes_no', value: true })
}

function confirmNo() {
  void submitConfirm({ type: 'cancel' })
}

function confirmSelect() {
  if (!pendingConfirm.value) return
  if (pendingConfirm.value.ui.type === 'single_select') {
    if (!selectedOptionId.value) {
      errorText.value = '请选择一项'
      return
    }
    void submitConfirm({ type: 'single_select', option_id: selectedOptionId.value })
    return
  }
  if (!selectedOptionIds.value.length) {
    errorText.value = '请至少选择一首'
    return
  }
  void submitConfirm({ type: 'multi_select', option_ids: [...selectedOptionIds.value] })
}

watch(focusSessionId, (id) => {
  if (id && id !== currentSessionId.value) void selectSession(id)
})

watch(pendingConfirm, (p) => {
  if (!p) return
  selectedOptionIds.value = (p.ui.options || []).map((o) => o.id)
  selectedOptionId.value = p.ui.options?.[0]?.id || ''
})

onMounted(() => {
  const unsub = aiChat.subscribe(bindStreamHandlers())
  onBeforeUnmount(unsub)
  void (async () => {
    await refreshSessions()
    if (focusSessionId.value) await selectSession(focusSessionId.value)
  })()
})
</script>

<style scoped>
.ai-page {
  display: grid;
  grid-template-columns: minmax(200px, 260px) minmax(0, 1fr);
  gap: 1rem;
  min-height: calc(100vh - 140px);
}
.ai-sidebar,
.ai-main {
  display: flex;
  flex-direction: column;
  min-height: 0;
  border: 1px solid rgb(229 231 235);
  border-radius: 0.75rem;
  background: rgb(255 255 255 / 0.7);
}
.dark .ai-sidebar,
.dark .ai-main {
  border-color: rgb(55 65 81);
  background: rgb(31 41 55 / 0.55);
}
.ai-sidebar__head {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 0.85rem 0.9rem 0.5rem;
}
.ai-session-list {
  list-style: none;
  margin: 0;
  padding: 0.35rem 0.5rem 0.75rem;
  overflow-y: auto;
  flex: 1;
}
.ai-session-row {
  display: flex;
  align-items: stretch;
  gap: 0.15rem;
  margin-bottom: 0.25rem;
}
.ai-session-item {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.15rem;
  padding: 0.55rem 0.65rem;
  border: 0;
  border-radius: 0.5rem;
  background: transparent;
  text-align: left;
  cursor: pointer;
  color: inherit;
}
.ai-session-item:hover,
.ai-session-item--active {
  background: rgb(243 244 246);
}
.dark .ai-session-item:hover,
.dark .ai-session-item--active {
  background: rgb(55 65 81);
}
.ai-session-item--voice .ai-session-item__title {
  font-weight: 600;
}
.ai-session-del {
  flex-shrink: 0;
  width: 1.75rem;
  border: 0;
  border-radius: 0.5rem;
  background: transparent;
  color: rgb(156 163 175);
  font-size: 1.1rem;
  line-height: 1;
  cursor: pointer;
}
.ai-session-del:hover:not(:disabled) {
  background: rgb(254 226 226);
  color: rgb(185 28 28);
}
.dark .ai-session-del:hover:not(:disabled) {
  background: rgb(127 29 29 / 0.45);
  color: rgb(252 165 165);
}
.ai-session-del:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
.ai-session-item__title {
  font-size: 0.9rem;
  font-weight: 500;
  line-height: 1.35;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}
.ai-session-item__meta {
  font-size: 0.75rem;
  color: rgb(156 163 175);
}
.ai-main__bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.85rem 1rem;
  border-bottom: 1px solid rgb(229 231 235);
}
.dark .ai-main__bar {
  border-bottom-color: rgb(55 65 81);
}
.ai-main__title {
  font-size: 0.95rem;
  font-weight: 600;
}
.ai-main__tokens {
  font-size: 0.8rem;
  color: rgb(107 114 128);
}
.ai-messages {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.ai-bubble {
  max-width: min(720px, 92%);
  padding: 0.7rem 0.85rem;
  border-radius: 0.75rem;
  background: rgb(243 244 246);
}
.dark .ai-bubble {
  background: rgb(55 65 81);
}
.ai-bubble--user {
  align-self: flex-end;
  background: rgb(219 234 254);
}
.dark .ai-bubble--user {
  background: rgb(30 64 175 / 0.35);
}
.ai-bubble--tool {
  align-self: flex-start;
  max-width: min(720px, 92%);
  background: transparent;
  border: 1px dashed rgb(209 213 219);
  padding: 0.45rem 0.7rem;
  opacity: 0.9;
}
.dark .ai-bubble--tool {
  border-color: rgb(75 85 99);
}
.ai-live-tools {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  align-self: flex-start;
}
.ai-live-tool {
  font-size: 0.8rem;
  color: rgb(107 114 128);
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  background: rgb(243 244 246);
}
.dark .ai-live-tool {
  background: rgb(55 65 81);
  color: rgb(209 213 219);
}
.ai-live-tool__name {
  font-weight: 600;
  margin-right: 0.5rem;
}
.ai-confirm {
  align-self: stretch;
  margin-top: 0.35rem;
  padding: 0.85rem 1rem;
  border: 1px solid rgb(147 197 253);
  border-radius: 0.75rem;
  background: rgb(239 246 255);
}
.dark .ai-confirm {
  border-color: rgb(59 130 246 / 0.45);
  background: rgb(30 58 138 / 0.25);
}
.ai-confirm__title {
  margin: 0 0 0.75rem;
  font-size: 0.9375rem;
  font-weight: 600;
}
.ai-confirm__options {
  list-style: none;
  margin: 0 0 0.75rem;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  max-height: 220px;
  overflow-y: auto;
}
.ai-confirm__opt {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  font-size: 0.875rem;
  cursor: pointer;
}
.ai-confirm__actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}
.ai-bubble__role {
  font-size: 0.72rem;
  font-weight: 600;
  color: rgb(107 114 128);
  margin-bottom: 0.25rem;
}
.ai-bubble__content {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 0.9375rem;
  line-height: 1.55;
}
.ai-bubble__content :deep(strong) {
  font-weight: 700;
}
.ai-bubble__tokens {
  margin-top: 0.35rem;
  font-size: 0.72rem;
  color: rgb(156 163 175);
}
.ai-composer {
  border-top: 1px solid rgb(229 231 235);
  padding: 0.75rem 1rem 1rem;
}
.dark .ai-composer {
  border-top-color: rgb(55 65 81);
}
.ai-composer__input {
  width: 100%;
  resize: vertical;
  min-height: 4.5rem;
  padding: 0.65rem 0.75rem;
  border: 1px solid rgb(209 213 219);
  border-radius: 0.5rem;
  background: rgb(255 255 255);
  color: inherit;
  font: inherit;
  outline: none;
}
.dark .ai-composer__input {
  border-color: rgb(75 85 99);
  background: rgb(17 24 39);
}
.ai-composer__input:focus {
  border-color: rgb(59 130 246);
}
.ai-composer__actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 0.55rem;
}
.ai-btn {
  height: 2rem;
  padding: 0 0.85rem;
  border: 1px solid rgb(209 213 219);
  border-radius: 0.5rem;
  background: #fff;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  color: inherit;
}
.dark .ai-btn {
  border-color: rgb(75 85 99);
  background: rgb(55 65 81);
}
.ai-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}
.ai-btn--primary {
  border-color: rgb(37 99 235);
  background: rgb(37 99 235);
  color: #fff;
}
.ai-error {
  margin: 0;
  flex: 1;
  font-size: 0.8rem;
  color: rgb(220 38 38);
  text-align: left;
}
@media (max-width: 900px) {
  .ai-page {
    grid-template-columns: 1fr;
    min-height: auto;
  }
  .ai-sidebar {
    max-height: 220px;
  }
  .ai-messages {
    min-height: 320px;
  }
}
</style>

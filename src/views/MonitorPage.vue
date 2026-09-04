<template>
  <PageContent>
    <div class="monitor-page">
      <div class="monitor-page__header">
        <h1 class="browse-section-title border-l-4 border-blue-500 pl-2">监控</h1>
        <button type="button" class="monitor-btn" :disabled="probing" @click="probeActiveTab">
          {{ probing ? '检测中…' : '刷新当前标签' }}
        </button>
      </div>

      <n-tabs v-model:value="activeTab" type="line" animated class="monitor-tabs">
        <n-tab-pane
          v-for="tab in MONITOR_TABS"
          :key="tab.id"
          :name="tab.id"
          :tab="tab.label"
        >
          <div class="monitor-cards">
            <article
              v-for="card in tab.cards"
              :key="card.id"
              class="monitor-card"
              :class="`monitor-card--${statusOf(tab.id, card.id)}`"
            >
              <div class="monitor-card__top">
                <div>
                  <h2 class="monitor-card__title">{{ card.title }}</h2>
                  <p class="monitor-card__url">{{ card.url }}</p>
                </div>
                <span class="monitor-card__badge">{{ statusLabel(statusOf(tab.id, card.id)) }}</span>
              </div>

              <dl class="monitor-card__meta">
                <div>
                  <dt>延迟</dt>
                  <dd>{{ latencyText(tab.id, card.id) }}</dd>
                </div>
                <div>
                  <dt>最近检测</dt>
                  <dd>{{ checkedAtText(tab.id, card.id) }}</dd>
                </div>
              </dl>

              <p v-if="errorOf(tab.id, card.id)" class="monitor-card__error">
                {{ errorOf(tab.id, card.id) }}
              </p>
              <pre v-else-if="bodyOf(tab.id, card.id)" class="monitor-card__body">{{ bodyOf(tab.id, card.id) }}</pre>

              <button
                type="button"
                class="monitor-btn monitor-btn--ghost"
                :disabled="probing"
                @click="probeCard(tab.id, card)"
              >
                单独检测
              </button>
            </article>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>
  </PageContent>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import { NTabPane, NTabs } from 'naive-ui'
import PageContent from '@/components/layout/PageContent.vue'

type ProbeStatus = 'idle' | 'loading' | 'ok' | 'fail'

type MonitorCard = {
  id: string
  title: string
  url: string
}

type MonitorTab = {
  id: string
  label: string
  cards: MonitorCard[]
}

type ProbeResult = {
  status: ProbeStatus
  latencyMs?: number
  checkedAt?: number
  error?: string
  body?: string
}

/** 标签与卡片写死在前端；后续可在此追加 whisper 等 */
const MONITOR_TABS: MonitorTab[] = [
  {
    id: 'onij.fun',
    label: 'onij.fun',
    cards: [
      {
        id: 'ping',
        title: 'Ping',
        url: 'https://onij.fun/go/ping',
      },
    ],
  },
  {
    id: '172.26.250.1',
    label: '172.26.250.1',
    cards: [
      {
        id: 'uvr',
        title: 'UVR',
        url: 'http://172.26.250.1:5555/health',
      },
      {
        id: 'whisper',
        title: 'Whisper',
        url: 'http://172.26.250.1:5555/api/v1/whisper/health',
      },
    ],
  },
]

const activeTab = ref(MONITOR_TABS[0]!.id)
const probing = ref(false)
const results = reactive<Record<string, ProbeResult>>({})

function resultKey(tabId: string, cardId: string) {
  return `${tabId}::${cardId}`
}

function statusOf(tabId: string, cardId: string): ProbeStatus {
  return results[resultKey(tabId, cardId)]?.status ?? 'idle'
}

function errorOf(tabId: string, cardId: string) {
  return results[resultKey(tabId, cardId)]?.error ?? ''
}

function bodyOf(tabId: string, cardId: string) {
  return results[resultKey(tabId, cardId)]?.body ?? ''
}

function latencyText(tabId: string, cardId: string) {
  const ms = results[resultKey(tabId, cardId)]?.latencyMs
  return typeof ms === 'number' ? `${ms} ms` : '—'
}

function checkedAtText(tabId: string, cardId: string) {
  const t = results[resultKey(tabId, cardId)]?.checkedAt
  if (!t) return '—'
  return new Date(t).toLocaleTimeString()
}

function statusLabel(status: ProbeStatus) {
  switch (status) {
    case 'loading':
      return '检测中'
    case 'ok':
      return '正常'
    case 'fail':
      return '异常'
    default:
      return '未检测'
  }
}

async function probeCard(tabId: string, card: MonitorCard) {
  const key = resultKey(tabId, card.id)
  results[key] = { status: 'loading' }
  const started = performance.now()
  try {
    const resp = await fetch(card.url, {
      method: 'GET',
      cache: 'no-store',
      signal: AbortSignal.timeout(8000),
    })
    const latencyMs = Math.round(performance.now() - started)
    const text = (await resp.text()).slice(0, 500)
    if (!resp.ok) {
      results[key] = {
        status: 'fail',
        latencyMs,
        checkedAt: Date.now(),
        error: `HTTP ${resp.status}`,
        body: text || undefined,
      }
      return
    }
    results[key] = {
      status: 'ok',
      latencyMs,
      checkedAt: Date.now(),
      body: text || undefined,
    }
  } catch (err) {
    let msg = err instanceof Error ? err.message : String(err)
    // 浏览器对跨域 / 防火墙拦截只抛 Failed to fetch
    if (/failed to fetch|networkerror|load failed/i.test(msg)) {
      msg = 'Failed to fetch（多半未放行）'
    }
    results[key] = {
      status: 'fail',
      latencyMs: Math.round(performance.now() - started),
      checkedAt: Date.now(),
      error: msg,
    }
  }
}

async function probeActiveTab() {
  const tab = MONITOR_TABS.find((t) => t.id === activeTab.value)
  if (!tab) return
  probing.value = true
  try {
    await Promise.all(tab.cards.map((card) => probeCard(tab.id, card)))
  } finally {
    probing.value = false
  }
}

watch(activeTab, () => {
  void probeActiveTab()
})

onMounted(() => {
  void probeActiveTab()
})
</script>

<style scoped>
.monitor-page {
  max-width: 960px;
}
.monitor-page__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1rem;
}
.monitor-tabs {
  margin-top: 0.25rem;
}
.monitor-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 1rem;
  padding: 1rem 0 1.5rem;
}
.monitor-card {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  padding: 1rem;
  border: 1px solid rgb(229 231 235);
  border-radius: 0.75rem;
  background: #fff;
}
.dark .monitor-card {
  border-color: rgb(55 65 81);
  background: rgb(31 41 55);
}
.monitor-card--ok {
  border-color: rgb(134 239 172);
}
.monitor-card--fail {
  border-color: rgb(252 165 165);
}
.monitor-card--loading {
  border-color: rgb(147 197 253);
}
.monitor-card__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 0.75rem;
}
.monitor-card__title {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 600;
}
.monitor-card__url {
  margin: 0.25rem 0 0;
  font-size: 0.75rem;
  color: rgb(107 114 128);
  word-break: break-all;
}
.monitor-card__badge {
  flex-shrink: 0;
  padding: 0.15rem 0.55rem;
  border-radius: 999px;
  font-size: 0.75rem;
  background: rgb(243 244 246);
  color: rgb(55 65 81);
}
.monitor-card--ok .monitor-card__badge {
  background: rgb(220 252 231);
  color: rgb(22 101 52);
}
.monitor-card--fail .monitor-card__badge {
  background: rgb(254 226 226);
  color: rgb(153 27 27);
}
.monitor-card--loading .monitor-card__badge {
  background: rgb(219 234 254);
  color: rgb(30 64 175);
}
.monitor-card__meta {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem;
  margin: 0;
}
.monitor-card__meta dt {
  font-size: 0.7rem;
  color: rgb(156 163 175);
}
.monitor-card__meta dd {
  margin: 0.1rem 0 0;
  font-size: 0.9rem;
}
.monitor-card__error {
  margin: 0;
  font-size: 0.8rem;
  color: rgb(185 28 28);
  word-break: break-word;
}
.monitor-card__body {
  margin: 0;
  max-height: 6rem;
  overflow: auto;
  padding: 0.5rem;
  border-radius: 0.4rem;
  background: rgb(249 250 251);
  font-size: 0.72rem;
  white-space: pre-wrap;
  word-break: break-word;
}
.dark .monitor-card__body {
  background: rgb(17 24 39);
}
.monitor-btn {
  height: 2rem;
  padding: 0 0.85rem;
  border: 0;
  border-radius: 0.5rem;
  background: rgb(37 99 235);
  color: #fff;
  font-size: 0.875rem;
  cursor: pointer;
}
.monitor-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.monitor-btn--ghost {
  align-self: flex-start;
  background: transparent;
  color: rgb(37 99 235);
  border: 1px solid rgb(191 219 254);
}
</style>

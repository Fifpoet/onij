import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import {
  GetGdstudioSign,
  ListPlaySources,
  SaveGdstudioSign,
  SavePlaySources,
  type PlaySourceDTO,
  type PlaySourceKind,
} from '@/api/playSource'
import {
  applyServerGdstudioSign,
  fetchAndApplyRemoteGdstudioSign,
  gdstudioSignVersion,
} from '@/api/netease/gdstudioSign'

export type { PlaySourceKind }

export type PlaySource = PlaySourceDTO

function envGdOrigin(): string {
  return (
    (import.meta.env.VITE_GDSTUDIO_MUSIC_URL as string | undefined)?.replace(/\/$/, '') ||
    (import.meta.env.DEV ? '/gdstudio-music' : '/gdstudio')
  )
}

function envGdApiOrigin(): string {
  return 'https://music-api.gdstudio.xyz'
}

export const DEFAULT_PLAY_SOURCES: PlaySource[] = [
  {
    id: 'gdapi',
    title: 'GD 公开 API',
    kind: 'gdapi',
    origin: 'https://music-api.gdstudio.xyz',
    enabled: true,
    builtin: true,
  },
  {
    id: 'gdstudio',
    title: 'GD 播放器',
    kind: 'gdstudio',
    origin: '',
    enabled: true,
    builtin: true,
  },
]

export function defaultOriginFor(kind: PlaySourceKind): string {
  return kind === 'gdstudio' ? envGdOrigin() : envGdApiOrigin()
}

export function resolvePlaySourceOrigin(source: PlaySource): string {
  const raw = source.origin.trim()
  return (raw || defaultOriginFor(source.kind)).replace(/\/$/, '')
}

function isPlaySourceKind(kind: string | undefined): kind is PlaySourceKind {
  return kind === 'gdstudio' || kind === 'gdapi'
}

function mergeSources(saved: PlaySource[] | undefined): PlaySource[] {
  const byId = new Map((saved ?? []).map((s) => [s.id, s]))
  const out: PlaySource[] = []
  for (const def of DEFAULT_PLAY_SOURCES) {
    const prev = byId.get(def.id)
    byId.delete(def.id)
    out.push({
      ...def,
      title: prev?.title?.trim() || def.title,
      origin: prev?.origin ?? def.origin,
      enabled: prev?.enabled ?? def.enabled,
    })
  }
  for (const extra of byId.values()) {
    if (!extra?.id || !isPlaySourceKind(extra.kind)) continue
    out.push({
      id: extra.id,
      title: extra.title?.trim() || '自定义源',
      kind: extra.kind,
      origin: extra.origin?.trim() || '',
      enabled: extra.enabled !== false,
      builtin: false,
    })
  }
  return out
}

function toDTO(list: PlaySource[]): PlaySourceDTO[] {
  return list.map((s) => ({
    id: s.id,
    title: s.title,
    kind: s.kind,
    origin: s.origin,
    enabled: s.enabled,
    builtin: !!s.builtin,
  }))
}

export const usePlaySourcesStore = defineStore('playSources', () => {
  const sources = ref<PlaySource[]>(mergeSources(undefined))
  const enabledSources = computed(() => sources.value.filter((s) => s.enabled))
  const signVersion = ref('')
  const signUpdatedAt = ref(0)

  let saveTimer: ReturnType<typeof setTimeout> | undefined
  let saveChain: Promise<void> = Promise.resolve()

  function applyServerList(next: PlaySourceDTO[] | undefined) {
    sources.value = mergeSources(next)
  }

  async function persistNow() {
    try {
      const resp = await SavePlaySources(toDTO(sources.value))
      if (Array.isArray(resp?.sources)) applyServerList(resp.sources)
    } catch {
      /* 保留本地列表，下次再试 */
    }
  }

  function scheduleSave(delayMs: number) {
    if (saveTimer) {
      clearTimeout(saveTimer)
      saveTimer = undefined
    }
    if (delayMs <= 0) {
      saveChain = saveChain.then(persistNow, persistNow)
      return
    }
    saveTimer = setTimeout(() => {
      saveTimer = undefined
      saveChain = saveChain.then(persistNow, persistNow)
    }, delayMs)
  }

  function flushSave() {
    scheduleSave(0)
  }

  function applySignMeta(version: string, updatedAt?: number) {
    signVersion.value = version.trim()
    signUpdatedAt.value = updatedAt || 0
  }

  async function loadSignFromServer() {
    try {
      const sign = await GetGdstudioSign()
      if (!sign?.version || !sign?.crc32) {
        applySignMeta(gdstudioSignVersion())
        return
      }
      applyServerGdstudioSign(sign.version, sign.crc32)
      applySignMeta(sign.version, sign.updatedAt)
    } catch {
      applySignMeta(gdstudioSignVersion())
    }
  }

  async function loadFromServer() {
    try {
      const [resp] = await Promise.all([ListPlaySources(), loadSignFromServer()])
      if (Array.isArray(resp?.sources)) applyServerList(resp.sources)
      else applyServerList(undefined)
    } catch {
      applyServerList(undefined)
    }
  }

  async function syncGdstudioSign(origin: string): Promise<string> {
    const remote = await fetchAndApplyRemoteGdstudioSign(origin)
    applySignMeta(remote.version)
    try {
      const saved = await SaveGdstudioSign({ version: remote.version, crc32: remote.crc32 })
      applyServerGdstudioSign(remote.version, remote.crc32)
      applySignMeta(saved?.version || remote.version, saved?.updatedAt)
    } catch (err) {
      const detail = err instanceof Error ? err.message : String(err)
      throw new Error(`本机已换到 ${remote.version}，写入服务端失败：${detail}`)
    }
    return remote.version
  }

  function patchSource(id: string, patch: Partial<Pick<PlaySource, 'title' | 'origin' | 'enabled'>>) {
    sources.value = sources.value.map((s) => (s.id === id ? { ...s, ...patch } : s))
    const onlyOrigin = Object.keys(patch).length === 1 && 'origin' in patch
    scheduleSave(onlyOrigin ? 400 : 0)
  }

  function addSource(input: { title: string; kind: PlaySourceKind; origin: string }) {
    const origin = input.origin.trim().replace(/\/$/, '')
    sources.value = [
      ...sources.value,
      {
        id: `custom-${Date.now()}`,
        title: input.title.trim() || '自定义源',
        kind: input.kind,
        origin,
        enabled: true,
        builtin: false,
      },
    ]
    scheduleSave(0)
  }

  function removeSource(id: string) {
    sources.value = sources.value.filter((s) => s.id !== id || s.builtin)
    scheduleSave(0)
  }

  function moveSource(id: string, dir: -1 | 1) {
    const idx = sources.value.findIndex((s) => s.id === id)
    const next = idx + dir
    if (idx < 0 || next < 0 || next >= sources.value.length) return
    const copy = sources.value.slice()
    const [item] = copy.splice(idx, 1)
    copy.splice(next, 0, item!)
    sources.value = copy
    scheduleSave(0)
  }

  return {
    sources,
    enabledSources,
    signVersion,
    signUpdatedAt,
    loadFromServer,
    syncGdstudioSign,
    flushSave,
    patchSource,
    addSource,
    removeSource,
    moveSource,
  }
})

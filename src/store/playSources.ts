import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export type PlaySourceKind = 'meting' | 'gdstudio'

export type PlaySource = {
  id: string
  title: string
  kind: PlaySourceKind
  /** 空字符串表示用环境默认（本地反代 / 生产同源） */
  origin: string
  enabled: boolean
  builtin?: boolean
}

function envGdOrigin(): string {
  return (
    (import.meta.env.VITE_GDSTUDIO_MUSIC_URL as string | undefined)?.replace(/\/$/, '') ||
    (import.meta.env.DEV ? '/gdstudio-music' : 'https://onij.fun/gdstudio')
  )
}

function envMetingOrigin(): string {
  return (
    (import.meta.env.VITE_METING_URL as string | undefined)?.replace(/\/$/, '') ||
    'https://api.qijieya.cn/meting'
  )
}

export const DEFAULT_PLAY_SOURCES: PlaySource[] = [
  {
    id: 'meting-qijieya',
    title: 'Meting',
    kind: 'meting',
    origin: '',
    enabled: true,
    builtin: true,
  },
  {
    id: 'gdstudio',
    title: 'GD Studio',
    kind: 'gdstudio',
    origin: '',
    enabled: true,
    builtin: true,
  },
]

export function defaultOriginFor(kind: PlaySourceKind): string {
  return kind === 'gdstudio' ? envGdOrigin() : envMetingOrigin()
}

export function resolvePlaySourceOrigin(source: PlaySource): string {
  const raw = source.origin.trim()
  return (raw || defaultOriginFor(source.kind)).replace(/\/$/, '')
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
    if (!extra?.id || (extra.kind !== 'meting' && extra.kind !== 'gdstudio')) continue
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

export const usePlaySourcesStore = defineStore(
  'playSources',
  () => {
    const sources = ref<PlaySource[]>(mergeSources(undefined))

    const enabledSources = computed(() => sources.value.filter((s) => s.enabled))

    function setSources(next: PlaySource[]) {
      sources.value = mergeSources(next)
    }

    function patchSource(id: string, patch: Partial<Pick<PlaySource, 'title' | 'origin' | 'enabled'>>) {
      sources.value = sources.value.map((s) => (s.id === id ? { ...s, ...patch } : s))
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
    }

    function removeSource(id: string) {
      sources.value = sources.value.filter((s) => s.id !== id || s.builtin)
    }

    function moveSource(id: string, dir: -1 | 1) {
      const idx = sources.value.findIndex((s) => s.id === id)
      const next = idx + dir
      if (idx < 0 || next < 0 || next >= sources.value.length) return
      const copy = sources.value.slice()
      const [item] = copy.splice(idx, 1)
      copy.splice(next, 0, item!)
      sources.value = copy
    }

    return {
      sources,
      enabledSources,
      setSources,
      patchSource,
      addSource,
      removeSource,
      moveSource,
    }
  },
  {
    persist: {
      paths: ['sources'],
      afterRestore: (ctx) => {
        ctx.store.setSources(ctx.store.sources)
      },
    },
  },
)

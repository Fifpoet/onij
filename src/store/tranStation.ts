import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { FileType } from '@/api/types/enums'
import { DownloadFiles } from '@/api/file'
import {
  AddTranFile,
  AddTranText,
  DeleteTran,
  ListTran,
  PinTran,
  type TranItemDTO,
} from '@/api/tran'
import { useFileTransfer } from '@/composables/useFileTransfer'
import { copyToClipboard } from '@/util/common'

type TranItemBase = {
  id: string
  createdAt: number
  pinned?: boolean
  /** 置顶时间，越大越靠前（后来者居上） */
  pinnedAt?: number
}

export type TranTextItem = TranItemBase & {
  kind: 'text'
  content: string
}

export type TranFileItem = TranItemBase & {
  kind: 'file'
  fileId: number
  name: string
  size: number
  /** 上传时写入；旧数据可从文件名推断 */
  format?: FileType
  /** 图片预览用签名链（不持久化，刷新后由页面重新拉取） */
  previewUrl?: string
}

export type TranItem = TranTextItem | TranFileItem

const LOCAL_MIGRATE_KEY = 'tranStation_migrated_v1'

function sortForDisplay(list: TranItem[]): TranItem[] {
  const pinned = list
    .filter((i) => i.pinned)
    .sort((a, b) => (b.pinnedAt ?? 0) - (a.pinnedAt ?? 0))
  const unpinned = list
    .filter((i) => !i.pinned)
    .sort((a, b) => b.createdAt - a.createdAt)
  return [...pinned, ...unpinned]
}

function fromDTO(dto: TranItemDTO): TranItem {
  if (dto.kind === 'file') {
    return {
      id: String(dto.id),
      kind: 'file',
      fileId: Number(dto.file_id ?? 0),
      name: dto.name ?? '',
      size: Number(dto.size ?? 0),
      format: dto.format,
      pinned: !!dto.pinned,
      pinnedAt: dto.pinned_at || undefined,
      createdAt: Number(dto.created_at) || Date.now(),
    }
  }
  return {
    id: String(dto.id),
    kind: 'text',
    content: dto.content ?? '',
    pinned: !!dto.pinned,
    pinnedAt: dto.pinned_at || undefined,
    createdAt: Number(dto.created_at) || Date.now(),
  }
}

function readLegacyLocalItems(): TranItem[] {
  try {
    const raw = localStorage.getItem('tranStation')
    if (!raw) return []
    const parsed = JSON.parse(raw) as { items?: TranItem[] }
    return Array.isArray(parsed.items) ? parsed.items : []
  } catch {
    return []
  }
}

function clearLegacyLocal() {
  try {
    localStorage.removeItem('tranStation')
  } catch {
    /* ignore */
  }
}

export const useTranStationStore = defineStore('tranStation', () => {
  const fileTransfer = useFileTransfer(0)
  const items = ref<TranItem[]>([])
  const loading = ref(false)

  const sortedItems = computed(() => sortForDisplay(items.value))

  async function fetchList() {
    loading.value = true
    try {
      const resp = await ListTran()
      items.value = (resp.items ?? []).map(fromDTO)
      await migrateLegacyIfNeeded()
    } finally {
      loading.value = false
    }
  }

  /** 一次性把旧 localStorage 中转内容迁到服务端 */
  async function migrateLegacyIfNeeded() {
    if (localStorage.getItem(LOCAL_MIGRATE_KEY) === '1') return
    const legacy = readLegacyLocalItems()
    if (!legacy.length) {
      localStorage.setItem(LOCAL_MIGRATE_KEY, '1')
      clearLegacyLocal()
      return
    }
    // 仅当服务端尚无数据时迁入，避免重复
    if (items.value.length > 0) {
      localStorage.setItem(LOCAL_MIGRATE_KEY, '1')
      clearLegacyLocal()
      return
    }
    // 按时间从旧到新写入，保持大致顺序
    const ordered = [...legacy].sort((a, b) => a.createdAt - b.createdAt)
    for (const item of ordered) {
      try {
        let newId = ''
        if (item.kind === 'text' && item.content?.trim()) {
          const resp = await AddTranText(item.content)
          newId = resp.item?.id ? String(resp.item.id) : ''
        } else if (item.kind === 'file' && item.fileId > 0) {
          const resp = await AddTranFile({
            file_id: item.fileId,
            name: item.name,
            size: item.size,
            format: item.format ?? FileType.FT_Unknown,
          })
          newId = resp.item?.id ? String(resp.item.id) : ''
        }
        if (item.pinned && newId) await PinTran(newId)
      } catch {
        /* 单条失败不阻断 */
      }
    }
    localStorage.setItem(LOCAL_MIGRATE_KEY, '1')
    clearLegacyLocal()
    const resp = await ListTran()
    items.value = (resp.items ?? []).map(fromDTO)
  }

  async function addText(content: string) {
    const text = content.trim()
    if (!text) return false
    const resp = await AddTranText(text)
    if (resp.item) items.value = [fromDTO(resp.item), ...items.value]
    else await fetchList()
    return true
  }

  async function addFile(file: File, onProgress?: (n: number) => void) {
    const meta = await fileTransfer.uploadOne(file, onProgress)
    let previewUrl: string | undefined
    if (meta.format === FileType.FT_Image) {
      try {
        const resp = await DownloadFiles({ file_ids: [meta.fileId] })
        previewUrl = resp.urls?.[0]
      } catch {
        /* 页面 mount 时会补链 */
      }
    }
    const resp = await AddTranFile({
      file_id: meta.fileId,
      name: meta.name,
      size: meta.size,
      format: meta.format,
    })
    const item = resp.item ? fromDTO(resp.item) : null
    if (item && item.kind === 'file') {
      item.previewUrl = previewUrl
      items.value = [item, ...items.value]
    } else {
      await fetchList()
    }
  }

  async function removeItem(item: TranItem) {
    await DeleteTran(item.id)
    items.value = items.value.filter((i) => i.id !== item.id)
  }

  async function copyText(item: TranTextItem) {
    const ok = await copyToClipboard(item.content)
    if (!ok) throw new Error('复制失败，请检查浏览器权限或使用 HTTPS')
  }

  async function downloadFile(item: TranFileItem) {
    await fileTransfer.downloadOne(item.fileId, item.name)
  }

  async function togglePin(item: TranItem) {
    const resp = await PinTran(item.id)
    if (!resp.item) {
      await fetchList()
      return
    }
    const next = fromDTO(resp.item)
    const idx = items.value.findIndex((i) => i.id === item.id)
    if (idx >= 0) {
      const prev = items.value[idx]
      if (prev.kind === 'file' && next.kind === 'file') {
        next.previewUrl = prev.previewUrl
      }
      items.value[idx] = next
    } else {
      await fetchList()
    }
  }

  return {
    items,
    sortedItems,
    loading,
    fetchList,
    addText,
    addFile,
    removeItem,
    copyText,
    downloadFile,
    togglePin,
  }
})

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useFileTransfer } from '@/composables/useFileTransfer'

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
}

export type TranItem = TranTextItem | TranFileItem

function sortForDisplay(list: TranItem[]): TranItem[] {
  const pinned = list
    .filter((i) => i.pinned)
    .sort((a, b) => (b.pinnedAt ?? 0) - (a.pinnedAt ?? 0))
  const unpinned = list
    .filter((i) => !i.pinned)
    .sort((a, b) => b.createdAt - a.createdAt)
  return [...pinned, ...unpinned]
}

function insertUnpinnedAtTop(list: TranItem[], item: TranItem) {
  const firstUnpinned = list.findIndex((i) => !i.pinned)
  if (firstUnpinned === -1) list.push(item)
  else list.splice(firstUnpinned, 0, item)
}

export const useTranStationStore = defineStore(
  'tranStation',
  () => {
    const fileTransfer = useFileTransfer(0)
    const items = ref<TranItem[]>([])

    const sortedItems = computed(() => sortForDisplay(items.value))

    function addText(content: string) {
      const text = content.trim()
      if (!text) return false
      const item: TranTextItem = {
        id: crypto.randomUUID(),
        kind: 'text',
        content: text,
        createdAt: Date.now(),
        pinned: false,
      }
      insertUnpinnedAtTop(items.value, item)
      return true
    }

    async function addFile(file: File, onProgress?: (n: number) => void) {
      const meta = await fileTransfer.uploadOne(file, onProgress)
      const item: TranFileItem = {
        id: String(meta.fileId),
        kind: 'file',
        fileId: meta.fileId,
        name: meta.name,
        size: meta.size,
        createdAt: Date.now(),
        pinned: false,
      }
      insertUnpinnedAtTop(items.value, item)
    }

    async function removeItem(item: TranItem) {
      if (item.kind === 'file') {
        await fileTransfer.removeOne(item.fileId)
      }
      items.value = items.value.filter((i) => i.id !== item.id)
    }

    async function copyText(item: TranTextItem) {
      await navigator.clipboard.writeText(item.content)
    }

    async function downloadFile(item: TranFileItem) {
      await fileTransfer.downloadOne(item.fileId)
    }

    function togglePin(item: TranItem) {
      const target = items.value.find((i) => i.id === item.id)
      if (!target) return
      if (target.pinned) {
        target.pinned = false
        target.pinnedAt = undefined
      } else {
        target.pinned = true
        target.pinnedAt = Date.now()
      }
    }

    return {
      items,
      sortedItems,
      addText,
      addFile,
      removeItem,
      copyText,
      downloadFile,
      togglePin,
    }
  },
  { persist: { paths: ['items'] } },
)

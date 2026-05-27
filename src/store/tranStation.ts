import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { FileType } from '@/api/types/enums'
import { DownloadFiles } from '@/api/file'
import { useFileTransfer } from '@/composables/useFileTransfer'
import { copyToClipboard } from '@/util/common'
import { randomUUID } from '@/util/crypto'

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
        id: randomUUID(),
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
      let previewUrl: string | undefined
      if (meta.format === FileType.FT_Image) {
        try {
          const resp = await DownloadFiles({ file_ids: [meta.fileId] })
          previewUrl = resp.urls?.[0]
        } catch {
          /* 页面 mount 时会补链 */
        }
      }
      const item: TranFileItem = {
        id: String(meta.fileId),
        kind: 'file',
        fileId: meta.fileId,
        name: meta.name,
        size: meta.size,
        format: meta.format,
        previewUrl,
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
      const ok = await copyToClipboard(item.content)
      if (!ok) throw new Error('复制失败，请检查浏览器权限或使用 HTTPS')
    }

    async function downloadFile(item: TranFileItem) {
      await fileTransfer.downloadOne(item.fileId, item.name)
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
  {
    persist: {
      paths: ['items'],
      serializer: {
        serialize: (state: { items: TranItem[] }) =>
          JSON.stringify({
            items: state.items.map((item) =>
              item.kind === 'file' ? { ...item, previewUrl: undefined } : item,
            ),
          }),
        deserialize: (value: string) => JSON.parse(value) as { items: TranItem[] },
      },
    },
  },
)

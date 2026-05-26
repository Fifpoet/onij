<template>
  <PageContent>
    <div class="tran-page">
      <h1 class="browse-section-title mb-6 border-l-4 border-blue-500 pl-2">中转站</h1>

      <div
        class="tran-drop"
        :class="{ 'tran-drop--active': dropActive, 'tran-drop--busy': uploading }"
        @dragover.prevent="dropActive = true"
        @dragleave.prevent="dropActive = false"
        @drop.prevent="onDrop"
        @click="openFilePicker"
      >
        <input
          ref="fileInputRef"
          type="file"
          class="sr-only"
          multiple
          @change="onFileInputChange"
        >
        <p class="tran-drop__title">拖入文件，或点击选择</p>
        <p class="tran-drop__hint">支持 Ctrl+V 粘贴文件 / 文本</p>
        <div v-if="uploading" class="tran-drop__progress">
          <div class="tran-drop__bar" :style="{ width: `${uploadPercent}%` }" />
        </div>
        <p v-if="uploading" class="tran-drop__label">{{ uploadLabel }}</p>
      </div>

      <div class="tran-text-row">
        <input
          v-model="textDraft"
          type="text"
          class="tran-text-input"
          placeholder="输入要中转的文本"
          @keydown.enter.prevent="submitText"
        >
        <button type="button" class="tran-btn tran-btn--primary" @click="submitText">确认</button>
      </div>

      <p v-if="!tran.sortedItems.length" class="browse-muted py-8 text-center">暂无中转内容</p>
      <ul v-else class="tran-list">
        <li
          v-for="item in tran.sortedItems"
          :key="item.id"
          class="tran-card"
          :class="{ 'tran-card--pinned': item.pinned }"
        >
          <div class="tran-card__main min-w-0">
            <span class="tran-card__tag">{{ itemTag(item) }}</span>
            <span v-if="item.pinned" class="tran-card__pin-badge">置顶</span>
            <p v-if="item.kind === 'text'" class="tran-card__text">{{ item.content }}</p>
            <div
              v-else-if="item.kind === 'file' && isTranImage(item)"
              class="tran-card__media"
            >
              <img
                v-if="imageSrc(item)"
                :src="imageSrc(item)"
                :alt="item.name"
                class="tran-card__img"
              >
              <p v-else class="tran-card__meta">图片加载中…</p>
              <p class="tran-card__meta">{{ item.name }} · {{ formatFileSize(item.size) }}</p>
            </div>
            <template v-else>
              <p class="tran-card__name">{{ item.name }}</p>
              <p class="tran-card__meta">{{ formatFileSize(item.size) }}</p>
            </template>
          </div>
          <div class="tran-card__actions">
            <button
              type="button"
              class="tran-icon-btn"
              :class="{ 'tran-icon-btn--pinned': item.pinned }"
              :aria-label="item.pinned ? '取消置顶' : '置顶'"
              @click="tran.togglePin(item)"
            >
              <n-icon :component="item.pinned ? PinSharp : PinOutline" :size="20" />
            </button>
            <button
              type="button"
              class="tran-icon-btn"
              :aria-label="item.kind === 'text' ? '复制' : '下载'"
              @click="onPrimaryAction(item)"
            >
              <n-icon
                :component="item.kind === 'text' ? CopyOutline : DownloadOutline"
                :size="20"
              />
            </button>
            <button
              type="button"
              class="tran-icon-btn tran-icon-btn--danger"
              aria-label="删除"
              @click="onRemove(item)"
            >
              <n-icon :component="TrashOutline" :size="20" />
            </button>
          </div>
        </li>
      </ul>
    </div>
  </PageContent>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { NIcon } from 'naive-ui'
import {
  CopyOutline,
  DownloadOutline,
  PinOutline,
  PinSharp,
  TrashOutline,
} from '@vicons/ionicons5'
import PageContent from '@/components/layout/PageContent.vue'
import { DownloadFiles } from '@/api/file'
import { FileType } from '@/api/types/enums'
import { useTranStationStore, type TranItem, type TranFileItem } from '@/store/tranStation'
import { formatFileSize, getFileTypeFromString, isImage } from '@/util/file'

const tran = useTranStationStore()
const imageUrlCache = ref<Record<number, string>>({})
let alive = true
let loadingUrls = false

const fileInputRef = ref<HTMLInputElement | null>(null)
const dropActive = ref(false)
const textDraft = ref('')
const uploading = ref(false)
const uploadPercent = ref(0)
const uploadLabel = ref('')

function fileFormat(item: TranFileItem): FileType {
  return item.format ?? getFileTypeFromString(item.name)
}

function isTranImage(item: TranItem): item is TranFileItem {
  return item.kind === 'file' && isImage(fileFormat(item))
}

function itemTag(item: TranItem) {
  if (item.kind === 'text') return '文本'
  if (isTranImage(item)) return '图片'
  return '文件'
}

function isSignedUrl(url: string): boolean {
  return /[?&]token=/.test(url)
}

function imageSrc(item: TranFileItem): string | undefined {
  const cached = imageUrlCache.value[item.fileId]
  if (cached) return cached
  if (item.previewUrl && isSignedUrl(item.previewUrl)) return item.previewUrl
  return undefined
}

async function loadMissingImageUrls() {
  if (!alive || loadingUrls) return
  const missing = tran.sortedItems.filter(
    (i): i is TranFileItem => isTranImage(i) && !imageSrc(i),
  )
  if (!missing.length) return
  loadingUrls = true
  try {
    const resp = await DownloadFiles({ file_ids: missing.map((i) => i.fileId) })
    if (!alive) return
    const urls = resp.urls ?? []
    const next = { ...imageUrlCache.value }
    missing.forEach((item, idx) => {
      if (urls[idx]) next[item.fileId] = urls[idx]
    })
    imageUrlCache.value = next
  } catch {
    /* ignore */
  } finally {
    loadingUrls = false
  }
}

function openFilePicker() {
  if (uploading.value) return
  fileInputRef.value?.click()
}

async function ingestFiles(files: FileList | File[]) {
  for (const file of Array.from(files)) {
    uploading.value = true
    uploadPercent.value = 0
    uploadLabel.value = `上传中 ${file.name}`
    try {
      await tran.addFile(file, (n) => {
        uploadPercent.value = n
      })
    } catch (e) {
      uploadLabel.value = e instanceof Error ? e.message : '上传失败'
      await new Promise((r) => setTimeout(r, 1200))
    }
  }
  uploading.value = false
  uploadLabel.value = ''
}

function onFileInputChange(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files?.length) void ingestFiles(input.files)
  input.value = ''
}

function onDrop(e: DragEvent) {
  dropActive.value = false
  const files = e.dataTransfer?.files
  if (files?.length) void ingestFiles(files)
}

function submitText() {
  if (tran.addText(textDraft.value)) textDraft.value = ''
}

async function onPrimaryAction(item: TranItem) {
  try {
    if (item.kind === 'text') await tran.copyText(item)
    else await tran.downloadFile(item)
  } catch {
    // ignore
  }
}

async function onRemove(item: TranItem) {
  try {
    await tran.removeItem(item)
  } catch {
    // ignore
  }
}

function onPaste(e: ClipboardEvent) {
  const dt = e.clipboardData
  if (!dt) return

  const files: File[] = []
  for (const item of dt.items) {
    if (item.kind === 'file') {
      const f = item.getAsFile()
      if (f) files.push(f)
    }
  }
  if (files.length) {
    e.preventDefault()
    void ingestFiles(files)
    return
  }
  const text = dt.getData('text/plain')?.trim()
  if (text) {
    e.preventDefault()
    tran.addText(text)
  }
}

onMounted(() => {
  window.addEventListener('paste', onPaste)
  void loadMissingImageUrls()
})

const stopItemsWatch = watch(
  () => tran.sortedItems.map((i) => (i.kind === 'file' ? i.fileId : i.id)).join(','),
  () => {
    void loadMissingImageUrls()
  },
)

onBeforeUnmount(() => {
  alive = false
  stopItemsWatch()
  window.removeEventListener('paste', onPaste)
})
</script>

<style scoped>
.tran-page {
  width: 100%;
  padding-bottom: 2rem;
}

.tran-drop {
  margin-bottom: 1.25rem;
  padding: 2rem 1.25rem 1.5rem;
  border: 2px dashed rgb(209 213 219);
  border-radius: 0.75rem;
  background: rgb(249 250 251);
  text-align: center;
  cursor: pointer;
  transition: border-color 0.15s, background-color 0.15s;
}

.dark .tran-drop {
  border-color: rgb(75 85 99);
  background: rgb(17 24 39);
}

.tran-drop:hover,
.tran-drop--active {
  border-color: rgb(59 130 246);
  background: rgb(239 246 255);
}

.dark .tran-drop:hover,
.dark .tran-drop--active {
  background: rgb(30 41 59);
}

.tran-drop--busy {
  pointer-events: none;
  opacity: 0.88;
}

.tran-drop__title {
  margin: 0 0 0.35rem;
  font-size: 1.0625rem;
  font-weight: 600;
  color: rgb(17 24 39);
}

.dark .tran-drop__title {
  color: rgb(243 244 246);
}

.tran-drop__hint {
  margin: 0;
  font-size: 0.9375rem;
  color: rgb(107 114 128);
}

.tran-drop__progress {
  margin: 1rem auto 0;
  max-width: 20rem;
  height: 6px;
  border-radius: 999px;
  background: rgb(229 231 235);
  overflow: hidden;
}

.tran-drop__bar {
  height: 100%;
  background: rgb(59 130 246);
  transition: width 0.12s ease;
}

.tran-drop__label {
  margin: 0.5rem 0 0;
  font-size: 0.875rem;
  color: rgb(107 114 128);
}

.tran-text-row {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1.75rem;
}

.tran-text-input {
  flex: 1;
  min-width: 0;
  height: 2.75rem;
  padding: 0 0.85rem;
  border: 1px solid rgb(209 213 219);
  border-radius: 0.5rem;
  font-size: 1rem;
  outline: none;
}

.tran-text-input:focus {
  border-color: rgb(59 130 246);
  box-shadow: 0 0 0 2px rgb(59 130 246 / 0.2);
}

.dark .tran-text-input {
  border-color: rgb(75 85 99);
  background: rgb(17 24 39);
  color: rgb(243 244 246);
}

.tran-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.tran-card {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 1.1rem;
  border: 1px solid rgb(229 231 235);
  border-radius: 0.625rem;
  background: rgb(255 255 255);
}

.dark .tran-card {
  border-color: rgb(55 65 81);
  background: rgb(31 41 55);
}

.tran-card--pinned {
  border-color: rgb(191 219 254);
  background: rgb(239 246 255);
}

.dark .tran-card--pinned {
  border-color: rgb(59 130 246);
  background: rgb(30 41 59);
}

.tran-card__tag {
  display: inline-block;
  margin-bottom: 0.35rem;
  margin-right: 0.35rem;
  padding: 0.1rem 0.45rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 600;
  color: rgb(37 99 235);
  background: rgb(219 234 254);
}

.tran-card__pin-badge {
  display: inline-block;
  margin-bottom: 0.35rem;
  padding: 0.1rem 0.45rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 600;
  color: rgb(29 78 216);
  background: rgb(191 219 254);
}

.dark .tran-card__pin-badge {
  color: rgb(147 197 253);
  background: rgb(30 64 175);
}

.tran-card__text {
  margin: 0;
  font-size: 1rem;
  line-height: 1.55;
  white-space: pre-wrap;
  word-break: break-word;
}

.tran-card__img {
  display: block;
  max-width: 100%;
  max-height: 20rem;
  margin-top: 0.35rem;
  border-radius: 0.5rem;
  object-fit: contain;
  background: rgb(243 244 246);
}

.dark .tran-card__img {
  background: rgb(17 24 39);
}

.tran-card__name {
  margin: 0;
  font-weight: 600;
  word-break: break-all;
}

.tran-card__meta {
  margin: 0.25rem 0 0;
  font-size: 0.875rem;
  color: rgb(107 114 128);
}

.tran-card__actions {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 0.25rem;
}

.tran-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  padding: 0;
  border: 0;
  border-radius: 0.5rem;
  background: transparent;
  color: rgb(107 114 128);
  cursor: pointer;
  transition: background-color 0.12s ease, color 0.12s ease;
}

.tran-icon-btn:hover:not(:disabled) {
  background: rgb(243 244 246);
  color: rgb(55 65 81);
}

.dark .tran-icon-btn {
  color: rgb(156 163 175);
}

.dark .tran-icon-btn:hover:not(:disabled) {
  background: rgb(55 65 81);
  color: rgb(243 244 246);
}

.tran-icon-btn--pinned {
  color: rgb(37 99 235);
}

.tran-icon-btn--pinned:hover {
  background: rgb(219 234 254);
  color: rgb(29 78 216);
}

.dark .tran-icon-btn--pinned {
  color: rgb(96 165 250);
}

.dark .tran-icon-btn--pinned:hover {
  background: rgb(30 64 175);
  color: rgb(191 219 254);
}

.tran-icon-btn--danger {
  color: rgb(220 38 38);
}

.tran-icon-btn--danger:hover:not(:disabled) {
  background: rgb(254 242 242);
  color: rgb(185 28 28);
}

.dark .tran-icon-btn--danger {
  color: rgb(248 113 113);
}

.dark .tran-icon-btn--danger:hover:not(:disabled) {
  background: rgb(127 29 29);
  color: rgb(254 202 202);
}

.tran-btn {
  height: 2.25rem;
  padding: 0 0.85rem;
  border-radius: 0.5rem;
  border: 1px solid transparent;
  font-size: 0.9375rem;
  font-weight: 500;
  cursor: pointer;
}

.tran-btn--primary {
  flex-shrink: 0;
  color: #fff;
  background: rgb(37 99 235);
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  border: 0;
}
</style>

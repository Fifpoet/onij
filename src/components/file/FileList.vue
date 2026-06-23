<template>
  <div class="file-list">
    <div class="px-3 pt-3 pb-3 flex flex-wrap items-center gap-3 lg:px-5 lg:pt-5">
      <n-button
        size="small"
        quaternary
        :disabled="currentParentId === 0"
        class="flex items-center h-8 min-w-[90px] justify-center"
        @click="goBack"
      >
        <template #icon>
          <n-icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m15 18-6-6 6-6" />
            </svg>
          </n-icon>
        </template>
        返回上级
      </n-button>

      <div class="flex items-center gap-1 shrink-0 min-w-0 flex-wrap">
        <span class="text-gray-300">|</span>
        <n-button
          size="small"
          text
          class="px-1"
          :type="currentParentId === 0 ? 'primary' : 'default'"
          @click="goToRoot()"
        >
          根目录
        </n-button>
        <template v-for="(name, idx) in folderNameHistory.slice(1)" :key="`${name}-${idx}`">
          <span class="text-gray-300">/</span>
          <span
            class="text-sm truncate max-w-[8rem]"
            :class="idx === folderNameHistory.length - 2 ? 'font-medium text-gray-900' : 'text-gray-500'"
          >
            {{ name }}
          </span>
        </template>
      </div>

      <div class="flex-1 flex justify-end min-w-[10rem]">
        <n-input
          v-model:value="searchKeyword"
          type="text"
          placeholder="搜索文件..."
          size="small"
          class="max-w-[220px]"
          @keydown.enter="handleSearch"
        >
          <template #prefix>
            <n-icon>
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="11" cy="11" r="8" />
                <line x1="21" y1="21" x2="16.65" y2="16.65" />
              </svg>
            </n-icon>
          </template>
        </n-input>
      </div>
    </div>

    <div class="flex justify-end flex-wrap gap-2 px-3 lg:px-5 mb-4">
      <n-button
        quaternary
        size="small"
        class="h-8 w-8 flex items-center justify-center"
        :type="viewMode === 'grid' ? 'default' : 'primary'"
        @click="toggleViewMode"
      >
        <n-icon>
          <svg v-if="viewMode === 'grid'" xmlns="http://www.w3.org/2000/svg" width="16" height="16"
            viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
            stroke-linejoin="round">
            <line x1="3" y1="6" x2="21" y2="6" />
            <line x1="3" y1="12" x2="21" y2="12" />
            <line x1="3" y1="18" x2="21" y2="18" />
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="3" width="7" height="7" />
            <rect x="14" y="3" width="7" height="7" />
            <rect x="14" y="14" width="7" height="7" />
            <rect x="3" y="14" width="7" height="7" />
          </svg>
        </n-icon>
      </n-button>
      <n-button size="small" quaternary class="h-8 px-3" @click="showUploadModal = true">
        <template #icon>
          <n-icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
              <polyline points="17 8 12 3 7 8" />
              <line x1="12" y1="3" x2="12" y2="15" />
            </svg>
          </n-icon>
        </template>
        上传文件
      </n-button>
      <n-button size="small" quaternary class="h-8 px-3" @click="showFolderModal = true">
        <template #icon>
          <n-icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path
                d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2Z" />
              <line x1="12" y1="10" x2="12" y2="16" />
              <line x1="9" y1="13" x2="15" y2="13" />
            </svg>
          </n-icon>
        </template>
        新建文件夹
      </n-button>
    </div>

    <div class="file-list__body p-3 lg:p-5">
      <div
        v-if="viewMode === 'grid'"
        class="grid gap-4"
        style="grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));"
      >
        <div
          v-for="file in fileList"
          :key="file.id"
          class="file-card"
          @click="handleFileClick(file)"
        >
          <div class="file-card__preview">
            <div v-if="showFileIcon(file.format)" class="flex justify-center items-center h-full">
              <n-icon size="40" class="text-gray-500">
                <component :is="getFileIcon(file.format)" />
              </n-icon>
            </div>
            <img v-else :src="fileImageSrc(file)" :alt="file.name" class="w-full h-full object-cover">
          </div>

          <n-ellipsis :line-clamp="2" class="file-card__name" :title="file.name">
            {{ file.name }}
          </n-ellipsis>

          <div v-if="!isFolder(file.format)" class="file-card__meta">
            <div class="truncate">{{ formatFileSize(file.size || 0) }}</div>
            <div class="truncate">{{ formatFileDisplayTime(file) }}</div>
          </div>

          <div class="file-card__actions" @click.stop>
            <template v-if="!isFolder(file.format)">
              <n-button size="tiny" type="primary" @click="downloadFile(file)">下载</n-button>
            </template>
            <n-button size="tiny" type="error" @click="deleteFile(file)">删除</n-button>
          </div>
        </div>
      </div>

      <div v-else class="space-y-2">
        <div
          v-for="file in fileList"
          :key="file.id"
          class="file-row"
          @click="handleFileClick(file)"
        >
          <div class="file-row__icon">
            <img
              v-if="isImage(file.format) && file.url"
              :src="fileImageSrc(file)"
              :alt="file.name"
              class="w-10 h-10 rounded object-cover"
            >
            <n-icon v-else size="24" class="text-gray-500">
              <component :is="getFileIcon(file.format)" />
            </n-icon>
          </div>

          <div class="flex-1 min-w-0">
            <n-ellipsis class="text-sm font-medium text-gray-900" :title="file.name">
              {{ file.name }}
            </n-ellipsis>
          </div>

          <div v-if="!isFolder(file.format)" class="file-row__size">
            {{ formatFileSize(file.size || 0) }}
          </div>
          <div v-if="!isFolder(file.format)" class="file-row__time">
            {{ formatFileDisplayTime(file) }}
          </div>

          <div class="file-row__actions" @click.stop>
            <n-button v-if="!isFolder(file.format)" size="tiny" type="primary" @click="downloadFile(file)">
              下载
            </n-button>
            <n-button size="tiny" type="error" @click="deleteFile(file)">删除</n-button>
          </div>
        </div>
      </div>

      <div v-if="loading" class="file-list__loading">
        <n-spin size="large" />
      </div>

      <div v-if="fileList.length === 0 && !loading" class="flex justify-center py-16">
        <n-empty description="暂无文件" />
      </div>
    </div>

    <div v-if="totalCount > 0" class="border-t border-gray-200 bg-white p-4">
      <div class="flex flex-col items-center space-y-2">
        <div class="text-sm text-gray-600">
          共 {{ totalCount }} 项，第 {{ currentPage }} / {{ totalPages || 1 }} 页
        </div>
        <n-pagination
          v-model:page="currentPage"
          :page-count="totalPages"
          :page-sizes="[10, 20, 50, 100]"
          :page-size="pageSize"
          show-size-picker
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </div>
  </div>

  <FileUploadModal
    v-model:show="showUploadModal"
    :parent-id="currentParentId"
    :folder-path="folderPathSegments"
    @uploaded="fetchFileList"
  />

  <FolderCreateModal
    v-model:show="showFolderModal"
    :parent-id="currentParentId"
    @created="fetchFileList"
  />

  <div
    v-if="previewVisible"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/75"
    @click="closePreview"
  >
    <div class="relative w-full h-full flex items-center justify-center overflow-hidden" @wheel="handleZoom" @click.stop>
      <img
        :src="previewImage"
        class="max-w-full max-h-full object-contain transition-transform duration-200"
        :style="{ transform: `scale(${zoomLevel})` }"
      >
      <button
        type="button"
        class="absolute top-5 right-5 bg-white/75 rounded-full p-2"
        @click.stop="closePreview"
      >
        ✕
      </button>
      <button
        type="button"
        class="absolute left-5 top-1/2 -translate-y-1/2 bg-white/75 rounded-full p-3 disabled:opacity-50"
        :disabled="!hasPrevImage"
        @click.stop="prevImage"
      >
        ‹
      </button>
      <button
        type="button"
        class="absolute right-5 top-1/2 -translate-y-1/2 bg-white/75 rounded-full p-3 disabled:opacity-50"
        :disabled="!hasNextImage"
        @click.stop="nextImage"
      >
        ›
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue'
import {
  NIcon,
  NButton,
  NSpin,
  NEmpty,
  NPagination,
  NEllipsis,
  useMessage,
  NInput,
} from 'naive-ui'
import { GetFileList, DeleteFileById } from '@/api'
import { type FileDetail, type GetFileListReq, type DeleteFileReq } from '@/api/types/file'
import {
  getFileIcon,
  formatFileSize,
  formatFileDisplayTime,
  fileImageSrc,
  fileDownloadUrl,
  isFolder,
  showFileIcon,
  isImage,
} from '@/util'
import { downloadByPrivateUrl } from '@/util/qiniu'
import FileUploadModal from './FileUploadModal.vue'
import FolderCreateModal from './FolderCreateModal.vue'

const message = useMessage()

const fileList = ref<FileDetail[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalCount = ref(0)
const totalPages = ref(0)
const searchKeyword = ref('')

const currentParentId = ref(0)
const cloudFolderId = ref<number | null>(null)
const initialized = ref(false)
const folderHistory = ref<number[]>([0])
const folderNameHistory = ref<string[]>(['根目录'])

const showUploadModal = ref(false)
const showFolderModal = ref(false)
const viewMode = ref<'grid' | 'list'>('grid')

const previewVisible = ref(false)
const previewImage = ref('')
const zoomLevel = ref(1)
const currentPreviewIndex = ref(-1)

const folderPathSegments = computed(() => folderNameHistory.value.slice(1))

const imageFiles = computed(() => fileList.value.filter((file) => isImage(file.format) && file.url))
const hasPrevImage = computed(() => currentPreviewIndex.value > 0)
const hasNextImage = computed(() => currentPreviewIndex.value < imageFiles.value.length - 1)

const prevImage = () => {
  if (!hasPrevImage.value) return
  currentPreviewIndex.value -= 1
  previewImage.value = fileImageSrc(imageFiles.value[currentPreviewIndex.value])
  zoomLevel.value = 1
}

const nextImage = () => {
  if (!hasNextImage.value) return
  currentPreviewIndex.value += 1
  previewImage.value = fileImageSrc(imageFiles.value[currentPreviewIndex.value])
  zoomLevel.value = 1
}

const fetchFileList = async (page: number = currentPage.value) => {
  if (loading.value) return

  loading.value = true
  try {
    const req: GetFileListReq = {
      parent_id: currentParentId.value,
      page,
      limit: pageSize.value,
      keyword: searchKeyword.value.trim() || undefined,
    }

    const response = await GetFileList(req)
    fileList.value = response.files ?? []
    currentPage.value = page
    totalCount.value = response.total || 0
    totalPages.value = Math.max(1, Math.ceil(totalCount.value / pageSize.value))
  } catch (error) {
    console.error('获取文件列表失败:', error)
    fileList.value = []
    totalCount.value = 0
    totalPages.value = 0
    message.error('获取文件列表失败')
  } finally {
    loading.value = false
  }
}

watch(currentParentId, () => {
  if (!initialized.value) return
  currentPage.value = 1
  fetchFileList(1)
})

const goToRoot = () => {
  folderHistory.value = [0]
  folderNameHistory.value = ['根目录']
  currentParentId.value = 0
}

const resolveCloudHome = async () => {
  loading.value = true
  try {
    const response = await GetFileList({ parent_id: 0, page: 1, limit: 100 })
    const cloud = response.files?.find((f) => isFolder(f.format) && f.name === 'cloud')
    if (cloud) {
      cloudFolderId.value = cloud.id
      folderHistory.value = [0, cloud.id]
      folderNameHistory.value = ['根目录', 'cloud']
      currentParentId.value = cloud.id
    } else {
      cloudFolderId.value = null
      goToRoot()
      message.warning('未找到 cloud 文件夹，已显示根目录。需初始化请运行 go run ./cmd/ensure_cloud_folder')
    }
  } catch (error) {
    console.error('定位 cloud 文件夹失败:', error)
    goToRoot()
  } finally {
    initialized.value = true
    loading.value = false
    currentPage.value = 1
    await fetchFileList(1)
  }
}

onMounted(() => {
  resolveCloudHome()
})

const enterFolder = (file: FileDetail) => {
  folderHistory.value.push(file.id)
  folderNameHistory.value.push(file.name)
  currentParentId.value = file.id
}

const handleFileClick = (file: FileDetail) => {
  if (isFolder(file.format)) {
    enterFolder(file)
    return
  }
  if (isImage(file.format) && file.url) {
    previewImage.value = fileImageSrc(file)
    previewVisible.value = true
    zoomLevel.value = 1
    currentPreviewIndex.value = imageFiles.value.findIndex((img) => img.id === file.id)
  }
}

const goBack = () => {
  if (folderHistory.value.length <= 1) return
  folderHistory.value.pop()
  folderNameHistory.value.pop()
  currentParentId.value = folderHistory.value[folderHistory.value.length - 1]
}

const handleSearch = () => {
  currentPage.value = 1
  fetchFileList(1)
}

const handlePageChange = (page: number) => {
  fetchFileList(page)
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchFileList(1)
}

const downloadFile = (file: FileDetail) => {
  if (import.meta.env.DEV) {
    window.location.href = fileDownloadUrl(file)
    return
  }
  downloadByPrivateUrl(fileDownloadUrl(file), file.name)
}

const deleteFile = async (file: FileDetail) => {
  const label = isFolder(file.format) ? '文件夹' : '文件'
  const confirmed = window.confirm(`确定要删除${label}「${file.name}」吗？此操作不可恢复。`)
  if (!confirmed) return

  try {
    const deleteData: DeleteFileReq = { file_id: file.id }
    const response = await DeleteFileById(deleteData)
    if (response.message === 'ok') {
      message.success('删除成功')
      fetchFileList()
    } else {
      message.error(response.message || '删除失败')
    }
  } catch {
    message.error('删除失败')
  }
}

const toggleViewMode = () => {
  viewMode.value = viewMode.value === 'grid' ? 'list' : 'grid'
}

const closePreview = () => {
  previewVisible.value = false
  previewImage.value = ''
  currentPreviewIndex.value = -1
}

const handleZoom = (e: WheelEvent) => {
  e.preventDefault()
  if (e.deltaY < 0) {
    zoomLevel.value = Math.min(zoomLevel.value + 0.1, 3)
  } else {
    zoomLevel.value = Math.max(zoomLevel.value - 0.1, 0.5)
  }
}
</script>

<style scoped>
.file-list__body {
  width: 100%;
  min-height: 420px;
  position: relative;
}

.file-list__loading {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgb(255 255 255 / 0.8);
  z-index: 10;
}

.file-card {
  background: #fff;
  border: 1px solid rgb(229 231 235);
  border-radius: 0.5rem;
  padding: 1rem;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  min-width: 0;
  aspect-ratio: 1 / 1.15;
  transition: box-shadow 0.2s, border-color 0.2s;
}

.file-card:hover {
  border-color: rgb(209 213 219);
  box-shadow: 0 4px 12px rgb(0 0 0 / 0.06);
}

.file-card__preview {
  width: 100%;
  height: 50%;
  margin-bottom: 0.75rem;
  overflow: hidden;
  border-radius: 0.375rem;
  background: rgb(249 250 251);
}

.file-card__name {
  text-align: center;
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(17 24 39);
  margin-bottom: 0.5rem;
}

.file-card__meta {
  text-align: center;
  font-size: 0.75rem;
  color: rgb(107 114 128);
  margin-bottom: 0.5rem;
}

.file-card__actions {
  margin-top: auto;
  display: flex;
  justify-content: center;
  gap: 0.35rem;
  flex-wrap: wrap;
}

.file-row {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  background: #fff;
  border: 1px solid rgb(229 231 235);
  border-radius: 0.5rem;
  cursor: pointer;
  transition: box-shadow 0.2s, border-color 0.2s;
}

.file-row:hover {
  border-color: rgb(209 213 219);
  box-shadow: 0 4px 12px rgb(0 0 0 / 0.06);
}

.file-row__icon {
  flex-shrink: 0;
  width: 2.5rem;
  display: flex;
  justify-content: center;
}

.file-row__size {
  flex-shrink: 0;
  width: 5rem;
  text-align: right;
  font-size: 0.875rem;
  color: rgb(107 114 128);
}

.file-row__time {
  flex-shrink: 0;
  width: 9rem;
  text-align: right;
  font-size: 0.875rem;
  color: rgb(107 114 128);
}

.file-row__actions {
  flex-shrink: 0;
  display: flex;
  gap: 0.35rem;
}

@media (max-width: 768px) {
  .file-row__time,
  .file-row__size {
    display: none;
  }
}
</style>

<template>
  <n-modal
    :show="show"
    preset="card"
    title="上传文件"
    style="width: 600px; max-width: calc(100vw - 2rem);"
    :z-index="1000"
    transform-origin="center"
    @update:show="handleClose"
  >
    <n-spin :show="loading">
      <n-form>
        <n-form-item label="文件选择">
          <n-upload
            ref="uploadRef"
            multiple
            :default-upload="false"
            list-type="image-card"
            @change="handleFileChange"
          >
            <div class="upload-trigger">
              <n-icon size="48" class="text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                  <polyline points="7,10 12,15 17,10"/>
                  <line x1="12" y1="15" x2="12" y2="3"/>
                </svg>
              </n-icon>
              <p class="text-sm text-gray-500 mt-2">点击或拖拽上传</p>
            </div>
          </n-upload>
        </n-form-item>

        <div v-if="uploadFileList.length > 0">
          <n-divider title-placement="left">
            <n-text type="primary" style="font-size: 14px;">文件信息</n-text>
          </n-divider>

          <n-space vertical size="small">
            <n-card
              v-for="fileInfo in uploadFileList"
              :key="fileInfo.id"
              size="small"
              :bordered="false"
              style="background-color: #f9fafb;"
            >
              <n-space vertical size="small">
                <n-space align="center" size="small">
                  <n-text style="font-weight: 500; color: #374151;">{{ fileInfo.name }}</n-text>
                  <n-text depth="3" style="font-size: 12px;">{{ formatFileSize(fileInfo.file?.size || 0) }}</n-text>
                </n-space>

                <n-space size="small" align="center" wrap>
                  <n-input
                    v-model:value="fileInfo.customName"
                    :placeholder="fileInfo.name"
                    size="small"
                    style="width: 200px;"
                  >
                    <template #prefix>
                      <n-text depth="3" style="font-size: 12px;">文件名</n-text>
                    </template>
                  </n-input>

                  <n-date-picker
                    v-model:value="fileInfo.customDate"
                    type="datetime"
                    placeholder="拍摄/创建时间"
                    size="small"
                    style="width: 200px;"
                    clearable
                  />
                  <n-text v-if="fileInfo.exifHint" depth="3" style="font-size: 12px;">{{ fileInfo.exifHint }}</n-text>
                </n-space>
              </n-space>
            </n-card>
          </n-space>
        </div>

        <n-form-item>
          <n-space justify="end">
            <n-button :loading="loading" type="primary" :disabled="uploadFileList.length === 0" @click="handleUpload">
              {{ loading ? '上传中...' : '开始上传' }}
            </n-button>
            <n-button @click="handleReset">重置</n-button>
            <n-button @click="handleClose">取消</n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-spin>
  </n-modal>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  NModal, NUpload, NButton, NIcon, NSpin, NSpace, NForm, NFormItem,
  NInput, NDatePicker, NDivider, NText, NCard, useMessage,
} from 'naive-ui'
import type { UploadFileInfo, UploadInst } from 'naive-ui'
import { UploadFile } from '@/api'
import { uploadToQiniu, getFileTypeFromFile } from '@/util/qiniu'
import { formatFileSize } from '@/util'
import { buildExifPayload, readPhotoOriginAt } from '@/util/exif'

type UploadItem = UploadFileInfo & {
  customName?: string
  customDate?: number | null
  exifHint?: string
}

const props = defineProps<{
  show: boolean
  parentId: number
  folderPath?: string[]
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  uploaded: []
}>()

const message = useMessage()
const loading = ref(false)
const uploadRef = ref<UploadInst | null>(null)
const uploadFileList = ref<UploadItem[]>([])

const handleFileChange = async (options: { fileList: UploadFileInfo[] }) => {
  const next: UploadItem[] = []

  for (const file of options.fileList) {
    const item: UploadItem = {
      ...file,
      customName: undefined,
      customDate: undefined,
      exifHint: '',
    }

    if (file.file) {
      const originAt = await readPhotoOriginAt(file.file)
      if (originAt) {
        item.customDate = originAt * 1000
        item.exifHint = '已从 EXIF 读取拍摄时间'
      }
    }

    next.push(item)
  }

  uploadFileList.value = next
}

const handleUpload = async () => {
  if (uploadFileList.value.length === 0) return

  loading.value = true
  try {
    for (const fileInfo of uploadFileList.value) {
      if (!fileInfo.file) continue

      const result = await uploadToQiniu(fileInfo.file, props.folderPath ?? [])
      const originAt = fileInfo.customDate ? Math.floor(fileInfo.customDate / 1000) : 0
      const exif = buildExifPayload(originAt)

      const response = await UploadFile({
        parent_id: props.parentId,
        files: [{
          filename: fileInfo.customName || fileInfo.name,
          store_key: result.key,
          hash: result.hash,
          size: fileInfo.file.size,
          format: getFileTypeFromFile(fileInfo.file),
          exif,
        }],
      })

      if (!response.file_ids?.length) {
        message.error(response.message || '上传失败')
        return
      }
    }

    message.success('上传成功')
    emit('uploaded')
    handleClose()
  } catch {
    message.error('上传失败')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  uploadFileList.value = []
  uploadRef.value?.clear()
}

const handleClose = () => {
  emit('update:show', false)
  handleReset()
}
</script>

<style scoped>
.upload-trigger {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100px;
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s;
}

.upload-trigger:hover {
  border-color: #3b82f6;
}
</style>

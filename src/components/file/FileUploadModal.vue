<template>
  <n-modal 
    v-model:show="props.show" 
    preset="card" 
    title="上传文件" 
    style="width: 600px;"
    :z-index="1000"
    transform-origin="center"
  >
    <n-form
      ref="uploadFormRef"
      :model="uploadForm"
      label-placement="left"
      label-width="auto"
      require-mark-placement="right-hanging"
      size="medium"
    >
      <n-form-item label="文件选择">
        <n-upload
          ref="uploadRef"
          :max="10"
          :show-file-list="true"
          accept="*/*"
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
            <p class="text-sm text-gray-500 mt-2">点击上传文件</p>
          </div>
        </n-upload>
      </n-form-item>
      <div v-if="uploadFileList.length > 0" class="mt-4">
        <h4 class="text-sm font-medium mb-3">文件信息设置</h4>
        <div class="space-y-3">
          <div 
            v-for="(file, index) in uploadFileList" 
            :key="index"
            class="border border-gray-200 rounded-lg p-3"
          >
            <div class="flex items-center space-x-3">
              <div class="flex-shrink-0">
                <n-icon size="32" class="text-gray-500">
                  <component :is="getFileIcon(getFileTypeFromString(file.type))" />
                </n-icon>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-900 truncate">{{ file.name }}</p>
                <p class="text-xs text-gray-500">{{ formatFileSize(file.size) }}</p>
              </div>
            </div>
            <div class="mt-3 grid grid-cols-2 gap-3">
              <n-form-item :label="`文件名 ${index + 1}`" :path="`fileNames.${index}`">
                <n-input
                  v-model:value="uploadFileList[index].customName"
                  :placeholder="file.name"
                  clearable
                />
              </n-form-item>
              <n-form-item :label="`创建时间 ${index + 1}`" :path="`fileDates.${index}`">
                <n-date-picker
                  v-model:value="uploadFileList[index].customDate"
                  type="datetime"
                  placeholder="选择时间"
                  clearable
                />
              </n-form-item>
            </div>
          </div>
        </div>
      </div>
      <n-form-item>
        <n-space>
          <n-button
            type="primary"
            :loading="uploading"
            :disabled="uploadFileList.length === 0"
            @click="handleUpload"
          >
            {{ uploading ? '上传中...' : '开始上传' }}
          </n-button>
          <n-button @click="handleUploadReset">重置</n-button>
          <n-button @click="close">取消</n-button>
        </n-space>
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch, defineProps, defineEmits, withDefaults } from 'vue';
import { 
  NModal,
  NIcon, 
  NButton, 
  NForm, 
  NFormItem, 
  NInput, 
  NDatePicker, 
  NUpload, 
  NSpace, 
  useMessage, 
  type FormInst, 
  type UploadInst 
} from 'naive-ui';
import { getFileIcon, getFileTypeFromString, formatFileSize } from '@/util';
import { UploadFile } from '@/api';
import { UploadFileReq, FileInfo } from '@/api/types/file';

const props = withDefaults(defineProps<{
  show: boolean;
  parentId: number;
}>(), {
  show: false,
  parentId: 0
});

const emit = defineEmits<{
  'update:show': [value: boolean];
  'uploaded': [];
}>();

const uploadFormRef = ref<FormInst | null>(null);
const uploadRef = ref<UploadInst | null>(null);
const uploadFileList = ref<Array<{ file: File; name: string; size: number; type: string; customName?: string; customDate?: number }>>([]);
const uploadForm = ref({ files: [] });
const uploading = ref(false);
const message = useMessage();

const handleFileChange = (options: any) => {
  const { fileList: newFileList } = options;
  uploadFileList.value = newFileList.map((item: any) => ({
    file: item.file,
    name: item.file.name,
    size: item.file.size,
    type: item.file.type,
    customName: undefined,
    customDate: undefined
  }));
};

const handleUpload = async () => {
  if (uploadFileList.value.length === 0) {
    message.error('请选择文件');
    return;
  }
  try {
    uploading.value = true;
    const { batchUploadToQiniu, getFileTypeFromFile } = await import('@/util/qiniu');
    const folderPath: string[] = [];
    const uploadResults = await batchUploadToQiniu(
      uploadFileList.value.map(item => item.file),
      folderPath
    );
    const files: FileInfo[] = uploadResults.map((result, index) => {
      const fileInfo = uploadFileList.value[index];
      return {
        filename: fileInfo.customName || fileInfo.name,
        store_key: result.key,
        hash: result.hash,
        format: getFileTypeFromFile(fileInfo.file),
        origin_at: fileInfo.customDate ? Math.floor(fileInfo.customDate / 1000) : 0
      };
    });
    const uploadData: UploadFileReq = {
      parent_id: props.parentId,
      files
    };
    const response = await UploadFile(uploadData);
    if (response.file_ids && response.file_ids.length > 0) {
      message.success(`成功上传 ${response.file_ids.length} 个文件`);
      emit('update:show', false);
      emit('uploaded');
      handleUploadReset();
    } else {
      message.error(response.message || '上传失败');
    }
  } catch (error) {
    message.error('上传失败: ' + (error instanceof Error ? error.message : '未知错误'));
  } finally {
    uploading.value = false;
  }
};

const handleUploadReset = () => {
  uploadFileList.value = [];
  if (uploadRef.value) {
    uploadRef.value.clear();
  }
};

const close = () => {
  emit('update:show', false);
  handleUploadReset();
};
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
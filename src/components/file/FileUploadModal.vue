<template>
  <n-modal
    :show="show"
    @update:show="handleClose"
    preset="card"
    title="上传文件"
    style="width: 600px;"
    :z-index="1000"
    transform-origin="center"
  >
    <n-spin :show="loading">
      <n-form>
        <n-form-item label="文件选择">
          <n-upload
            ref="uploadRef"
            multiple
            :default-upload="false"
            @change="handleFileChange"
            list-type="image-card"
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

        <n-form-item>
          <n-space justify="end">
            <n-button :loading="loading" type="primary" @click="handleUpload" :disabled="uploadFileList.length === 0">
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
import { ref } from 'vue';
import { NModal, NUpload, NButton, NIcon, NSpin, NSpace, NForm, NFormItem, useMessage } from 'naive-ui';
import type { UploadFileInfo, UploadInst } from 'naive-ui';
import { UploadFile } from '@/api';
import { uploadToQiniu, getFileTypeFromFile } from '@/util/qiniu';
import { formatFileSize } from '@/util';

const props = defineProps<{
  show: boolean;
  parentId: number;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
  'uploaded': [];
}>();

const message = useMessage();
const loading = ref(false);
const uploadRef = ref<UploadInst | null>(null);
const uploadFileList = ref<UploadFileInfo[]>([]);

// 处理文件选择
const handleFileChange = (options: { fileList: UploadFileInfo[] }) => {
  uploadFileList.value = options.fileList;
};

// 处理上传
const handleUpload = async () => {
  if (uploadFileList.value.length === 0) return;
  
  loading.value = true;
  try {
    for (const fileInfo of uploadFileList.value) {
      if (!fileInfo.file) continue;
      
      // 上传到七牛云
      const result = await uploadToQiniu(fileInfo.file);
      
      // 调用后端接口
      const response = await UploadFile({
        parent_id: props.parentId,
        files: [{
          filename: fileInfo.name,
          store_key: result.key,
          hash: result.hash,
          format: getFileTypeFromFile(fileInfo.file),
          origin_at: 0
        }]
      });

      if (response.code === 0) {
        message.success('上传成功');
        emit('uploaded');
      } else {
        message.error(response.message || '上传失败');
      }
    }
    
    // 上传完成后关闭弹窗
    handleClose();
  } catch (error) {
    console.error('上传失败:', error);
    message.error('上传失败');
  } finally {
    loading.value = false;
  }
};

// 处理重置
const handleReset = () => {
  uploadFileList.value = [];
  if (uploadRef.value) {
    uploadRef.value.clear();
  }
};

// 处理关闭
const handleClose = () => {
  emit('update:show', false);
  handleReset();
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
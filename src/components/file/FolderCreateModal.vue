<template>
  <n-modal
    :show="show"
    @update:show="handleClose"
    preset="card"
    title="新建文件夹"
    style="width: 400px;"
    :z-index="1000"
    transform-origin="center"
  >
    <n-spin :show="loading">
      <n-form>
        <n-form-item label="文件夹名称">
          <n-input v-model:value="folderName" placeholder="请输入文件夹名称" />
        </n-form-item>

        <n-form-item>
          <n-space justify="end">
            <n-button :loading="loading" type="primary" @click="handleCreate" :disabled="!folderName">
              {{ loading ? '创建中...' : '创建' }}
            </n-button>
            <n-button @click="handleClose">取消</n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-spin>
  </n-modal>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { NModal, NInput, NButton, NIcon, NSpin, NSpace, NForm, NFormItem, useMessage } from 'naive-ui';
import { UploadFile } from '@/api';
import { FileType } from '@/api/types/file';

const props = defineProps<{
  show: boolean;
  parentId: number;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
  'created': [];
}>();

const message = useMessage();
const loading = ref(false);
const folderName = ref('');

// 处理创建
const handleCreate = async () => {
  if (!folderName.value) return;
  
  loading.value = true;
  try {
    const response = await UploadFile({
      parent_id: props.parentId,
      files: [{
        filename: folderName.value,
        store_key: '',
        hash: '',
        format: FileType.FT_Folder,
        origin_at: 0
      }]
    });

    if (response.file_ids.length > 0 ) {
      emit('created');
      handleClose();
    } else {
      message.error(response.message || '创建失败');
    }
  } catch (error) {
    message.error('创建失败');
  } finally {
    loading.value = false;
  }
};

// 处理关闭
const handleClose = () => {
  emit('update:show', false);
  folderName.value = '';
};
</script> 
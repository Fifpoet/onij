<template>
  <n-modal 
    v-model:show="props.show" 
    preset="card" 
    title="新建文件夹" 
    style="width: 400px;"
    :z-index="1000"
    transform-origin="center"
  >
    <n-form
      ref="folderFormRef"
      :model="folderForm"
      :rules="folderRules"
      label-placement="left"
      label-width="auto"
      require-mark-placement="right-hanging"
      size="medium"
    >
      <n-form-item label="文件夹名称" path="filename">
        <n-input
          v-model:value="folderForm.filename"
          placeholder="请输入文件夹名称"
          clearable
        />
      </n-form-item>
      <n-form-item>
        <n-space>
          <n-button
            type="primary"
            :loading="creatingFolder"
            :disabled="!folderForm.filename"
            @click="handleCreateFolder"
          >
            {{ creatingFolder ? '创建中...' : '创建文件夹' }}
          </n-button>
          <n-button @click="handleFolderReset">重置</n-button>
          <n-button @click="close">取消</n-button>
        </n-space>
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, defineProps, defineEmits } from 'vue';
import { 
  NModal,
  NForm, 
  NFormItem, 
  NInput, 
  NButton, 
  NSpace, 
  useMessage, 
  type FormInst 
} from 'naive-ui';
import { UploadFile } from '@/api';
import { FileType, UploadFileReq } from '@/api/types/file';

const props = defineProps<{
  show: boolean;
  parentId: number;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
  'created': [];
}>();

const folderFormRef = ref<FormInst | null>(null);
const folderForm = ref({ filename: '' });
const creatingFolder = ref(false);
const folderRules = {
  filename: {
    required: true,
    message: '请输入文件夹名称',
    trigger: 'blur'
  }
};

const message = useMessage();

const handleCreateFolder = async () => {
  if (!folderFormRef.value) return;
  try {
    await folderFormRef.value.validate();
    creatingFolder.value = true;
    const folderData: UploadFileReq = {
      parent_id: props.parentId,
      files: [{
        filename: folderForm.value.filename,
        store_key: '',
        hash: `folder_${Date.now()}_${folderForm.value.filename}`,
        format: FileType.FT_Folder,
        origin_at: Math.floor(Date.now() / 1000)
      }]
    };
    const response = await UploadFile(folderData);
    if (response.file_ids && response.file_ids.length > 0) {
      message.success('文件夹创建成功');
      emit('update:show', false);
      emit('created');
      handleFolderReset();
    } else {
      message.error(response.message || '创建文件夹失败');
    }
  } catch (error) {
    message.error('创建文件夹失败');
  } finally {
    creatingFolder.value = false;
  }
};

const handleFolderReset = () => {
  folderForm.value.filename = '';
};

const close = () => {
  emit('update:show', false);
  handleFolderReset();
};
</script> 
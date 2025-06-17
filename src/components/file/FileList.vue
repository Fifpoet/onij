<template>
  <Transition name="fade">
    <div v-if="musicStore.midShowWhat == MidShowWhat.ShowFileList">
      <!-- 面包屑导航 -->
      <div class="px-5 pt-5 pb-3">
        <div v-if="parentId !== 0" class="flex items-center space-x-3">
          <n-button 
            size="small" 
            type="primary" 
            ghost 
            @click="goBack"
            class="flex items-center"
          >
            <template #icon>
              <n-icon>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m15 18-6-6 6-6"/>
                </svg>
              </n-icon>
            </template>
            返回上级
          </n-button>
          <span class="text-gray-400">|</span>
          <span class="text-sm text-gray-600">{{ currentFolderName }}</span>
        </div>
        <div v-else class="h-8 flex items-center">
          <span class="text-sm text-gray-600">根目录</span>
        </div>
      </div>

      <!-- 文件列表区域 -->
      <div class="p-5" style="width: 100%; min-width: 900px; margin: 0 auto; position: relative;">
        <!-- 操作按钮区域 -->
        <div class="flex justify-end space-x-2 mb-4">
          <n-button type="primary" @click="showUploadModal = true">
            <template #icon>
              <n-icon>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/>
                  <polyline points="14,2 14,8 20,8"/>
                </svg>
              </n-icon>
            </template>
            上传文件
          </n-button>
          <n-button type="info" @click="showFolderModal = true">
            <template #icon>
              <n-icon>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2Z"/>
                </svg>
              </n-icon>
            </template>
            新建文件夹
          </n-button>
        </div>

        <!-- 文件网格 -->
        <div class="grid gap-4" style="grid-template-columns: repeat(auto-fit, minmax(200px, 200px));">
          <div 
            v-for="file in fileList" 
            :key="file.id" 
            class="bg-white rounded-lg border border-gray-200 p-4 cursor-pointer transition-all duration-200 hover:shadow-md hover:border-gray-300 h-48 flex flex-col min-w-0 max-w-full"
            style="width: 100%; box-sizing: border-box;"
            @click="handleFileClick(file)"
          >
            <!-- 文件图标 -->
            <div class="flex justify-center mb-3 flex-shrink-0">
              <n-icon size="48" class="text-gray-500">
                <component :is="getFileIcon(file.format)" />
              </n-icon>
            </div>
            
            <!-- 文件名称 - 固定高度，最多两行 -->
            <div class="text-center mb-2 flex-1 flex flex-col justify-center min-h-0 min-w-0 overflow-hidden">
              <n-ellipsis 
                :line-clamp="2" 
                class="text-sm font-medium text-gray-900 w-full"
                :title="file.name"
              >
                {{ file.name }}
              </n-ellipsis>
            </div>
            
            <!-- 文件信息 -->
            <div class="text-center text-xs text-gray-500 space-y-1 flex-shrink-0 overflow-hidden">
              <div v-if="!isFolder(file.format)" class="truncate px-1">{{ getFileSize(file.format) }}</div>
              <div v-if="!isFolder(file.format)" class="truncate px-1">{{ formatFileTime(file.origin_at) }}</div>
            </div>
            
            <!-- 操作按钮 -->
            <div v-if="!isFolder(file.format)" class="flex justify-center mt-3 space-x-2 flex-shrink-0" @click.stop>
              <n-button size="tiny" type="primary" @click="downloadFile(file)">
                <template #icon>
                  <n-icon>
                    <DownloadOutline />
                  </n-icon>
                </template>
                下载
              </n-button>
              <n-button size="tiny" type="error" @click="deleteFile(file)">
                <template #icon>
                  <n-icon>
                    <TrashOutline />
                  </n-icon>
                </template>
                删除
              </n-button>
            </div>
          </div>
        </div>
        
        <!-- 加载状态 - 覆盖在网格上方 -->
        <div v-if="loading" class="absolute inset-0 bg-white bg-opacity-80 flex justify-center items-center z-10">
          <n-spin size="large" />
        </div>
        
        <!-- 空状态 -->
        <div v-if="fileList.length === 0 && !loading" class="flex justify-center py-16">
          <n-empty description="暂无文件" />
        </div>
      </div>
      
      <!-- 分页控件 -->
      <div v-if="fileList.length > 0" class="border-t border-gray-200 bg-white p-4">
        <div class="flex flex-col items-center space-y-2">
          <!-- 分页信息 -->
          <div class="text-sm text-gray-600">
            共 {{ totalCount }} 个文件，
            第 {{ currentPage }} 页，
            每页 {{ pageSize }} 个，
            <span v-if="isLastPage">本页 {{ fileList.length }} 个</span>
            <span v-else>本页 {{ pageSize }} 个</span>
          </div>
          
          <!-- 分页控件 -->
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
  </Transition>

  <!-- 文件上传弹框 -->
  <n-modal v-model:show="showUploadModal" preset="card" title="上传文件" style="width: 600px">
    <n-form
      ref="uploadFormRef"
      :model="uploadForm"
      :rules="uploadRules"
      label-placement="left"
      label-width="auto"
      require-mark-placement="right-hanging"
      size="medium"
    >
      <n-form-item label="文件选择" path="files">
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

      <!-- 文件信息编辑区域 -->
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
          <n-button @click="showUploadModal = false">取消</n-button>
        </n-space>
      </n-form-item>
    </n-form>
  </n-modal>

  <!-- 新建文件夹弹框 -->
  <n-modal v-model:show="showFolderModal" preset="card" title="新建文件夹" style="width: 400px">
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
          <n-button @click="showFolderModal = false">取消</n-button>
        </n-space>
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { 
  NIcon, 
  NButton, 
  NSpin, 
  NEmpty,
  NPagination,
  NEllipsis,
  NModal,
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
import { 
  DownloadOutline,
  TrashOutline
} from '@vicons/ionicons5';
import { GetFileList, DownloadFile, UploadFile } from '@/api';
import { FileType, FileDetail, GetFileListReq, UploadFileReq, FileInfo } from '@/api/types/file';
import { MidShowWhat } from '@/api/types';
import { useMusicStore } from '@/store/music';
import { 
  getFileIcon, 
  getFileSize, 
  formatFileTime, 
  isFolder,
  getFileTypeFromString,
  formatFileSize
} from '@/util';

// 使用全局状态
const musicStore = useMusicStore();

// Props
interface Props {
  parentId?: number;
}

const props = withDefaults(defineProps<Props>(), {
  parentId: 0
});

// Emits
const emit = defineEmits<{
  fileClick: [file: FileDetail];
}>();

// 响应式数据
const fileList = ref<FileDetail[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(10);
const hasMore = ref(true);
const totalCount = ref(0);
const totalPages = ref(0);
const isLastPage = ref(false);

// 面包屑导航相关
const currentFolderName = ref('根目录');
const folderHistory = ref<number[]>([0]);
const folderNameHistory = ref<string[]>(['根目录']);

// 弹框相关状态
const showUploadModal = ref(false);
const showFolderModal = ref(false);
const uploading = ref(false);
const creatingFolder = ref(false);

// 表单引用
const uploadFormRef = ref<FormInst | null>(null);
const folderFormRef = ref<FormInst | null>(null);
const uploadRef = ref<UploadInst | null>(null);

// 文件列表（用于上传）
const uploadFileList = ref<Array<{
  file: File;
  name: string;
  size: number;
  type: string;
  customName?: string;
  customDate?: number;
}>>([]);

// 上传表单
const uploadForm = ref({
  files: []
});

// 文件夹表单
const folderForm = ref({
  filename: ''
});

// 表单验证规则
const uploadRules = {
  files: {
    required: true,
    message: '请选择文件',
    trigger: 'change'
  }
};

const folderRules = {
  filename: {
    required: true,
    message: '请输入文件夹名称',
    trigger: 'blur'
  }
};

// 消息提示
const message = useMessage();

// 获取文件列表
const fetchFileList = async (page: number = 1, append: boolean = false) => {
  if (loading.value) return;
  
  loading.value = true;
  try {
    const req: GetFileListReq = {
      parent_id: props.parentId,
      page,
      limit: pageSize.value
    };
    
    const response = await GetFileList(req);
    
    // 直接处理数据，不检查code
    if (response.files) {
      if (append) {
        fileList.value.push(...response.files);
      } else {
        fileList.value = response.files;
      }
      
      hasMore.value = response.files.length === pageSize.value;
      currentPage.value = page;
      totalCount.value = response.total || 0;
      totalPages.value = Math.ceil(totalCount.value / pageSize.value);
      isLastPage.value = currentPage.value >= totalPages.value;
    } else {
      // 如果没有files字段，清空列表
      fileList.value = [];
      hasMore.value = false;
      totalCount.value = 0;
      totalPages.value = 0;
      isLastPage.value = true;
    }
  } catch (error) {
    console.error('获取文件列表失败:', error);
    // 出错时清空列表
    fileList.value = [];
    hasMore.value = false;
    totalCount.value = 0;
    totalPages.value = 0;
    isLastPage.value = true;
  } finally {
    loading.value = false;
  }
};

// 监听parentId变化
watch(() => props.parentId, () => {
  fetchFileList();
}, { immediate: true });

// 处理文件点击
const handleFileClick = (file: FileDetail) => {
  if (isFolder(file.format)) {
    // 进入文件夹
    folderHistory.value.push(file.id);
    folderNameHistory.value.push(file.name);
    currentFolderName.value = file.name;
    console.log('进入文件夹:', file.name, 'ID:', file.id);
  }
  emit('fileClick', file);
};

// 返回上级文件夹
const goBack = () => {
  if (folderHistory.value.length > 1) {
    folderHistory.value.pop();
    folderNameHistory.value.pop();
    const newParentId = folderHistory.value[folderHistory.value.length - 1];
    const newFolderName = folderNameHistory.value[folderNameHistory.value.length - 1];
    currentFolderName.value = newFolderName;
    
    // 通知父组件更新parentId
    emit('fileClick', { 
      id: newParentId, 
      format: FileType.FT_Folder, 
      name: newFolderName 
    } as FileDetail);
  }
};

// 处理页码变化
const handlePageChange = (page: number) => {
  currentPage.value = page;
  fetchFileList(page);
};

// 处理每页数量变化
const handlePageSizeChange = (size: number) => {
  pageSize.value = size;
  currentPage.value = 1;
  fetchFileList(1);
};

// 下载文件
const downloadFile = async (file: FileDetail) => {
  try {
    const response = await DownloadFile({ file_ids: [file.id] });
    
    if (response.urls && response.urls.length > 0) {
      const fileUrl = response.urls[0];
      const fileName = file.name;
      
      // 创建一个隐藏的iframe来强制下载
      const iframe = document.createElement('iframe');
      iframe.style.display = 'none';
      iframe.src = fileUrl;
      document.body.appendChild(iframe);
      
      // 同时创建一个下载链接作为备用方案
      const link = document.createElement('a');
      link.href = fileUrl;
      link.download = fileName;
      link.target = '_blank';
      link.rel = 'noopener noreferrer';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      
      // 延迟移除iframe
      setTimeout(() => {
        document.body.removeChild(iframe);
      }, 1000);
    }
  } catch (error) {
    message.error('下载文件失败');
  }
};

// 删除文件
const deleteFile = async (file: FileDetail) => {
  // TODO: 实现删除文件接口
  console.log('删除文件:', file.name);
};

// 处理文件选择
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

// 处理文件上传
const handleUpload = async () => {
  if (!uploadFormRef.value || uploadFileList.value.length === 0) return;
  
  try {
    await uploadFormRef.value.validate();
    
    uploading.value = true;
    
    const files: FileInfo[] = uploadFileList.value.map(file => ({
      filename: file.customName || file.name,
      file: file.file,
      origin_at: file.customDate ? Math.floor(file.customDate / 1000) : undefined
    }));
    
    const uploadData: UploadFileReq = {
      parent_id: props.parentId,
      files
    };
    
    const response = await UploadFile(uploadData);
    
    if (response.code === 0) {
      message.success(`成功上传 ${response.file_ids.length} 个文件`);
      showUploadModal.value = false;
      handleUploadReset();
      fetchFileList(); // 刷新文件列表
    } else {
      message.error(response.message || '上传失败');
    }
  } catch (error) {
    console.error('上传失败:', error);
    message.error('上传失败');
  } finally {
    uploading.value = false;
  }
};

// 处理创建文件夹
const handleCreateFolder = async () => {
  if (!folderFormRef.value) return;
  
  try {
    await folderFormRef.value.validate();
    creatingFolder.value = true;
    // 使用UploadFile接口创建文件夹
    const folderData: UploadFileReq = {
      parent_id: props.parentId,
      files: [{
        filename: folderForm.value.filename,
        file: new File([], folderForm.value.filename, { type: 'application/octet-stream' }), // 空文件表示创建文件夹
        origin_at: Math.floor(Date.now() / 1000)
      }]
    };
    
    const response = await UploadFile(folderData);
    console.log(response);

    if (response.file_ids.length > 0) {
      message.success('文件夹创建成功');
      showFolderModal.value = false;
      handleFolderReset();
      fetchFileList(); // 刷新文件列表
    } else {
      message.error(response.message || '创建文件夹失败');
    }
  } catch (error) {
    message.error('创建文件夹失败');
  } finally {
    creatingFolder.value = false;
  }
};

// 重置上传表单
const handleUploadReset = () => {
  uploadFileList.value = [];
  if (uploadRef.value) {
    uploadRef.value.clear();
  }
};

// 重置文件夹表单
const handleFolderReset = () => {
  folderForm.value.filename = '';
};

// 暴露刷新方法
const refresh = () => {
  fetchFileList();
};

defineExpose({
  refresh
});
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

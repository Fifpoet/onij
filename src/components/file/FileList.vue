<template>
  <Transition name="fade">
    <div v-if="musicStore.midShowWhat == MidShowWhat.ShowFileList">
      <!-- 面包屑导航 -->
      <div class="px-5 pt-5 pb-3 flex items-center space-x-3">
        <!-- 返回按钮 -->
        <n-button size="small" quaternary :disabled="parentId === 0" @click="goBack"
          class="flex items-center h-8 min-w-[90px] justify-center">
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

        <!-- 分隔符和当前目录 -->
        <div class="flex items-center space-x-3 shrink-0">
          <span class="text-gray-300">|</span>
          <n-button size="small" text class="px-2 !text-gray-600">
            <template #icon>
              <n-icon>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path
                    d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2Z" />
                </svg>
              </n-icon>
            </template>
            {{ currentFolderName }}
          </n-button>
        </div>

        <!-- 搜索框 -->
        <div class="flex-1 flex justify-end">
          <n-input v-model:value="searchKeyword" type="text" placeholder="搜索文件..." @keydown.enter="handleSearch"
            class="max-w-[200px]" size="small">
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

      <!-- 操作按钮区域 -->
      <div class="flex justify-end space-x-3 mb-4">
        <n-button quaternary size="small" class="h-8 w-8 flex items-center justify-center" @click="toggleViewMode"
          :type="viewMode === 'grid' ? 'default' : 'primary'">
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
        <n-button size="small" quaternary class="h-8 flex items-center space-x-1 px-3" @click="showUploadModal = true">
          <n-icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
              <polyline points="17 8 12 3 7 8" />
              <line x1="12" y1="3" x2="12" y2="15" />
            </svg>
          </n-icon>
          <span>上传文件</span>
        </n-button>
        <n-button size="small" quaternary class="h-8 flex items-center space-x-1 px-3" @click="showFolderModal = true">
          <n-icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path
                d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2Z" />
              <line x1="12" y1="10" x2="12" y2="16" />
              <line x1="9" y1="13" x2="15" y2="13" />
            </svg>
          </n-icon>
          <span>新建文件夹</span>
        </n-button>
      </div>

      <!-- 文件列表区域 -->
      <div class="p-5" style="width: 100%; min-width: 900px; min-height: 720px; margin: 0 auto; position: relative;">
        <!-- 文件网格 -->
        <div v-if="viewMode === 'grid'" class="grid gap-4"
          style="grid-template-columns: repeat(auto-fit, minmax(160px, 120px));">
          <div v-for="file in fileList" :key="file.id"
            class="bg-white rounded-lg border border-gray-200 p-4 cursor-pointer transition-all duration-200 hover:shadow-md hover:border-gray-300 flex flex-col min-w-0 max-w-full"
            style="width: 100%; aspect-ratio: 1/1.2; box-sizing: border-box;" @click="handleFileClick(file)">
            
            <!-- 预览区域 - 30% -->
            <div class="w-full mb-3" style="height: 50%;">
              <!-- 文件图标或图片 -->
              <div v-if="showFileIcon(file.format)" class="flex justify-center items-center h-full">
                <n-icon size="40" class="text-gray-500">
                  <component :is="getFileIcon(file.format)" />
                </n-icon>
              </div>
              <div v-else class="w-full h-full overflow-hidden">
                <img :src="file.url" :alt="file.name" class="w-full h-full object-cover">
              </div>
            </div>
            
            <!-- 文件名称 - 25% -->
            <div class="text-center mb-2 flex flex-col justify-center" style="height: 15%;">
              <n-ellipsis :line-clamp="2" class="text-sm font-medium text-gray-900 w-full" :title="file.name">
                {{ file.name }}
              </n-ellipsis>
            </div>

            <!-- 文件信息 - 20% -->
            <div class="text-center text-xs text-gray-500 mb-2" style="height: 10%;">
              <template v-if="!isFolder(file.format)">
                <div class="truncate">{{ getFileSize(file.format) }}</div>
                <div class="truncate">{{ formatFileTime(file.origin_at) }}</div>
              </template>
            </div>

            <!-- 操作按钮 - 25% -->
            <div class="flex justify-center items-end mt-auto" style="height: 15%;" @click.stop>
              <template v-if="!isFolder(file.format)">
                <n-button size="tiny" type="primary" class="mx-1" @click="downloadFile(file)">
                  <template #icon><n-icon><DownloadOutline /></n-icon></template>
                  下载
                </n-button>
                <n-button size="tiny" type="error" class="mx-1" @click="deleteFile(file)">
                  <template #icon><n-icon><TrashOutline /></n-icon></template>
                  删除
                </n-button>
              </template>
            </div>
          </div>
        </div>

        <!-- 文件列表 -->
        <div v-else class="space-y-2">
          <div v-for="file in fileList" :key="file.id"
            class="bg-white rounded-lg border border-gray-200 px-4 py-3 cursor-pointer transition-all duration-200 hover:shadow-md hover:border-gray-300 flex items-center space-x-4"
            @click="handleFileClick(file)">
            <!-- 文件图标 -->
            <n-icon size="24" class="text-gray-500 flex-shrink-0">
              <component :is="getFileIcon(file.format)" />
            </n-icon>

            <!-- 文件名称 -->
            <div class="flex-1 min-w-0">
              <n-ellipsis class="text-sm font-medium text-gray-900" :title="file.name">
                {{ file.name }}
              </n-ellipsis>
            </div>

            <!-- 文件信息 -->
            <div v-if="!isFolder(file.format)" class="text-sm text-gray-500 flex-shrink-0 w-32 text-right">
              {{ getFileSize(file.format) }}
            </div>
            <div v-if="!isFolder(file.format)" class="text-sm text-gray-500 flex-shrink-0 w-48 text-right">
              {{ formatFileTime(file.origin_at) }}
            </div>

            <!-- 操作按钮 -->
            <div v-if="!isFolder(file.format)" class="flex space-x-2 flex-shrink-0" @click.stop>
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
          <n-pagination v-model:page="currentPage" :page-count="totalPages" :page-sizes="[10, 20, 50, 100]"
            :page-size="pageSize" show-size-picker @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange" />
        </div>
      </div>
    </div>
  </Transition>

  <!-- 文件上传弹框 -->
  <FileUploadModal v-model:show="showUploadModal" :parent-id="props.parentId" @uploaded="fetchFileList" />

  <!-- 新建文件夹弹框 -->
  <FolderCreateModal v-model:show="showFolderModal" :parent-id="props.parentId" @created="fetchFileList" />
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
  useMessage,
  NInput
} from 'naive-ui';
import {
  DownloadOutline,
  TrashOutline
} from '@vicons/ionicons5';
import { GetFileList, DownloadFile, DeleteFileById } from '@/api';
import { FileType, FileDetail, GetFileListReq, DeleteFileReq } from '@/api/types/file';
import { MidShowWhat } from '@/api/types';
import { useMusicStore } from '@/store/music';
import {
  getFileIcon,
  getFileSize,
  formatFileTime,
  isFolder,
  showFileIcon
} from '@/util';
import FileUploadModal from './FileUploadModal.vue';
import FolderCreateModal from './FolderCreateModal.vue';

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
const searchKeyword = ref('');

// 面包屑导航相关
const currentFolderName = ref('根目录');
const folderHistory = ref<number[]>([0]);
const folderNameHistory = ref<string[]>(['根目录']);

// 弹框相关状态
const showUploadModal = ref(false);
const showFolderModal = ref(false);

// 视图模式
const viewMode = ref<'grid' | 'list'>('grid');

// 获取文件列表
const fetchFileList = async (page: number = 1, append: boolean = false) => {
  if (loading.value) return;

  loading.value = true;
  try {
    const req: GetFileListReq = {
      parent_id: props.parentId,
      page,
      limit: pageSize.value,
      keyword: searchKeyword.value.trim()
    };

    const response = await GetFileList(req);

    // 直接判断返回的files数组
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

// 新增搜索处理
const handleSearch = () => {
  currentPage.value = 1;
  fetchFileList(1);
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

    // 直接判断返回的urls数组
    const fileUrl = file.url;

    // 创建一个隐藏的iframe来强制下载
    const iframe = document.createElement('iframe');
    iframe.style.display = 'none';
    iframe.src = fileUrl;
    document.body.appendChild(iframe);

    // 延迟移除iframe
    setTimeout(() => {
      document.body.removeChild(iframe);
    }, 1000);

  } catch (error) {
    message.error('下载文件失败');
  }
};

// 删除文件
const deleteFile = async (file: FileDetail) => {
  try {
    // 显示确认对话框
    const confirmed = window.confirm(`确定要删除文件 "${file.name}" 吗？此操作不可恢复。`);
    if (!confirmed) {
      return;
    }

    const deleteData: DeleteFileReq = {
      file_id: file.id
    };

    const response = await DeleteFileById(deleteData);

    // 直接判断返回数据，不依赖状态码
    if (response.message === 'ok') {
      message.success('文件删除成功');
      fetchFileList(); // 刷新文件列表
    } else {
      message.error(response.message || '删除文件失败');
    }
  } catch (error) {
    message.error('删除文件失败');
  }
};

// 切换视图模式
const toggleViewMode = () => {
  viewMode.value = viewMode.value === 'grid' ? 'list' : 'grid';
};

const message = useMessage();

defineExpose({
  fetchFileList,
  handleFileClick,
  goBack,
  handlePageChange,
  handlePageSizeChange,
  downloadFile,
  deleteFile
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

<template>
  <Transition name="fade">
    <div v-if="musicStore.midShowWhat == MidShowWhat.ShowFileList">
      <!-- 文件列表区域 -->
      <div class="p-5" style="width: 100%; min-width: 900px; margin: 0 auto; position: relative;">
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
              <div v-if="file.format !== FileType.FT_Folder" class="truncate px-1">{{ getFileSize(file.format) }}</div>
              <div v-if="file.format !== FileType.FT_Folder" class="truncate px-1">{{ formatTime(file.origin_at) }}</div>
            </div>
            
            <!-- 操作按钮 -->
            <div v-if="file.format !== FileType.FT_Folder" class="flex justify-center mt-3 space-x-2 flex-shrink-0" @click.stop>
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
  useMessage
} from 'naive-ui';
import { 
  FolderOutline, 
  MusicalNoteOutline, 
  DocumentTextOutline, 
  ImageOutline,
  DownloadOutline,
  TrashOutline
} from '@vicons/ionicons5';
import { GetFileList, DownloadFile } from '@/api';
import { FileType, FileDetail, GetFileListReq } from '@/api/types/file';
import { MidShowWhat } from '@/api/types';
import { useMusicStore } from '@/store/music';

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

// 获取文件图标
const getFileIcon = (fileType: FileType) => {
  switch (fileType) {
    case FileType.FT_Folder:
      return FolderOutline;
    case FileType.FT_Mp3:
      return MusicalNoteOutline;
    case FileType.FT_Lyrics:
      return DocumentTextOutline;
    case FileType.FT_Png:
      return ImageOutline;
    default:
      return DocumentTextOutline;
  }
};

// 获取文件大小（暂时写死）
const getFileSize = (fileType: FileType): string => {
  if (fileType === FileType.FT_Folder) {
    return '';
  }
  
  // 暂时写死的大小，后续接口返回后可以替换
  const sizeMap = {
    [FileType.FT_Mp3]: '3.2 MB',
    [FileType.FT_Lyrics]: '0.1 MB',
    [FileType.FT_Png]: '1.5 MB',
  };
  
  return sizeMap[fileType] || '0.1 MB';
};

// 格式化时间
const formatTime = (timestamp: number): string => {
  const date = new Date(timestamp * 1000);
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '';
  }
  
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// 处理文件点击
const handleFileClick = (file: FileDetail) => {
  emit('fileClick', file);
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
      // 创建下载链接
      const link = document.createElement('a');
      link.href = response.urls[0];
      link.download = file.name;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
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

// 暴露刷新方法
const refresh = () => {
  fetchFileList();
};

defineExpose({
  refresh
});
</script>

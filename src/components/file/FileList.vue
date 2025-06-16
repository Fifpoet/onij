<template>
  <div class="p-5">
    <n-card title="文件列表" class="w-full max-w-4xl mx-auto">
      <n-list hoverable>
        <n-list-item 
          v-for="file in fileList" 
          :key="file.id" 
          class="border-b border-gray-100 cursor-pointer transition-colors duration-200 hover:bg-gray-50 last:border-b-0"
          clickable
          @click="handleFileClick(file)"
        >
          <template #prefix>
            <n-icon size="24" class="text-gray-500">
              <component :is="getFileIcon(file.format)" />
            </n-icon>
          </template>
          
          <n-thing :title="file.name" content-style="margin-top: 4px;">
            <template #description>
              <div class="flex gap-4 text-gray-500 text-xs">
                <span class="min-w-15">{{ getFileSize(file.format) }}</span>
                <span class="min-w-30">{{ formatTime(file.origin_at) }}</span>
              </div>
            </template>
          </n-thing>
          
          <template #suffix>
            <div v-if="file.format !== FileType.FT_Folder" class="flex items-center" @click.stop>
              <n-button size="small" type="primary" @click="downloadFile(file)">
                <template #icon>
                  <n-icon>
                    <DownloadOutline />
                  </n-icon>
                </template>
                下载
              </n-button>
              <n-button size="small" type="error" @click="deleteFile(file)" class="ml-2">
                <template #icon>
                  <n-icon>
                    <TrashOutline />
                  </n-icon>
                </template>
                删除
              </n-button>
            </div>
          </template>
        </n-list-item>
      </n-list>
      
      <div v-if="loading" class="flex justify-center py-5">
        <n-spin size="large" />
      </div>
      
      <div v-if="fileList.length === 0 && !loading" class="py-10">
        <n-empty description="暂无文件" />
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { 
  NCard, 
  NList, 
  NListItem, 
  NThing, 
  NIcon, 
  NButton, 
  NSpin, 
  NEmpty,
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
const pageSize = ref(20);
const hasMore = ref(true);

// 消息提示
const message = useMessage();

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
    
    if (response.code === 0 && response.files) {
      if (append) {
        fileList.value.push(...response.files);
      } else {
        fileList.value = response.files;
      }
      
      hasMore.value = response.files.length === pageSize.value;
      currentPage.value = page;
    } else {
      message.error(response.message || '获取文件列表失败');
    }
  } catch (error) {
    console.error('获取文件列表失败:', error);
    message.error('获取文件列表失败');
  } finally {
    loading.value = false;
  }
};

// 下载文件
const downloadFile = async (file: FileDetail) => {
  try {
    const response = await DownloadFile({ file_ids: [file.id] });
    
    if (response.code === 0 && response.urls && response.urls.length > 0) {
      // 创建下载链接
      const link = document.createElement('a');
      link.href = response.urls[0];
      link.download = file.name;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      
      message.success('开始下载文件');
    } else {
      message.error(response.message || '下载失败');
    }
  } catch (error) {
    console.error('下载文件失败:', error);
    message.error('下载文件失败');
  }
};

// 删除文件
const deleteFile = async (file: FileDetail) => {
  // TODO: 实现删除文件接口
  message.warning('删除功能待实现');
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
.file-list-container {
  padding: 20px;
}

.file-item {
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.2s;
}

.file-item:hover {
  background-color: #f8f9fa;
}

.file-item:last-child {
  border-bottom: none;
}

.file-icon {
  color: #666;
}

.file-info {
  display: flex;
  gap: 16px;
  color: #999;
  font-size: 12px;
}

.file-size {
  min-width: 60px;
}

.file-time {
  min-width: 120px;
}

.file-actions {
  display: flex;
  align-items: center;
}

.loading-container {
  display: flex;
  justify-content: center;
  padding: 20px;
}

.empty-container {
  padding: 40px;
}
</style>

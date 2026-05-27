<template>
  <n-modal 
    :show="visible"
    preset="card" 
    style="width: 600px" 
    :mask-closable="false"
    :show-close="false"
    class="album-upload-modal"
    @update:show="onUpdateShow"
  >
    <template #header>
      <div class="flex items-center gap-2">
        <n-icon size="20" class="text-primary">
          <i class="icon-album"></i>
        </n-icon>
        <span class="text-lg font-semibold">上传专辑</span>
      </div>
    </template>

    <n-tabs v-model:value="mode" type="line" animated>
      <n-tab-pane name="search" tab="搜索已有专辑">
        <div class="space-y-4">
          <n-input-group>
            <n-input 
              v-model:value="searchKeyword"
              placeholder="输入专辑名称搜索..."
              @keyup.enter="searchAlbums"
              clearable
            />
            <n-button type="primary" @click="searchAlbums" :loading="searching">
              <template #icon>
                <n-icon><i class="icon-search"></i></n-icon>
              </template>
              搜索
            </n-button>
          </n-input-group>

          <div v-if="searchResults.length > 0" class="space-y-2">
            <n-card 
              v-for="(album, index) in searchResults" 
              :key="album.id"
              class="cursor-pointer hover:shadow-md transition-all duration-300 hover:scale-105"
              @click="selectAlbum(album)"
              :style="{ animationDelay: `${index * 0.1}s` }"
            >
              <div class="flex items-center gap-3">
                <n-image
                  :src="album.cover_file_url"
                  :alt="album.name"
                  width="48"
                  height="48"
                  object-fit="cover"
                  class="rounded shadow-sm"
                  :preview-disabled="true"
                />
                <div class="flex-1">
                  <div class="font-medium text-base">{{ album.name }}</div>
                  <div class="text-sm text-gray-500 mt-1">点击选择此专辑</div>
                </div>
                <n-icon size="20" class="text-gray-400">
                  <i class="icon-arrow-right"></i>
                </n-icon>
              </div>
            </n-card>
          </div>

          <n-empty v-else-if="hasSearched" description="未找到相关专辑">
            <template #icon>
              <n-icon size="48" class="text-gray-400">
                <i class="icon-search"></i>
              </n-icon>
            </template>
          </n-empty>
        </div>
      </n-tab-pane>

      <n-tab-pane name="create" tab="创建新专辑">
        <n-form 
          ref="formRef"
          :model="formData"
          :rules="rules"
          label-placement="left"
          label-width="80"
          require-mark-placement="right-hanging"
        >
          <n-form-item label="专辑名称" path="name">
            <n-input 
              v-model:value="formData.name"
              placeholder="请输入专辑名称"
              clearable
              maxlength="50"
              show-count
            />
          </n-form-item>

          <n-form-item label="发行时间" path="issue_time">
            <n-input-number 
              v-model:value="formData.issue_time"
              placeholder="请输入发行年份"
              :min="1900"
              :max="new Date().getFullYear()"
              clearable
            />
          </n-form-item>

          <n-form-item label="专辑封面" path="cover_file_id">
            <div class="flex items-center gap-4">
              <n-image
                v-if="coverPreview"
                :src="coverPreview"
                width="80"
                height="80"
                object-fit="cover"
                class="rounded shadow-md"
                :preview-disabled="true"
              />
              <div 
                v-else
                class="w-20 h-20 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded flex items-center justify-center hover:border-primary transition-colors"
              >
                <n-icon size="24" class="text-gray-400">
                  <i class="icon-image"></i>
                </n-icon>
              </div>
              <n-button 
                type="primary"
                @click="uploadCover"
                :loading="uploading"
                ghost
              >
                <template #icon>
                  <n-icon><i class="icon-upload"></i></n-icon>
                </template>
                上传封面
              </n-button>
            </div>
          </n-form-item>
        </n-form>
      </n-tab-pane>
    </n-tabs>

    <template #footer>
      <div class="flex justify-end gap-2">
        <n-button @click="close" ghost>
          取消
        </n-button>
        <n-button 
          v-if="mode === 'create'"
          type="primary"
          @click="createAlbum"
          :disabled="!canCreate"
          :loading="creating"
        >
          创建专辑
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { 
  NModal, 
  NTabs, 
  NTabPane, 
  NInput, 
  NInputGroup, 
  NButton, 
  NIcon, 
  NImage, 
  NCard, 
  NEmpty, 
  NForm, 
  NFormItem, 
  NInputNumber,
  useMessage
} from 'naive-ui';
import { UploadAlbum, GetAlbumDetail, SearchAlbum } from '@/api/album';
import { GetUploadToken } from '@/api/file';
import { AlbumProfile, UploadAlbumReq } from '@/api/types/album';

interface Props {
  visible: boolean;
  artistId: number;
  musicId: number;
}

interface Emits {
  (e: 'close'): void;
  (e: 'success', albumId: number): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();
const message = useMessage();

const onUpdateShow = (val: boolean) => {
  if (!val) emit('close');
};

const mode = ref<'search' | 'create'>('search');
const searchKeyword = ref('');
const searchResults = ref<AlbumProfile[]>([]);
const hasSearched = ref(false);
const coverPreview = ref('');
const searching = ref(false);
const uploading = ref(false);
const creating = ref(false);
const formRef = ref();

const formData = ref({
  name: '',
  issue_time: new Date().getFullYear(),
  cover_file_id: 0
});

const rules = {
  name: [
    { required: true, message: '请输入专辑名称', trigger: 'blur' }
  ],
  issue_time: [
    { required: true, message: '请输入发行时间', trigger: 'blur' }
  ],
  cover_file_id: [
    { required: true, message: '请上传专辑封面', trigger: 'change' }
  ]
};

const canCreate = computed(() => {
  return formData.value.name.trim() && formData.value.cover_file_id > 0;
});

const close = () => {
  emit('close');
  resetForm();
};

const resetForm = () => {
  mode.value = 'search';
  searchKeyword.value = '';
  searchResults.value = [];
  hasSearched.value = false;
  coverPreview.value = '';
  formData.value = {
    name: '',
    issue_time: new Date().getFullYear(),
    cover_file_id: 0
  };
};

const searchAlbums = async () => {
  if (!searchKeyword.value.trim()) {
    message.warning('请输入搜索关键词');
    return;
  }
  
  searching.value = true;
  try {
    const response = await SearchAlbum({
      keyword: searchKeyword.value,
      artist_id: props.artistId,
      page: 1,
      limit: 20
    });
    searchResults.value = response.albums;
    hasSearched.value = true;
    
    if (response.albums.length === 0) {
      message.info('未找到相关专辑');
    }
  } catch (error) {
    console.error('搜索专辑失败:', error);
    message.error('搜索专辑失败，请稍后重试');
  } finally {
    searching.value = false;
  }
};

const selectAlbum = async (album: AlbumProfile) => {
  try {
    message.success(`已选择专辑：${album.name}`);
    emit('success', album.id);
    close();
  } catch (error) {
    console.error('关联专辑失败:', error);
    message.error('关联专辑失败，请稍后重试');
  }
};

const uploadCover = async () => {
  const fileInput = document.createElement('input');
  fileInput.type = 'file';
  fileInput.accept = 'image/*';
  fileInput.onchange = async (e) => {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (file) {
      uploading.value = true;
      try {
        const tokenResponse = await GetUploadToken();
        const formData = new FormData();
        formData.append('file', file);
        formData.append('token', tokenResponse.upload_token);
        
        const uploadUrl =
          tokenResponse.upload_url ??
          tokenResponse.uploadUrl ??
          'https://up-z0.qiniup.com'
        const response = await fetch(uploadUrl, {
          method: 'POST',
          body: formData
        });
        
        if (response.ok) {
          const result = await response.json();
          formData.value.cover_file_id = result.key;
          coverPreview.value = `${tokenResponse.domain}/${result.key}`;
          message.success('封面上传成功');
        } else {
          throw new Error('上传失败');
        }
      } catch (error) {
        console.error('上传文件失败:', error);
        message.error('上传文件失败，请稍后重试');
      } finally {
        uploading.value = false;
      }
    }
  };
  fileInput.click();
};

const createAlbum = async () => {
  creating.value = true;
  try {
    const response = await UploadAlbum({
      name: formData.value.name,
      artist_id: props.artistId,
      cover_file_id: formData.value.cover_file_id,
      issue_time: formData.value.issue_time,
      music_id: props.musicId
    });
    
    message.success('专辑创建成功');
    emit('success', response.album_id);
    close();
  } catch (error) {
    console.error('创建专辑失败:', error);
    message.error('创建专辑失败，请稍后重试');
  } finally {
    creating.value = false;
  }
};
</script>

<style scoped>
.album-upload-modal :deep(.n-card) {
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style> 
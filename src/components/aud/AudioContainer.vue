<template>
  <div
      ref="audioContainer"
      class="audio-container fixed bg-[rgb(245,245,245)] rounded-lg shadow-lg flex items-center bottom-5 left-5 w-[400px] h-[60px]"
  >

    <!-- 左侧两竖排小点 -->
    <div
        class="drag-handle flex justify-between p-2 cursor-grab"
        @mousedown="startDragging"
    >
      <!-- 第一排小点 -->
      <div class="flex flex-col">
        <div class="dot w-[4px] h-[4px] bg-gray-500 rounded-full mb-1" v-for="n in 3" :key="'left' + n"></div>
      </div>
      <!-- 第二排小点 -->
      <div class="flex flex-col ml-1">
        <div class="dot w-[4px] h-[4px] bg-gray-500 rounded-full mb-1" v-for="n in 3" :key="'right' + n"></div>
      </div>
    </div>

    <!-- 音乐播放器主体，显示当前播放的音乐 -->
    <div class="audio-content flex-grow pl-5 group relative">
      <div v-if="currentMusicDetail" class="flex items-center">
        <div class="hidden group-hover:flex justify-center items-center">
          <button @click="playPrevious" class="w-[30px] h-[30px] rounded-full bg-gray-300 mr-2">⬅️</button>
          <button @click="togglePlay" class="w-[30px] h-[30px] rounded-full bg-gray-300 mx-2">▶️</button>
          <button @click="playNext" class="w-[30px] h-[30px] rounded-full bg-gray-300 ml-2">➡️</button>
        </div>
        <strong class="transition-all duration-300 group-hover:hidden">{{ currentMusicDetail.title }}</strong> <span
          class="group-hover:hidden">{{ currentMusicDetail.artist_name }}</span>
      </div>
      <div v-else>
        <p>请选择一首音乐播放</p>
      </div>
    </div>


    <!-- 左侧歌曲详情图标 -->
    <div class="song-detail-toggle p-2 cursor-pointer" @click="toggleSongDetail">
      📝 <!-- 可以替换为你想要的图标 -->
    </div>
    <!-- 右侧展示音乐列表的图标 -->
    <div class="music-list-toggle p-2 cursor-pointer" @click="toggleMusicList">
      🎵 <!-- 可以替换为你想要的图标 -->
    </div>

    <!-- 音乐列表展示 -->
    <div v-if="showMusicList && !showSongDetail"
         class="music-list absolute bg-white shadow-lg rounded-lg p-4 w-[400px] bottom-[70px] left-0">
      <n-list hoverable clickable>
        <div
            v-for="music in musicStore.MusicList"
            :key="music.id"
            class="mb-2 cursor-pointer"
            @dblclick="playMusic(music.id)">
          <n-list-item>
            <n-thing :title="music.title" content-style="margin-top: 10px;">
              <template #description>
                <n-space size="small" style="margin-top: 4px">
                  <n-tag :bordered="false" type="info" size="small">
                    暑夜
                  </n-tag>
                  <n-tag :bordered="false" type="info" size="small">
                    晚春
                  </n-tag>
                </n-space>
              </template>
            </n-thing>
          </n-list-item>
        </div>
      </n-list>
    </div>

    <!-- 歌曲详情展示 -->


    <div v-if="showSongDetail && currentMusicDetail" class="song-detail fixed bg-[#00000000] rounded-lg p-4 w-[400px]">
      <input type="hidden" v-model="currentMusicDetail.id">
      <input type="hidden" v-model="currentMusicDetail.root_id">
      <n-space vertical>
        <n-form>
          <n-form-item label="歌手" :show-label="true">
            <n-select
                v-model:value="selectedSingerValues"
                multiple
                filterable
                placeholder="搜索歌手"
                :options="singerOptions"
                :loading="loadingSinger"
                clearable
                remote
                :clear-filter-after-select="false"
                @search="handleSearchSinger"
            />
          </n-form-item>


          <n-form-item label="作曲" :show-label="true">
            <n-select
                v-model:value="selectedComposerValues"
                filterable
                placeholder="搜索作曲"
                :options="composerOptions"
                :loading="loadingComposer"
                clearable
                remote
                @search="handleSearchComposer"
            />
          </n-form-item>

          <n-form-item label="作词" :show-label="true">
            <n-select
                v-model:value="selectedWriterValues"
                filterable
                placeholder="搜索作词"
                :options="writerOptions"
                :loading="loadingWriter"
                clearable
                remote
                @search="handleSearchWriter"
            />
          </n-form-item>

          <n-form-item label="发布年份" path="issue_year">
            <n-input v-model:value="musicForm.issue_year" placeholder="Input Name"/>
          </n-form-item>

          <n-form-item label="地区" path="language">
            <n-input v-model:value="musicForm.language" placeholder="Input Name"/>
          </n-form-item>

          <n-form-item label="表演形式" path="perform_type">
            <n-select v-model:value="musicForm.perform_type" :options="performTypeOptions"/>
          </n-form-item>

          <n-form-item label="演唱会" path="concert">
            <n-input v-model:value="musicForm.concert" placeholder="Input Name"/>
          </n-form-item>

          <n-form-item label="表演时间" path="concert_year">
            <n-input v-model:value="musicForm.concert_year" placeholder="Input Name"/>
          </n-form-item>

          <n-form-item label="MV链接" path="concert_year">
            <n-input v-model:value="musicForm.mv_url" placeholder="Input Name"/>
          </n-form-item>

          <n-upload :custom-request="uploadMp3">
            <n-button>上传MP3</n-button>
          </n-upload>
          <n-upload :custom-request="uploadLyric">
            <n-button>上传歌词</n-button>
          </n-upload>
          <n-upload :custom-request="uploadCover">
            <n-button>上传封面</n-button>
          </n-upload>
          <n-upload :custom-request="uploadSheet">
            <n-button>上传曲谱</n-button>
          </n-upload>


        </n-form>
        <n-button @click="submitMusicForm">确定</n-button>

      </n-space>

    </div>


  </div>
</template>


<script lang="ts" setup>

import {onMounted, Ref, ref} from 'vue';
import apiClient from '@/util/http.ts'; // 引入 axios 实例
import type {MusicDetail} from "@/store/music.ts";
import {convertToUpsertMusicReq, useMusicStore} from "@/store/music.ts";
import type {SelectOption, UploadCustomRequestOptions} from 'naive-ui'
import {
  NButton,
  NForm,
  NFormItem,
  NInput,
  NList,
  NListItem,
  NSelect,
  NSpace,
  NTag,
  NThing,
  NUpload,
  useMessage
} from "naive-ui"
import {performTypeOptions} from "@/util/enum.ts";

const message = useMessage()
// *************************************************** 音乐列表展示逻辑 *************************************************** //
const showMusicList = ref(false);
const showSongDetail = ref(false);
const currentMusicDetail = ref<MusicDetail>(); // 用于保存当前音乐详情
const musicStore = useMusicStore(); // 获取 Pinia store


const listMusicReq = {
  "title": "",
  "artist": 0,
  "perform_type": 0,
  "page": 1,
  "size": 5
};

// *************** music表单相关控件 *************** //
const selectedSingerValues = ref(null)
const loadingSinger = ref(false)
const singerOptions = ref<SelectOption[]>([])
const selectedComposerValues = ref(null)
const loadingComposer = ref(false)
const composerOptions = ref<SelectOption[]>([])
const selectedWriterValues = ref(null)
const loadingWriter = ref(false)
const writerOptions = ref<SelectOption[]>([])

const musicForm = ref<MusicForm | null>({
  issue_year: "",
  language: "",
  perform_type: "",
  concert: "",
  concert_year: "",
  mv_url: "",
  mp3: null,
  lyric: null,
  sheet: null,
  cover: null,
})

interface MusicForm {
  issue_year: string,
  language: string,
  perform_type: string,
  concert: string,
  concert_year: string,
  mv_url: string,
  mp3: File | null,
  lyric: File | null,
  sheet: File | null,
  cover: File | null,
}

interface SingerModel {
  id: number;
  name: string;
}

const handleSearchSinger = async (query: string) => {
  await handleSearch(query, 1, singerOptions, loadingSinger); // 传入typ为1
};

const handleSearchComposer = async (query: string) => {
  await handleSearch(query, 2, composerOptions, loadingComposer); // 传入typ为2
};

const handleSearchWriter = async (query: string) => {
  await handleSearch(query, 3, writerOptions, loadingWriter); // 传入typ为2
};

const uploadMp3 = ({file}: UploadCustomRequestOptions) => {
  if (musicForm.value) {
    if (file.file) {
      musicForm.value.mp3 = file.file; // 确保 file.file 不是 null
      console.log("上传mp3文件暂存: ", file.name)
    } else {
      message.error("文件无效");
      return;
    }
  } else {
    message.error("请先选择歌曲");
    return;
  }
  message.info("上传mp3成功");
}
const uploadLyric = ({file}: UploadCustomRequestOptions) => {
  if (musicForm.value) {
    if (file.file) {
      musicForm.value.lyric = file.file; // 确保 file.file 不是 null
      console.log("上传歌词文件暂存: ", file.name)
    } else {
      message.error("文件无效");
      return;
    }
  } else {
    message.error("请先选择歌曲");
    return;
  }
  message.info("上传歌词成功");
}
const uploadCover = ({file}: UploadCustomRequestOptions) => {
  if (musicForm.value) {
    if (file.file) {
      musicForm.value.cover = file.file; // 确保 file.file 不是 null
      console.log("上传封面文件暂存: ", file.name)
    } else {
      message.error("文件无效");
      return;
    }
  } else {
    message.error("请先选择歌曲");
    return;
  }
}
const uploadSheet = ({file}: UploadCustomRequestOptions) => {
  if (musicForm.value) {
    if (file.file) {
      musicForm.value.sheet = file.file; // 确保 file.file 不是 null
      console.log("上传曲谱文件暂存: ", file.name)
    } else {
      message.error("文件无效");
      return;
    }
  } else {
    message.error("请先选择歌曲");
    return;
  }
}

const submitMusicForm = async () => {
  const response = await apiClient.post('/music/upsert', musicForm.value);
  if (response.status === 200) {
    message.success("上传成功");
  } else {
    message.error("上传失败");
  }
};

const handleSearch = async (query: string, typ: number, options: Ref<SelectOption[]>, loading: Ref<boolean>) => {
  if (!query.length) {
    return
  }
  loading.value = true
  // 搜索singer
  try {
    const response = await apiClient.get('/performer/get', {
      params: {
        type: typ,
        name: query,
      }
    });
    const singerList = response.data.data;
    options.value = singerList.map((s: SingerModel) => ({
      label: s.name,
      value: s.id
    }));
    loading.value = false
  } catch (error) {
    console.error('搜索失败', error);
  } finally {
    loading.value = false; // 无论成功与否，都需要结束加载状态
  }
};

// 处理文件上传
const handleFileUpload = (field, event) => {
  const file = event.target.files[0];
  if (file && currentMusicDetail.value) {
    currentMusicDetail.value[field] = file;
  }
};

const fetchMusicList = async () => {
  const response = await apiClient.post('/music/list', listMusicReq);
  const musicList = response.data.data;
  musicStore.setMusicList(musicList);
};

// 播放音乐，获取音乐详情
const playMusic = async (id: number) => {
  try {
    const response = await apiClient.get(`/music/get/${id}`); // 获取音乐详情的 API
    currentMusicDetail.value = response.data.data; // 更新当前音乐详情
    musicStore.setCurrentMusic(id); // 更新 Pinia store 中的当前播放的音乐 id
  } catch (error) {
    console.error('获取音乐详情失败', error);
  }
};

// *************************************************** 拖动操作 *************************************************** //
const audioContainer = ref<HTMLDivElement | null>(null);
let isDragging = false;
let offset = {x: 0, y: 0};

// 切换音乐列表展示
const toggleMusicList = () => {
  showMusicList.value = !showMusicList.value;
  showSongDetail.value = false; // 隐藏歌曲详情窗口
};

// 切换歌曲详情展示
const toggleSongDetail = () => {
  showSongDetail.value = !showSongDetail.value;
  showMusicList.value = false; // 隐藏音乐列表
};

// 使用转换函数保存音乐详情
const saveMusicDetails = async () => {
  try {
    if (currentMusicDetail.value) {
      const upsertMusicReq = convertToUpsertMusicReq(currentMusicDetail.value);
      console.log(upsertMusicReq);
      // 发送更新音乐详情的请求
      const response = await apiClient.post(`/music/upsert`, upsertMusicReq, {
        headers: {
          'Content-Type': 'multipart/form-data', // 指定请求类型为 multipart/form-data
        },
      });
      console.log('音乐详情已更新', response.data);
      // 可以选择刷新音乐列表或其他操作
    }
  } catch (error) {
    console.error('保存音乐详情失败', error);
  }
};

// 播放上一曲
const playPrevious = () => {
  // 实现上一曲的逻辑
};

// 暂停或播放音乐
const togglePlay = () => {
  // 实现暂停/播放的逻辑
};

// 播放下一曲
const playNext = () => {
  // 实现下一曲的逻辑
};

// 获取音乐列表
onMounted(() => {
  fetchMusicList();
});


const startDragging = (e: MouseEvent) => {
  if (!audioContainer.value) return;

  isDragging = true;
  offset.x = e.clientX - audioContainer.value.getBoundingClientRect().left;
  offset.y = e.clientY - audioContainer.value.getBoundingClientRect().top;

  // 禁用文本选择
  document.body.style.userSelect = 'none';

  document.addEventListener('mousemove', drag);
  document.addEventListener('mouseup', stopDragging);
};

const stopDragging = () => {
  isDragging = false;

  // 恢复文本选择
  document.body.style.userSelect = '';

  document.removeEventListener('mousemove', drag);
  document.removeEventListener('mouseup', stopDragging);
};

const drag = (e: MouseEvent) => {
  if (!isDragging || !audioContainer.value) return;

  // 计算新的位置
  let newLeft = e.clientX - offset.x;
  let newTop = e.clientY - offset.y;

  // 获取窗口宽高
  const windowWidth = window.innerWidth;
  const windowHeight = window.innerHeight;

  // 限制新位置不能超出窗口范围
  const containerWidth = audioContainer.value.offsetWidth;
  const containerHeight = audioContainer.value.offsetHeight;

  // 限制左、右、上、下的最大最小值
  if (newLeft < 0) newLeft = 0;
  if (newLeft + containerWidth > windowWidth) newLeft = windowWidth - containerWidth;
  if (newTop < 0) newTop = 0;
  if (newTop + containerHeight > windowHeight) newTop = windowHeight - containerHeight;

  // 设置新位置
  audioContainer.value.style.left = `${newLeft}px`;
  audioContainer.value.style.top = `${newTop}px`;
};


</script>


<style scoped>
.song-detail {
  /* 设置为fixed，使其相对于屏幕定位 */
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1000; /* 确保位于最前面 */
  width: 400px;
}
</style>
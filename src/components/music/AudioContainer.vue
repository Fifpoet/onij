<template>
  <div ref="audioContainer"
    class="audio-container fixed bg-[rgb(245,245,245)] flex items-center bottom-5 left-5 w-[300px] h-[60px]">

    <!-- 拖拽手柄 -->
    <div class="drag-handle flex justify-between p-2 cursor-grab" @mousedown="startDragging">
      <div class="flex flex-col">
        <div class="dot w-[3px] h-[3px] bg-gray-500 rounded-full mb-1" v-for="n in 4" :key="'left' + n"></div>
      </div>
      <div class="flex flex-col ml-1">
        <div class="dot w-[3px] h-[3px] bg-gray-500 rounded-full mb-1" v-for="n in 4" :key="'right' + n"></div>
      </div>
    </div>

    <!-- 音乐播放器主体，显示当前播放的音乐 -->
    <div class="audio-content flex-grow pl-5 group relative">
      <div v-if="musicStore.current.detail" class="flex items-center">
        <div class="hidden group-hover:flex justify-center items-center">
          <button @click="playPrevious" class="w-[30px] h-[30px] rounded-full bg-gray-300 mr-2">⬅️</button>
          <button @click="startOrPause" class="w-[30px] h-[30px] rounded-full bg-gray-300 mx-2">
            {{ isPlaying ? '⏸️' : '▶️' }}
          </button>
          <button @click="playNext" class="w-[30px] h-[30px] rounded-full bg-gray-300 ml-2">➡️</button>
        </div>
        <span class="transition-all duration-300 group-hover:hidden">{{ musicStore.currentMusicName }} - {{
          musicStore.currentMusicArtistName }}</span>
      </div>
      <div v-else>
        <p>请选择一首音乐播放</p>
      </div>
    </div>

    <!-- 音量按钮 -->
    <div class="volume-control p-2 cursor-pointer relative" @click="toggleMute" @mouseenter="showVolumeSlider = true"
      @mouseleave="showVolumeSlider = false">
      {{ isMuted ? '🔇' : volume > 50 ? '🔊' : volume > 0 ? '🔉' : '🔇' }}

      <!-- 音量滑块 -->
      <div v-if="showVolumeSlider"
        class="volume-slider absolute bottom-[40px] left-0 bg-white shadow-lg rounded-lg p-2 flex flex-col items-center">
        <input type="range" class="h-[80px] w-[20px] bg-gray-300 outline-none appearance-none" min="0" max="100"
          v-model="volume" @input="adjustVolume" orient="vertical" />
      </div>
    </div>

    <div class="music-list-toggle p-2 cursor-pointer" @click="toggleMusicList">
      🎵
    </div>
    <!-- 进度条 -->
    <div class="absolute bottom-[-1px] left-0 w-full p-0 m-0 flex items-center">
      <span class="text-xs text-gray-500 ml-2">{{ formattedCurrentTime }}</span>
      <input type="range" class="flex-grow h-[2px] bg-gray-300 outline-none appearance-none p-0 m-0 ml-2 mr-2" min="0"
        max="100" v-model="currentProgress" @input="onPullTime" />
    </div>

    <!-- music列表 -->
    <MusicList v-if="showMusicList && !showSongDetail" @play-music="playMusic" />

    <!-- 歌曲详情展示 -->
    <div v-if="showSongDetail && currentMusicDetail" class="song-detail fixed bg-[#00000000] rounded-lg p-4 w-[400px]">
      <input type="hidden" v-model="currentMusicDetail.id">
      <input type="hidden" v-model="currentMusicDetail.root_id">
      <n-space vertical>
        <n-form>
          <n-form-item label="歌手" :show-label="true">
            <n-select v-model:value="selectedSingerValues" multiple filterable placeholder="搜索歌手"
              :options="singerOptions" :loading="loadingSinger" clearable remote :clear-filter-after-select="false"
              @search="handleSearchSinger" />
          </n-form-item>


          <n-form-item label="作曲" :show-label="true">
            <n-select v-model:value="selectedComposerValues" filterable placeholder="搜索作曲" :options="composerOptions"
              :loading="loadingComposer" clearable remote @search="handleSearchComposer" />
          </n-form-item>

          <n-form-item label="作词" :show-label="true">
            <n-select v-model:value="selectedWriterValues" filterable placeholder="搜索作词" :options="writerOptions"
              :loading="loadingWriter" clearable remote @search="handleSearchWriter" />
          </n-form-item>

          <n-form-item label="发布年份" path="issue_year">
            <n-input v-model:value="currentMusicDetail.issue_year" placeholder="Input Name" />
          </n-form-item>

          <n-form-item label="地区" path="language">
            <n-select v-model:value="currentMusicDetail.language" :options="languageOptions" />
          </n-form-item>

          <n-form-item label="表演形式" path="perform_type">
            <n-select v-model:value="currentMusicDetail.perform_type" :options="performTypeOptions" />
          </n-form-item>

          <n-form-item label="演唱会" path="concert">
            <n-input v-model:value="currentMusicDetail.concert" placeholder="Input Name" />
          </n-form-item>

          <n-form-item label="表演时间" path="concert_year">
            <n-input v-model:value="currentMusicDetail.concert_year" placeholder="Input Name" />
          </n-form-item>

          <n-form-item label="MV链接" path="concert_year">
            <n-input v-model:value="currentMusicDetail.mv_url" placeholder="Input Name" />
          </n-form-item>


        </n-form>
        <n-button @click="submitMusicForm">确定</n-button>

      </n-space>

    </div>


  </div>
</template>


<script lang="ts" setup>

import { computed, onMounted, ref } from 'vue';
import MusicList from './MusicList.vue';
import { useMusicStore } from "@/store/music.ts";
import type { SelectOption } from 'naive-ui'
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
import { performTypeOptions } from "@/util/enum.ts";
import { GetMusicDetail, GetMusicList } from '@/api';
import { GetMusicListReq, MusicSortType, MusicDetail } from '@/api/types';
import { formatTime } from '@/util/time';

const message = useMessage()
// *************************************************** 音乐列表展示逻辑 *************************************************** //
const showMusicList = ref(false);
const showSongDetail = ref(false);
const musicStore = useMusicStore(); // 获取 Pinia store

// 音量控制相关
const volume = ref(80); // 音量值，范围0-100
const isMuted = ref(false); // 是否静音
const showVolumeSlider = ref(false); // 是否显示音量滑块
const previousVolume = ref(100); // 存储静音前的音量值

// 调整音量
const adjustVolume = () => {
  if (audio.value) {
    isMuted.value = volume.value == 0
    audio.value.volume = volume.value / 100;
  }
};

// 切换静音状态
const toggleMute = () => {
  if (audio.value) {
    if (isMuted.value) {
      // 取消静音
      isMuted.value = false;
      volume.value = previousVolume.value;
      audio.value.volume = volume.value / 100;
    } else {
      // 静音
      isMuted.value = true;
      previousVolume.value = volume.value;
      volume.value = 0;
      audio.value.volume = 0;
    }
  }
};

// *************************************************** 音乐播放逻辑 *************************************************** //
const audioContainer = ref<HTMLDivElement | null>(null);
const audio = ref<HTMLAudioElement | null>(null); // 用于音频控制的全局 Audio 实例
const isPlaying = ref(false); // 控制播放状态
const duration = ref(0); // 音频总时长
const currentTime = ref(0); // 添加一个响应式变量来跟踪当前时间


let isDragging = false;
let offset = { x: 0, y: 0 };

const formattedCurrentTime = computed(() => {
  if (!audio.value) return '00:00'; // TODO这里不会重新计算

  return formatTime(currentTime.value);
});

// 播放音乐，获取音乐详情
const playMusic = async (id: number) => {
  try {
    const response = await GetMusicDetail({ music_id: id });
    useMusicStore.setCurrentMusic(response.detail)
    handleMp3()
  } catch (error) {
    console.error('获取音乐详情失败', error);
  }
};

// 切换音乐列表展示
const toggleMusicList = () => {
  showMusicList.value = !showMusicList.value;
};

// 使用转换函数保存音乐详情

// 处理音频文件
const handleMp3 = () => {
  console.log("开始处理mp3并开始播放: ", musicStore.currentMusicName)
  if (!musicStore.current.detail) {
    message.error("没有找到音频文件");
    return;
  }
  // 先干掉之前的Audio对象, 初始化新Audio对象
  if (audio.value) {
    audio.value.pause();
    audio.value = null;
  }
  isPlaying.value = false;
  for (let i = 0; i < 3; i++) {
    if (!audio.value) {
      audio.value = new Audio(musicStore.current.detail?.mp3_file_url);
      audio.value.onerror = () => {
        audio.value = null;
        isPlaying.value = false;
      };
    } else {
      break;
    }
  }
  if (!audio.value) {
    message.error("音频文件加载失败，请检查链接");
    return;
  }

  // 设置音量
  audio.value.volume = volume.value / 100;

  audio.value.ontimeupdate = () => {
    currentProgress.value = (audio.value!.currentTime / audio.value!.duration) * 100;
    currentTime.value = audio.value!.currentTime;
  };
  audio.value.onloadedmetadata = () => {
    duration.value = audio.value!.duration;
  };

  audio.value.play();
  isPlaying.value = true;
};

// 开始播放或暂停
const startOrPause = () => {
  if (audio.value) {
    if (isPlaying.value) {
      audio.value.pause();
      isPlaying.value = false;
    } else {
      audio.value.play();
      isPlaying.value = true;
    }
  }
}
const playPrevious = () => {
  if (!musicStore.musicList || musicStore.musicList.length === 0) return;

  const currentIndex = musicStore.musicList.findIndex(music => music?.id === musicStore.current.detail?.id);
  if (currentIndex > 0) {
    playMusic(musicStore.musicList[currentIndex - 1]?.id);
  } else {
    playMusic(musicStore.musicList[musicStore.musicList.length - 1]?.id);
  }
};
const playNext = () => {
  if (!musicStore.musicList || musicStore.musicList.length === 0) return;

  const currentIndex = musicStore.musicList.findIndex(music => music?.id === musicStore.current.detail?.id);
  let nextIndex;
  if (musicStore.musicList.length > 1) {
    do {
      nextIndex = Math.floor(Math.random() * musicStore.musicList.length);
    } while (nextIndex === currentIndex);
  } else {
    nextIndex = 0;
  }
  playMusic(musicStore.musicList[nextIndex]?.id);
};

const currentProgress = computed({
  get: () => musicStore.current.progress,
  set: (value) => {
    musicStore.current.progress = value;
  }
});

const onPullTime = () => {
  if (!audio.value || !musicStore.current.detail) return;
  const time = (musicStore.current.progress / 100) * duration.value;
  audio.value.currentTime = time;
};

// 添加时间更新监听
onMounted(() => {
  if (audio.value) {
    audio.value.addEventListener('timeupdate', () => {
      if (!isDragging && audio.value) {
        musicStore.current.progress = (audio.value.currentTime / duration.value) * 100;
      }
    });
  }
});



// ************************* 播放器拖动 START ************************* //
const startDragging = (e: MouseEvent) => {
  if (!audioContainer.value) return;
  isDragging = true;
  offset.x = e.clientX - audioContainer.value.getBoundingClientRect().left;
  offset.y = e.clientY - audioContainer.value.getBoundingClientRect().top;

  document.body.style.userSelect = 'none';
  document.addEventListener('mousemove', drag);
  document.addEventListener('mouseup', stopDragging);
};

const stopDragging = () => {
  isDragging = false;

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


onMounted(async () => {

});

</script>


<style scoped>
.song-detail {
  /* 设置为fixed，使其相对于屏幕定位 */
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1000;
  /* 确保位于最前面 */
  width: 400px;
}

/* 调整进度条滑块圆圈的大小 */
input[type="range"]::-webkit-slider-thumb {
  width: 8px;
  /* 圆圈宽度 */
  height: 8px;
  /* 圆圈高度 */
  background-color: #333;
  /* 圆圈颜色 */
  border-radius: 50%;
  cursor: pointer;
  -webkit-appearance: none;
  /* 移除默认样式 */
}

input[type="range"]::-moz-range-thumb {
  width: 8px;
  height: 8px;
  background-color: #333;
  border-radius: 50%;
  cursor: pointer;
}

/* 垂直方向的音量滑块样式 */
input[type="range"][orient="vertical"] {
  -webkit-appearance: slider-vertical;
  writing-mode: bt-lr;
}

input[type="range"][orient="vertical"]::-webkit-slider-thumb {
  width: 12px;
  height: 12px;
}

input[type="range"][orient="vertical"]::-moz-range-thumb {
  width: 12px;
  height: 12px;
}
</style>

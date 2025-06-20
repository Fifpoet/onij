<template>
  <div ref="audioContainer"
    class="audio-container fixed bg-[rgb(245,245,245)] flex items-center 
    // 移动端样式：底部固定，宽度100%，左对齐
    bottom-0 left-0 right-0 w-full h-[80px]
    // 桌面端还原原样式
    md:bottom-5 md:left-5 md:right-auto md:w-[300px] md:h-[60px]">

    <!-- 拖拽手柄 - 移动端隐藏 -->
    <div class="hidden md:flex drag-handle justify-between p-2 cursor-grab" @mousedown="startDragging">
      <div class="flex flex-col">
        <div class="dot w-[3px] h-[3px] bg-gray-500 rounded-full mb-1" v-for="n in 4" :key="'left' + n"></div>
      </div>
      <div class="flex flex-col ml-1">
        <div class="dot w-[3px] h-[3px] bg-gray-500 rounded-full mb-1" v-for="n in 4" :key="'right' + n"></div>
      </div>
    </div>

    <!-- 音乐播放器主体 -->
    <div class="audio-content flex-grow px-2 md:px-5 group relative">
      <div v-if="musicStore.current.detail" class="flex items-center justify-between w-full">
        <!-- 左侧：播放控制和歌曲信息 -->
        <div class="flex items-center flex-1 min-w-0 mr-2">
          <!-- 控制按钮 - 移动端常驻显示，桌面端hover显示 -->
          <div class="flex md:hidden justify-center items-center shrink-0">
            <NButton text @click="playPrevious" class="w-[30px] h-[30px] rounded-full mr-1">
              <NIcon :component="PlaySkipBackOutline" size="18" />
            </NButton>
            <NButton text @click="startOrPause" class="w-[30px] h-[30px] rounded-full mx-1">
              <NIcon :component="isPlaying ? PauseCircleOutline : PlayCircleOutline" size="18" />
            </NButton>
            <NButton text @click="playNext" class="w-[30px] h-[30px] rounded-full ml-1">
              <NIcon :component="PlaySkipForwardOutline" size="18" />
            </NButton>
          </div>
          <div class="hidden md:group-hover:flex justify-center items-center shrink-0">
            <NButton text @click="playPrevious" class="w-[30px] h-[30px] rounded-full mr-1">
              <NIcon :component="PlaySkipBackOutline" size="18" />
            </NButton>
            <NButton text @click="startOrPause" class="w-[30px] h-[30px] rounded-full mx-1">
              <NIcon :component="isPlaying ? PauseCircleOutline : PlayCircleOutline" size="18" />
            </NButton>
            <NButton text @click="playNext" class="w-[30px] h-[30px] rounded-full ml-1">
              <NIcon :component="PlaySkipForwardOutline" size="18" />
            </NButton>
          </div>
          <!-- 歌曲信息 - 移动端常驻显示，桌面端hover隐藏 -->
          <span class="truncate md:transition-all md:duration-300 md:group-hover:hidden ml-2">
            {{ musicStore.currentMusicName }} - {{ musicStore.currentMusicArtistName }}
          </span>
        </div>

        <!-- 右侧：音量和操作按钮 -->
        <div class="flex items-center shrink-0">
          <!-- 音量按钮 - 移动端隐藏 -->
          <div class="hidden md:block volume-control cursor-pointer relative" @mouseenter="showVolumeSlider = true"
            @mouseleave="showVolumeSlider = false">
            <NButton text @click="toggleMute">
              <NIcon
                :component="isMuted ? VolumeMuteOutline : (volume > 50 ? VolumeHighOutline : (volume > 0 ? VolumeLowOutline : VolumeMuteOutline))"
                size="18" />
            </NButton>

            <!-- 音量滑块 -->
            <div v-if="showVolumeSlider"
              class="volume-slider absolute bottom-[40px] left-0 bg-white shadow-lg rounded-lg p-2 flex flex-col items-center">
              <input type="range" class="h-[80px] w-[20px] bg-gray-300 outline-none appearance-none" min="0" max="100"
                v-model="volume" @input="adjustVolume" orient="vertical" />
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="like-toggle cursor-pointer px-2" @click="toggleMusicLike">
            <NIcon :component="HeartOutline" size="18" />
          </div>
          <div class="music-list-toggle cursor-pointer px-2" @click="toggleMusicList">
            <NIcon :component="MusicalNoteOutline" size="18" />
          </div>
        </div>
      </div>
      <div v-else>
        <p>请选择一首音乐播放</p>
      </div>
    </div>

    <!-- 进度条 -->
    <div class="absolute bottom-[-1px] left-0 w-full p-0 m-0 flex items-center">
      <span class="text-xs text-gray-500 ml-2">{{ formattedCurrentTime }}</span>
      <input type="range" class="flex-grow h-[2px] bg-gray-300 outline-none appearance-none p-0 m-0 ml-2 mr-2" min="0"
        max="100" v-model="currentProgress" @input="onPullTime" />
    </div>

    <!-- music列表 -->
    <MusicListComponent v-if="showMusicList && !showSongDetail" @play-music="playMusic" />
  </div>
</template>


<script lang="ts" setup>

import { computed, onMounted, ref } from 'vue';
import MusicListComponent from './MusicList.vue';
import { useMusicStore } from "@/store/music.ts";
import {
  NButton,
  useMessage,
  NIcon // 引入 NIcon
} from "naive-ui"
import {
  PlayCircleOutline,
  PauseCircleOutline,
  PlaySkipBackOutline,
  PlaySkipForwardOutline,
  VolumeHighOutline,
  VolumeLowOutline,
  VolumeMuteOutline,
  HeartOutline,
  MusicalNoteOutline
} from '@vicons/ionicons5'; // 引入图标
import { GetMusicDetail } from '@/api';
import { MidShowWhat } from '@/api/types';
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
    // 直接判断返回的detail对象
    if (response?.detail) {
      // 切换新歌曲时重置进度
      musicStore.current.progress = 0;
      musicStore.setCurrentMusic(response.detail)
      musicStore.setMidShowWhat(MidShowWhat.ShowMusicDetail)
      handleMp3()
    } else {
      message.error('获取音乐详情失败');
    }
  } catch (error) {
    console.error('获取音乐详情失败', error);
    message.error('获取音乐详情失败');
  }
};

// 切换音乐列表展示
const toggleMusicList = () => {
  showMusicList.value = !showMusicList.value;
};

const toggleMusicLike = () => {
  if (musicStore.midShowWhat == MidShowWhat.ShowMusicLike) {
    musicStore.setMidShowWhat(MidShowWhat.ShowMusicDetail)
  } else {
    musicStore.setMidShowWhat(MidShowWhat.ShowMusicLike)
  }
};


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
    if (!audio.value) return;
    duration.value = audio.value.duration;
    // 页面刷新时会自动使用store中保存的进度
    if (musicStore.current.progress > 0) {
      audio.value.currentTime = (musicStore.current.progress / 100) * audio.value.duration;
    }
  };
  audio.value.onended = () => {
    isPlaying.value = false;
    playNext();
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
  // 如果 store 中存在音乐详情，说明之前在播放，需要恢复播放状态
  if (musicStore.current.detail) {
    await handleMp3();
  }
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

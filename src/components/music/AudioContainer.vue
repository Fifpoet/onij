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
        <span class="transition-all duration-300 group-hover:hidden">{{ musicStore.currentMusicName }} - {{ musicStore.currentMusicArtistName }}</span>
      </div>
      <div v-else>
        <p>请选择一首音乐播放</p>
      </div>
    </div>

    <div class="song-detail-toggle p-2 cursor-pointer" @click="toggleSongDetail">
      📝
    </div>
    <div class="music-list-toggle p-2 cursor-pointer" @click="toggleMusicList">
      🎵
    </div>
    <!-- 进度条 -->
    <div class="absolute bottom-[-1px] left-0 w-full p-0 m-0 flex items-center">
      <span class="text-xs text-gray-500 ml-2">{{ formattedCurrentTime }}</span>
      <input type="range" class="flex-grow h-[2px] bg-gray-300 outline-none appearance-none p-0 m-0 ml-2 mr-2" min="0"
             max="100" v-model="currentProgress" @input="onSeek"/>
    </div>

    <!-- music列表 -->
    <div v-if="showMusicList && !showSongDetail" @scroll="handleScroll"
         class="music-list absolute bg-white shadow-lg rounded-lg p-4 w-[400px] bottom-[70px] left-0 max-h-[300px] overflow-y-auto">
      <n-list hoverable clickable>
        <div v-for="music in musicStore.musicList" :key="music?.id" class="mb-0.5 cursor-pointer"
             @dblclick="playMusic(music?.id)">
          <n-list-item>
            <n-thing :title="music?.name" content-style="margin-top: 10px;">
              <template #description>
                <n-space size="small" style="margin-top: 2px">
                  <n-tag :bordered="false" type="info" size="small">
                    暑夜
                  </n-tag>
                  <n-tag :bordered="false" type="info" size="small">
                    晚春
                  </n-tag>
                </n-space>
              </template>
              <template #header-extra>
                {{music?.artist }}
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
            <n-select v-model:value="selectedSingerValues" multiple filterable placeholder="搜索歌手"
                      :options="singerOptions"
                      :loading="loadingSinger" clearable remote :clear-filter-after-select="false"
                      @search="handleSearchSinger"/>
          </n-form-item>


          <n-form-item label="作曲" :show-label="true">
            <n-select v-model:value="selectedComposerValues" filterable placeholder="搜索作曲"
                      :options="composerOptions"
                      :loading="loadingComposer" clearable remote @search="handleSearchComposer"/>
          </n-form-item>

          <n-form-item label="作词" :show-label="true">
            <n-select v-model:value="selectedWriterValues" filterable placeholder="搜索作词" :options="writerOptions"
                      :loading="loadingWriter" clearable remote @search="handleSearchWriter"/>
          </n-form-item>

          <n-form-item label="发布年份" path="issue_year">
            <n-input v-model:value="currentMusicDetail.issue_year" placeholder="Input Name"/>
          </n-form-item>

          <n-form-item label="地区" path="language">
            <n-select v-model:value="currentMusicDetail.language" :options="languageOptions"/>
          </n-form-item>

          <n-form-item label="表演形式" path="perform_type">
            <n-select v-model:value="currentMusicDetail.perform_type" :options="performTypeOptions"/>
          </n-form-item>

          <n-form-item label="演唱会" path="concert">
            <n-input v-model:value="currentMusicDetail.concert" placeholder="Input Name"/>
          </n-form-item>

          <n-form-item label="表演时间" path="concert_year">
            <n-input v-model:value="currentMusicDetail.concert_year" placeholder="Input Name"/>
          </n-form-item>

          <n-form-item label="MV链接" path="concert_year">
            <n-input v-model:value="currentMusicDetail.mv_url" placeholder="Input Name"/>
          </n-form-item>


        </n-form>
        <n-button @click="submitMusicForm">确定</n-button>

      </n-space>

    </div>


  </div>
</template>


<script lang="ts" setup>

import {computed, onMounted, ref} from 'vue';
import { useMusicStore} from "@/store/music.ts";
import type {SelectOption} from 'naive-ui'
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
import {languageOptions, performTypeOptions} from "@/util/enum.ts";
import {GetMusicDetail, GetMusicList} from '@/api';
import {GetMusicListReq, MusicSortType, MusicDetail} from '@/api/types';

const message = useMessage()
// *************************************************** 音乐列表展示逻辑 *************************************************** //
const showMusicList = ref(false);
const showSongDetail = ref(false);
const musicStore = useMusicStore(); // 获取 Pinia store

// *************** music表单相关控件 *************** //
const selectedSingerValues = ref<number[]>([]);
const loadingSinger = ref(false)
const singerOptions = ref<SelectOption[]>([])
const selectedComposerValues = ref<number[]>([]);
const loadingComposer = ref(false)
const composerOptions = ref<SelectOption[]>([])
const selectedWriterValues = ref<number[]>([]);
const loadingWriter = ref(false)
const writerOptions = ref<SelectOption[]>([])
const fetching = ref(false); // 防止重复请求

const handleScroll = (event: Event) => {
  const target = event.target as HTMLElement;
  // 检查是否滚动到底部
  if (target.scrollHeight - target.scrollTop <= target.clientHeight + 10) {
    if (!fetching.value && musicStore.musicListContinue) {
      fetching.value = true; // 设置为正在加载状态
      fetchMusicList().finally(() => {
        fetching.value = false; // 重置加载状态
      });
    }
  }
};


const fetchMusicList = async () => {
  try {
    const response = await GetMusicList({
      sort_type: MusicSortType.MST_Created_At_Desc,
      page: musicStore.page,
      limit: 100
    });
    if (response?.musics) {
      musicStore.append(response.musics);
      musicStore.incrPage()
    } else {
      musicStore.stopPage()
    }
  } catch (error) {
    console.error("加载音乐列表失败:", error);
  }
};

// *************************************************** 音乐播放逻辑 *************************************************** //
const audioContainer = ref<HTMLDivElement | null>(null);
const audio = ref<HTMLAudioElement | null>(null); // 用于音频控制的全局 Audio 实例
const isPlaying = ref(false); // 控制播放状态
const currentTime = ref(0);
const duration = ref(0); // 音频总时长
const playMod = ref(2); // 播放模式


let isDragging = false;
let offset = {x: 0, y: 0};

const formattedCurrentTime = computed(() => {
  const minutes = Math.floor(currentTime.value / 60)
      .toString()
      .padStart(2, "0");
  const seconds = Math.floor(currentTime.value % 60)
      .toString()
      .padStart(2, "0");
  return `${minutes}:${seconds}`;
});

// 播放音乐，获取音乐详情
const playMusic = async (id: number) => {
  try {
    const response = await GetMusicDetail({ music_id: id }); 
    musicStore.setCurrentMusic(response.detail)
    handleMp3()
  } catch (error) {
    console.error('获取音乐详情失败', error);
  }
};

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
  for (let i = 0; i < 5; i++) {
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

  // 设置下一曲
  if (audio.value) {
    switch (playMod.value) {
      case 1:
        audio.value.addEventListener('ended', () => {
          playNext(); // 播放下一首
        });
        break;
      case 2:
        audio.value.addEventListener('ended', () => {
          playRandom(); // 播放随机一首
        });
        break;
      default:
        audio.value.addEventListener('ended', () => {
          playNext(); // 播放下一首
        });
        break;
    }
  }

  audio.value.ontimeupdate = () => {
    currentTime.value = audio.value!.currentTime;
    currentProgress.value = (audio.value!.currentTime / audio.value!.duration) * 100;
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

const currentProgress = computed({
  get: () => musicStore.current.progress,
  set: (value) => {
    musicStore.current.progress = value;
    onSeek();
  }
});
// TODO
const onSeek = () => {
  if (!audio.value || !musicStore.current.detail) return;
  const time = (musicStore.current.progress / 100) * duration.value;
  audio.value.currentTime = time;
  currentTime.value = time;
};

// 添加时间更新监听
onMounted(() => {
  if (audio.value) {
    audio.value.addEventListener('timeupdate', () => {
      if (!isDragging && audio.value) {
        currentTime.value = audio.value.currentTime;
        musicStore.current.progress = (audio.value.currentTime / duration.value) * 100;
      }
    });
  }
});
const playPrevious = () => {
  if (!musicStore.CurrentMusic || !musicStore.MusicList.length) {
    message.error("没有可播放的音乐");
    return;
  }

  // 获取当前音乐在列表中的索引
  const currentIndex = musicStore.MusicList.findIndex(
      (music: MusicDetail) => music.id === musicStore.CurrentMusic?.id
  );

  let previousMusic = musicStore.MusicList[currentIndex]
  if (currentIndex > 0) {
    previousMusic = musicStore.MusicList[currentIndex - 1];
  } else {
    previousMusic = musicStore.MusicList[musicStore.MusicList.length - 1];
  }
  playMusic(previousMusic.id);
};

const playNext = () => {
  if (!musicStore.CurrentMusic || !musicStore.MusicList.length) {
    message.error("没有可播放的音乐");
    return;
  }

  // 获取当前音乐在列表中的索引
  const currentIndex = musicStore.MusicList.findIndex(
      (music: MusicDetail) => music.id === musicStore.CurrentMusic?.id
  );

  let nextMusic = musicStore.MusicList[currentIndex];
  if (currentIndex < musicStore.MusicList.length - 1) {
    // 如果有下一曲，播放下一曲
     nextMusic = musicStore.MusicList[currentIndex + 1];
  } else {
    nextMusic = musicStore.MusicList[0];
  }
  playMusic(nextMusic.id);
};

const playRandom = () => {
  if (!musicStore.CurrentMusic || !musicStore.MusicList.length) {
    message.error("没有可播放的音乐");
    return;
  }

  // 随机选择下一首音乐
   const nextIndex = Math.floor(Math.random() * musicStore.MusicList.length);

  const nextMusic = musicStore.MusicList[nextIndex];
  playMusic(nextMusic.id);
};



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
  if (musicStore.musicListContinue) {
    fetchMusicList()
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
</style>


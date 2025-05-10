import { defineStore } from 'pinia';
import type { CurrentMusicState, MusicDetail, MusicProfile } from '@/api/types/music';
import { ref, computed } from 'vue';

export const useMusicStore = defineStore('music', () => {
    // 定义状态
    const musicList = ref<MusicProfile[]>([]);
    const current = ref<CurrentMusicState>({
        detail: null,
        isPlaying: false,
        progressTime: 0,
        progress: 0,
        volume: 80
    });
    const musicListPage = ref(1);
    const musicListContinue = ref(true);

    // 计算属性
    const count = computed(() => musicList.value.length);
    const currentMusicName = computed(() => current.value.detail?.name || '');
    const currentMusicArtistName = computed(() => {
        if (!current.value.detail) return '';
        const artists = current.value.detail.artist_names || [];
        return Array.isArray(artists) ? artists.join(' & ') : '';
    });

    // 方法
    function append(newMusics: MusicProfile | MusicProfile[]) {
        musicList.value.push(...(Array.isArray(newMusics) ? newMusics : [newMusics]));
    }

    function setCurrentMusic(detail: MusicDetail) {
        current.value = {
            ...current.value,
            detail: { ...detail }
        };
    }

    function incrPage() {
        musicListPage.value += 1;
    }

    function stopPage() {
        musicListContinue.value = false;
    }

    return {
        musicList,
        current,
        musicListPage,
        musicListContinue,
        count,
        currentMusicName,
        currentMusicArtistName,
        append,
        setCurrentMusic,
        incrPage,
        stopPage
    };
}, {
    persist: {
        enabled: true,
    },
});


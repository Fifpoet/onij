import { defineStore } from 'pinia';
import type { CurrentMusicState, MusicDetail, MusicProfile } from '@/api/types/music';
import { ref, computed } from 'vue';
import { loadEnv } from 'vite';
import { MidShowWhat, TagGroup } from '@/api/types';

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
    const midShowWhat = ref(0);

    // 计算属性
    const count = computed(() => musicList.value.length);
    const currentMusicId = computed(() => current.value.detail?.id || 0);
    const currentMusicName = computed(() => current.value.detail?.name || '');
    const currentMusicArtistName = computed(() => {
        if (!current.value.detail) return '';
        const artists = current.value.detail.artist_names || [];
        return Array.isArray(artists) ? artists.join(' & ') : '';
    });
    const likeMap = computed(() => {
        if (!current.value.detail) return false;
        const mp = new Map<number, boolean>();
        for (const tag of current.value.detail.tag_details || []) {
            switch (tag.tag_group) {
                case TagGroup.TG_ArtistStar:
                    mp.set(tag.resource_id, true);
                    break;
                default:
                    mp.set(tag.tag_type, true)
                    break;
            }
        }
        return mp;
    });

    // 方法
    function append(newMusics: MusicProfile[]) {
        musicList.value.push(...newMusics);
        if (newMusics.length == 0) {
            musicListContinue.value = false;
        } else {
            musicListPage.value += 1;
        }
    }
    function setCurrentMusic(detail: MusicDetail) {
        current.value = {
            ...current.value,
            detail: { ...detail }
        };
    }
    function resetMusicList(list: MusicProfile[]) {
        musicList.value = list;
        musicListPage.value = 1;
        musicListContinue.value = true;
    }
    function setMidShowWhat(what: MidShowWhat) {
        midShowWhat.value = what;
    }

    return {
        musicList,
        current,
        musicListPage,
        musicListContinue,
        midShowWhat,
        likeMap,
        count,
        currentMusicId,
        currentMusicName,
        currentMusicArtistName,
        append,
        setCurrentMusic,
        setMidShowWhat,
        resetMusicList,
    };
}, {
    persist: {
        enabled: true,
    },
});


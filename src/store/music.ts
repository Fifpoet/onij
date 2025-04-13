import { defineStore } from 'pinia';
import type {CurrentMusicState,MusicDetail,  MusicProfile} from '@/api/types/music';


export const useMusicStore = defineStore('music', {
    // 定义状态
    state: () => ({
        musicList: [] as MusicProfile[], // 音乐列表
        current: {
            detail: null,
            isPlaying: false,
            progress: 0,
            volume: 80
        } as CurrentMusicState,
        page : 1,
    }),

    getters: {
        // 获取音乐列表数量
        count: (state) =>  { return state.musicList.length},
        currentMusicName: (state) =>   {
            return state.current.detail?.name || '';
        },
        currentMusicArtistName: (state) =>  {
            if (!state.current.detail) return '';
            const artists = this.current.detail.artist_names ||
                state.current.detail.artist_names || [];
            return Array.isArray(artists) ? artists.join(' & ') : '';
        }
    },

    actions: {
        append(newMusics: MusicProfile | MusicProfile[]) {
            this.musicList.push(...(Array.isArray(newMusics) ? newMusics : [newMusics]));
        },

        setCurrentMusic(detail: MusicDetail) {
            this.current = {
                ...this.current, // 保留其他字段
                detail: { ...detail } // 深拷贝新数据
            };
            // this.play(); TODO
        },

        incrPage() {
            this.page += 1;
        }


    },

    persist: {
        enabled: true,
    },
});


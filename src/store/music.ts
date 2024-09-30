import { defineStore } from 'pinia';

export interface Music {
    id: string;
    title: string;
    artist: string;
    composer: string;
    writer: string;
    concert: string;
    sequence: string;
    mv_url: string;
}

export const useMusicStore = defineStore('music', {
    // 定义状态
    state: () => ({
        MusicList: [] as Music[], // 保存音乐列表
        CurrentMusicId: null as string | null, // 当前播放的音乐 ID
    }),

    // 定义 getters，如果需要从 state 中派生数据，可以使用 getters

    // 定义 actions，用于修改状态
    actions: {
        // 更新音乐列表
        setMusicList(musicList: Music[]) {
            this.MusicList = musicList;
        },

        // 设置当前播放的音乐 ID
        setCurrentMusic(id: string) {
            this.CurrentMusicId = id;
        },
    },
});


import { defineStore } from 'pinia';

export interface Music {
    id: string;
    title: string;
    artist: string;
    url: string;
    // 添加更多字段根据实际的音乐数据结构
}

export const useMusicStore = defineStore('music', {
    // 定义状态
    state: () => ({
        MusicList: [] as Music[], // 保存音乐列表
        CurrentMusicId: null as string | null, // 当前播放的音乐 ID
    }),

    // 定义 getters，如果需要从 state 中派生数据，可以使用 getters
    getters: {
        currentMusic(state) {
            return state.MusicList.find((music) => music.id === state.CurrentMusicId);
        },
    },

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


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
export interface MusicDetail {
    id: number;
    root_id: number;
    title: string;
    artist_ids: string;
    composer: number;
    writer: number;
    issue_year: number;
    language: number;
    perform_type: number;
    concert: string;
    concert_year: number;
    sequence: number;
    mv_url: string;
    cover_oss: number;
    mp_oss: number;
    lyric_oss: number;
    sheet_oss: number;

    artist_name: string;
    composer_name: string;
    writer_name: string;
    cover_url: string;
    mp_url: string;
    lyric_url: string;
    sheet_url: string;

    cover?: File;
    mp?: File;
    lyric?: File;
    sheet?: File;
}
export interface UpsertMusicReq {
    id: number;
    root_id: number;
    title: string;
    artist_ids: number[];
    composer: number;
    writer: number;
    issue_year: number;
    language: number;
    perform_type: number;
    concert: string;
    concert_year: number;
    sequence: number;
    mv_url: string;
    cover_oss: number;
    mp_oss: number;
    lyric_oss: number;
    sheet_oss: number;
    cover?: File;
    mp?: File;
    lyric?: File;
    sheet?: File;
}

// 转换函数，将 MusicDetail 转换为 UpsertMusicReq
export const convertToUpsertMusicReq = (musicDetail: MusicDetail): UpsertMusicReq => {
    return {
        id: musicDetail.id,
        root_id: musicDetail.root_id,
        title: musicDetail.title,
        artist_ids: musicDetail.artist_ids.split(',').map(Number),  // 转换为数组
        composer: musicDetail.composer,
        writer: musicDetail.writer,
        issue_year: musicDetail.issue_year,
        language: musicDetail.language,
        perform_type: musicDetail.perform_type,
        concert: musicDetail.concert,
        concert_year: musicDetail.concert_year,
        sequence: musicDetail.sequence,
        mv_url: musicDetail.mv_url,
        cover_oss: musicDetail.cover_oss,
        mp_oss: musicDetail.mp_oss,
        lyric_oss: musicDetail.lyric_oss,
        sheet_oss: musicDetail.sheet_oss,

        // 如果有文件上传则添加，没有则跳过
        cover: musicDetail.cover,  // 或者你可以在某个输入表单中获取 File 对象
        mp: musicDetail.mp,     // 同上，处理文件
        lyric: musicDetail.lyric,  // 同上
        sheet: musicDetail.sheet,   // 同上
    };
};


export const useMusicStore = defineStore('music', {
    // 定义状态
    state: () => ({
        MusicList: [] as Music[], // 保存音乐列表
        CurrentMusic: null as MusicDetail | null, // 当前播放的音乐 ID
    }),

    // 定义 getters，如果需要从 state 中派生数据，可以使用 getters

    // 定义 actions，用于修改状态
    actions: {
        // 更新音乐列表
        setMusicList(musicList: Music[]) {
            this.MusicList = musicList;
        },

        // 设置当前播放的音乐 ID
        setCurrentMusic(detail: MusicDetail) {
            this.CurrentMusic = detail;
        },
    },
});


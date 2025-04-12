// src/api/types/music.ts

/* 请求参数类型 */
export interface GetMusicListReq {
    sort_type?: MusicSortType;
    keyword?: string;
    album_id?: number;
    artist_id?: number;
    page: number;
    limit: number;
}

/* 响应数据类型 */
export interface GetMusicListResp {
    code: number;
    message: string;
    musics: MusicProfile[];
}

export interface MusicProfile {
    id: number;
    name: string;
}

export enum MusicSortType {
    MST_Unknown = 0,
    MST_Created_At_Desc = 1,
    MST_Created_At_Asc = 2,
    MST_Issue_Time_Desc = 3,
    MST_Issue_Time_Asc = 4,
}

/* 通用响应结构 */
export interface BaseResponse<T> {
    code: number;
    message: string;
    data?: T;
}
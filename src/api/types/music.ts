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

export interface GetMusicDetailReq {
    music_id: number;
}

// 响应数据类型
export interface GetMusicDetailResp {
    code: number;
    message: string;
    detail: MusicDetail;
}

// 添加上传音乐的接口定义

/* 请求参数类型 */
export interface UploadMusicReq {
  // 基础信息
  name: string;
  artist_ids: number[];
  mp3_file_id: number;
  lyrics_file_id: number;

  // 表单信息
  composer_id?: number;
  writer_id?: number;
  album_id?: number;
  mv_url?: string;
  root_music_id?: number;
  issue_time?: number;

  // 更新
  music_id?: number;
}

/* 响应数据类型 */
export interface UploadMusicResp {
  code: number;
  message: string;
  music_id: number;
}

// 更新 MusicDetail 接口，添加缺少的字段
export interface MusicDetail {
  id: number;
  name: string;
  artist_ids: number[];
  artist_names: string[];
  composer_id: number;
  composer_name: string;
  writer_id: number;
  writer_name: string;
  issue_time: number;
  mv_url: string;

  mp3_file_url: string;
  lyrics_file_url: string;

  // 更新专辑信息结构
  album_profile: {
    id: number;
    name: string;
    cover_file_url: string;
  };
  // 添加标签信息
  tag_details: Array<{
    resource_id: number;
    resource_type: number;
    tag_biz: number;
    tag_group: number;
    tag_type: number;
    target_id?: number;
    target_type?: number;
    extra: string;
    list_show: boolean;
  }>;
}

export interface MusicProfile {
  id: number;
  name: string;
  singer_profiles: Array<{
    id: number;
    name: string;
    artist_type: number;
  }>;
  tags: Array<{
    resource_id: number;
    resource_type: number;
    tag_biz: number;
    tag_group: number;
    tag_type: number;
    target_id?: number;
    target_type?: number;
    extra: string;
    list_show: boolean;
  }>;
}

// ********* 自定model

export interface CurrentMusicState {
    detail: MusicDetail | null;
    isPlaying: boolean;
    progress: number; // 当前播放进度 (0-100)
    progressTime: number; // 当前播放时间 (秒)
    volume: number;   // 音量 (0-100)
}
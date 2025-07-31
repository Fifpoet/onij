// src/api/types/album.ts

import { MusicProfile } from './music';

/* 请求参数类型 */
export interface UploadAlbumReq {
  name: string;
  artist_id: number;
  cover_file_id: number;
  issue_time: number;
  music_id?: number;
  related_album_id?: number;
}

/* 响应数据类型 */
export interface UploadAlbumResp {
  code: number;
  message: string;
  album_id: number;
}

export interface GetAlbumDetailReq {
  album_id: number;
}

export interface GetAlbumDetailResp {
  code: number;
  message: string;
  detail: AlbumDetail;
}

export interface SearchAlbumReq {
  keyword: string;
  artist_id?: number;
  page: number;
  limit: number;
}

export interface SearchAlbumResp {
  code: number;
  message: string;
  albums: AlbumProfile[];
}

export interface AlbumDetail {
  id: number;
  name: string;
  artist_id: number;
  artist_name: string;
  issue_time: number;
  
  cover_file_url: string;
  mv_url: string;
  
  related_musics: MusicProfile[];
}

export interface AlbumProfile {
  id: number;
  name: string;
  cover_file_url: string;
}
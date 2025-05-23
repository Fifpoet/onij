// src/api/types/artist.ts

/* 枚举类型 */
export enum ArtistType {
  AT_Unknown = 0,
  AT_Singer = 1,
  AT_Writer = 2,
  AT_Composer = 3,
}

/* 请求参数类型 */
export interface UploadArtistReq {
  name: string;
  artist_type: ArtistType;
}

/* 响应数据类型 */
export interface UploadArtistResp {
  code: number;
  message: string;
  artist_id: number;
}

/* 搜索艺人请求参数 */
export interface SearchArtistReq {
  keyword: string;
  tag_group?: number;
  tag_type?: number;
  page: number;
  limit: number;
}

/* 搜索艺人响应数据 */
export interface SearchArtistResp {
  code: number;
  message: string;
  artists: Array<{
    id: number;
    name: string;
    artist_type: ArtistType;
  }>;
}
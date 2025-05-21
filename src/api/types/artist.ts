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
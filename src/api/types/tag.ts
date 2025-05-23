// src/api/types/tag.ts

export enum TagBiz {
  TB_Unknown = 0,
  TB_Music = 1,
}

export enum TagGroup {
  TG_Unknown = 0,
  TG_MusicTheme = 101,
  TG_MusicStyle = 102,
  TG_MusicEmotion = 103,
  TG_MusicHighlight = 104,
  TG_MusicCollection = 105,
  TG_ArtistStar = 120,
}

export enum TagType {
  TT_Unknown = 0,
  TT_NcBridge = 10401,
  TT_Star = 12001,
}

export interface TagDetail {
  resource_id: number;
  resource_type: number;
  tag_biz: TagBiz;
  tag_group: TagGroup;
  tag_type: TagType;
  target_id?: number;
  target_type?: number;
  extra: string;
  list_show: boolean;
}

export interface UploadTagReq {
  tag_detail: TagDetail;
}

export interface UploadTagResp {
  code: number;
  message: string;
}
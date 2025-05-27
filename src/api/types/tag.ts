// src/api/types/tag.ts
import { NList, NListItem, NThing, NSpace, NTag, NSelect, NIcon, SelectGroupOption } from 'naive-ui';

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

// 选择tag的value值结构, 用于区分不同tag类型
export interface SelectTagExtra {
  artist_id?: number;
  artist_type?: number;
  tag_type?: number;
}

// select tag
export const tagOpts: SelectGroupOption[] = [
  {
    type: 'group',
    label: '主题',
    key: 'theme',
    children: [{
      label: "第三者",
      value: JSON.stringify({tag_type: 124325}),
      type: 'success'
    }]
  },
  {
    type: 'group',
    label: '风格',
    key: 'style',
    children: [{
      label: "hhh",
      value: JSON.stringify({tag_type: 124331425}),
      type: 'success'
    }]
  },
  {
    type: 'group',
    label: '情感',
    key: 'emotion',
    children: [{
      label: "痛",
      value: JSON.stringify({tag_type: 1332452}),
      type: 'success'
    }]
  },
  {
    type: 'group',
    label: 'highlight',
    key: 'highlight',
    children: [{
      label: "神级live",
      value: JSON.stringify({tag_type: 24531}),
      type: 'success'
    }]
  },
  {
    type: 'group',
    label: '合集',
    key: 'collection',
    children: [{
      label: "垃圾三部曲",
      value: JSON.stringify({tag_type: "1351341"}),
      type: 'success'
    }]
  }
]

export interface UploadTagReq {
  tag_detail: TagDetail;
}

export interface UploadTagResp {
  code: number;
  message: string;
}
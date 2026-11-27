/**
 * 前端专用的 UI 常量与类型。
 *
 * 这里只放「后端 proto 里没有、但界面需要」的东西。
 * 凡是 proto 有定义的，一律走 gen:types 生成的 enums / models / api，不要在这里重抄。
 */
import type { SelectGroupOption } from 'naive-ui'
import { ArtistType, FileType, PracticeType, ResourceType, TagBiz, TagGroup, TagType } from './enums'

export const FILE_TYPE_LABELS: Record<FileType, string> = {
  [FileType.FT_Unknown]: '未知',
  [FileType.FT_Pdf]: 'PDF',
  [FileType.FT_Word]: 'Word',
  [FileType.FT_Excel]: 'Excel',
  [FileType.FT_PowerPoint]: 'PPT',
  [FileType.FT_Image]: '图片',
  [FileType.FT_Video]: '视频',
  [FileType.FT_Audio]: '音频',
  [FileType.FT_Text]: '文本',
  [FileType.FT_Archive]: '压缩包',
  [FileType.FT_Code]: '代码',
  [FileType.FT_CSV]: 'CSV',
  [FileType.FT_Executable]: '可执行文件',
  [FileType.FT_Folder]: '文件夹',
}

export const PRACTICE_TYPE_LABELS: Record<PracticeType, string> = {
  [PracticeType.PCT_Unknown]: '未分类',
  [PracticeType.PCT_AnaerobicExercise]: '无氧',
  [PracticeType.PCT_AerobicExercise]: '有氧',
  [PracticeType.PCT_Reading]: '阅读',
  [PracticeType.PCT_Language]: '语言',
  [PracticeType.PCT_Game]: '游戏',
  [PracticeType.PCT_Piano]: '钢琴',
  [PracticeType.PCT_Flute]: '长笛',
  [PracticeType.PCT_Singing]: '歌唱',
  [PracticeType.PCT_Other]: '其他',
}

/** 音乐详情弹层的展示模式 */
export enum MidShowWhat {
  ShowMusicDetail = 1,
  ShowMusicLike = 2,
  ShowMemoList = 3,
  ShowFileList = 4,
  ShowChatWindow = 5,
}

/** 标签详情（后端 Tag 的展示侧扩展） */
export interface TagDetail {
  resource_id: number
  resource_type: ResourceType
  tag_biz: TagBiz
  tag_group: TagGroup
  tag_type: TagType
  target_id?: number
  target_type?: number
  extra: string
  list_show: boolean
}

/** 选择 tag 时 value 里序列化的结构，用于区分不同类型的标签 */
export interface SelectTagExtra {
  artist_id?: number
  artist_type?: ArtistType
  tag_type?: TagType
}

/**
 * 文件下载接口的请求 / 响应。
 * 走 server/handler/download_file.go 的 customizedRegister 手写路由，不在 proto 里。
 */
export interface DownloadFileReq {
  file_ids: number[]
}

export interface DownloadFileResp {
  code: number
  message: string
  urls: string[]
}

/** 文件列表项的别名（历史命名，等同生成的 File） */
export type { File as FileDetail } from './models'

/** 搜索栏里的固定标签分组 */
export const tagOpts: SelectGroupOption[] = [
  {
    type: 'group',
    label: '主题',
    key: String(TagGroup.TG_MusicTheme),
    children: [
      {
        label: '第三者',
        value: JSON.stringify({ tag_type: TagType.TT_Star }),
      },
    ],
  },
  {
    type: 'group',
    label: '风格',
    key: String(TagGroup.TG_MusicStyle),
    children: [
      {
        label: '流行',
        value: JSON.stringify({ tag_type: TagType.TT_NcBridge }),
      },
    ],
  },
]

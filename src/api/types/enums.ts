/**
 * 与 server/_proto/onij/all_enum.proto 保持一致
 */

/** 文件类型 */
export enum FileType {
  FT_Unknown = 0,
  FT_Pdf = 1,
  FT_Word = 2,
  FT_Excel = 3,
  FT_PowerPoint = 4,
  FT_Image = 5,
  FT_Video = 6,
  FT_Audio = 7,
  FT_Text = 8,
  FT_Archive = 9,
  FT_Code = 10,
  FT_CSV = 11,
  FT_Executable = 12,
  FT_Folder = 99,
}

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

export enum ArtistType {
  ATT_Unknown = 0,
  ATT_Artist = 1,
  ATT_Writer = 2,
  ATT_Composer = 3,
}

export enum AlbumType {
  ALT_Unknown = 0,
  ALT_Studio = 1,
  ALT_Live = 2,
  ALT_OST = 3,
  ALT_Cover = 4,
  ALT_Other = 99,
}

export enum AudioQuality {
  AQ_Unknown = 0,
  AQ_Low = 1,
}

export enum PerformType {
  PT_Unknown = 0,
}

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

export enum ResourceType {
  RT_Unknown = 0,
  RT_Artist = 1,
  RT_Music = 2,
  RT_Album = 3,
}

/** 练习类型，与 all_enum.proto PracticeType 一致 */
export enum PracticeType {
  PCT_Unknown = 0,
  PCT_AnaerobicExercise = 1,
  PCT_AerobicExercise = 2,
  PCT_Reading = 3,
  PCT_Language = 4,
  PCT_Game = 5,
  PCT_Piano = 6,
  PCT_Flute = 7,
  PCT_Singing = 8,
  PCT_Other = 99,
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

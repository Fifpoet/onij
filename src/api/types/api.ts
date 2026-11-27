// 此文件由 scripts/gen_ts_types.py 自动生成，请勿手工修改。
// 源：server/_proto/onij/all.proto
// 重新生成：pnpm gen:types

import type {
  Album,
  AlbumMusic,
  Artist,
  File,
  Music,
  Practice,
  Tag,
} from './models'
import type { AlbumType } from './enums'
import type { AudioQuality } from './enums'
import type { FileType } from './enums'
import type { PerformType } from './enums'
import type { PracticeType } from './enums'
import type { TagType } from './enums'

export interface UploadFileReq {
  parent_id: number
  files: UploadFileReq_FileInfo[]
  url?: string
}

export interface UploadFileReq_FileInfo {
  filename: string
  store_key: string
  hash: string
  format: FileType
  exif: string
  size: number
}

export interface UploadFileResp {
  code: number
  message: string
  file_ids: number[]
}

export interface UploadMusicReq {
  name: string
  sub_name: string
  artist_ids: number[]
  audio_file_id: number
  lyric_file_id: number
  time_length: number
  composer_ids: number[]
  writer_ids: number[]
  album_id?: number
  mv_url?: string
  root_music_id?: number
  issue_time?: number
  priority?: number
  lyric_content?: string
  audio_quality: AudioQuality
  perform_type: PerformType
  third_id?: number
  music_id?: number
}

export interface UploadMusicResp {
  code: number
  message: string
  music_id: number
}

export interface UploadArtistReq {
  name: string
  avatar_file_id?: number
  artist_id?: string
}

export interface UploadArtistResp {
  code: number
  message: string
  artist_id: number
}

export interface UploadTagReq {
  tag: Tag
}

export interface UploadTagResp {
  code: number
  message: string
}

export interface UploadAlbumReq {
  name: string
  profile: string
  artist_ids: number[]
  cover_file_id: number
  issue_time: number
  album_type: AlbumType
  live_url: string
  third_id?: number
  album_musics: UploadAlbumReq_AlbumMusic[]
  album_id?: number
}

export interface UploadAlbumReq_AlbumMusic {
  name: string
  time_length: number
  music_id: number
  artist_name: string[]
  third_id?: number
  is_available: boolean
}

export interface UploadAlbumResp {
  code: number
  message: string
  album_id: number
}

export interface GetMusicDetailReq {
  music_id: number
}

export interface GetMusicDetailResp {
  code: number
  message: string
  music: Music
}

export interface UpdateMusicMvReq {
  music_id?: number
  third_id?: number
  mv_url: string
}

export interface UpdateMusicMvResp {
  code: number
  message: string
  music_id: number
}

export interface GetMusicMvBatchReq {
  third_ids: number[]
}

export interface MusicMvItem {
  third_id: number
  music_id: number
  mv_url: string
}

export interface GetMusicMvBatchResp {
  code: number
  message: string
  items: MusicMvItem[]
}

export interface GetAlbumDetailReq {
  album_id: number
}

export interface GetAlbumDetailResp {
  code: number
  message: string
  album: Album
  album_musics: AlbumMusic[]
}

export interface GetMusicListReq {
  keyword: string
  artist_ids: number[]
  tag_types: TagType[]
  page: number
  limit: number
}

export interface GetMusicListResp {
  code: number
  message: string
  musics: Music[]
}

export interface GetFileListReq {
  parent_id?: number
  keyword?: string
  page: number
  limit: number
}

export interface GetFileListResp {
  code: number
  message: string
  files: File[]
  total: number
}

export interface GetArtistListReq {
  keyword: string
  name: string
  page: number
  limit: number
}

export interface GetArtistListResp {
  code: number
  message: string
  artists: Artist[]
}

export interface GetAlbumListReq {
  keyword: string
  artist_id?: number
  page: number
  limit: number
}

export interface GetAlbumListResp {
  code: number
  message: string
  albums: Album[]
}

export interface DeleteFileReq {
  file_id: number
}

export interface DeleteFileResp {
  code: number
  message: string
}

export interface DeleteTagReq {
  resource_id: number
  tag_type: TagType
}

export interface DeleteTagResp {
  code: number
  message: string
}

export interface GetUploadTokenReq {
}

export interface GetUploadTokenResp {
  code: number
  message: string
  upload_token: string
  domain: string
}

export interface UploadPracticeReq {
  practice_type: PracticeType
  content: string
  practice_at: number
  duration: number
  id?: number
}

export interface UploadPracticeResp {
  code: number
  message: string
  practice_id: number
}

export interface DeletePracticeReq {
  id: number
}

export interface DeletePracticeResp {
  code: number
  message: string
}

export interface GetMonthlyPracticeReq {
  year: number
  month: number
}

export interface PracticeDay {
  day: number
  practices: Practice[]
}

export interface GetMonthlyPracticeResp {
  code: number
  message: string
  days: PracticeDay[]
}

// 此文件由 scripts/gen_ts_types.py 自动生成，请勿手工修改。
// 源：server/_proto/onij/all_model.proto
// 重新生成：pnpm gen:types

import type { AlbumType } from './enums'
import type { ArtistType } from './enums'
import type { FileType } from './enums'
import type { PracticeType } from './enums'
import type { ResourceType } from './enums'
import type { TagBiz } from './enums'
import type { TagGroup } from './enums'
import type { TagType } from './enums'

export type { AlbumType } from './enums'
export type { ArtistType } from './enums'
export type { FileType } from './enums'
export type { PracticeType } from './enums'
export type { ResourceType } from './enums'
export type { TagBiz } from './enums'
export type { TagGroup } from './enums'
export type { TagType } from './enums'

export interface Music {
  id: number
  name: string
  artists: Artist[]
  composers: Artist[]
  writers: Artist[]
  issue_time: number
  mv_url: string
  audio_file_url: string
  lyrics_file_url: string
  album: Album[]
  tags: Tag[]
}

export interface Tag {
  id: number
  resource_id: number
  resource_type: ResourceType
  tag_biz: TagBiz
  tag_group: TagGroup
  tag_type: TagType
  target_id: number
  extra: string
}

export interface Artist {
  id: number
  name: string
  avatar_file_url: string
  artist_types: ArtistType[]
}

export interface Album {
  id: number
  name: string
  profile: string
  album_type: AlbumType
  artists: Artist[]
  issue_time: number
  cover_file_url: string
  live_url: string
  album_musics: AlbumMusic[]
}

export interface AlbumMusic {
  name: string
  time_length: number
  music_id: number
  artists: Artist[]
  is_available: boolean
}

export interface File {
  id: number
  name: string
  format: FileType
  size: number
  parent_id: number
  extra: string
  created_at: number
  updated_at: number
  url: string
}

export interface Practice {
  id: number
  practice_type: PracticeType
  content: string
  practice_at: number
  duration: number
  created_at: number
  updated_at: number
}

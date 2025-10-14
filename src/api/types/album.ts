import { AlbumType } from './enums'
import { Artist } from './artist'

export interface AlbumMusic {
  name: string
  timeLength: number
  musicId: number
  artists: Artist[]
  isAvailable: boolean
}

export interface Album {
  id: number
  name: string
  profile: string
  albumType: AlbumType
  artists: Artist[]
  issueTime: number
  coverFileUrl: string
  liveUrl: string
  albumMusics: AlbumMusic[]
}

// 请求和响应类型
export interface UploadAlbumReq {
  name: string
  profile: string
  artistIds: number[]
  coverFileId: number
  issueTime: number
  albumType: AlbumType
  liveUrl: string
  thirdId?: number
  albumMusics: Array<{
    name: string
    timeLength: number
    musicId: number
    artistName: string[]
    thirdId?: number
    isAvailable: boolean
  }>
  albumId?: number
}

export interface UploadAlbumResp {
  code: number
  message: string
  albumId: number
}

export interface GetAlbumDetailReq {
  albumId: number
}

export interface GetAlbumDetailResp {
  code: number
  message: string
  album: Album
  albumMusics: AlbumMusic[]
}

export interface GetAlbumListReq {
  keyword: string
  artistId?: number
  page: number
  limit: number
}

export interface GetAlbumListResp {
  code: number
  message: string
  albums: Album[]
}
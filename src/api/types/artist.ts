import { ArtistType } from './enums'

export interface Artist {
  id: number
  name: string
  avatarFileUrl: string
  artistTypes: ArtistType[]
}

// 请求和响应类型
export interface UploadArtistReq {
  name: string
  avatarFileId?: number
  artistId?: string
}

export interface UploadArtistResp {
  code: number
  message: string
  artistId: number
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
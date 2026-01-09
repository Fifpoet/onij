// 通用响应结构
import { SearchResult } from "@/api/netease/search.ts";
import {ArtistDetail, ArtistHotSong, SongPrivilege} from "@/api/netease/artists.ts";

export interface BaseResponse<T = any> {
    result: T
    code: number
    message?: string
}

// 歌曲搜索响应
export interface SearchResponse extends BaseResponse<SearchResult> {}

// 歌手详情 - 完整结构
export interface ArtistDetailResponse {
    artist: ArtistDetail
    hotSongs: (ArtistHotSong & { privilege: SongPrivilege })[]
    more: boolean
    code: number
}

// 歌曲详情 - 质量信息
export interface SongDetailQuality {
    br: number
    fid: number
    size: number
    vd: number
    sr: number
}

// 歌曲详情 - 艺术家
export interface SongDetailArtist {
    id: number
    name: string
    tns: string[]
    alias: string[]
}

// 歌曲详情 - 专辑
export interface SongDetailAlbum {
    id: number
    name: string
    picUrl: string
    tns: string[]
    pic_str: string
    pic: number
}

// 歌曲详情
export interface SongDetail {
    name: string
    mainTitle: string | null
    additionalTitle: string | null
    id: number
    pst: number
    t: number
    ar: SongDetailArtist[]
    alia: string[]
    pop: number
    st: number
    rt: string
    fee: number
    v: number
    crbt: string | null
    cf: string
    al: SongDetailAlbum
    dt: number
    h: SongDetailQuality | null
    m: SongDetailQuality | null
    l: SongDetailQuality | null
    sq: SongDetailQuality | null
    hr: SongDetailQuality | null
    a: any | null
    cd: string
    no: number
    rtUrl: string | null
    ftype: number
    rtUrls: string[]
    djId: number
    copyright: number
    s_id: number
    mark: number
    originCoverType: number
    originSongSimpleData: any | null
    tagPicList: any | null
    resourceState: boolean
    version: number
    songJumpInfo: any | null
    entertainmentTags: any | null
    awardTags: any | null
    displayTags: string[]
    markTags: string[]
    single: number
    noCopyrightRcmd: any | null
    mv: number
    rtype: number
    rurl: string | null
    mst: number
    cp: number
    publishTime: number
}

// 歌曲详情响应
export interface SongDetailResponse {
    songs: SongDetail[]
    code: number
}
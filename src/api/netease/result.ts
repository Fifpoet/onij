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

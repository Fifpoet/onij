// 通用响应结构
import { SearchResult } from "@/api/netease/search.ts";

export interface BaseResponse<T = any> {
    result: T
    code: number
    message?: string
}

// 歌曲搜索响应
export interface SearchResponse extends BaseResponse<SearchResult> {}

// 请求参数类型
export interface UploadMemoReq {
    title: string;
    content: string;
    origin_at: number;
    cover_file_id: number;
    profile_content: string;
}

export interface GetMemoDetailReq {
    memo_id: number;
}

export interface GetMemoListReq {
    keyword?: string;
    page: number;
    limit: number;
}

// 响应类型
export interface UploadMemoResp {
    code: number;
    message: string;
    memo_id: number;
}

export interface GetMemoDetailResp {
    code: number;
    message: string;
    memo_detail: MemoDetail;
}

export interface GetMemoListResp {
    code: number;
    message: string;
    memo_profiles: MemoProfile[];
}

// 数据类型
export interface MemoDetail {
    id: number;
    title: string;
    content: string;
    origin_at: number;
    cover_url: string;
    profile_content: string;
}

export interface MemoProfile {
    id: number;
    title: string;
    cover_url: string;
    profile_content: string;
    created_at: number;
    updated_at: number;
    origin_at: number;
}
import type {UploadCustomRequestOptions} from "naive-ui";

import type {
    GetMusicListReq,
    GetMusicListResp,
    GetMusicDetailReq,
    GetMusicDetailResp,
    MusicSortType
} from '@/api/types/music';
import apiClient from "@/util/http.ts";


export const GetMusicList = async (
    params: GetMusicListReq
): Promise<GetMusicListResp> => {
    try {
        // 转换枚举值为后端需要的数值
        const requestParams = {
            ...params,
            sort_type: params.sort_type
        };
        const response = await apiClient.post<GetMusicListResp>('/music/list', {
            ...requestParams
        });
        if (response.data.code !== 200) {
            throw new Error(response.data.message || '获取音乐列表失败');
        }
        return response.data;
    } catch (error) {
        console.error('getMusicList error:', error);
        throw error;
    }
};

export const GetMusicDetail = async (
    params: GetMusicDetailReq
): Promise<GetMusicDetailResp> => {
    try {
        // 注意：根据IDL定义，这是POST请求
        const response = await apiClient.post<GetMusicDetailResp>('/music/detail', params);

        if (response.data.code !== 200) {
            throw new Error(response.data.message || '获取音乐详情失败');
        }

        return response.data;
    } catch (error) {
        console.error('getMusicDetail error:', error);
        throw error;
    }
};




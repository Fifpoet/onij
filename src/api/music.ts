import type {UploadCustomRequestOptions} from "naive-ui";
import {convertToUpsertMusicReq} from "@/store/music.ts";

import type {
    GetMusicListReq,
    GetMusicListResp,
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


export function uploadMp3 ({file}: UploadCustomRequestOptions) {
    if (currentMusicDetail.value) {
        if (file.file) {
            currentMusicDetail.value.mp = file.file; // 确保 file.file 不是 null
            console.log("上传mp3文件暂存: ", file.name)
        } else {
            message.error("文件无效");
            return;
        }
    } else {
        message.error("请先选择歌曲");
        return;
    }
    message.info("上传mp3成功");
}

export  async function  submitMusicForm  () {
    if (currentMusicDetail.value) {
        if (selectedSingerValues.value) {
            currentMusicDetail.value.artist_ids = selectedSingerValues.value;
        } else {
            currentMusicDetail.value.artist_ids = [];
        }


        const upsertMusicReq = convertToUpsertMusicReq(currentMusicDetail.value);
        console.log(upsertMusicReq);
        // 发送更新音乐详情的请求
        const response = await apiClient.post(`/music/upsert`, upsertMusicReq, {
            headers: {
                'Content-Type': 'multipart/form-data', // 指定请求类型为 multipart/form-data
            },
        });
        console.log('音乐详情已更新', response.data);
        // 可以选择刷新音乐列表或其他操作
        if (response.status === 200) {
            message.success("上传成功");
        } else {
            message.error("上传失败");
        }
    }

}

// src/api/music.ts
import {post} from '@/util';
import type {
  GetMusicDetailReq,
  GetMusicDetailResp,
  GetMusicListReq,
  GetMusicListResp,
  GetMusicMvBatchReq,
  GetMusicMvBatchResp,
  UpdateMusicMvReq,
  UpdateMusicMvResp,
  UploadMusicReq,
  UploadMusicResp,
} from '@/api/types'

// 获取音乐列表
export const GetMusicList = async (params: GetMusicListReq): Promise<GetMusicListResp> => {
    return (await post<GetMusicListResp>('/music/list', params)).data;
};

// 获取音乐详情
export const GetMusicDetail = async (params: GetMusicDetailReq): Promise<GetMusicDetailResp> => {
    return (await post<GetMusicDetailResp>('/music/detail', params)).data;
};

// 上传音乐
export const UploadMusic = async (params: UploadMusicReq): Promise<UploadMusicResp> => {
    return (await post<UploadMusicResp>('/music/upload', params)).data;
};

/** 更新 MV 内嵌地址（持久化至 music.mv_url） */
export const UpdateMusicMv = async (params: UpdateMusicMvReq): Promise<UpdateMusicMvResp> => {
    return (await post<UpdateMusicMvResp>('/music/mv', params)).data;
};

/** 播放前批量查询 MV（third_id = 网易云歌曲 id） */
export const GetMusicMvBatch = async (params: GetMusicMvBatchReq): Promise<GetMusicMvBatchResp> => {
    return (await post<GetMusicMvBatchResp>('/music/mv/batch', params)).data;
};


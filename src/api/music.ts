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
    return post<GetMusicListResp>('/music/list', params);
};

// 获取音乐详情
export const GetMusicDetail = async (params: GetMusicDetailReq): Promise<GetMusicDetailResp> => {
    return post<GetMusicDetailResp>('/music/detail', params);
};

// 上传音乐
export const UploadMusic = async (params: UploadMusicReq): Promise<UploadMusicResp> => {
    return post<UploadMusicResp>('/music/upload', params);
};

/** 更新 MV 内嵌地址（持久化至 music.mv_url） */
export const UpdateMusicMv = async (params: UpdateMusicMvReq): Promise<UpdateMusicMvResp> => {
    return post<UpdateMusicMvResp>('/music/mv', params);
};

/** 播放前批量查询 MV（third_id = 网易云歌曲 id） */
export const GetMusicMvBatch = async (params: GetMusicMvBatchReq): Promise<GetMusicMvBatchResp> => {
    return post<GetMusicMvBatchResp>('/music/mv/batch', params);
};


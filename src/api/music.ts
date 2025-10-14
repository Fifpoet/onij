// src/api/music.ts
import {post} from '@/util';
import {
  GetMusicDetailReq,
  GetMusicDetailResp,
  GetMusicListReq,
  GetMusicListResp,
  UploadMusicReq,
  UploadMusicResp
} from '@/api/types';

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




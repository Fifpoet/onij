// src/api/music.ts
import { post } from '../util/http';
import { 
  GetMusicListReq, 
  GetMusicListResp,
  GetMusicDetailReq,
  GetMusicDetailResp,
  UploadMusicReq,
  UploadMusicResp
} from './types/music';

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




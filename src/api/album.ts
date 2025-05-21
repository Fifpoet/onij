// src/api/album.ts
import { post } from '../util/http';
import { 
  UploadAlbumReq, 
  UploadAlbumResp,
  GetAlbumDetailReq,
  GetAlbumDetailResp
} from './types/album';

// 上传专辑
export const UploadAlbum = async (params: UploadAlbumReq): Promise<UploadAlbumResp> => {
  return post<UploadAlbumResp>('/album/upload', params);
};

// 获取专辑详情
export const GetAlbumDetail = async (params: GetAlbumDetailReq): Promise<GetAlbumDetailResp> => {
  return post<GetAlbumDetailResp>('/album/detail', params);
};
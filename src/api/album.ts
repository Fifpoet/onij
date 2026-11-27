import { post } from '@/util';
import {
  UploadAlbumReq,
  UploadAlbumResp,
  GetAlbumDetailReq,
  GetAlbumDetailResp,
  GetAlbumListReq,
  GetAlbumListResp
} from '@/api/types';

export const UploadAlbum = async (params: UploadAlbumReq): Promise<UploadAlbumResp> => {
  return post<UploadAlbumResp>('/album/upload', params);
};

export const GetAlbumDetail = async (params: GetAlbumDetailReq): Promise<GetAlbumDetailResp> => {
  return post<GetAlbumDetailResp>('/album/detail', params);
};

export const SearchAlbum = async (params: GetAlbumListReq): Promise<GetAlbumListResp> => {
  return post<GetAlbumListResp>('/album/search', params);
};
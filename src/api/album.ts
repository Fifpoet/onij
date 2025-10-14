// src/api/album.ts
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
  return (await post<UploadAlbumResp>('/album/upload', params)).data;
};

export const GetAlbumDetail = async (params: GetAlbumDetailReq): Promise<GetAlbumDetailResp> => {
  return (await post<GetAlbumDetailResp>('/album/detail', params)).data;
};

export const SearchAlbum = async (params: GetAlbumListReq): Promise<GetAlbumListResp> => {
  return (await post<GetAlbumListResp>('/album/search', params)).data;
};
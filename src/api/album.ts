// src/api/album.ts
import { post } from '../util/http';
import {
  UploadAlbumReq,
  UploadAlbumResp,
  GetAlbumDetailReq,
  GetAlbumDetailResp,
  SearchAlbumReq,
  SearchAlbumResp
} from './types/album';

export const UploadAlbum = async (params: UploadAlbumReq): Promise<UploadAlbumResp> => {
  return (await post<UploadAlbumResp>('/album/upload', params)).data;
};

export const GetAlbumDetail = async (params: GetAlbumDetailReq): Promise<GetAlbumDetailResp> => {
  return (await post<GetAlbumDetailResp>('/album/detail', params)).data;
};

export const SearchAlbum = async (params: SearchAlbumReq): Promise<SearchAlbumResp> => {
  return (await post<SearchAlbumResp>('/album/search', params)).data;
};
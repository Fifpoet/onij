import { post } from '../util/http';
import type {
  UploadMemoReq,
  UploadMemoResp,
  GetMemoDetailReq,
  GetMemoDetailResp,
  GetMemoListReq,
  GetMemoListResp,
} from './types/memo';

// 上传Memo
export const UploadMemo = async (params: UploadMemoReq): Promise<UploadMemoResp> => {
  return (await post<UploadMemoResp>('/memo/upload', params)).data;
};

// 获取Memo详情
export const GetMemoDetail = async (params: GetMemoDetailReq): Promise<GetMemoDetailResp> => {
  return (await post<GetMemoDetailResp>('/memo/detail', params)).data;
};

// 获取Memo列表
export const GetMemoList = async (params: GetMemoListReq): Promise<GetMemoListResp> => {
  return (await post<GetMemoListResp>('/memo/list', params)).data;
};
// src/api/file.ts
import { post, uploadFile } from '../util/http';
import { 
  UploadFileReq, 
  UploadFileResp, 
  GetFileListReq, 
  GetFileListResp,
  DownloadFileReq,
  DownloadFileResp
} from './types/file';

// 上传文件
export const UploadFile = async (params: UploadFileReq): Promise<UploadFileResp> => {
  const formData = new FormData();
  formData.append('parent_id', params.parent_id.toString());
  
  // 使用简单的字段名格式
  params.files.forEach((fileInfo, index) => {
    formData.append(`filename_${index}`, fileInfo.filename || '');
    formData.append(`file_${index}`, fileInfo.file);
    if (fileInfo.origin_at !== undefined) {
      formData.append(`origin_at_${index}`, fileInfo.origin_at.toString());
    }
  });
  
  return (await uploadFile<UploadFileResp>('/file/upload', formData));
};

// 获取文件列表
export const GetFileList = async (params: GetFileListReq): Promise<GetFileListResp> => {
  return (await post<GetFileListResp>('/file/list', params)).data;
};

// 下载文件
export const DownloadFile = async (params: DownloadFileReq): Promise<DownloadFileResp> => {
  return (await post<DownloadFileResp>('/file/download', params)).data;
};
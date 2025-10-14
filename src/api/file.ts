// src/api/file.ts
import {post} from '../util/http';
import {
    DeleteFileReq,
    DeleteFileResp,
    GetFileListReq,
    GetFileListResp,
    GetUploadTokenResp,
    UploadFileReq,
    UploadFileResp
} from '@/api/types';

export const UploadFile = async (params: UploadFileReq): Promise<UploadFileResp> => {
    const response = await post<UploadFileResp>('/file/upload', params);
    return response.data;
};

// 获取文件列表
export const GetFileList = async (params: GetFileListReq): Promise<GetFileListResp> => {
    return (await post<GetFileListResp>('/file/list', params)).data;
};

// 删除文件
export const DeleteFileById = async (params: DeleteFileReq): Promise<DeleteFileResp> => {
    return (await post<DeleteFileResp>('/file/delete', params)).data;
};

// 获取上传token
export const GetUploadToken = async (): Promise<GetUploadTokenResp> => {
    return (await post<GetUploadTokenResp>('/file/upload_token', {})).data;
};
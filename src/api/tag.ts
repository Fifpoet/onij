import {post} from '@/util';
import {DeleteTagReq, DeleteTagResp, UploadTagReq, UploadTagResp} from '@/api/types';

export const UploadTag = async (params: UploadTagReq): Promise<UploadTagResp> => {
    return post<UploadTagResp>('/tag/upload', params);
};

export const DeleteTag = async (params: DeleteTagReq): Promise<DeleteTagResp> => {
    return post<DeleteTagResp>('/tag/delete', params);
};
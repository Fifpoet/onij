import { post } from '../util/http'
import type {
  DeletePracticeReq,
  DeletePracticeResp,
  GetMonthlyPracticeReq,
  GetMonthlyPracticeResp,
  UploadPracticeReq,
  UploadPracticeResp,
} from '@/api/types/practice'

export const UploadPractice = async (params: UploadPracticeReq): Promise<UploadPracticeResp> => {
  return post('/practice/upload', params) as Promise<UploadPracticeResp>
}

export const DeletePractice = async (params: DeletePracticeReq): Promise<DeletePracticeResp> => {
  return post('/practice/delete', params) as Promise<DeletePracticeResp>
}

export const GetMonthlyPractice = async (
  params: GetMonthlyPracticeReq,
): Promise<GetMonthlyPracticeResp> => {
  return post('/practice/monthly', params) as Promise<GetMonthlyPracticeResp>
}

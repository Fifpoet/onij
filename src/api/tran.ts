import { post } from '../util/http'
import type { FileType } from '@/api/types/enums'

export type TranItemDTO = {
  id: string
  kind: 'text' | 'file'
  content?: string
  file_id?: number
  name?: string
  size?: number
  format?: FileType
  pinned: boolean
  pinned_at?: number
  created_at: number
}

type ListTranResp = {
  code: number
  message: string
  items: TranItemDTO[]
}

type ItemTranResp = {
  code: number
  message: string
  item: TranItemDTO
}

type OkTranResp = {
  code: number
  message: string
}

export const ListTran = async (): Promise<ListTranResp> => {
  return post('/tran/list', {}) as Promise<ListTranResp>
}

export const AddTranText = async (content: string): Promise<ItemTranResp> => {
  return post('/tran/text', { content }) as Promise<ItemTranResp>
}

export const AddTranFile = async (params: {
  file_id: number
  name: string
  size: number
  format: FileType
}): Promise<ItemTranResp> => {
  return post('/tran/file', params) as Promise<ItemTranResp>
}

export const PinTran = async (id: string): Promise<ItemTranResp> => {
  return post('/tran/pin', { id }) as Promise<ItemTranResp>
}

export const DeleteTran = async (id: string): Promise<OkTranResp> => {
  return post('/tran/delete', { id }) as Promise<OkTranResp>
}

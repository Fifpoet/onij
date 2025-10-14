import { TagBiz, TagGroup, TagType, ResourceType } from './enums'

export interface Tag {
  id: number
  resourceId: number
  resourceType: ResourceType
  tagBiz: TagBiz
  tagGroup: TagGroup
  tagType: TagType
  targetId: number
  extra: string
}

// 请求和响应类型
export interface UploadTagReq {
  tag: Tag
}

export interface UploadTagResp {
  code: number
  message: string
}

export interface DeleteTagReq {
  resourceId: number
  tagType: TagType
}

export interface DeleteTagResp {
  code: number
  message: string
}
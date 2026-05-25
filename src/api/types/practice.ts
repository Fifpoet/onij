import { PracticeType } from './enums'

export { PracticeType, PRACTICE_TYPE_LABELS } from './enums'

export interface Practice {
  id: number
  practice_type: PracticeType
  content: string
  practice_at: number
  duration: number
  created_at: number
  updated_at: number
}

export interface PracticeDay {
  day: number
  practices: Practice[]
}

export interface UploadPracticeReq {
  practice_type: PracticeType
  content: string
  practice_at: number
  duration?: number
  id?: number
}

export interface UploadPracticeResp {
  code: number
  message: string
  practice_id: number
}

export interface DeletePracticeReq {
  id: number
}

export interface DeletePracticeResp {
  code: number
  message: string
}

export interface GetMonthlyPracticeReq {
  year: number
  month: number
}

export interface GetMonthlyPracticeResp {
  code: number
  message: string
  days: PracticeDay[]
}

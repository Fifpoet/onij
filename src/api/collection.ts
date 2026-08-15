import { post } from '../util/http'

export type CollectionDTO = {
  id: string
  name: string
  cover_url?: string
  song_ids?: number[]
  song_count: number
  updated_at: number
}

type ListResp = { code: number; message: string; items: CollectionDTO[] }
type ItemResp = { code: number; message: string; item: CollectionDTO }
type OkResp = { code: number; message: string }

export type CollectionSongMeta = {
  song_id: number
  song_name: string
  artist_names: string
}

export const ListCollections = async (keyword = ''): Promise<ListResp> => {
  return post('/collection/list', { keyword }) as Promise<ListResp>
}

export const CreateCollection = async (params: {
  name: string
  cover_url?: string
}): Promise<ItemResp> => {
  return post('/collection/create', params) as Promise<ItemResp>
}

export const UpdateCollection = async (params: {
  id: string
  name?: string
  cover_url?: string
}): Promise<ItemResp> => {
  return post('/collection/update', params) as Promise<ItemResp>
}

export const DetailCollection = async (id: string): Promise<ItemResp> => {
  return post('/collection/detail', { id }) as Promise<ItemResp>
}

export const DeleteCollection = async (id: string): Promise<OkResp> => {
  return post('/collection/delete', { id }) as Promise<OkResp>
}

export const AddCollectionSongs = async (params: {
  id: string
  songs: CollectionSongMeta[]
}): Promise<ItemResp> => {
  return post('/collection/add_songs', params) as Promise<ItemResp>
}

export const RemoveCollectionSong = async (params: {
  id: string
  song_id: number
}): Promise<ItemResp> => {
  return post('/collection/remove_song', params) as Promise<ItemResp>
}

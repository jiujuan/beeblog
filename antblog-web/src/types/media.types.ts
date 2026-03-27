import type { PageQuery } from './api.types'

export interface Media {
  id: number
  uploader_id: number
  article_id: number | null
  original_name: string
  storage_path: string
  url: string
  mime_type: string
  file_size: number
  file_size_human: string
  width: number
  height: number
  created_at: string
}

export interface AdminListMediaQuery extends PageQuery {
  uploader_id?: number
  article_id?: number
  mime_type?: string
}

export interface BindArticleReq {
  article_id: number | null
}

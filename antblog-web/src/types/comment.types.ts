import type { PageQuery } from './api.types'

export interface Comment {
  id: number
  article_id: number
  user_id: number | null
  parent_id: number | null
  root_id: number | null
  reply_to_id: number | null
  nickname: string
  avatar: string
  content: string
  status: number // 1=待审 2=已通过 3=已拒绝
  like_count: number
  created_at: string
  children?: Comment[]
}

export interface CreateCommentReq {
  article_id: number
  parent_id?: number
  reply_to_id?: number
  content: string
  nickname?: string
  email?: string
}

export interface ListCommentQuery extends PageQuery {
  article_id: number
}

export interface AdminListCommentQuery extends PageQuery {
  article_id?: number
  status?: number
  keyword?: string
}

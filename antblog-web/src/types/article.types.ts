import type { PageQuery } from './api.types'

export interface TagSimple {
  id: number
  name: string
  slug: string
  color: string
}

/** 文章列表项（不含正文） */
export interface ArticleListItem {
  id: number
  uuid: string
  author_id: number
  category_id: number | null
  title: string
  slug: string
  summary: string
  cover: string
  status: number
  is_top: boolean
  is_featured: boolean
  view_count: number
  like_count: number
  comment_count: number
  bookmark_count: number
  word_count: number
  published_at: string | null
  created_at: string
  tags: TagSimple[]
  liked?: boolean
  bookmarked?: boolean
}

/** 文章详情（含正文） */
export interface Article extends ArticleListItem {
  content: string
  content_html: string
  allow_comment: boolean
  updated_at: string
}

/** 归档时间线 */
export interface ArchiveItem {
  year: number
  month: number
  article_count: number
}

/** 文章列表请求参数 */
export interface ListArticleQuery extends PageQuery {
  category_id?: number
  tag_id?: number
  keyword?: string
}

/** 归档详情请求 */
export interface ArchiveDetailQuery extends PageQuery {
  year: number
  month: number
}

/** 创建文章 */
export interface CreateArticleReq {
  category_id?: number
  tag_ids?: number[]
  title: string
  slug?: string
  summary?: string
  content: string
  cover?: string
  is_top?: boolean
  is_featured?: boolean
  allow_comment?: boolean
}

/** 更新文章 */
export interface UpdateArticleReq extends Partial<CreateArticleReq> {}

/** 更新状态 */
export interface UpdateStatusReq {
  status: number // 1=草稿 2=已发布 3=已归档
}

/** 后台文章列表请求 */
export interface AdminListArticleQuery extends PageQuery {
  keyword?: string
  status?: number
  category_id?: number
  tag_id?: number
}

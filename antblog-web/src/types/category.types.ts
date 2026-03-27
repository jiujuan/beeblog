export interface Category {
  id: number
  name: string
  slug: string
  description: string
  cover: string
  sort_order: number
  article_count: number
}

export interface CreateCategoryReq {
  name: string
  slug?: string
  description?: string
  cover?: string
  sort_order?: number
}

export interface UpdateCategoryReq extends Partial<CreateCategoryReq> {}

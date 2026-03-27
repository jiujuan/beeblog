export interface Tag {
  id: number
  name: string
  slug: string
  color: string
  article_count: number
}

export interface CreateTagReq {
  name: string
  slug?: string
  color?: string
}

export interface UpdateTagReq extends Partial<CreateTagReq> {}

export interface BatchCreateTagReq {
  tags: CreateTagReq[]
}

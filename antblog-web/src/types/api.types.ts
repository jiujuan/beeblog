/** 统一 API 响应结构 */
export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data: T
}

/** 分页结果 */
export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

/** 通用分页请求参数 */
export interface PageQuery {
  page?: number
  page_size?: number
}

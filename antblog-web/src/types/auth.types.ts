export interface LoginReq {
  email: string
  password: string
}

export interface RegisterReq {
  username: string
  email: string
  password: string
  nickname?: string
}

export interface UserInfo {
  id: number
  uuid: string
  username: string
  email: string
  nickname: string
  avatar: string
  bio: string
  role: number   // 1=普通用户 2=管理员
  status: number
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_at: string
}

export interface LoginResp {
  access_token: string
  refresh_token: string
  expires_at: string
  user: UserInfo
}

export interface UpdateProfileReq {
  nickname?: string
  avatar?: string
  bio?: string
}

export interface ChangePasswordReq {
  old_password: string
  new_password: string
}

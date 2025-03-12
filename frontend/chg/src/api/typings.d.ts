declare namespace API {
  type DeleteRequest = {
    id: string
  }

  type getGetParams = {
    /** 用户的ID */
    id: number
  }

  type getGetVoParams = {
    /** 用户的ID */
    id: number
  }

  type ListUserVOResponse = {
    /** 当前页数 */
    current?: number
    /** 页面大小 */
    pages?: number
    records?: UserVO[]
    /** 页面大小 */
    size?: number
    /** 总记录数 */
    total?: number
  }

  type Response = {
    code?: number
    data?: Record<string, any>
    msg?: string
  }

  type User = {
    createTime?: string
    editTime?: string
    id?: string
    updateTime?: string
    userAccount?: string
    userAvatar?: string
    userName?: string
    userPassword?: string
    userProfile?: string
    userRole?: string
  }

  type UserAddRequest = {
    /** 用户账号 */
    userAccount: string
    /** 用户头像 */
    userAvatar?: string
    /** 用户昵称 */
    userName?: string
    /** 用户简介 */
    userProfile?: string
    /** 用户权限 */
    userRole?: string
  }

  type UserLoginRequest = {
    userAccount: string
    userPassword: string
  }

  type UserLoginVO = {
    createTime?: string
    editTime?: string
    id?: string
    updateTime?: string
    userAccount?: string
    userAvatar?: string
    userName?: string
    userProfile?: string
    userRole?: string
  }

  type UserQueryRequest = {
    /** 当前页数 */
    current?: number
    /** 用户ID */
    id?: number
    /** 页面大小 */
    pageSize?: number
    /** 排序字段 */
    sortField?: string
    /** 排序顺序（默认升序） */
    sortOrder?: string
    /** 用户账号 */
    userAccount?: string
    /** 用户昵称 */
    userName?: string
    /** 用户简介 */
    userProfile?: string
    /** 用户权限 */
    userRole?: string
  }

  type UserRegsiterRequest = {
    checkPassword: string
    userAccount: string
    userPassword: string
  }

  type UserUpdateRequest = {
    /** 用户ID */
    id?: number
    /** 用户头像 */
    userAvatar?: string
    /** 用户昵称 */
    userName?: string
    /** 用户简介 */
    userProfile?: string
    /** 用户权限 */
    userRole?: string
  }

  type UserVO = {
    createTime?: string
    id?: string
    userAccount?: string
    userAvatar?: string
    userName?: string
    userProfile?: string
    userRole?: string
  }
}

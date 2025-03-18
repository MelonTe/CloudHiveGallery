declare namespace API {
  type DeleteRequest = {
    id: string
  }

  type getFileTestDownloadParams = {
    /** 文件存储在 COS 的 KEY */
    key: string
  }

  type getPictureGetParams = {
    /** 图片的ID */
    id: string
  }

  type getPictureGetVoParams = {
    /** 图片的ID */
    id: string
  }

  type getUserGetParams = {
    /** 用户的ID */
    id: string
  }

  type getUserGetVoParams = {
    /** 用户的ID */
    id: string
  }

  type ListPictureResponse = {
    /** 当前页数 */
    current?: number
    /** 总页数 */
    pages?: number
    records?: Picture[]
    /** 页面大小 */
    size?: number
    /** 总记录数 */
    total?: number
  }

  type ListPictureVOResponse = {
    /** 当前页数 */
    current?: number
    /** 总页数 */
    pages?: number
    records?: PictureVO[]
    /** 页面大小 */
    size?: number
    /** 总记录数 */
    total?: number
  }

  type ListUserVOResponse = {
    /** 当前页数 */
    current?: number
    /** 总页数 */
    pages?: number
    records?: UserVO[]
    /** 页面大小 */
    size?: number
    /** 总记录数 */
    total?: number
  }

  type Picture = {
    category?: string
    createTime?: string
    editTime?: string
    id?: string
    introduction?: string
    name?: string
    picFormat?: string
    picHeight?: number
    picScale?: number
    picSize?: number
    picWidth?: number
    /** 存储的格式：["golang","java","c++"] */
    tags?: string
    updateTime?: string
    url?: string
    userId?: string
  }

  type PictureEditRequest = {
    category?: string
    id?: string
    introduction?: string
    name?: string
    tags?: string[]
  }

  type PictureQueryRequest = {
    category?: string
    /** 当前页数 */
    current?: number
    /** 图片ID */
    id?: string
    introduction?: string
    name?: string
    /** 页面大小 */
    pageSize?: number
    picFormat?: string
    picHeight?: number
    picScale?: number
    picSize?: number
    picWidth?: number
    /** 搜索词 */
    searchText?: string
    /** 排序字段 */
    sortField?: string
    /** 排序顺序（默认升序） */
    sortOrder?: string
    tags?: string[]
    userId?: string
  }

  type PictureTagCategory = {
    categoryList?: string[]
    tagList?: string[]
  }

  type PictureUpdateRequest = {
    category?: string
    id?: string
    introduction?: string
    name?: string
    tags?: string[]
  }

  type PictureVO = {
    category?: string
    createTime?: string
    editTime?: string
    id?: string
    introduction?: string
    name?: string
    picFormat?: string
    picHeight?: number
    picScale?: number
    picSize?: number
    picWidth?: number
    tags?: string[]
    updateTime?: string
    url?: string
    user?: UserVO
    userId?: string
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
    id?: string
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
    id?: string
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

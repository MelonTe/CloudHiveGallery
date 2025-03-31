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

  type getSpaceGetVoParams = {
    /** 空间的ID */
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

  type ListSpaceResponse = {
    /** 当前页数 */
    current?: number
    /** 总页数 */
    pages?: number
    records?: Space[]
    /** 页面大小 */
    size?: number
    /** 总记录数 */
    total?: number
  }

  type ListSpaceVOResponse = {
    /** 当前页数 */
    current?: number
    /** 总页数 */
    pages?: number
    records?: SpaceVO[]
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
    reviewMessage?: string
    reviewStatus?: number
    reviewTime?: string
    reviewerId?: string
    spaceId?: string
    /** 存储的格式：["golang","java","c++"] */
    tags?: string
    thumbnailUrl?: string
    updateTime?: string
    url?: string
    userId?: string
  }

  type PictureEditRequest = {
    category?: string
    id?: string
    introduction?: string
    name?: string
    /** 空间ID */
    spaceId?: string
    tags?: string[]
  }

  type PictureQueryRequest = {
    category?: string
    /** 当前页数 */
    current?: number
    /** 图片ID */
    id?: string
    introduction?: string
    /** 是否查询空间ID为空的图片 */
    isNullSpaceId?: boolean
    name?: string
    /** 页面大小 */
    pageSize?: number
    picFormat?: string
    picHeight?: number
    picScale?: number
    picSize?: number
    picWidth?: number
    reviewMessage?: string
    /** 新增审核字段 */
    reviewStatus?: string
    /** 审核人ID */
    reviewerId?: string
    /** 搜索词 */
    searchText?: string
    /** 排序字段 */
    sortField?: string
    /** 排序顺序（默认升序） */
    sortOrder?: string
    /** 新增空间筛选字段 */
    spaceId?: string
    tags?: string[]
    /** 图片上传人信息 */
    userId?: string
  }

  type PictureReviewRequest = {
    /** 图片ID */
    id?: string
    /** 审核信息 */
    reviewMessage?: string
    /** 审核状态 */
    reviewStatus?: string
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
    /** 空间ID */
    spaceId?: string
    tags?: string[]
  }

  type PictureUploadByBatchRequest = {
    /** 图片数量 */
    count?: number
    /** 图片名称前缀，默认为SearchText */
    namePrefix?: string
    /** 搜索词 */
    searchText?: string
  }

  type PictureUploadRequest = {
    /** 图片地址 */
    fileUrl?: string
    /** 图片ID */
    id?: string
    /** 图片名称 */
    picName?: string
    /** 空间ID */
    spaceId?: string
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
    spaceId?: string
    tags?: string[]
    thumbnailUrl?: string
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

  type Space = {
    createTime?: string
    editTime?: string
    id?: string
    maxCount?: number
    maxSize?: number
    spaceLevel?: number
    spaceName?: string
    totalCount?: number
    totalSize?: number
    updateTime?: string
    userId?: string
  }

  type SpaceAddRequest = {
    /** 空间级别：0-普通版 1-专业版 2-旗舰版 */
    spaceLevel?: number
    /** 空间名称 */
    spaceName?: string
  }

  type SpaceEditRequest = {
    /** Space ID */
    id?: string
    /** Space name */
    spaceName?: string
  }

  type SpaceLevelResponse = {
    /** 空间图片的最大数量 */
    maxCount?: number
    /** 空间图片的最大总大小 */
    maxSize?: number
    /** 空间的等级名称 */
    text?: string
    /** 空间的等级 */
    value?: number
  }

  type SpaceQueryRequest = {
    /** 当前页数 */
    current?: number
    /** 空间 ID */
    id?: string
    /** 页面大小 */
    pageSize?: number
    /** 排序字段 */
    sortField?: string
    /** 排序顺序（默认升序） */
    sortOrder?: string
    /** 空间级别：0-普通版 1-专业版 2-旗舰版 使用指针来区分0和未传参 */
    spaceLevel?: number
    /** 空间名称 */
    spaceName?: string
    /** 用户 ID */
    userId?: string
  }

  type SpaceUpdateRequest = {
    /** Space ID */
    id?: string
    /** Maximum number of space images */
    maxCount?: number
    /** Maximum total size of space images */
    maxSize?: number
    /** Space level: 0-普通版 1-专业版 2-旗舰版 */
    spaceLevel?: number
    /** Space name */
    spaceName?: string
  }

  type SpaceVO = {
    createTime?: string
    editTime?: string
    /** Space ID */
    id?: string
    maxCount?: number
    maxSize?: number
    spaceLevel?: number
    spaceName?: string
    totalCount?: number
    totalSize?: number
    updateTime?: string
    user?: UserVO
    /** User ID */
    userId?: string
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

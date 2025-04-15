declare namespace API {
  type CreateOutPaintingTaskRequest = {
    /** 图像处理任务的参数 */
    parameters?: ImageParameters
    /** 图片ID */
    pictureId?: string
  }

  type CreateOutPaintingTaskResponse = {
    /** 错误码（失败时返回） */
    code?: string
    /** 错误信息（失败时返回） */
    message?: string
    /** 任务输出信息（成功时返回） */
    output?: Output
    /** 请求唯一标识符 */
    requestId?: string
  }

  type DeleteRequest = {
    id: string
  }

  type getFileTestDownloadParams = {
    /** 文件存储在 COS 的 KEY */
    key: string
  }

  type GetOutPaintingResponse = {
    /** 任务输出信息（一定包含） */
    output?: TaskDetailOutput
    /** 请求唯一标识 */
    requestId?: string
    /** 图像统计信息（仅在成功时返回） */
    usage?: Usage
  }

  type getPictureGetParams = {
    /** 图片的ID */
    id: string
  }

  type getPictureGetVoParams = {
    /** 图片的ID */
    id: string
  }

  type getPictureOutPaintingCreateTaskParams = {
    /** 任务的ID */
    taskId: string
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

  type ImageParameters = {
    /** 是否添加水印 */
    addWatermark?: boolean
    /** 图像旋转角度（单位：度） */
    angle?: number
    /** 是否启用最佳质量 */
    bestQuality?: boolean
    /** 图像底部的偏移量 */
    bottomOffset?: number
    /** 图像左侧的偏移量 */
    leftOffset?: number
    /** 是否限制图像大小 */
    limitImageSize?: boolean
    /** 输出图像的宽高比（例如："16:9"） */
    outputRatio?: string
    /** 图像右侧的偏移量 */
    rightOffset?: number
    /** 图像顶部的偏移量 */
    topOffset?: number
    /** 图像的水平缩放比例，范围在1.0 ~ 3.0 */
    xScale?: number
    /** 图像的垂直缩放比例，范围在1.0 ~ 3.0 */
    yScale?: number
  }

  type ImageSearchResult = {
    /** 来源地址 */
    fromURL?: string
    /** 缩略图地址 */
    thumbURL?: string
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

  type Output = {
    /** 任务的唯一标识符 */
    taskId?: string
    /** 任务状态：PENDING、RUNNING、SUSPENDED、SUCCEEDED、FAILED、UNKNOWN */
    taskStatus?: string
  }

  type Picture = {
    category?: string
    createTime?: string
    editTime?: string
    id?: string
    introduction?: string
    name?: string
    picColor?: string
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

  type PictureEditByBatchRequest = {
    /** 分类 */
    category?: string
    /** 名称规则，暂时只支持“名称{序号}的形式，序号将会自动递增” */
    nameRule?: string
    /** 图片ID列表 */
    pictureIdList?: string[]
    /** 空间ID */
    spaceId?: string
    /** 标签 */
    tags?: string[]
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
    /** 结束编辑时间 */
    endEditTime?: string
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
    /** 开始编辑时间 */
    startEditTime?: string
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
    reviewStatus?: number
  }

  type PictureSearchByColorRequest = {
    /** 图片颜色 */
    picColor?: string
    /** 空间ID */
    spaceId?: string
  }

  type PictureSearchByPictureRequest = {
    /** 图片ID */
    pictureId?: string
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
    /** 空间的权限列表 */
    permissionList?: string[]
    picColor?: string
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
    spaceType?: number
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
    /** 空间类型：0-个人空间 1-团队空间 */
    spaceType?: number
  }

  type SpaceCategoryAnalyzeRequest = {
    /** 是否查询所有空间 */
    queryAll?: boolean
    /** 是否查询公开空间 */
    queryPublic?: boolean
    /** 空间ID */
    spaceId?: string
  }

  type SpaceCategoryAnalyzeResponse = {
    /** 分类名称 */
    category?: string
    /** 分类数量 */
    count?: number
    /** 分类总大小 */
    totalSize?: number
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
    /** 空间类型：0-个人空间 1-团队空间 使用指针来区分0和未传参 */
    spaceType?: number
    /** 用户 ID */
    userId?: string
  }

  type SpaceRankAnalyzeRequest = {
    /** 排名前N的空间 */
    top_n?: number
  }

  type SpaceSizeAnalyzeRequest = {
    /** 是否查询所有空间 */
    queryAll?: boolean
    /** 是否查询公开空间 */
    queryPublic?: boolean
    /** 空间ID */
    spaceId?: string
  }

  type SpaceSizeAnalyzeResponse = {
    /** 分类数量 */
    count?: number
    /** 大小范围，格式为"<100KB","100KB-500KB","500KB-1MB",">1MB" */
    sizeRange?: string
  }

  type SpaceTagAnalyzeRequest = {
    /** 是否查询所有空间 */
    queryAll?: boolean
    /** 是否查询公开空间 */
    queryPublic?: boolean
    /** 空间ID */
    spaceId?: string
  }

  type SpaceTagAnalyzeResponse = {
    /** 标签数量 */
    count?: number
    /** 标签名称 */
    tag?: string
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

  type SpaceUsageAnalyzeRequest = {
    /** 是否查询所有空间 */
    queryAll?: boolean
    /** 是否查询公开空间 */
    queryPublic?: boolean
    /** 空间ID */
    spaceId?: string
  }

  type SpaceUsageAnalyzeResponse = {
    /** 资源数量使用比例 */
    countUsageRatio?: number
    /** 最大资源数量 */
    maxCount?: number
    /** 最大空间大小 */
    maxSize?: number
    /** 空间使用比例 */
    sizeUsageRatio?: number
    /** 已使用的资源数量 */
    usedCount?: number
    /** 已使用的空间大小 */
    usedSize?: number
  }

  type SpaceUser = {
    createTime?: string
    id?: string
    spaceId?: string
    spaceRole?: string
    updateTime?: string
    userId?: string
  }

  type SpaceUserAddRequest = {
    /** 空间ID */
    spaceId?: string
    /** 空间角色：viewer-查看者 editor-编辑者 admin-管理员 */
    spaceRole?: string
    /** 用户ID */
    userId?: string
  }

  type SpaceUserAnalyzeRequest = {
    /** 是否查询所有空间 */
    queryAll?: boolean
    /** 是否查询公开空间 */
    queryPublic?: boolean
    /** 空间ID */
    spaceId?: string
    /** 时间维度：day/week/month */
    timeDimension?: string
    /** 用户ID */
    userId?: string
  }

  type SpaceUserAnalyzeResponse = {
    /** 周期内上传的图片数量 */
    count?: number
    /** 时间周期 */
    period?: string
  }

  type SpaceUserEditRequest = {
    /** 表的元组ID */
    Id?: string
    /** 空间角色：viewer-查看者 editor-编辑者 admin-管理员 */
    spaceRole?: string
  }

  type SpaceUserQueryRequest = {
    /** 表的元组ID */
    Id?: string
    /** 空间ID */
    spaceId?: string
    /** 空间角色：viewer-查看者 editor-编辑者 admin-管理员 */
    spaceRole?: string
    /** 用户ID */
    userId?: string
  }

  type SpaceUserRemoveRequest = {
    /** 表的元组ID */
    Id?: string
  }

  type SpaceUserVO = {
    createTime?: string
    id?: string
    /** 空间信息 */
    space?: SpaceVO
    spaceId?: string
    spaceRole?: string
    updateTime?: string
    /** 用户信息 */
    user?: UserVO
    userId?: string
  }

  type SpaceVO = {
    createTime?: string
    editTime?: string
    /** Space ID */
    id?: string
    maxCount?: number
    maxSize?: number
    /** 空间的权限列表 */
    permissionList?: string[]
    spaceLevel?: number
    spaceName?: string
    /** Space type: 0 - 私人空间, 1 - 团队空间 */
    spaceType?: number
    totalCount?: number
    totalSize?: number
    updateTime?: string
    user?: UserVO
    /** User ID */
    userId?: string
  }

  type TaskDetailOutput = {
    /** 错误码（失败时返回） */
    code?: string
    /** 任务完成时间（成功或失败时返回） */
    endTime?: string
    /** 错误信息（失败时返回） */
    message?: string
    /** 输出图像的 URL（成功时返回） */
    outputImageUrl?: string
    /** 任务调度时间（成功或失败时返回） */
    scheduledTime?: string
    /** 任务提交时间（成功或失败时返回） */
    submitTime?: string
    /** 任务的唯一标识符 */
    taskId?: string
    /** 任务结果统计（执行中时返回） */
    taskMetrics?: TaskMetrics
    /** 任务状态：PENDING、RUNNING、SUSPENDED、SUCCEEDED、FAILED、UNKNOWN */
    taskStatus?: string
  }

  type TaskMetrics = {
    /** 失败任务数 */
    failed?: number
    /** 成功任务数 */
    succeeded?: number
    /** 总任务数 */
    total?: number
  }

  type Usage = {
    /** 生成的图片数量 */
    imageCount?: number
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

  type UserEditRequest = {
    /** 用户ID */
    id?: string
    /** 用户昵称 */
    userName?: string
    /** 用户简介 */
    userProfile?: string
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

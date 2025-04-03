// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 根据ID软删除图片「登录校验」 POST /v1/picture/delete */
export async function postPictureOpenApiDelete(
  body: API.DeleteRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: boolean }>('/v1/picture/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 编辑图片 若图片不存在，则返回false POST /v1/picture/edit */
export async function postPictureEdit(
  body: API.PictureEditRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: boolean }>('/v1/picture/edit', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 批量更新图片请求「登录校验」 POST /v1/picture/edit/batch */
export async function postPictureEditBatch(
  body: API.PictureEditByBatchRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: boolean }>('/v1/picture/edit/batch', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 根据ID获取图片「管理员」 GET /v1/picture/get */
export async function getPictureGet(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getPictureGetParams,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.Picture }>('/v1/picture/get', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 根据ID获取脱敏的图片 GET /v1/picture/get/vo */
export async function getPictureGetVo(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getPictureGetVoParams,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.PictureVO }>('/v1/picture/get/vo', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 分页获取一系列图片信息「管理员」 POST /v1/picture/list/page */
export async function postPictureListPage(
  body: API.PictureQueryRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.ListPictureResponse }>('/v1/picture/list/page', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 分页获取一系列图片信息 POST /v1/picture/list/page/vo */
export async function postPictureListPageVo(
  body: API.PictureQueryRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.ListPictureVOResponse }>('/v1/picture/list/page/vo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 带有缓存的分页获取一系列图片信息 POST /v1/picture/list/page/vo/cache */
export async function postPictureListPageVoCache(
  body: API.PictureQueryRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.ListPictureVOResponse }>(
    '/v1/picture/list/page/vo/cache',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      data: body,
      ...(options || {}),
    }
  )
}

/** 获取AI扩图任务信息「登录校验」 GET /v1/picture/out_painting/create_task */
export async function getPictureOutPaintingCreateTask(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getPictureOutPaintingCreateTaskParams,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.GetOutPaintingResponse }>(
    '/v1/picture/out_painting/create_task',
    {
      method: 'GET',
      params: {
        ...params,
      },
      ...(options || {}),
    }
  )
}

/** 创建AI扩图任务请求「登录校验」 POST /v1/picture/out_painting/create_task */
export async function postPictureOutPaintingCreateTask(
  body: API.CreateOutPaintingTaskRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.CreateOutPaintingTaskResponse }>(
    '/v1/picture/out_painting/create_task',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      data: body,
      ...(options || {}),
    }
  )
}

/** 执行图片审核「管理员」 POST /v1/picture/review */
export async function postPictureReview(
  body: API.PictureReviewRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: boolean }>('/v1/picture/review', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 根据图片的颜色搜索相似图片「登录校验」 POST /v1/picture/search/color */
export async function postPictureSearchColor(
  body: API.PictureSearchByColorRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.PictureVO[] }>('/v1/picture/search/color', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 根据图片ID搜索图片 POST /v1/picture/search/picture */
export async function postPictureSearchPicture(
  body: API.PictureSearchByPictureRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.ImageSearchResult[] }>('/v1/picture/search/picture', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取图片的标签和分类（固定） GET /v1/picture/tag_category */
export async function getPictureTagCategory(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.PictureTagCategory }>('/v1/picture/tag_category', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 更新图片「登录校验」 若图片不存在，则返回false POST /v1/picture/update */
export async function postPictureUpdate(
  body: API.PictureUpdateRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: boolean }>('/v1/picture/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 上传图片接口「需要登录校验」 根据是否存在ID来上传图片或者修改图片信息，返回图片信息视图 POST /v1/picture/upload */
export async function postPictureUpload(
  body: {
    /** 图片的ID，非必需 */
    id?: string
    /** 图片的上传空间ID，非必需 */
    spaceId?: string
  },
  file?: File,
  options?: { [key: string]: any }
) {
  const formData = new FormData()

  if (file) {
    formData.append('file', file)
  }

  Object.keys(body).forEach((ele) => {
    const item = (body as any)[ele]

    if (item !== undefined && item !== null) {
      if (typeof item === 'object' && !(item instanceof File)) {
        if (item instanceof Array) {
          item.forEach((f) => formData.append(ele, f || ''))
        } else {
          formData.append(ele, JSON.stringify(item))
        }
      } else {
        formData.append(ele, item)
      }
    }
  })

  return request<API.Response & { data?: API.PictureVO }>('/v1/picture/upload', {
    method: 'POST',
    data: formData,
    requestType: 'form',
    ...(options || {}),
  })
}

/** 批量抓取图片「管理员」 POST /v1/picture/upload/batch */
export async function postPictureUploadBatch(
  body: API.PictureUploadByBatchRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: number }>('/v1/picture/upload/batch', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 根据URL上传图片接口「需要登录校验」 POST /v1/picture/upload/url */
export async function postPictureUploadUrl(
  body: API.PictureUploadRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.PictureVO }>('/v1/picture/upload/url', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

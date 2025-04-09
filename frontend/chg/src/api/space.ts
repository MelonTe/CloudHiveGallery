// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 增加空间「需要登录」 POST /v1/space/add */
export async function postSpaceAdd(body: API.SpaceAddRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: string }>('/v1/space/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 编辑空间昵称 若空间不存在，则返回false POST /v1/space/edit */
export async function postSpaceEdit(body: API.SpaceEditRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: boolean }>('/v1/space/edit', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取当个空间的视图信息「登录校验」 GET /v1/space/get/vo */
export async function getSpaceGetVo(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getSpaceGetVoParams,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.SpaceVO }>('/v1/space/get/vo', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 获取所有的空间等级信息 GET /v1/space/list/level */
export async function getSpaceListLevel(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.SpaceLevelResponse[] }>('/v1/space/list/level', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 分页获取一系列空间信息「管理员」 POST /v1/space/list/page */
export async function postSpaceListPage(
  body: API.SpaceQueryRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.ListSpaceResponse }>('/v1/space/list/page', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 分页获取一系列空间视图信息 POST /v1/space/list/page/vo */
export async function postSpaceListPageVo(
  body: API.SpaceQueryRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.ListSpaceVOResponse }>('/v1/space/list/page/vo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 更新空间「管理员」 若空间不存在，则返回false POST /v1/space/update */
export async function postSpaceUpdate(
  body: API.SpaceUpdateRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: boolean }>('/v1/space/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 增加成员到空间 POST /v1/spaceUser/add */
export async function postSpaceUserAdd(
  body: API.SpaceUserAddRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: string }>('/v1/spaceUser/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 从空间移除成员 POST /v1/spaceUser/delete */
export async function postSpaceUserOpenApiDelete(
  body: API.SpaceUserRemoveRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: boolean }>('/v1/spaceUser/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 编辑成员权限 POST /v1/spaceUser/edit */
export async function postSpaceUserEdit(
  body: API.SpaceUserEditRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: boolean }>('/v1/spaceUser/edit', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 查询某个成员在某个空间的信息 POST /v1/spaceUser/get */
export async function postSpaceUserGet(
  body: API.SpaceUserQueryRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.SpaceUser }>('/v1/spaceUser/get', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 查询成员信息列表 POST /v1/spaceUser/list */
export async function postSpaceUserList(
  body: API.SpaceUserQueryRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.SpaceUserVO[] }>('/v1/spaceUser/list', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 查询我加入的团队空间列表 POST /v1/spaceUser/list/my */
export async function postSpaceUserListMy(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.SpaceUserVO[] }>('/v1/spaceUser/list/my', {
    method: 'POST',
    ...(options || {}),
  })
}

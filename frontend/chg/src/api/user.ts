// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 创建一个用户「管理员」 默认密码为12345678 POST /v1/user/add */
export async function postUserAdd(body: API.UserAddRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: string }>('/v1/user/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 上传用户头像「需要登录校验」 根据ID，将头像保存到数据库，返回是否成功 POST /v1/user/avatar */
export async function postUserAvatar(body: {}, file?: File, options?: { [key: string]: any }) {
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

  return request<API.Response & { data?: boolean }>('/v1/user/avatar', {
    method: 'POST',
    data: formData,
    requestType: 'form',
    ...(options || {}),
  })
}

/** 根据ID软删除用户「管理员」 POST /v1/user/delete */
export async function postUserOpenApiDelete(
  body: API.DeleteRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: boolean }>('/v1/user/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 更新用户个人资料 若用户不存在，则返回失败 POST /v1/user/edit */
export async function postUserEdit(body: API.UserEditRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: boolean }>('/v1/user/edit', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 根据ID获取用户「管理员」 GET /v1/user/get */
export async function getUserGet(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getUserGetParams,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.User }>('/v1/user/get', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 获取登录的用户信息 GET /v1/user/get/login */
export async function getUserGetLogin(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.UserLoginVO }>('/v1/user/get/login', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 根据ID获取简略信息用户 GET /v1/user/get/vo */
export async function getUserGetVo(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getUserGetVoParams,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.UserVO }>('/v1/user/get/vo', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 分页获取一系列用户信息「管理员」 根据用户关键信息进行模糊查询 POST /v1/user/list/page/vo */
export async function postUserListPageVo(
  body: API.UserQueryRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.ListUserVOResponse }>('/v1/user/list/page/vo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 用户登录 根据账号密码进行登录 POST /v1/user/login */
export async function postUserLogin(body: API.UserLoginRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.UserLoginVO }>('/v1/user/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 执行用户注销（退出） POST /v1/user/logout */
export async function postUserLogout(options?: { [key: string]: any }) {
  return request<API.Response & { data?: boolean }>('/v1/user/logout', {
    method: 'POST',
    ...(options || {}),
  })
}

/** 注册用户 根据账号密码进行注册 POST /v1/user/register */
export async function postUserRegister(
  body: API.UserRegsiterRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: string }>('/v1/user/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 更新用户信息「管理员」 若用户不存在，则返回失败 POST /v1/user/update */
export async function postUserUpdate(
  body: API.UserUpdateRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: boolean }>('/v1/user/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

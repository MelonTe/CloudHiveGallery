// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 获取空间图片分类分析「登录校验」 POST /v1/space/analyze/category */
export async function postSpaceAnalyzeCategory(
  body: API.SpaceCategoryAnalyzeRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.SpaceCategoryAnalyzeResponse[] }>(
    '/v1/space/analyze/category',
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

/** 获取空间使用情况排名「管理员」 POST /v1/space/analyze/rank */
export async function postSpaceAnalyzeRank(
  body: API.SpaceRankAnalyzeRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.Space[] }>('/v1/space/analyze/rank', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取空间图片大小范围统计分析「登录校验」 POST /v1/space/analyze/size */
export async function postSpaceAnalyzeSize(
  body: API.SpaceSizeAnalyzeRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.SpaceSizeAnalyzeResponse[] }>(
    '/v1/space/analyze/size',
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

/** 获取空间标签出现量分析「登录校验」 POST /v1/space/analyze/tag */
export async function postSpaceAnalyzeTag(
  body: API.SpaceTagAnalyzeRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.SpaceTagAnalyzeResponse[] }>('/v1/space/analyze/tag', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取空间使用分析「登录校验」 POST /v1/space/analyze/usage */
export async function postSpaceAnalyzeUsage(
  body: API.SpaceUsageAnalyzeRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.SpaceUsageAnalyzeResponse }>(
    '/v1/space/analyze/usage',
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

/** 获取用户上传图片统计分析，支持分析特定用户「登录校验」 POST /v1/space/analyze/user */
export async function postSpaceAnalyzeUser(
  body: API.SpaceUserAnalyzeRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.SpaceUserAnalyzeResponse[] }>(
    '/v1/space/analyze/user',
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

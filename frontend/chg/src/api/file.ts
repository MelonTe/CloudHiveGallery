// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 测试文件下载接口「管理员」 GET /v1/file/test/download */
export async function getFileTestDownload(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getFileTestDownloadParams,
  options?: { [key: string]: any }
) {
  return request<string>('/v1/file/test/download', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 测试文件上传接口「管理员」 POST /v1/file/test/upload */
export async function postFileTestUpload(body: {}, file?: File, options?: { [key: string]: any }) {
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

  return request<API.Response & { data?: string }>('/v1/file/test/upload', {
    method: 'POST',
    data: formData,
    requestType: 'form',
    ...(options || {}),
  })
}

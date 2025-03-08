// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** Say Hello World This is a simple API endpoint that returns a "Hello, World!" message. GET /v1/hello */
export async function getV1Hello(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/v1/hello', {
    method: 'GET',
    ...(options || {}),
  })
}

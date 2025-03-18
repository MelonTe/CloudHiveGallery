import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { getUserGetLogin } from '@/api/user'
/**
 * 存储登录用户信息的状态
 */
export const useLoginUserStore = defineStore('loginUser', () => {
  const loginUser = ref<API.UserLoginVO>({
    userName: '未登录',
  })

  /**
   * 远程获取登录用户信息
   */
  async function fetchLoginUser() {
    //todo : 获取登录用户信息
    const res = await getUserGetLogin()
    /* 响应码为0并且正常响应 */
    if (res.data.code === 0 && res.data.data) {
      loginUser.value = res.data.data
    }
    /* setTimeout(() => {
      {
        loginUser.value = { userName: '测试用户', id: 1 }
      }
    }, 3000) */
  }

  /**
   * 设置登录用户
   * @param newLoginUser
   */
  function setLoginUser(newLoginUser: any) {
    loginUser.value = newLoginUser
  }

  // 返回
  return { loginUser, fetchLoginUser, setLoginUser }
})

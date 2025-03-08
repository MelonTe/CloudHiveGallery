import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
/**
 * 存储登录用户信息的状态
 */
export const useLoginUserStore = defineStore('loginUser', () => {
  const loginUser = ref<any>({
    userName: '未登录',
  })

  /**
   * 远程获取登录用户信息
   */
  async function fetchLoginUser() {
    //todo : 获取登录用户信息
    //测试，3s之后自动登录(假登入)
    setTimeout(() => {
      {
        loginUser.value = { userName: '测试用户', id: 1 }
      }
    }, 3000)
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

<!-- 主页 -->
<template>
  <div id="userRegisterPage">
    <h2 class="title">云巢画廊 - 用户注册</h2>
    <div class="desc">高效协同画廊</div>
    <a-form :model="formState" name="basic" autocomplete="off" @finish="handleSubmit">
      <a-form-item name="userAccount" :rules="[{ required: true, message: '请输入账号' }]">
        <a-input v-model:value="formState.userAccount" placeholder="请输入账号" />
      </a-form-item>
      <a-form-item
        name="userPassword"
        :rules="[
          { required: true, message: '请输入密码' },
          { min: 8, message: '密码长度不能小于 8 位' },
        ]"
      >
        <a-input-password v-model:value="formState.userPassword" placeholder="请输入密码" />
      </a-form-item>
      <a-form-item
        name="checkPassword"
        :rules="[
          { required: true, message: '请输入确认密码' },
          { min: 8, message: '确认密码长度不能小于 8 位' },
        ]"
      >
        <a-input-password v-model:value="formState.checkPassword" placeholder="请输入密码" />
      </a-form-item>
      <div class="tips">
        <!-- 引导跳转到注册页面 -->
        已有账号？
        <RouterLink to="/user/login">去登录</RouterLink>
      </div>
      <a-form-item>
        <a-button type="primary" html-type="submit" style="width: 100%">注册</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
import { postUserRegister } from '@/api/user'
import router from '@/router'
import { useLoginUserStore } from '@/stores/useLoginUserStore'
import { message } from 'ant-design-vue'
import { reactive } from 'vue'

interface FormState {
  username: string
  password: string
  remember: boolean
}

const loginUserStore = useLoginUserStore()
//获取，供全局使用
loginUserStore.fetchLoginUser()

const formState = reactive<API.UserRegsiterRequest>({
  userAccount: '',
  userPassword: '',
  checkPassword: '',
})
/* 提交表单 */
const handleSubmit = async (values: any) => {
  //校验密码是否一致
  if (values.userPassword !== values.checkPassword) {
    message.error('两次密码不一致')
    return
  }
  /* 传入表单项 */
  const res = await postUserRegister(values)
  if (res.data.code === 0 && res.data.data) {
    //注册成功，跳转到登录页面
    message.success('注册成功')
    /* 跳转登录页 */
    router.push({
      path: '/user/login',
      replace: true /* 覆盖掉注册页 */,
    })
  } else {
    message.error('注册失败，' + res.data.msg)
  }
}
</script>

<style scoped>
#userRegisterPage {
  max-width: 360px;
  /* 宽度 */
  margin: 0 auto;
  margin-top: 10%;
  /* 居中 */
}

.title {
  text-align: center;
  margin-bottom: 16px;
}

.desc {
  text-align: center;
  color: #bbb;
  margin-bottom: 16px;
}

.tips {
  color: #bbb;
  text-align: right;
  font-size: 13px;
  margin-bottom: 16px;
}
</style>

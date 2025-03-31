<!-- 主页 -->
<template>
    <div id="userLoginPage">
        <h2 class="title">云巢画廊 - 用户登录</h2>
        <div class="desc">高效协同画廊</div>
        <a-form :model="formState" name="basic" autocomplete="off" @finish="handleSubmit">
            <a-form-item name="userAccount" :rules="[{ required: true, message: '请输入账号' }]">
                <a-input v-model:value="formState.userAccount" placeholder="请输入账号" />
            </a-form-item>
            <a-form-item name="userPassword" :rules="[
                { required: true, message: '请输入密码' },
                { min: 8, message: '密码长度不能小于 8 位' },
            ]">
                <a-input-password v-model:value="formState.userPassword" placeholder="请输入密码" />
            </a-form-item>
            <div class="tips"> <!-- 引导跳转到注册页面 -->
                没有账号？
                <RouterLink to="/user/register">去注册</RouterLink>
            </div>
            <a-form-item>
                <a-button type="primary" html-type="submit" style="width: 100%">登录</a-button>
            </a-form-item>
        </a-form>
    </div>
</template>

<script lang="ts" setup>
import { postUserLogin } from '@/api/user'
import router from '@/router';
import { useLoginUserStore } from '@/stores/useLoginUserStore';
import { message } from 'ant-design-vue';
import { reactive } from 'vue';

interface FormState {
    username: string;
    password: string;
    remember: boolean;
}

const loginUserStore = useLoginUserStore()
//获取，供全局使用
loginUserStore.fetchLoginUser()

const formState = reactive<API.UserLoginRequest>({
    userAccount: '',
    userPassword: '',
});
/* 提交表单 */
const handleSubmit = async (values: any) => {
    /* 传入表单项 */
    const res = await postUserLogin(values);
    if (res.data.code === 0 && res.data.data) {
        /* 把登录态保存到全局状态中 */
        await loginUserStore.fetchLoginUser()
        message.success('登录成功')
        /* 跳转回主页 */
        router.push({
            path: "/",
            replace: true /* 覆盖掉登录页 */
        })
    } else {
        message.error('登录失败，' + res.data.msg)
    }
};

</script>

<style scoped>
#userLoginPage {
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

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
    width: 100%;
    max-width: 360px;
    margin: 0 auto;
    margin-top: 8%;
    background: rgba(255, 255, 255, 0.9);
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

#userLoginPage :deep(.ant-input-affix-wrapper) {
    border-radius: 8px;
    border: 1px solid #d9d9d9;
    background: transparent;
    height: 40px;
    transition: all 0.3s;
}

#userLoginPage :deep(.ant-input-affix-wrapper:hover) {
    border-color: #4CAF50;
}

#userLoginPage :deep(.ant-input-affix-wrapper-focused) {
    border-color: #4CAF50;
    box-shadow: 0 0 0 2px rgba(76, 175, 80, 0.2);
}

#userLoginPage :deep(.ant-input) {
    background: transparent;
    border: none;
    height: 38px;
}

#userLoginPage :deep(.ant-input:focus) {
    box-shadow: none;
}

#userLoginPage :deep(.ant-input-password) {
    background: transparent;
    border: none;
}

#userLoginPage :deep(.ant-input-password .ant-input) {
    background: transparent;
    border: none;
    height: 38px;
}

#userLoginPage :deep(.ant-input-password .ant-input-suffix) {
    background: transparent;
}

#userLoginPage :deep(.ant-input-password .ant-input-suffix .anticon) {
    color: #4CAF50;
    font-size: 16px;
}

.title {
    text-align: center;
    margin-bottom: 16px;
    font-size: 24px;
    background: linear-gradient(45deg, #43b16a, #64d487);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
}

.desc {
    text-align: center;
    color: #888;
    margin-bottom: 24px;
}

.tips {
    color: #888;
    text-align: right;
    font-size: 13px;
    margin-bottom: 16px;
}

:deep(.ant-btn-primary) {
    height: 44px;
    background: linear-gradient(135deg, #64d487, #43b16a);
    border: none;
    box-shadow: 0 4px 8px rgba(100, 212, 135, 0.2);
    transition: all 0.3s;
}

:deep(.ant-btn-primary:hover) {
    transform: translateY(-2px);
    box-shadow: 0 6px 12px rgba(100, 212, 135, 0.25);
    background: linear-gradient(135deg, #5bc77b, #3ba05f);
}
</style>

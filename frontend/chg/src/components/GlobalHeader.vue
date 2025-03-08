<template>
    <div id="globalHeader">
        <a-row :wrap="false"> <!-- 关闭自动换行 -->
            <a-col flex="200px">
                <!-- 固定大小的内容 -->
                <!-- 第一列为网站图标 -->
                <router-link to="/">
                    <div class="title-bar">
                        <img class="logo" src="../assets/logo.png" alt="logo" />
                        <div class="title">云巢画廊</div>
                    </div>
                </router-link>
            </a-col>
            <a-col flex="auto">
                <!-- 菜单列 -->
                <a-menu v-model:selectedKeys="current" mode="horizontal" :items="items" @click="doMenuClick" />
            </a-col>
            <a-col flex="120px">
                <div class="user-login-status">
                    <div v-if="loginUserStore.loginUser.id">
                        {{ loginUserStore.loginUser.userName ?? '无名' }}
                    </div>
                    <div v-else>
                        <a-button type="primary" href="/user/login">登录</a-button>
                    </div>
                </div>
            </a-col>
        </a-row>



    </div>
</template>
<script lang="ts" setup>
import { h, ref } from 'vue';
import { MailOutlined, AppstoreOutlined, SettingOutlined, HomeOutlined } from '@ant-design/icons-vue';
import type { MenuProps } from 'ant-design-vue';
import { useRouter } from 'vue-router';
import { useLoginUserStore } from '../stores/useLoginUserStore';
const loginUserStore = useLoginUserStore()
const items = ref<MenuProps['items']>([
    {
        key: '/',
        icon: () => h(HomeOutlined),
        label: '主页',
        title: '主页',
    },
    {
        key: '/about',
        icon: () => h(AppstoreOutlined),
        label: '关于',
        title: '关于',
    },
]);

const router = useRouter();
/* 路由跳转事件 */
const doMenuClick = ({ key }) => {
    /* 跳转到key的页面 */
    router.push({
        path: key
    })
}

/* current决定菜单项高亮 */
const current = ref<string[]>([]);
/* 钩子函数，每次跳转到新页面都会执行 */
router.afterEach((to, from, next) => {
    /* 把渲染current的值，改成url中的地址，表现为在哪个路由里，menu中的选型标记为选中 */
    current.value = [to.path]
})
</script>

<style scoped>
.title {
    flex-grow: 1;
    flex-shrink: 0;
    color: black;
    font-size: 18px;
    margin-left: 16px;
    /* 离左边logo远一点 */
}

.logo {
    height: 40px;
    width: 40px;
    display: inline-block;
    flex-shrink: 0;
}

#globalHeader .title-bar {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    /* 确保元素按顺序排列 */
}
</style>
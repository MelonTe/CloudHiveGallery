<template>
    <div id="globalHeader">
        <router-link to="/">
            <div id="title-bar">
                <img class="logo" src="../assets/logo.png" alt="logo">
                <div class="title">云巢画廊</div>
            </div>
        </router-link>

        <a-menu v-model:selectedKeys="current" mode="horizontal" :items="items" />
    </div>
</template>
<script lang="ts" setup>
import { h, ref } from 'vue';
import { MailOutlined, AppstoreOutlined, SettingOutlined, HomeOutlined } from '@ant-design/icons-vue';
import type { MenuProps } from 'ant-design-vue';
const current = ref<string[]>(['mail']);
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
</script>
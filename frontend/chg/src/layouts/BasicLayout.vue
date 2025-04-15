<template>
    <a-layout class="basic-layout">
        <GlobalHeader />
        <a-layout>
            <GlobalSider />
            <a-layout-content class="main-content">
                <router-view />
            </a-layout-content>
        </a-layout>
    </a-layout>
</template>

<script setup lang="ts">
import GlobalHeader from '../components/GlobalHeader.vue'
import GlobalSider from '../components/GlobalSider.vue'
</script>

<style scoped>
.basic-layout {
    min-height: 100vh;
}

.main-content {
    padding: 24px;
    background: #f5f5f5;
    min-height: calc(100vh - 64px);
    transition: margin-left 0.3s ease;
}

/* 侧边栏和主内容区域的关系 */
:deep(.ant-layout-has-sider) {
    display: flex;
    flex-direction: row;
    min-height: calc(100vh - 64px);
}

/* 侧边栏展开状态 */
:deep(.ant-layout-sider:not(.ant-layout-sider-collapsed))+.main-content {
    margin-left: 0;
    width: calc(100% - 200px);
}

/* 侧边栏折叠状态 */
:deep(.ant-layout-sider-collapsed)+.main-content {
    margin-left: 0;
    width: calc(100% - 80px);
}

@media (max-width: 992px) {
    .main-content {
        margin-left: 80px;
    }
}

@media (max-width: 576px) {
    .main-content {
        margin-left: 0;
    }
}
</style>
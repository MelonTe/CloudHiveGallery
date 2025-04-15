<template>
  <div id="basicLayout">
    <a-layout style="min-height: 100vh;" class="main-layout" :class="{ 'not-logged-in': !isLoggedIn }">
      <a-layout-header class="header glass-effect">
        <GlobalHeader />
      </a-layout-header>

      <!-- 登录状态下的布局 -->
      <a-layout class="main-content-layout" v-if="isLoggedIn">
        <GlobalSider class="sider" />
        <a-layout-content class="content">
          <div class="content-inner">
            <router-view />
          </div>
        </a-layout-content>
      </a-layout>

      <!-- 未登录状态下的简化布局 -->
      <a-layout-content class="content full-width" v-else>
        <div class="content-inner">
          <router-view />
        </div>
      </a-layout-content>

      <a-layout-footer class="footer glass-effect">
        云巢画廊 made By：MelonTe.粤ICP备2024329874号-2
      </a-layout-footer>
    </a-layout>
  </div>
</template>

<script setup lang="ts">
import GlobalHeader from '@/components/GlobalHeader.vue';
import GlobalSider from '@/components/GlobalSider.vue'
import { useLoginUserStore } from '../stores/useLoginUserStore';
import { computed } from 'vue';

const loginUserStore = useLoginUserStore();
const isLoggedIn = computed(() => !!loginUserStore.loginUser.id);

</script>

<style scoped>
#basicLayout .main-layout {
  background: #ffffff;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

#basicLayout .header {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  padding: 0 20px;
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.03);
  height: 64px;
  line-height: 64px;
  z-index: 1000;
  position: sticky;
  top: 0;
}

#basicLayout .main-content-layout {
  flex: 1;
  display: flex;
  min-height: calc(100vh - 130px);
}

#basicLayout .sider {
  background: rgba(255, 255, 255, 0.9);
  padding-top: 16px;
  border-right: 1px solid rgba(240, 240, 240, 0.5);
  transition: all 0.3s;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  position: sticky;
  top: 64px;
  height: calc(100vh - 130px);
  overflow-y: auto;
  z-index: 900;
}

#basicLayout .content {
  flex: 1;
  padding: 24px;
  background: transparent;
  margin-bottom: 50px;
  overflow-y: auto;
}

/* 未登录状态下全宽内容 */
#basicLayout .content.full-width {
  width: 100%;
  padding: 24px 16px;
}

#basicLayout .content-inner {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.02);
  margin: 16px;
  min-height: calc(100vh - 182px);
  padding: 24px;
}

/* 未登录状态下移除毛玻璃效果 */
#basicLayout .not-logged-in .content-inner {
  backdrop-filter: none;
  -webkit-backdrop-filter: none;
  background: rgba(255, 255, 255, 0.95);
}

#basicLayout .footer {
  background: rgba(245, 245, 245, 0.9);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  padding: 16px;
  text-align: center;
  box-shadow: 0 -1px 4px rgba(0, 0, 0, 0.02);
  z-index: 1000;
}

#basicLayout :deep(.ant-menu-root) {
  border-bottom: none !important;
  border-inline-end: none !important;
  background: transparent;
}

/* 加入新的悬浮效果 */
#basicLayout :deep(.ant-menu-item) {
  border-radius: 4px;
  margin: 4px 8px;
  transition: all 0.3s;
}

#basicLayout :deep(.ant-menu-item:hover) {
  background-color: rgba(76, 175, 80, 0.15);
}

#basicLayout :deep(.ant-menu-item-selected) {
  background-color: rgba(76, 175, 80, 0.85) !important;
  color: white !important;
  font-weight: 500;
}
</style>

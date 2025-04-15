<template>
    <div id="basicLayout">
        <a-layout style="min-height: 100vh;">
            <a-layout-header class="header">
                <GlobalHeader />
            </a-layout-header>
          <a-layout>
            <GlobalSider class="sider" />
            <a-layout-content class="content">
              <router-view />
            </a-layout-content>
          </a-layout>
          <a-layout-footer class="footer">
                云巢画廊 made By：MelonTe.粤ICP备2024329874号-2
            </a-layout-footer>
        </a-layout>
    </div>
</template>

<script setup lang="ts">
import GlobalHeader from '@/components/GlobalHeader.vue';
import GlobalSider from '@/components/GlobalSider.vue'


</script>

<style scoped>
#basicLayout .content {
    padding: 20px;
    background: linear-gradient(to right, #fcfcfc, #ffffff);
    margin-bottom: 30px;
}

#basicLayout .footer {
    background: #efefef;
    padding: 16px;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    text-align: center;
}

#basicLayout .header {
  margin-bottom: 1px;
  background: white;
  color: unset;
  padding-inline: 20px;
}

#basicLayout .content {
  padding: 28px;
}

#basicLayout .sider {
  background: #fff;
  padding-top: 20px;
  border-right: 0.5px solid #eee;
}

#basicLayout :deep(.ant-menu-root) {
  border-bottom: none !important;
  border-inline-end: none !important;
}

</style>

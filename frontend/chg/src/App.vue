<template>
  <div id="app">
    <BasicLayOut></BasicLayOut>
  </div>
</template>



<script setup lang="ts">
import BasicLayOut from './layout/BasicLayOut.vue';
import { useLoginUserStore } from './stores/useLoginUserStore';
const loginUserStore = useLoginUserStore()
//获取，供全局使用
loginUserStore.fetchLoginUser()
</script>

<style scoped></style>

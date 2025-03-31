<template>
  <div id="globalHeader">
    <a-row :wrap="false">
      <!-- 关闭自动换行 -->
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
        <a-menu
          v-model:selectedKeys="current"
          mode="horizontal"
          :items="items"
          @click="doMenuClick"
        />
      </a-col>
      <!-- 用户信息展示 -->
      <a-col flex="120px">
        <div class="user-login-status">
          <div v-if="loginUserStore.loginUser.id">
            <a-dropdown :trigger="['click']">
              <a-space>
                <a-avatar :src="loginUserStore.loginUser.userAvatar" />
                {{ loginUserStore.loginUser.userName ?? '无名' }}
              </a-space>
              <a class="ant-dropdown-link" @click.prevent>
                <DownOutlined />
              </a>
              <!-- 插槽 -->
              <template #overlay>
                <a-menu>
                  <a-menu-item key="0" @click="doLogout">
                    <icon-font type="icon-dengchu" />
                    退出登录
                  </a-menu-item>
                  <a-menu-item>
                    <router-link to="/my_space">
                      <UserOutlined />
                      我的空间
                    </router-link>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
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
import { computed, h, ref } from 'vue'
import {
  HomeOutlined,
} from '@ant-design/icons-vue'
import { message, type MenuProps } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { useLoginUserStore } from '../stores/useLoginUserStore'
const loginUserStore = useLoginUserStore()

//未经处理的原始菜单
const originItmes = [
  {
    key: '/',
    icon: () => h(HomeOutlined),
    label: '主页',
    title: '主页',
  },
  {
    key: '/admin/userManage',
    label: '用户管理',
    title: '用户管理',
  },
  {
    key: '/admin/pictureManage',
    label: '图片管理',
    title: '图片管理',
  },
  {
    key: '/add_picture',
    label: '创建图片',
    title: '创建图片',
  },
  {
    key: '/admin/spaceManage',
    label: '空间管理',
    title: '空间管理',
  },
]
// 根据权限过滤菜单项
const filterMenus = (menus = [] as MenuProps['items']) => {
  return menus?.filter((menu) => {
    // 管理员才能看到 /admin 开头的菜单
    if (typeof menu?.key === 'string' && menu.key.startsWith('/admin')) {
      const loginUser = loginUserStore.loginUser
      if (!loginUser || loginUser.userRole !== 'admin') {
        return false
      }
    }
    return true
  })
}
//过滤后的菜单
const items = computed(() => {
  return filterMenus(originItmes)
})
const router = useRouter()
/* 路由跳转事件 */
const doMenuClick = ({ key }) => {
  /* 跳转到key的页面 */
  router.push({
    path: key,
  })
}

/* current决定菜单项高亮 */
const current = ref<string[]>([])
/* 钩子函数，每次跳转到新页面都会执行 */
router.afterEach((to, from, next) => {
  /* 把渲染current的值，改成url中的地址，表现为在哪个路由里，menu中的选型标记为选中 */
  current.value = [to.path]
})

/* 头像下拉菜单 */
import { DownOutlined } from '@ant-design/icons-vue'

import { createFromIconfontCN } from '@ant-design/icons-vue'
import { postUserLogout } from '@/api/user'

/* 项目图标导入 */
const IconFont = createFromIconfontCN({
  scriptUrl: '//at.alicdn.com/t/c/font_4855251_hyb7n0qsrh7.js',
})

/* 注销 */
const doLogout = async () => {
  const res = await postUserLogout()
  if (res.data.code === 0) {
    /* 重置未登录 */
    loginUserStore.setLoginUser({
      userName: '未登录',
    })
    message.success('登出成功')
    router.push({
      path: '/user/login',
    })
  } else {
    message.error('登出失败，' + res.data.msg)
  }
}
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

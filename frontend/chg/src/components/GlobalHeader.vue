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
        <div class="menu-container">
          <a-menu v-model:selectedKeys="current" mode="horizontal" :items="items" @click="doMenuClick"
            class="custom-menu" />
          <!-- 跟随线条 -->
          <div class="follow-line" :style="followLineStyle"></div>
        </div>
      </a-col>
      <!-- 用户信息展示 -->
      <a-col flex="120px">
        <div class="user-login-status">
          <div v-if="loginUserStore.loginUser.id">
            <a-dropdown :trigger="['click']">
              <a-space class="user-info">
                <a-avatar :src="loginUserStore.loginUser.userAvatar" class="user-avatar" />
                <span class="user-name">{{ loginUserStore.loginUser.userName ?? '无名' }}</span>
              </a-space>
              <a class="ant-dropdown-link" @click.prevent>
                <DownOutlined />
              </a>
              <!-- 插槽 -->
              <template #overlay>
                <a-menu class="user-dropdown glass-effect">
                  <a-menu-item>
                    <router-link to="/profile">
                      <UserOutlined />
                      我的主页
                    </router-link>
                  </a-menu-item>
                  <a-menu-item>
                    <router-link to="/my_space">
                      <UserOutlined />
                      我的空间
                    </router-link>
                  </a-menu-item>
                  <a-menu-item key="0" @click="doLogout">
                    <icon-font type="icon-dengchu" />
                    退出登录
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
          <div v-else>
            <a-button type="primary" href="/user/login" class="login-btn">登录</a-button>
          </div>
        </div>
      </a-col>
    </a-row>
  </div>
</template>
<script lang="ts" setup>
import { computed, h, onMounted, onUnmounted, ref } from 'vue'
import {
  HomeOutlined,
  UserOutlined
} from '@ant-design/icons-vue'
import { message, type MenuProps } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { useLoginUserStore } from '../stores/useLoginUserStore'
const loginUserStore = useLoginUserStore()

// 跟随线条样式
const followLineStyle = ref({
  width: '0px',
  left: '0px',
  opacity: 0
})

// 未经处理的原始菜单
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
const doMenuClick = (menuInfo: any) => {
  /* 跳转到key的页面 */
  router.push({
    path: menuInfo.key,
  })
}

/* current决定菜单项高亮 */
const current = ref<string[]>([])
/* 钩子函数，每次跳转到新页面都会执行 */
router.afterEach((to, from, next) => {
  /* 把渲染current的值，改成url中的地址，表现为在哪个路由里，menu中的选型标记为选中 */
  current.value = [to.path]
  // 当路由变化时，更新跟随线条位置
  setTimeout(() => {
    updateFollowLine(document.querySelector('.ant-menu-item-selected'))
  }, 100)
})

/* 头像下拉菜单 */
import { DownOutlined } from '@ant-design/icons-vue'

import { createFromIconfontCN } from '@ant-design/icons-vue'
import { postUserLogout } from '../api/user'

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

// 设置跟随线条位置
const updateFollowLine = (el: Element | null) => {
  if (!el) return

  const rect = el.getBoundingClientRect()
  const menuRect = (document.querySelector('.custom-menu') as Element).getBoundingClientRect()

  followLineStyle.value = {
    width: rect.width + 'px',
    left: (rect.left - menuRect.left) + 'px',
    opacity: 1
  }
}

// 设置菜单项鼠标移入移出事件
const setupMenuEvents = () => {
  const menuItems = document.querySelectorAll('.custom-menu .ant-menu-item')

  menuItems.forEach(item => {
    item.addEventListener('mouseenter', () => {
      updateFollowLine(item)
    })
  })

  const menu = document.querySelector('.custom-menu')
  menu?.addEventListener('mouseleave', () => {
    const selectedItem = document.querySelector('.ant-menu-item-selected')
    if (selectedItem) {
      updateFollowLine(selectedItem)
    } else {
      followLineStyle.value.opacity = 0
    }
  })
}

// 组件挂载和卸载时的处理
onMounted(() => {
  setTimeout(setupMenuEvents, 200)

  // 初始化跟随线条位置
  setTimeout(() => {
    const selectedItem = document.querySelector('.ant-menu-item-selected')
    if (selectedItem) {
      updateFollowLine(selectedItem)
    }
  }, 300)
})

onUnmounted(() => {
  // 清理事件监听
  const menuItems = document.querySelectorAll('.custom-menu .ant-menu-item')
  menuItems.forEach(item => {
    item.removeEventListener('mouseenter', () => { })
  })

  const menu = document.querySelector('.custom-menu')
  menu?.removeEventListener('mouseleave', () => { })
})
</script>

<style scoped>
#globalHeader {
  width: 100%;
  height: 64px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.03);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.title {
  flex-grow: 1;
  flex-shrink: 0;
  color: #333;
  font-size: 18px;
  font-weight: 500;
  margin-left: 16px;
  transition: all 0.3s ease;
  background: linear-gradient(45deg, #43b16a, #64d487);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.logo {
  height: 40px;
  width: 40px;
  display: inline-block;
  flex-shrink: 0;
  transition: all 0.3s ease;
  border-radius: 0;
  box-shadow: none;
}

#globalHeader .title-bar {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  height: 64px;
}

.menu-container {
  position: relative;
  height: 64px;
}

.custom-menu {
  line-height: 64px;
  height: 64px;
  border-bottom: none;
  background: transparent;
}

.follow-line {
  position: absolute;
  height: 3px;
  background: linear-gradient(90deg, #43b16a, #64d487);
  bottom: 0;
  border-radius: 3px;
  transition: all 0.3s ease;
  z-index: 1;
  box-shadow: 0 0 8px rgba(100, 212, 135, 0.5);
}

:deep(.ant-menu-item) {
  padding: 0 20px;
  color: #333;
  margin: 0 4px;
  transition: all 0.3s ease;
}

:deep(.ant-menu-item:hover) {
  color: #64d487;
  background-color: transparent;
}

:deep(.ant-menu-item-selected) {
  color: #64d487 !important;
  font-weight: 500;
}

:deep(.ant-menu-item::after) {
  display: none !important;
}

.user-login-status {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  height: 100%;
}

.user-info {
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.3s;
}

.user-info:hover {
  background-color: rgba(100, 212, 135, 0.05);
}

.user-avatar {
  border: 1px solid #e0e0e0;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.04);
}

.user-name {
  font-weight: 500;
  color: #333;
}

.login-btn {
  border-radius: 6px;
  box-shadow: 0 2px 0 rgba(0, 0, 0, 0.02);
  background-color: #64d487;
  border-color: #64d487;
  transition: all 0.3s ease;
}

.login-btn:hover {
  background-color: #43b16a;
  border-color: #43b16a;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
}

:deep(.ant-dropdown-menu) {
  border-radius: 6px;
  padding: 4px 0;
}

:deep(.ant-dropdown-menu-item) {
  padding: 8px 16px;
  border-radius: 4px;
  margin: 2px 4px;
  transition: all 0.3s ease;
}

:deep(.ant-dropdown-menu-item:hover) {
  background-color: rgba(100, 212, 135, 0.05);
}
</style>

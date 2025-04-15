<template>
  <div id="globalSider">
    <a-layout-sider v-if="loginUserStore.loginUser.id" class="sider glass-effect" :width="200" collapsible
      :collapsed="collapsed" :trigger="null" breakpoint="lg" @mouseenter="toggleCollapsed(false)"
      @mouseleave="toggleCollapsed(true)">
      <div class="sider-inner">
        <a-menu mode="inline" v-model:selectedKeys="current" :items="menuItems" @click="doMenuClick"
          class="custom-sider-menu" :inline-collapsed="collapsed" />
      </div>
    </a-layout-sider>
  </div>
</template>

<script lang="ts" setup>
import { computed, h, ref, watchEffect } from 'vue'
import { useRouter } from 'vue-router'
import { useLoginUserStore } from '../stores/useLoginUserStore'
import { PictureOutlined, UserOutlined } from '@ant-design/icons-vue'
import { SPACE_TYPE_ENUM } from '../constants/space'
import { TeamOutlined } from '@ant-design/icons-vue'
import { postSpaceUserListMy } from '../api/spaceUser'
import { message } from 'ant-design-vue'

// 侧边栏收缩状态
const collapsed = ref(true);

// 切换侧边栏展开/收缩
const toggleCollapsed = (value: boolean) => {
  collapsed.value = value;
}

// 固定的菜单列表
const fixedMenuItems = [
  {
    key: '/',
    label: '公共图库',
    icon: () => h(PictureOutlined),
  },
  {
    key: '/my_space',
    label: '我的空间',
    icon: () => h(UserOutlined),
  },
  {
    key: '/add_space?type=' + SPACE_TYPE_ENUM.TEAM,
    label: '创建团队',
    icon: () => h(TeamOutlined),
  },
]


const teamSpaceList = ref<any[]>([])
const menuItems = computed(() => {
  // 没有团队空间，只展示固定菜单
  if (teamSpaceList.value.length < 1) {
    return fixedMenuItems;
  }
  // 展示团队空间分组
  const teamSpaceSubMenus = teamSpaceList.value.map((spaceUser) => {
    const space = spaceUser.space
    return {
      key: '/space/' + spaceUser.spaceId,
      label: space?.spaceName,
    }
  })
  const teamSpaceMenuGroup = {
    type: 'group',
    label: '我的团队',
    key: 'teamSpace',
    children: teamSpaceSubMenus,
  }
  return [...fixedMenuItems, teamSpaceMenuGroup]
})

// 加载团队空间列表
const fetchTeamSpaceList = async () => {
  const res = await postSpaceUserListMy()
  if (res.data.code === 0 && res.data.data) {
    teamSpaceList.value = res.data.data
  } else {
    message.error('加载我的团队空间失败，' + res.data.msg)
  }
}
const router = useRouter()
const loginUserStore = useLoginUserStore()
/**
 * 监听变量，改变时触发数据的重新加载
 */
watchEffect(() => {
  // 登录才加载
  if (loginUserStore.loginUser.id) {
    fetchTeamSpaceList()
  }
})
// 当前选中菜单
const current = ref<string[]>([])
// 监听路由变化，更新当前选中菜单
router.afterEach((to, from, failure) => {
  current.value = [to.path]
})

// 路由跳转事件
const doMenuClick = ({ key }: { key: string }) => {
  router.push(key)
}
</script>

<style scoped>
#globalSider {
  height: 100%;
  position: relative;
  z-index: 900;
  transition: all 0.3s ease;
}

.sider {
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.02);
  transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1);
  overflow: hidden;
  background: transparent !important;
  height: 100%;
}

.sider-inner {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px 0;
}

.custom-sider-menu {
  background: transparent;
  border-right: none;
  flex: 1;
}

:deep(.ant-layout-sider) {
  background: transparent !important;
  transition: all 0.3s ease !important;
}

:deep(.ant-layout-sider-collapsed) {
  width: 80px !important;
  min-width: 80px !important;
  max-width: 80px !important;
}

:deep(.ant-layout-sider:not(.ant-layout-sider-collapsed)) {
  width: 200px !important;
  min-width: 200px !important;
  max-width: 200px !important;
}

/* 修复折叠侧边栏背景 */
:deep(.ant-layout-sider-collapsed) {
  background: transparent !important;
}

:deep(.ant-menu-inline-collapsed) {
  width: 80px;
  background: transparent !important;
}

:deep(.ant-menu-inline-collapsed > .ant-menu-item) {
  padding: 0 calc(50% - 16px / 2) !important;
}

:deep(.ant-menu-inline-collapsed .ant-menu-item-icon) {
  font-size: 18px;
}

:deep(.ant-menu-item) {
  margin: 4px 8px;
  height: 40px;
  line-height: 40px;
  padding: 0 16px;
  color: #333;
  border-radius: 6px;
  transition: all 0.3s;
}

:deep(.ant-menu-item:hover) {
  color: #64d487;
  background-color: rgba(100, 212, 135, 0.08);
}

:deep(.ant-menu-item-selected) {
  background-color: rgba(100, 212, 135, 0.15) !important;
  color: #64d487 !important;
  font-weight: 500;
}

:deep(.ant-menu-item .anticon) {
  margin-right: 10px;
  font-size: 16px;
}

:deep(.ant-menu-item-group-title) {
  padding: 12px 24px 8px;
  color: #999;
  font-size: 12px;
}

:deep(.ant-menu-inline .ant-menu-item::after) {
  display: none;
}

:deep(.ant-layout-sider-trigger) {
  display: none;
}

/* 确保菜单组背景透明 */
:deep(.ant-menu-item-group) {
  background: transparent !important;
}

/* 确保子菜单背景透明 */
:deep(.ant-menu-submenu) {
  background: transparent !important;
}

:deep(.ant-menu-submenu-title) {
  background: transparent !important;
}

:deep(.ant-menu-submenu-popup) {
  background: transparent !important;
}

:deep(.ant-menu-submenu-popup > .ant-menu) {
  background: rgba(255, 255, 255, 0.95) !important;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

/* 滚动条样式 */
#globalSider::-webkit-scrollbar {
  width: 6px;
}

#globalSider::-webkit-scrollbar-track {
  background: transparent;
}

#globalSider::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
}

#globalSider::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.2);
}
</style>

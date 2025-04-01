<template>
  <div class="picture-list">
    <!-- 图片列表 -->
    <a-list
      :grid="{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 5, xxl: 6 }"
      :data-source="dataList"
      :loading="loading"
    >
      <template #renderItem="{ item: picture }">
        <a-list-item style="padding: 0">
          <!-- 单张图片 -->
          <a-card hoverable @click="doClickPicture(picture)" class="custom-card">
            <!-- 图片封面 -->
            <template #cover>
              <div class="cover-container">
                <img
                  :alt="picture.name"
                  :src="picture.thumbnailUrl || picture.url"
                  class="cover-image"
                />
              </div>
            </template>

            <!-- 卡片内容 -->
            <a-card-meta :title="picture.name">
              <template #description>
                <a-flex wrap="wrap" gap="small">
                  <a-tag color="green">{{ picture.category || '默认' }}</a-tag>
                  <a-tag v-for="tag in picture.tags" :key="tag">{{ tag }}</a-tag>
                </a-flex>
              </template>
            </a-card-meta>

            <!-- 操作按钮 -->
            <template v-if="showOp" #actions>
              <a-space class="card-actions" @click="e => doSearch(picture, e)">
                <search-outlined />
                搜索
              </a-space>
              <a-space class="card-actions" @click="e => doEdit(picture, e)">
                <edit-outlined />
                编辑
              </a-space>
              <a-space class="card-actions" @click="e => doDelete(picture, e)">
                <delete-outlined />
                删除
              </a-space>
            </template>
          </a-card>

        </a-list-item>
      </template>
    </a-list>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { DeleteOutlined, EditOutlined, SearchOutlined } from '@ant-design/icons-vue'
import { postPictureOpenApiDelete } from '@/api/picture.ts'
import { message } from 'ant-design-vue'

interface Props {
  dataList?: API.PictureVO[]
  loading?: boolean
  showOp?: boolean
  onReload?: () => void
}


const props = withDefaults(defineProps<Props>(), {
  dataList: () => [],
  loading: false,
  showOp: false,
})

// 跳转至图片详情
const router = useRouter()
const doClickPicture = (picture) => {
  router.push({
    path: `/picture/${picture.id}`,
  })
}

// 搜索
const doSearch = (picture, e) => {
  e.stopPropagation()
  window.open(`/search_picture?pictureId=${picture.id}`)
}


// 编辑
const doEdit = (picture, e) => {
  e.stopPropagation()
  router.push({
    path: '/add_picture',
    query: {
      id: picture.id,
      spaceId: picture.spaceId,
    },
  })
}

// 删除
const doDelete = async (picture, e) => {
  e.stopPropagation()
  const id = picture.id
  if (!id) {
    return
  }
  const res = await postPictureOpenApiDelete({ id })
  if (res.data.code === 0) {
    message.success('删除成功')
    // 让外层刷新
    props.onReload?.()
  } else {
    message.error('删除失败')
  }
}

</script>

<style scoped>
/* 卡片美化 */
.custom-card {
  width: 100%;
  max-width: 240px;
  transition:
    transform 0.3s ease-in-out,
    box-shadow 0.3s ease-in-out;
  border-radius: 10px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

/* 鼠标悬停时的效果 */
.custom-card:hover {
  transform: translateY(-5px);
  box-shadow: 0px 8px 16px rgba(0, 0, 0, 0.12);
}

/* 覆盖容器（图片区域） */
.cover-container {
  width: 100%;
  height: 160px; /* ⬆️ 增加高度，让图片占比更多 */
  overflow: hidden;
  border-radius: 10px 10px 0 0;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 图片样式 */
.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease-in-out;
}

.custom-card:hover .cover-image {
  transform: scale(1.05);
}

/* 操作按钮 */
.card-actions {
  display: flex;
  justify-content: center; /* 居中 */
  align-items: center;
  gap: 16px; /* 按钮之间的间距 */
  padding: 6px 0;
  font-size: 14px;
  color: #333;
  cursor: pointer;
  transition: color 0.2s ease-in-out;
}

.card-actions:hover {
  color: #1890ff;
}

</style>

<!-- 主页 -->
<template>
  <div id="pictureDetailPage">
    <a-row :gutter="[16, 16]">
      <!-- 图片展示区 -->
      <a-col :sm="24" :md="16" :xl="18">
        <a-card title="图片预览">
          <a-image style="max-height: 600px; object-fit: contain" :src="picture.url" />
        </a-card>
      </a-col>
      <!-- 图片信息区 -->
      <a-col :sm="24" :md="8" :xl="6">
        <a-card title="图片信息">
          <a-descriptions :column="1">
            <a-descriptions-item label="作者">
              <a-space>
                <a-avatar :size="24" :src="picture.user?.userAvatar" />
                <div>{{ picture.user?.userName }}</div>
              </a-space>
            </a-descriptions-item>
            <a-descriptions-item label="名称">
              {{ picture.name ?? '未命名' }}
            </a-descriptions-item>
            <a-descriptions-item label="简介">
              {{ picture.introduction ?? '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="分类">
              {{ picture.category ?? '默认' }}
            </a-descriptions-item>
            <a-descriptions-item label="标签">
              <a-tag v-for="tag in picture.tags" :key="tag">
                {{ tag }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="格式">
              {{ picture.picFormat ?? '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="宽度">
              {{ picture.picWidth ?? '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="高度">
              {{ picture.picHeight ?? '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="宽高比">
              {{ picture.picScale ?? '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="大小">
              {{ formatSize(picture.picSize) }}
            </a-descriptions-item>
            <a-descriptions-item label="主色调">
              <a-space>
                {{ picture.picColor ?? '-' }}
                <div v-if="picture.picColor" :style="{
                  backgroundColor: toHexColor(picture.picColor),
                  width: '16px',
                  height: '16px',
                }" />
              </a-space>
            </a-descriptions-item>
          </a-descriptions>
          <a-space wrap>
            <a-button type="primary" ghost @click="doShare">
              分享
              <template #icon>
                <share-alt-outlined />
              </template>
            </a-button>
            <a-button type="primary" @click="doDownload">
              免费下载
              <template #icon>
                <DownloadOutlined />
              </template>
            </a-button>
            <a-button v-if="canEdit" type="default" @click="doEdit">
              编辑
              <template #icon>
                <EditOutlined />
              </template>
            </a-button>
            <a-button v-if="canDelete" danger @click="doDelete">
              删除
              <template #icon>
                <DeleteOutlined />
              </template>
            </a-button>
          </a-space>
        </a-card>
      </a-col>
    </a-row>
    <ShareModal ref="shareModalRef" :link="shareLink" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { getPictureGetVo, postPictureOpenApiDelete } from '@/api/picture.ts'
import { message } from 'ant-design-vue'
import { downloadImage, formatSize, toHexColor } from '../utils'
import { useLoginUserStore } from '@/stores/useLoginUserStore.ts'
import { useRouter } from 'vue-router'
import ShareModal from '@/components/ShareModal.vue'
import { SPACE_PERMISSION_ENUM } from '@/constants/space.ts'
const props = defineProps<{
  id: string | number
}>()
const picture = ref<API.PictureVO>({})
const router = useRouter()
// 获取图片详情
const fetchPictureDetail = async () => {
  try {
    const res = await getPictureGetVo({
      id: props.id,
    })
    if (res.data.code === 0 && res.data.data) {
      picture.value = res.data.data
    } else {
      message.error('获取图片详情失败，' + res.data.msg)
    }
  } catch (e: any) {
    message.error('获取图片详情失败：' + e.message)
  }
}

// 通用权限检查函数
function createPermissionChecker(permission: string) {
  return computed(() => {
    return (picture.value.permissionList ?? []).includes(permission)
  })
}

// 定义权限检查
const canEdit = createPermissionChecker(SPACE_PERMISSION_ENUM.PICTURE_EDIT)
const canDelete = createPermissionChecker(SPACE_PERMISSION_ENUM.PICTURE_DELETE)


onMounted(() => {
  fetchPictureDetail()
})

// 编辑
const doEdit = () => {
  router.push({
    path: '/add_picture',
    query: {
      id: picture.value.id,
      ...(picture.value.spaceId && picture.value.spaceId !== "0" ? { spaceId: picture.value.spaceId } : {}),
    },
  })
}

// 删除
const doDelete = async () => {
  const id = picture.value.id
  if (!id) {
    return
  }
  const res = await postPictureOpenApiDelete({ id })
  if (res.data.code === 0) {
    message.success('删除成功')
  } else {
    message.error('删除失败')
  }
}

const loginUserStore = useLoginUserStore()

// 处理下载
const doDownload = () => {
  downloadImage(picture.value.url)
}

// 分享弹窗引用
const shareModalRef = ref()
// 分享链接
const shareLink = ref<string>()

// 分享
const doShare = () => {
  shareLink.value = `${window.location.protocol}//${window.location.host}/picture/${picture.value.id}`
  if (shareModalRef.value) {
    shareModalRef.value.openModal()
  }
}

</script>

<style scoped>
#pictureDetailPage {
  margin-bottom: 16px;
}

:deep(.ant-btn-primary) {
  background-color: #64d487;
  border-color: #64d487;
  color: #fff;
}

:deep(.ant-btn-primary:hover),
:deep(.ant-btn-primary:focus) {
  background-color: #4bc072;
  border-color: #4bc072;
  color: #fff;
}

:deep(.ant-btn-primary[disabled]),
:deep(.ant-btn-primary[disabled]:hover),
:deep(.ant-btn-primary[disabled]:focus),
:deep(.ant-btn-primary[disabled]:active) {
  background-color: rgba(0, 0, 0, 0.04);
  border-color: #d9d9d9;
  color: rgba(0, 0, 0, 0.25);
}

:deep(.ant-btn-link) {
  color: #64d487;
}

:deep(.ant-btn-link:hover),
:deep(.ant-btn-link:focus) {
  color: #4bc072;
}

:deep(.ant-btn-default) {
  border-color: #64d487;
  color: #64d487;
}

:deep(.ant-btn-default:hover),
:deep(.ant-btn-default:focus) {
  border-color: #4bc072;
  color: #4bc072;
}
</style>

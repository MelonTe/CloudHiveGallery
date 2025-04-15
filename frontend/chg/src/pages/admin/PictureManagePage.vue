<template>
  <div id="pictureManagePage">
    <a-flex justify="space-between">
      <h2>图片管理</h2>
      <a-space>
        <a-button type="primary" href="/add_picture" target="_blank">+ 创建图片</a-button>
        <a-button type="primary" href="/add_picture/batch" target="_blank" ghost>+ 批量创建图片</a-button>
      </a-space>
    </a-flex>
    <div style="margin-bottom: 16px"></div>
    <!-- 搜索表单 -->
    <a-form layout="inline" :model="searchParams" @finish="doSearch">
      <a-form-item label="关键词" name="searchText">
        <a-input v-model:value="searchParams.searchText" placeholder="从名称和简介搜索" allow-clear />
      </a-form-item>
      <a-form-item label="类型" name="category">
        <a-input v-model:value="searchParams.category" placeholder="请输入类型" allow-clear />
      </a-form-item>
      <a-form-item label="标签" name="tags">
        <a-select v-model:value="searchParams.tags" mode="tags" placeholder="请输入标签" style="min-width: 180px"
          allow-clear />
      </a-form-item>
      <a-form-item label="审核状态" name="reviewStatus">
        <a-select v-model:value="searchParams.reviewStatus" :options="PIC_REVIEW_STATUS_OPTIONS" placeholder="请输入审核状态"
          style="min-width: 180px" allow-clear />
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit">搜索</a-button>
      </a-form-item>
    </a-form>
    <div style="margin-bottom: 16px" />
    <!-- 表格 -->
    <a-table :columns="columns" :data-source="dataList" :pagination="pagination" @change="doTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'url'">
          <a-image :src="record.url" :width="120" />
        </template>
        <!-- 标签 -->
        <template v-if="column.dataIndex === 'tags'">
          <a-space wrap>
            <a-tag v-for="tag in JSON.parse(record.tags || '[]')" :key="tag">{{ tag }}</a-tag>
          </a-space>
        </template>
        <!-- 图片信息 -->
        <template v-if="column.dataIndex === 'picInfo'">
          <div>格式：{{ record.picFormat }}</div>
          <div>宽度：{{ record.picWidth }}</div>
          <div>高度：{{ record.picHeight }}</div>
          <div>宽高比：{{ record.picScale }}</div>
          <div>大小：{{ (record.picSize / 1024).toFixed(2) }}KB</div>
        </template>
        <!-- 审核信息 -->
        <template v-if="column.dataIndex === 'reviewMessage'">
          <div>审核状态：{{ PIC_REVIEW_STATUS_MAP[record.reviewStatus] }}</div>
          <div>审核信息：{{ record.reviewMessage }}</div>
          <div>审核人：{{ record.reviewerId }}</div>
        </template>
        <template v-else-if="column.dataIndex === 'createTime'">
          {{ dayjs(record.createTime).format('YYYY-MM-DD HH:mm:ss') }}
        </template>
        <template v-else-if="column.dataIndex === 'editTime'">
          {{ dayjs(record.editTime).format('YYYY-MM-DD HH:mm:ss') }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space wrap>
            <a-button v-if="record.reviewStatus !== PIC_REVIEW_STATUS_ENUM.PASS" type="link"
              @click="handleReview(record, PIC_REVIEW_STATUS_ENUM.PASS)">
              通过
            </a-button>
            <a-button v-if="record.reviewStatus !== PIC_REVIEW_STATUS_ENUM.REJECT" type="link" danger
              @click="handleReview(record, PIC_REVIEW_STATUS_ENUM.REJECT)">
              拒绝
            </a-button>
            <a-button type="link" :href="`/add_picture?id=${record.id}`" target="_blank">编辑
            </a-button>
            <a-button type="link" danger @click="doDelete(record.id)">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>
  </div>
</template>
<script lang="ts" setup>
import dayjs from 'dayjs'
import {
  postPictureListPage,
  postPictureOpenApiDelete,
  postPictureReview,
} from '@/api/picture'
import { message } from 'ant-design-vue'
import { computed, onMounted, reactive, ref } from 'vue'
import {
  PIC_REVIEW_STATUS_MAP,
  PIC_REVIEW_STATUS_ENUM,
  PIC_REVIEW_STATUS_OPTIONS,
} from '../../constants/picture.ts'
const columns = [
  {
    title: 'id',
    dataIndex: 'id',
    width: 80,
  },
  {
    title: '图片',
    dataIndex: 'url',
  },
  {
    title: '名称',
    dataIndex: 'name',
  },
  {
    title: '简介',
    dataIndex: 'introduction',
    ellipsis: true,
  },
  {
    title: '类型',
    dataIndex: 'category',
  },
  {
    title: '标签',
    dataIndex: 'tags',
  },
  {
    title: '图片信息',
    dataIndex: 'picInfo',
  },
  {
    title: '用户 id',
    dataIndex: 'userId',
    width: 80,
  },
  {
    title: '空间 id',
    dataIndex: 'spaceId',
    width: 80,
  },
  {
    title: '审核信息',
    dataIndex: 'reviewMessage',
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
  },
  {
    title: '编辑时间',
    dataIndex: 'editTime',
  },
  {
    title: '操作',
    key: 'action',
  },
]
// 定义数据
// ref是数据整体更新，reactive是数据局部更新
const dataList = ref<API.Picture[]>()
const total = ref(0)

/* 表格变换，重新获取数据 */
const doTableChange = (page: any) => {
  searchParams.current = page.current
  searchParams.pageSize = page.pageSize
  fetchData()
}

//定义搜索条件
const searchParams = reactive<API.PictureQueryRequest>({
  current: 1,
  pageSize: 10,
  sortField: 'create_time',
  sortOrder: 'descend',
})
//分页器
const pagination = computed(() => {
  return {
    current: searchParams.current,
    pageSize: searchParams.pageSize,
    total: total.value,
    showSizeChanger: true,
    showTotal: (total: number) => `共 ${total} 条`,
  }
})
//获取数据
const fetchData = async () => {
  const res = await postPictureListPage({
    ...searchParams,
    isNullSpaceId: true,
  })
  if (res.data.code === 0 && res.data.data) {
    dataList.value = res.data.data.records ?? []
    total.value = res.data.data.total ?? 0
  } else {
    message.error('获取数据失败，' + res.data.msg)
  }
}

//生命周期，开始时加载
onMounted(() => {
  fetchData()
})

// 搜索数据
const doSearch = () => {
  // 重置页码
  searchParams.current = 1
  fetchData()
}

// 删除数据
const doDelete = async (id: string) => {
  if (!id) {
    return
  }
  const res = await postPictureOpenApiDelete({ id })
  if (res.data.code === 0) {
    message.success('删除成功')
    // 刷新数据
    fetchData()
  } else {
    message.error('删除失败')
  }
}
//审核函数
const handleReview = async (record: API.Picture, reviewStatus: number) => {
  const reviewMessage =
    reviewStatus === PIC_REVIEW_STATUS_ENUM.PASS ? '管理员操作通过' : '管理员操作拒绝'
  const res = await postPictureReview({
    id: record.id,
    reviewStatus,
    reviewMessage,
  })
  if (res.data.code === 0) {
    message.success('审核操作成功')
    // 重新获取列表
    fetchData()
  } else {
    message.error('审核操作失败，' + res.data.msg)
  }
}
</script>

<style scoped>
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

:deep(.ant-btn-default) {
  border-color: #64d487;
  color: #64d487;
}

:deep(.ant-btn-default:hover),
:deep(.ant-btn-default:focus) {
  border-color: #4bc072;
  color: #4bc072;
}

:deep(.ant-btn-link) {
  color: #64d487;
}

:deep(.ant-btn-link:hover),
:deep(.ant-btn-link:focus) {
  color: #4bc072;
}

:deep(.ant-btn-link.ant-btn-dangerous),
:deep(.ant-btn-link.ant-btn-dangerous:hover),
:deep(.ant-btn-link.ant-btn-dangerous:focus) {
  color: #ff4d4f;
}
</style>

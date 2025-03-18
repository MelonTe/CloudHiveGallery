<template>
    <div id="pictureManagePage">
        <!-- 搜索表单 -->
      <a-form layout="inline" :model="searchParams" @finish="doSearch">
        <a-form-item label="关键词" name="searchText">
          <a-input
            v-model:value="searchParams.searchText"
            placeholder="从名称和简介搜索"
            allow-clear
          />
        </a-form-item>
        <a-form-item label="类型" name="category">
          <a-input v-model:value="searchParams.category" placeholder="请输入类型" allow-clear />
        </a-form-item>
        <a-form-item label="标签" name="tags">
          <a-select
            v-model:value="searchParams.tags"
            mode="tags"
            placeholder="请输入标签"
            style="min-width: 180px"
            allow-clear
          />
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
            <template v-else-if="column.dataIndex === 'createTime'">
              {{ dayjs(record.createTime).format('YYYY-MM-DD HH:mm:ss') }}
            </template>
            <template v-else-if="column.dataIndex === 'editTime'">
              {{ dayjs(record.editTime).format('YYYY-MM-DD HH:mm:ss') }}
            </template>
            <template v-else-if="column.key === 'action'">
              <a-space>
                <a-button type="link" :href="`/add_picture?id=${record.id}`" target="_blank">
                  编辑
                </a-button>
                <a-button danger @click="doDelete(record.id)">删除</a-button>
              </a-space>

            </template>
          </template>

        </a-table>
    </div>
</template>
<script lang="ts" setup>
import dayjs from 'dayjs';
import { postPictureListPage, postPictureListPageVo, postPictureOpenApiDelete } from '@/api/picture'
import { message } from 'ant-design-vue';
import { computed, onMounted, reactive, ref } from 'vue';
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
  sortOrder:'descend',
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
        ...searchParams
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
    const res = await postPictureOpenApiDelete({id})
    if (res.data.code === 0) {
        message.success('删除成功')
        // 刷新数据
        fetchData()
    } else {
        message.error('删除失败')
    }
}
</script>

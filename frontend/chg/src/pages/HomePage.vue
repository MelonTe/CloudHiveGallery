<!-- 主页 -->
<template>
  <div id="homePage">
    <!-- 搜索框 -->
    <div class="search-bar">
      <a-input-search
        placeholder="从海量图片中搜索"
        v-model:value="searchParams.searchText"
        enter-button="搜索"
        size="large"
        @search="doSearch"
      />
    </div>
    <!-- 分类 + 标签 -->
    <a-tabs v-model:activeKey="selectedCategory" @change="doSearch">
      <a-tab-pane key="all" tab="全部" />
      <a-tab-pane v-for="category in categoryList" :key="category" :tab="category" />
    </a-tabs>
    <div class="tag-bar">
      <span style="margin-right: 8px">标签：</span>
      <a-space :size="[0, 8]" wrap>
        <a-checkable-tag
          v-for="(tag, index) in tagList"
          :key="tag"
          v-model:checked="selectedTagList[index]"
          @change="doSearch"
        >
          {{ tag }}
        </a-checkable-tag>
      </a-space>
    </div>

    <!-- 图片列表 -->
    <a-list
      :grid="{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 5, xxl: 6 }"
      :data-source="dataList"
      :pagination="pagination"
      :loading="loading"
    >
      <template #renderItem="{ item: picture }">
        <a-list-item style="padding: 0">
          <!-- 单张图片 -->
          <a-card
            hoverable
            @click="doClickPicture(picture)"
            class="custom-card"
            style="height: 270px"
          >
            <template #cover>
              <div class="cover-container">
                <img :alt="picture.name" :src="picture.thumbnailUrl || picture.url" class="cover-image" />
              </div>
            </template>
            <a-card-meta :title="picture.name">
              <template #description>
                <a-flex wrap="wrap">
                  <a-tag color="green">{{ picture.category || '默认' }}</a-tag>
                  <a-tag v-for="tag in picture.tags" :key="tag">{{ tag }}</a-tag>
                </a-flex>
              </template>
            </a-card-meta>
          </a-card>
        </a-list-item>
      </template>
    </a-list>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { getPictureTagCategory, postPictureListPageVo } from '@/api/picture.ts'
import { message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
// 数据
const dataList = ref([])
const total = ref(0)
const loading = ref(true)
const router = useRouter()
// 跳转至图片详情
const doClickPicture = (picture) => {
  router.push({
    path: `/picture/${picture.id}`,
  })
}

// 搜索条件
const searchParams = reactive<API.PictureQueryRequest>({
  current: 1,
  pageSize: 12,
  sortField: 'create_time',
  sortOrder: 'descend',
})

// 分页参数
const pagination = computed(() => {
  return {
    current: searchParams.current ?? 1,
    pageSize: searchParams.pageSize ?? 10,
    total: total.value,
    // 切换页号时，会修改搜索参数并获取数据
    onChange: (page, pageSize) => {
      searchParams.current = page
      searchParams.pageSize = pageSize
      fetchData()
    },
  }
})

// 获取数据
const fetchData = async () => {
  loading.value = true
  // 转换搜索参数
  const params = {
    ...searchParams,
    tags: [] as string[],
  }
  if (selectedCategory.value !== 'all') {
    params.category = selectedCategory.value
  }
  selectedTagList.value.forEach((useTag, index) => {
    if (useTag) {
      params.tags.push(tagList.value[index])
    }
  })
  const res = await postPictureListPageVo(params)
  if (res.data.data) {
    dataList.value = res.data.data.records ?? []
    total.value = res.data.data.total ?? 0
  } else {
    message.error('获取数据失败，' + res.data.msg)
  }
  loading.value = false
}

// 页面加载时请求一次
onMounted(() => {
  fetchData()
})

const doSearch = () => {
  // 重置搜索条件
  searchParams.current = 1
  fetchData()
}

const categoryList = ref<string[]>([])
const selectedCategory = ref<string>('all')
const tagList = ref<string[]>([])
const selectedTagList = ref<boolean[]>([])

// 获取标签和分类选项
const getTagCategoryOptions = async () => {
  const res = await getPictureTagCategory()
  if (res.data.code === 0 && res.data.data) {
    // 转换成下拉选项组件接受的格式
    categoryList.value = res.data.data.categoryList ?? []
    tagList.value = res.data.data.tagList ?? []
  } else {
    message.error('加载分类标签失败，' + res.data.msg)
  }
}

onMounted(() => {
  getTagCategoryOptions()
})
</script>

<style scoped>
#homePage {
  margin-bottom: 16px;
}
#homePage .search-bar {
  max-width: 480px;
  margin: 0 auto 16px;
}
#homePage .tag-bar {
  margin-bottom: 16px;
}

/* 卡片整体美化 */
.custom-card {
  transition:
    transform 0.3s ease-in-out,
    box-shadow 0.3s ease-in-out;
  border-radius: 8px;
}

.custom-card:hover {
  transform: translateY(-5px); /* 轻微浮起 */
  box-shadow: 0px 10px 20px rgba(0, 0, 0, 0.15); /* 柔和阴影 */
}

/* 覆盖容器 */
.cover-container {
  position: relative;
  overflow: hidden;
  border-radius: 8px 8px 0 0;
}

/* 图片放大动画 */
.cover-image {
  width: 100%;
  height: 180px;
  object-fit: cover;
  transition: transform 0.3s ease-in-out;
}

.custom-card:hover .cover-image {
  transform: scale(1.08); /* 轻微放大，保持美观 */
}

/* 昵称显示优化 */
.image-name {
  position: absolute;
  left: 0;
  bottom: 0;
  width: 100%;
  padding: 6px 12px;
  font-size: 18px;
  color: white;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.2), transparent); /* 渐变背景 */
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  opacity: 0;
  transition: opacity 0.3s ease-in-out;
}

.custom-card:hover .image-name {
  opacity: 1;
  font-size: 18px; /* 变大 */
}
</style>

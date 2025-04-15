<!-- 主页 -->
<template>
  <div id="homePage">
    <!-- 搜索框 -->
    <div class="search-bar">
      <a-input-search placeholder="从海量图片中搜索" v-model:value="searchParams.searchText" enter-button="搜索" size="large"
        @search="doSearch" class="glass-search" />
    </div>

    <!-- 现代化筛选器 -->
    <div class="filter-container glass-effect">
      <div class="filter-section">
        <div class="filter-title">
          <span>分类</span>
        </div>
        <div class="filter-options">
          <div v-for="(cat, index) in ['all', ...categoryList]" :key="cat"
            :class="['filter-chip', selectedCategory === cat ? 'active' : '']"
            @click="selectedCategory = cat; doSearch()">
            <span>{{ cat === 'all' ? '全部' : cat }}</span>
            <span v-if="selectedCategory === cat" class="check-icon">✓</span>
          </div>
        </div>
      </div>

      <a-divider style="margin: 12px 0" />

      <div class="filter-section">
        <div class="filter-title">
          <span>标签</span>
          <a-button type="link" size="small" @click="clearTags" v-if="hasSelectedTag" class="clear-btn">
            清除
          </a-button>
        </div>
        <div class="filter-options">
          <div v-for="(tag, index) in tagList" :key="tag"
            :class="['filter-tag', selectedTagList[index] ? 'active' : '']" @click="tagClick(index)">
            {{ tag }}
          </div>
        </div>
      </div>
    </div>

    <!-- 图片列表 -->
    <div class="picture-container">
      <PictureList :dataList="dataList" :loading="loading" />
      <a-pagination style="text-align: right; margin-top: 24px" v-model:current="searchParams.current"
        v-model:pageSize="searchParams.pageSize" :total="total" @change="onPageChange" show-quick-jumper
        :show-size-changer="true" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { getPictureTagCategory, postPictureListPageVo } from '@/api/picture.ts'
import { message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import PictureList from '@/components/PictureList.vue'
// 数据
const dataList = ref([])
const total = ref(0)
const loading = ref(true)
const router = useRouter()
// 搜索条件
const searchParams = reactive<API.PictureQueryRequest>({
  current: 1,
  pageSize: 12,
  sortField: 'create_time',
  sortOrder: 'descend',
})

const onPageChange = (page, pageSize) => {
  searchParams.current = page
  searchParams.pageSize = pageSize
  fetchData()
}

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

// 是否有选中的标签
const hasSelectedTag = computed(() => {
  return selectedTagList.value.some(tag => tag === true)
})

// 点击标签
const tagClick = (index: number) => {
  selectedTagList.value[index] = !selectedTagList.value[index]
  doSearch()
}

// 清除所有标签
const clearTags = () => {
  selectedTagList.value = selectedTagList.value.map(() => false)
  doSearch()
}

// 获取标签和分类选项
const getTagCategoryOptions = async () => {
  const res = await getPictureTagCategory()
  if (res.data.code === 0 && res.data.data) {
    // 转换成下拉选项组件接受的格式
    categoryList.value = res.data.data.categoryList ?? []
    tagList.value = res.data.data.tagList ?? []
    // 初始化标签选择状态
    selectedTagList.value = new Array(tagList.value.length).fill(false)
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
  margin-bottom: 24px;
  max-width: 1800px;
  margin: 0 auto;
}

#homePage .search-bar {
  max-width: 600px;
  margin: 0 auto 24px;
}

.glass-search :deep(.ant-input) {
  border-radius: 50px 0 0 50px;
  padding-left: 20px;
  height: 48px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid rgba(200, 200, 200, 0.3);
}

.glass-search :deep(.ant-input-search-button) {
  height: 48px;
  border-radius: 0 50px 50px 0;
  width: 100px;
  background: #64d487;
  border-color: #64d487;
  box-shadow: 0 2px 8px rgba(100, 212, 135, 0.15);
}

.glass-search :deep(.ant-input-search-button:hover) {
  background: #43b16a;
  border-color: #43b16a;
}

/* 现代化筛选器样式 */
.filter-container {
  padding: 20px;
  margin-bottom: 24px;
  border-radius: 12px;
}

.filter-section {
  margin-bottom: 8px;
}

.filter-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-weight: 500;
  color: #333;
}

.clear-btn {
  padding: 0;
  height: auto;
  color: #999;
}

.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.filter-chip {
  display: inline-flex;
  align-items: center;
  padding: 6px 16px;
  border-radius: 50px;
  background: rgba(245, 245, 245, 0.8);
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
  border: 1px solid transparent;
}

.filter-chip:hover {
  background: rgba(240, 240, 240, 0.9);
  transform: translateY(-1px);
}

.filter-chip.active {
  background: rgba(100, 212, 135, 0.1);
  color: #64d487;
  border-color: rgba(100, 212, 135, 0.3);
  font-weight: 500;
}

.check-icon {
  margin-left: 6px;
  font-size: 12px;
}

.filter-tag {
  display: inline-flex;
  align-items: center;
  padding: 5px 12px;
  border-radius: 4px;
  background: rgba(245, 245, 245, 0.8);
  cursor: pointer;
  transition: all 0.3s;
  font-size: 13px;
  border: 1px solid transparent;
}

.filter-tag:hover {
  background: rgba(240, 240, 240, 0.9);
  transform: translateY(-1px);
}

.filter-tag.active {
  background: rgba(100, 212, 135, 0.1);
  color: #64d487;
  border-color: rgba(100, 212, 135, 0.3);
  font-weight: 500;
}

.picture-container {
  margin-top: 24px;
  min-height: 300px;
  display: flex;
  flex-direction: column;
}

/* 响应式调整 - 确保在小屏幕上也有良好展示 */
@media (max-width: 768px) {
  .picture-container {
    margin-top: 16px;
  }

  .filter-container {
    padding: 16px;
  }

  #homePage .search-bar {
    margin-bottom: 16px;
  }
}

/* 卡片整体美化 */
.custom-card {
  transition:
    transform 0.3s ease-in-out,
    box-shadow 0.3s ease-in-out;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  border: none;
}

.custom-card:hover {
  transform: translateY(-5px);
  /* 轻微浮起 */
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  /* 柔和阴影 */
}

/* 覆盖容器 */
.cover-container {
  position: relative;
  overflow: hidden;
  border-radius: 12px 12px 0 0;
}

/* 图片放大动画 */
.cover-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
  transition: transform 0.3s ease-in-out;
}

.custom-card:hover .cover-image {
  transform: scale(1.08);
  /* 轻微放大，保持美观 */
}

.custom-card:hover .image-name {
  opacity: 1;
  font-size: 18px;
  /* 变大 */
}

:deep(.ant-pagination-item) {
  border-color: #64d487;
}

:deep(.ant-pagination-item:hover) {
  border-color: #64d487;
}

:deep(.ant-pagination-item-active) {
  background-color: #64d487;
  border-color: #64d487;
}

:deep(.ant-pagination-item-active a) {
  color: #fff;
}

:deep(.ant-pagination-prev .ant-pagination-item-link),
:deep(.ant-pagination-next .ant-pagination-item-link) {
  border-color: #64d487;
  color: #64d487;
}

:deep(.ant-pagination-prev:hover .ant-pagination-item-link),
:deep(.ant-pagination-next:hover .ant-pagination-item-link) {
  border-color: #64d487;
  color: #64d487;
}

/* 按钮样式统一 */
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
</style>

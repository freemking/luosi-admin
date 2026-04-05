<template>
  <div class="dashboard">
    <a-page-header title="Dashboard" style="padding: 0 0 24px 0;" />
    <a-row :gutter="16">
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card :bordered="true" class="stat-card">
          <template #title>
            <a-space>
              <AppstoreOutlined style="color: #52c41a; font-size: 24px;" />
              <span>产品总数</span>
            </a-space>
          </template>
          <div class="stat-number">{{ stats.products }}</div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card :bordered="true" class="stat-card">
          <template #title>
            <a-space>
              <MessageOutlined style="color: #faad14; font-size: 24px;" />
              <span>反馈总数</span>
            </a-space>
          </template>
          <div class="stat-number">{{ stats.feedbacks }}</div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card :bordered="true" class="stat-card">
          <template #title>
            <a-space>
              <FileTextOutlined style="color: #1890ff; font-size: 24px;" />
              <span>新闻总数</span>
            </a-space>
          </template>
          <div class="stat-number">{{ stats.news }}</div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card :bordered="true" class="stat-card">
          <template #title>
            <a-space>
              <LayoutOutlined style="color: #722ed1; font-size: 24px;" />
              <span>广告位总数</span>
            </a-space>
          </template>
          <div class="stat-number">{{ stats.adPositions }}</div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card :bordered="true" class="stat-card">
          <template #title>
            <a-space>
              <PictureOutlined style="color: #eb2f96; font-size: 24px;" />
              <span>广告总数</span>
            </a-space>
          </template>
          <div class="stat-number">{{ stats.ads }}</div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { AppstoreOutlined, MessageOutlined, FileTextOutlined, LayoutOutlined, PictureOutlined } from '@ant-design/icons-vue'
import { useProductStore, useFeedbackStore, useNewsStore } from '../stores/auth'
import config from '../config'

const productStore = useProductStore()
const feedbackStore = useFeedbackStore()
const newsStore = useNewsStore()

const stats = ref({
  products: 0,
  feedbacks: 0,
  news: 0,
  adPositions: 0,
  ads: 0
})

const fetchAdPositionCount = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`${config.API_BASE_URL}/ad-positions/count`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    return data.count || 0
  } catch (err) {
    console.error('Failed to fetch ad position count:', err)
    return 0
  }
}

const fetchAdCount = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`${config.API_BASE_URL}/ads/count`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    return data.count || 0
  } catch (err) {
    console.error('Failed to fetch ad count:', err)
    return 0
  }
}

const fetchStats = async () => {
  try {
    const [productCount, feedbackCount, newsCount, adPositionCount, adCount] = await Promise.all([
      productStore.getProductCount(),
      feedbackStore.getFeedbackCount(),
      newsStore.getNewsCount(),
      fetchAdPositionCount(),
      fetchAdCount()
    ])
    stats.value.products = productCount
    stats.value.feedbacks = feedbackCount
    stats.value.news = newsCount
    stats.value.adPositions = adPositionCount
    stats.value.ads = adCount
  } catch (err) {
    console.error('Failed to fetch stats:', err)
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped lang="less">
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700&family=Source+Sans+3:wght@400;500;600&display=swap');

@primary: #1e3a5f;
@primary-light: #2d5a8a;
@accent: #c77b30;
@success: #2e7d5a;
@steel-gradient: linear-gradient(135deg, #2d3e50 0%, #1e3a5f 50%, #0f2538 100%);

.dashboard {
  width: 100%;

  :deep(.ant-page-header-heading-title) {
    font-family: 'Outfit', sans-serif;
    font-weight: 700;
    font-size: 28px;
    letter-spacing: -0.5px;
  }
}

.stat-card {
  margin-bottom: 16px;
  border-radius: 8px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid #e8ecf0;
  overflow: hidden;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: @steel-gradient;
    opacity: 0;
    transition: opacity 0.25s ease;
  }

  &:hover {
    box-shadow: 0 4px 16px rgba(21, 28, 36, 0.1);
    transform: translateY(-4px);
    border-color: transparent;

    &::before {
      opacity: 1;
    }
  }

  :deep(.ant-card-head) {
    border-bottom: none;
    padding: 16px 20px 0 20px;

    .ant-card-head-title {
      font-family: 'Outfit', sans-serif;
      font-weight: 500;
      font-size: 14px;
      color: #5a6572;
      letter-spacing: 0.3px;
    }
  }

  :deep(.ant-card-body) {
    padding: 20px;
  }
}

.stat-number {
  font-family: 'Outfit', sans-serif;
  font-size: 48px;
  font-weight: 700;
  text-align: center;
  color: #151c24;
  padding: 20px 0 12px 0;
  line-height: 1.1;
  letter-spacing: -1px;
  background: linear-gradient(135deg, #151c24 0%, #3d4a58 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* Icon colors */
:deep(.anticon-appstore) {
  background: linear-gradient(135deg, @success 0%, #3d9970 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

:deep(.anticon-message) {
  background: linear-gradient(135deg, @accent 0%, #e5a55d 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

:deep(.anticon-file-text) {
  background: linear-gradient(135deg, #1890ff 0%, #40a9ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

:deep(.anticon-layout) {
  background: linear-gradient(135deg, #722ed1 0%, #9254de 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

:deep(.anticon-picture) {
  background: linear-gradient(135deg, #eb2f96 0%, #f759ab 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

@media (max-width: 768px) {
  .stat-number {
    font-size: 36px;
  }

  .stat-card {
    :deep(.ant-card-head) {
      padding: 14px 16px 0 16px;
    }

    :deep(.ant-card-body) {
      padding: 16px;
    }
  }
}
</style>
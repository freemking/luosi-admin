<template>
  <a-layout>
    <a-layout-sider 
      v-model="collapsed"
      :trigger="null"
      breakpoint="md"
      collapsed-width="0"
      :width="collapsed ? 0 : 208"
    >
      <div class="logo">
        <template v-if="!collapsed">
          <span class="logo-text">Admin Panel</span>
        </template>
      </div>
      <a-menu
        :selected-keys="selectedKeys"
        mode="inline"
        theme="dark"
        :inline-indent="collapsed ? 0 : 16"
        @click="handleMenuClick"
      >
        <a-menu-item key="dashboard">
          <router-link to="/">
            <template #icon><DashboardOutlined /></template>
            <span>Dashboard</span>
          </router-link>
        </a-menu-item>
        <a-menu-item key="users">
          <router-link to="/users">
            <template #icon><UserOutlined /></template>
            <span>用户管理</span>
          </router-link>
        </a-menu-item>
        <a-menu-item key="products">
          <router-link to="/products">
            <template #icon><AppstoreOutlined /></template>
            <span>产品管理</span>
          </router-link>
        </a-menu-item>
        <a-menu-item key="feedbacks">
          <router-link to="/feedbacks">
            <template #icon><MessageOutlined /></template>
            <span>反馈管理</span>
          </router-link>
        </a-menu-item>
        <a-menu-item key="news">
          <router-link to="/news">
            <template #icon><FileTextOutlined /></template>
            <span>新闻管理</span>
          </router-link>
        </a-menu-item>
        <a-menu-item key="ad-positions">
          <router-link to="/ad-positions">
            <template #icon><LayoutOutlined /></template>
            <span>广告位管理</span>
          </router-link>
        </a-menu-item>
        <a-menu-item key="ads">
          <router-link to="/ads">
            <template #icon><PictureOutlined /></template>
            <span>广告管理</span>
          </router-link>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header class="layout-header">
        <div class="header-left">
          <a-button 
            type="text" 
            @click="toggle"
          >
            <template #icon>
              <MenuUnfoldOutlined v-if="collapsed" />
              <MenuFoldOutlined v-else />
            </template>
          </a-button>
        </div>
        <div class="header-right">
          <a-space size="middle">
            <a-avatar :style="{ backgroundColor: '#1890ff' }">
              <template #icon><UserOutlined /></template>
            </a-avatar>
            <span class="username">{{ user?.username }}</span>
            <a-divider type="vertical" />
            <a-button @click="handleLogout" type="primary" danger size="small">
              退出登录
            </a-button>
          </a-space>
        </div>
      </a-layout-header>
      <a-layout-content class="layout-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  DashboardOutlined, 
  UserOutlined, 
  AppstoreOutlined, 
  MessageOutlined,
  FileTextOutlined,
  PictureOutlined,
  LayoutOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined
} from '@ant-design/icons-vue'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const collapsed = ref(false)
const selectedKeys = ref([
  window.location.pathname === '/' || window.location.pathname === '' 
    ? 'dashboard' 
    : window.location.pathname.substring(1)
])

const user = computed(() => authStore.user)
const isSuperAdmin = computed(() => authStore.user && authStore.user.role === 'super')

const toggle = () => {
  collapsed.value = !collapsed.value
}

const handleMenuClick = (e) => {
  selectedKeys.value = [e.key]
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const handleResize = () => {
  if (window.innerWidth < 768) {
    collapsed.value = true
  }
}

onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)
})

watch(() => window.location.pathname, () => {
  selectedKeys.value = [
    window.location.pathname === '/' || window.location.pathname === '' 
      ? 'dashboard' 
      : window.location.pathname.substring(1)
  ]
})
</script>

<style scoped lang="less">
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700&family=Source+Sans+3:wght@400;500;600&display=swap');

@primary: #1e3a5f;
@primary-light: #2d5a8a;
@primary-dark: #0f2538;
@sider-bg: #0a1628;
@sider-dark: #060e1c;
@accent: #c77b30;

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(180deg, #0c1a2d 0%, @sider-bg 100%);
  padding: 0 16px;
  position: relative;

  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 16px;
    right: 16px;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(255,255,255,0.08), transparent);
  }

  .logo-text {
    font-family: 'Outfit', sans-serif;
    color: #fff;
    font-size: 16px;
    font-weight: 600;
    letter-spacing: 0.5px;
  }
}

.layout-header {
  background: #fff;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(21, 28, 36, 0.06);
  height: 64px;
  line-height: 64px;
  position: relative;
  z-index: 10;

  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(30, 58, 95, 0.1), transparent);
  }
}

.header-left {
  display: flex;
  align-items: center;

  :deep(.ant-btn) {
    color: @primary;
    border-radius: 4px;
    transition: all 0.2s ease;

    &:hover {
      background: rgba(30, 58, 95, 0.05);
      color: @primary-light;
    }
  }
}

.header-right {
  display: flex;
  align-items: center;

  .username {
    font-family: 'Outfit', sans-serif;
    color: #151c24;
    font-size: 14px;
    font-weight: 500;
  }

  :deep(.ant-avatar) {
    background: linear-gradient(135deg, @primary 0%, @primary-light 100%) !important;
  }

  :deep(.ant-divider) {
    border-color: #e8ecf0;
  }

  :deep(.ant-btn-dangerous.ant-btn-primary) {
    font-family: 'Outfit', sans-serif;
    font-weight: 500;
    letter-spacing: 0.3px;
    border-radius: 4px;
    transition: all 0.2s ease;
    background: #5a6572;
    border-color: #5a6572;

    &:hover {
      background: #3d4a58;
      border-color: #3d4a58;
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(90, 101, 114, 0.3);
    }
  }
}

.layout-content {
  margin: 24px;
  padding: 28px;
  background: #fff;
  min-height: calc(100vh - 64px - 48px);
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(21, 28, 36, 0.04);
  border: 1px solid #e8ecf0;
}

/* Sider menu styling */
:deep(.ant-layout-sider) {
  background: @sider-bg !important;

  .ant-menu {
    background: transparent;
  }

  .ant-menu-item {
    margin: 4px 8px;
    border-radius: 4px;
    font-family: 'Source Sans 3', sans-serif;
    transition: all 0.2s ease;

    &:hover {
      background: rgba(255, 255, 255, 0.08);
    }

    &.ant-menu-item-selected {
      background: @primary !important;
      box-shadow: 0 2px 8px rgba(30, 58, 95, 0.3);
    }
  }

  .ant-menu-item a {
    color: rgba(255, 255, 255, 0.7);
    transition: color 0.2s ease;

    &:hover {
      color: #fff;
    }
  }

  .ant-menu-item-selected a {
    color: #fff;
  }
}

@media (max-width: 768px) {
  .layout-content {
    margin: 16px;
    padding: 20px;
  }

  .header-right {
    .username {
      display: none;
    }
  }
}
</style>

<template>
  <div class="product-management">
    <a-page-header title="产品管理" style="padding: 0 0 16px 0;" />
    <div class="toolbar">
      <a-button type="primary" @click="showAddModal">
        <PlusOutlined />
        新建产品
      </a-button>
    </div>
    <a-card :bordered="false">
      <a-skeleton :loading="loading" active>
        <a-table
          :data-source="products"
          :loading="loading"
          :columns="columns"
          rowKey="id"
          :scroll="{ x: 800 }"
          :pagination="{
            current: currentPage,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条`,
            size: 'middle',
            onChange: handlePaginationChange
          }"
          :row-hover="true"
          :bordered="false"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'action'">
              <a-space size="small">
                <a-button size="small" @click="handleView(record)">查看</a-button>
                <a-button size="small" @click="handleEdit(record)">编辑</a-button>
                <a-button size="small" danger @click="showDeleteModal(record)">删除</a-button>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-skeleton>
    </a-card>

    <a-modal
      v-model:open="deleteModalVisible"
      title="确认删除"
      @ok="handleDelete"
      :confirmLoading="deleting"
      ok-text="确认删除"
      cancel-text="取消"
    >
      <a-alert
        message="警告"
        description="确定要删除此产品吗？删除后无法恢复。"
        type="warning"
        show-icon
      />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { useProductStore } from '../stores/auth'
import { message } from 'ant-design-vue'

const router = useRouter()
const route = useRoute()
const productStore = useProductStore()
const loading = ref(true)
const deleteModalVisible = ref(false)
const deleting = ref(false)
const currentId = ref(null)

// 从URL查询参数中读取分页信息
const currentPage = ref(parseInt(route.query.page) || 1)
const pageSize = ref(parseInt(route.query.pageSize) || 10)
const total = ref(0)

const productForm = ref({
  name: '',
  description: '',
  category_name: '',
  category_slug: '',
  standard: '',
  material: '',
  images: []
})

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 60,
    fixed: 'left'
  },
  {
    title: '产品名称',
    dataIndex: 'name',
    key: 'name',
    width: 150
  },
  {
    title: '分类',
    dataIndex: 'category_name',
    key: 'category_name',
    width: 120
  },
  {
    title: '修改时间',
    dataIndex: 'updated_at',
    key: 'updated_at',
    width: 170
  },
  {
    title: '操作',
    key: 'action',
    width: 120,
    fixed: 'right'
  }
]

const products = computed(() => productStore.products)

const fetchProducts = async () => {
  try {
    loading.value = true
    const result = await productStore.getProducts(currentPage.value, pageSize.value)
    total.value = result.total
  } catch (err) {
    message.error('获取产品列表失败')
  } finally {
    loading.value = false
  }
}

const handleView = (record) => {
  window.open(record.view_url, '_blank')
}

const showAddModal = () => {
  router.push('/products/create')
}

const handleEdit = (record) => {
  router.push(`/products/${record.id}`)
}

const showDeleteModal = (record) => {
  currentId.value = record.id
  deleteModalVisible.value = true
}

const handleDelete = async () => {
  try {
    deleting.value = true
    await productStore.deleteProduct(currentId.value)
    message.success('产品删除成功')
    deleteModalVisible.value = false
    fetchProducts()
  } catch (err) {
    message.error(productStore.error || '删除产品失败')
  } finally {
    deleting.value = false
  }
}

// 分页变化处理函数
const handlePaginationChange = (current, size) => {
  currentPage.value = current
  pageSize.value = size
  // 更新URL查询参数
  router.replace({
    query: {
      page: current,
      pageSize: size
    }
  })
  fetchProducts()
}

onMounted(() => {
  fetchProducts()
})
</script>

<style scoped lang="less">
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700&family=Source+Sans+3:wght@400;500;600&display=swap');

@primary: #1e3a5f;
@primary-light: #2d5a8a;
@accent: #c77b30;

.product-management {
  width: 100%;

  :deep(.ant-page-header-heading-title) {
    font-family: 'Outfit', sans-serif;
    font-weight: 700;
    font-size: 24px;
    letter-spacing: -0.3px;
  }

  :deep(.ant-card) {
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(21, 28, 36, 0.04);
  }

  :deep(.ant-table-thead > tr > th) {
    font-family: 'Outfit', sans-serif;
    font-weight: 600;
    background: #f8fafb !important;
    color: #151c24;
    font-size: 13px;
    letter-spacing: 0.3px;
  }

  :deep(.ant-table-tbody > tr:hover > td) {
    background: #f4f6f8 !important;
  }

  :deep(.ant-btn-primary) {
    background: @primary;
    border-color: @primary;
    border-radius: 4px;
    font-family: 'Outfit', sans-serif;
    font-weight: 500;
    transition: all 0.2s ease;

    &:hover {
      background: @primary-light;
      border-color: @primary-light;
      transform: translateY(-1px);
    }
  }

  :deep(.ant-btn-dangerous) {
    border-radius: 4px;
    font-family: 'Outfit', sans-serif;
    font-weight: 500;
  }

  :deep(.ant-modal-content) {
    border-radius: 8px;
    overflow: hidden;
  }
}

.toolbar {
  margin-bottom: 20px;
  text-align: right;

  :deep(.ant-btn) {
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }
}

.image-list {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.image-thumb {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 4px;
  border: 1px solid #e8ecf0;
  transition: all 0.2s ease;
  cursor: pointer;

  &:hover {
    border-color: @primary;
    transform: scale(1.1);
    box-shadow: 0 2px 8px rgba(21, 28, 36, 0.12);
  }
}

.image-item {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  align-items: center;
}

@media (max-width: 768px) {
  .toolbar {
    text-align: left;
  }

  .image-item {
    flex-direction: column;
    align-items: stretch;

    .ant-input {
      width: 100% !important;
    }
  }
}
</style>

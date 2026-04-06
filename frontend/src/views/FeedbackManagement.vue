<template>
  <div class="feedback-management">
    <a-page-header title="反馈管理" style="padding: 0 0 16px 0;" />
    <a-card :bordered="false">
      <a-skeleton :loading="loading" active>
        <a-table
          :data-source="feedbacks"
          :loading="loading"
          :columns="columns"
          rowKey="id"
          :scroll="{ x: 700 }"
          :pagination="{
            pageSize: 10,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条`,
            size: 'middle'
          }"
          :row-hover="true"
          :bordered="false"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'action'">
              <a-space size="small">
                <a-button size="small" type="primary" @click="viewFeedback(record)">
                  查看
                </a-button>
                <a-button size="small" danger @click="showDeleteModal(record)">
                  删除
                </a-button>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-skeleton>
    </a-card>

    <a-modal
      v-model:open="modalVisible"
      title="反馈详情"
      width="640px"
      destroyOnClose
    >
      <template #footer>
        <a-button @click="closeModal">关闭</a-button>
      </template>
      <div v-if="selectedFeedback">
        <a-descriptions :column="1" bordered :colon="false">
          <a-descriptions-item label="ID">{{ selectedFeedback.id }}</a-descriptions-item>
          <a-descriptions-item label="姓名">{{ selectedFeedback.name }}</a-descriptions-item>
          <a-descriptions-item label="邮箱">{{ selectedFeedback.email }}</a-descriptions-item>
          <a-descriptions-item label="电话">{{ selectedFeedback.phone }}</a-descriptions-item>
          <a-descriptions-item label="公司">{{ selectedFeedback.company }}</a-descriptions-item>
          <a-descriptions-item label="产品">{{ selectedFeedback.product }}</a-descriptions-item>
          <a-descriptions-item label="留言">
            <pre style="margin: 0; white-space: pre-wrap; word-break: break-word; font-family: inherit;">{{ selectedFeedback.message }}</pre>
          </a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ selectedFeedback.created_at }}</a-descriptions-item>
        </a-descriptions>
      </div>
    </a-modal>

    <a-modal
      v-model:open="deleteModalVisible"
      title="确认删除"
      @ok="handleDelete"
      :confirmLoading="deleting"
      ok-text="确认"
      cancel-text="取消"
    >
      <a-alert
        message="警告"
        description="确定要删除此留言吗？删除后无法恢复。"
        type="warning"
        show-icon
      />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useFeedbackStore } from '../stores/auth'
import { message } from 'ant-design-vue'
import config from '../config'

const feedbackStore = useFeedbackStore()

const feedbacks = ref([])
const loading = ref(true)
const modalVisible = ref(false)
const deleteModalVisible = ref(false)
const selectedFeedback = ref(null)
const currentId = ref(null)
const deleting = ref(false)

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 60,
    fixed: 'left'
  },
  {
    title: '姓名',
    dataIndex: 'name',
    key: 'name',
    width: 120
  },
  {
    title: '邮箱',
    dataIndex: 'email',
    key: 'email',
    width: 180
  },
  {
    title: '产品',
    dataIndex: 'product',
    key: 'product',
    width: 120
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    width: 180
  },
  {
    title: '操作',
    key: 'action',
    width: 120,
    fixed: 'right'
  }
]

const fetchFeedbacks = async () => {
  try {
    loading.value = true
    await feedbackStore.getFeedbacks()
    feedbacks.value = feedbackStore.feedbacks
  } catch (err) {
    message.error('获取反馈列表失败')
  } finally {
    loading.value = false
  }
}

const viewFeedback = async (feedback) => {
  try {
    const data = await feedbackStore.getFeedback(feedback.id)
    selectedFeedback.value = data
    modalVisible.value = true
  } catch (err) {
    message.error('获取反馈详情失败')
  }
}

const closeModal = () => {
  modalVisible.value = false
  selectedFeedback.value = null
}

const showDeleteModal = (record) => {
  currentId.value = record.id
  deleteModalVisible.value = true
}

const handleDelete = async () => {
  deleting.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`${config.API_BASE_URL}/feedbacks/${currentId.value}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    if (response.ok) {
      message.success(data.message || '删除成功')
      deleteModalVisible.value = false
      fetchFeedbacks()
    } else {
      message.error(data.error || '删除失败')
    }
  } catch (error) {
    message.error('网络错误')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchFeedbacks()
})
</script>

<style scoped lang="less">
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700&family=Source+Sans+3:wght@400;500;600&display=swap');

@primary: #1e3a5f;
@primary-light: #2d5a8a;
@accent: #c77b30;

.feedback-management {
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

  :deep(.ant-modal-content) {
    border-radius: 8px;
    overflow: hidden;
  }

  :deep(.ant-modal-header) {
    border-bottom: 1px solid #e8ecf0;
    padding: 16px 24px;
  }

  :deep(.ant-modal-title) {
    font-family: 'Outfit', sans-serif;
    font-weight: 600;
  }

  :deep(.ant-modal-body) {
    padding: 24px;
  }

  :deep(.ant-descriptions-item-label) {
    font-family: 'Outfit', sans-serif;
    font-weight: 500;
    color: #151c24;
    font-size: 13px;
  }
}

@media (max-width: 768px) {
  :deep(.ant-descriptions-item-label) {
    width: 80px !important;
  }
}
</style>

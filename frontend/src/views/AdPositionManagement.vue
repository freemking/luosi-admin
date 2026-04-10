<template>
  <div class="ad-position-management">
    <a-page-header title="广告位管理" style="padding: 0 0 16px 0;" />
    <div class="toolbar">
      <a-button type="primary" @click="showAddModal">
        <PlusOutlined />
        新建广告位
      </a-button>
    </div>
    <a-card :bordered="false">
      <a-skeleton :loading="loading" active>
        <a-table
          :data-source="adPositions"
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
            <template v-if="column.key === 'status'">
              <a-tag :color="record.status === 1 ? 'green' : 'red'">
                {{ record.status === 1 ? '启用' : '禁用' }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'action'">
              <a-space size="small">
                <a-button size="small" @click="handleEdit(record)">编辑</a-button>
                <a-button size="small" danger @click="showDeleteModal(record)">删除</a-button>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-skeleton>
    </a-card>

    <a-modal
      v-model:open="modalVisible"
      :title="isEditing ? '编辑广告位' : '新建广告位'"
      @ok="handleSubmit"
      :confirmLoading="submitting"
      ok-text="确认"
      cancel-text="取消"
      width="600px"
    >
      <a-form :model="adPositionForm" layout="vertical">
        <a-form-item label="编码" name="code" :rules="[{ required: true, message: '请输入编码' }]">
          <a-input v-model:value="adPositionForm.code" placeholder="例如：home" />
        </a-form-item>
        <a-form-item label="名称" name="name" :rules="[{ required: true, message: '请输入名称' }]">
          <a-input v-model:value="adPositionForm.name" placeholder="例如：首页轮播" />
        </a-form-item>
        <a-form-item label="描述" name="description">
          <a-textarea v-model:value="adPositionForm.description" placeholder="描述广告位的位置和用途" :rows="3" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="宽度" name="width">
              <a-input-number v-model:value="adPositionForm.width" :min="0" placeholder="像素宽度" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="高度" name="height">
              <a-input-number v-model:value="adPositionForm.height" :min="0" placeholder="像素高度" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="状态" name="status">
          <a-radio-group v-model:value="adPositionForm.status">
            <a-radio :value="1">启用</a-radio>
            <a-radio :value="0">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>

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
        description="确定要删除此广告位吗？删除后该位置下的所有广告也会被软删除。"
        type="warning"
        show-icon
      />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import config from '../config'

const loading = ref(true)
const modalVisible = ref(false)
const deleteModalVisible = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const isEditing = ref(false)
const currentId = ref(null)

const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const adPositions = ref([])

const adPositionForm = ref({
  code: '',
  name: '',
  description: '',
  width: '',
  height: '',
  status: 1
})

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 60
  },
  {
    title: '编码',
    dataIndex: 'code',
    key: 'code',
    width: 120
  },
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    width: 150
  },
  {
    title: '尺寸',
    key: 'size',
    width: 120,
    customRender: ({ record }) => {
      return `${record.width}×${record.height}`
    }
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 80
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

const fetchAdPositions = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`${config.API_BASE_URL}/ad-positions?page=${currentPage.value}&page_size=${pageSize.value}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    if (response.ok) {
      adPositions.value = data.data
      total.value = data.total
    } else {
      message.error(data.error || '获取数据失败')
    }
  } catch (error) {
    message.error('网络错误')
  } finally {
    loading.value = false
  }
}

const showAddModal = () => {
  isEditing.value = false
  currentId.value = null
  adPositionForm.value = {
    code: '',
    name: '',
    description: '',
    width: '',
    height: '',
    status: 1
  }
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEditing.value = true
  currentId.value = record.id
  adPositionForm.value = {
    code: record.code,
    name: record.name,
    description: record.description || '',
    width: record.width || '',
    height: record.height || '',
    status: record.status
  }
  modalVisible.value = true
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    let url = `${config.API_BASE_URL}/ad-positions`
    let method = 'POST'
    
    if (isEditing.value && currentId.value) {
      url = `${url}/${currentId.value}`
      method = 'PUT'
    }
    
    const response = await fetch(url, {
      method,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(adPositionForm.value)
    })
    
    const data = await response.json()
    if (response.ok) {
      message.success(data.message || (isEditing.value ? '更新成功' : '创建成功'))
      modalVisible.value = false
      fetchAdPositions()
    } else {
      message.error(data.error || '操作失败')
    }
  } catch (error) {
    message.error('网络错误')
  } finally {
    submitting.value = false
  }
}

const showDeleteModal = (record) => {
  currentId.value = record.id
  deleteModalVisible.value = true
}

const handleDelete = async () => {
  deleting.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`${config.API_BASE_URL}/ad-positions/${currentId.value}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    if (response.ok) {
      message.success(data.message || '删除成功')
      deleteModalVisible.value = false
      fetchAdPositions()
    } else {
      message.error(data.error || '删除失败')
    }
  } catch (error) {
    message.error('网络错误')
  } finally {
    deleting.value = false
  }
}

const handlePaginationChange = (page, pageSizeVal) => {
  currentPage.value = page
  if (pageSizeVal !== pageSize.value) {
    pageSize.value = pageSizeVal
    currentPage.value = 1
  }
  fetchAdPositions()
}

onMounted(() => {
  fetchAdPositions()
})
</script>

<style scoped>
.toolbar {
  margin-bottom: 16px;
}
</style>
<template>
  <div class="ad-management">
    <a-page-header title="广告管理" style="padding: 0 0 16px 0;" />
    <div class="toolbar">
      <a-space>
        <a-select
          v-model:value="filterPositionId"
          placeholder="选择广告位"
          :options="positionOptions"
          allow-clear
          style="width: 200px"
          @change="handleFilterChange"
        />
        <a-button type="primary" @click="showAddModal">
          <PlusOutlined />
          新建广告
        </a-button>
      </a-space>
    </div>
    <a-card :bordered="false">
      <a-skeleton :loading="loading" active>
        <a-table
          :data-source="ads"
          :loading="loading"
          :columns="columns"
          rowKey="id"
          :scroll="{ x: 900 }"
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
            <template v-if="column.key === 'image'">
              <a-image
                v-if="record.image_url"
                :src="record.image_url"
                :width="80"
                :height="40"
                :preview="{ visible: false }"
                @click="handlePreview(record)"
                style="object-fit: cover; cursor: pointer"
              />
              <div v-else class="no-image-placeholder">
                <span>暂无图片</span>
              </div>
            </template>
            <template v-else-if="column.key === 'status'">
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
      :title="isEditing ? '编辑广告' : '新建广告'"
      @ok="handleSubmit"
      :confirmLoading="submitting"
      ok-text="确认"
      cancel-text="取消"
      width="700px"
    >
      <a-form :model="adForm" layout="vertical">
        <a-form-item label="广告位" name="position_id" :rules="[{ required: true, message: '请选择广告位' }]">
          <a-select v-model:value="adForm.position_id" placeholder="选择广告位" :options="positionOptions" style="width: 100%" />
        </a-form-item>
        <a-form-item label="标题" name="title">
          <a-input v-model:value="adForm.title" placeholder="广告标题（可选）" />
        </a-form-item>
        <a-form-item label="副标题" name="subtitle">
          <a-textarea v-model:value="adForm.subtitle" placeholder="广告副标题（可选）" :rows="3" />
        </a-form-item>
        <a-form-item label="图片" name="image_url" :rules="[{ required: true, message: '请上传图片' }]">
          <a-upload
            v-model:file-list="fileList"
            :action="config.getUploadUrl('ads')"
            list-type="picture-card"
            :max-count="1"
            :headers="uploadHeaders"
            name="image"
            @change="handleUploadChange"
          >
            <div v-if="fileList.length < 1">
              <PlusOutlined />
              <div style="margin-top: 8px">上传</div>
            </div>
          </a-upload>
        </a-form-item>
        <a-form-item label="链接地址" name="link_url">
          <a-input v-model:value="adForm.link_url" placeholder="点击跳转链接（可选）" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="排序" name="order">
              <a-input-number v-model:value="adForm.order" :min="0" style="width: 100%" />
              <div class="ant-form-item-extra">数字越小越靠前</div>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="状态" name="status">
              <a-radio-group v-model:value="adForm.status">
                <a-radio :value="1">启用</a-radio>
                <a-radio :value="0">禁用</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="开始时间" name="start_time">
              <a-date-picker v-model:value="adForm.start_time" show-time placeholder="选择开始时间" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="结束时间" name="end_time">
              <a-date-picker v-model:value="adForm.end_time" show-time placeholder="选择结束时间" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
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
        description="确定要删除此广告吗？删除后无法恢复。"
        type="warning"
        show-icon
      />
    </a-modal>

    <a-modal :open="previewVisible" :footer="null" @cancel="previewVisible = false">
      <img :src="previewImage" alt="预览" style="width: 100%" />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'
import config from '../config'

const loading = ref(true)
const modalVisible = ref(false)
const deleteModalVisible = ref(false)
const previewVisible = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const isEditing = ref(false)
const currentId = ref(null)
const previewImage = ref('')

const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const ads = ref([])
const positions = ref([])
const filterPositionId = ref(null)
const fileList = ref([])

const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`
}))

const adForm = ref({
  position_id: null,
  title: '',
  subtitle: '',
  image_url: '',
  link_url: '',
  order: 0,
  status: 1,
  start_time: null,
  end_time: null
})

const positionOptions = computed(() => {
  return positions.value.map(p => ({
    label: p.name,
    value: p.id
  }))
})

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 60
  },
  {
    title: '广告位',
    dataIndex: 'position_name',
    key: 'position_name',
    width: 150
  },
  {
    title: '图片',
    key: 'image',
    width: 100
  },
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title',
    width: 150
  },
  {
    title: '排序',
    dataIndex: 'order',
    key: 'order',
    width: 80
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 80
  },
  {
    title: '开始时间',
    dataIndex: 'start_time',
    key: 'start_time',
    width: 170
  },
  {
    title: '结束时间',
    dataIndex: 'end_time',
    key: 'end_time',
    width: 170
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    width: 170
  },
  {
    title: '操作',
    key: 'action',
    width: 120,
    fixed: 'right'
  }
]

const fetchPositions = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`${config.API_BASE_URL}/ad-positions/all`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    if (response.ok) {
      positions.value = data.data
    }
  } catch (error) {
    console.error('获取广告位列表失败', error)
  }
}

const fetchAds = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    let url = `${config.API_BASE_URL}/ads?page=${currentPage.value}&page_size=${pageSize.value}`
    if (filterPositionId.value) {
      url += `&position_id=${filterPositionId.value}`
    }
    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    if (response.ok) {
      ads.value = data.data
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

const handleFilterChange = () => {
  currentPage.value = 1
  fetchAds()
}

const handleUploadChange = ({ fileList: newFileList }) => {
  fileList.value = newFileList
  if (newFileList.length > 0) {
    adForm.value.image_url = newFileList[0].response?.url || newFileList[0].url || ''
  } else {
    adForm.value.image_url = ''
  }
}

const showAddModal = () => {
  isEditing.value = false
  currentId.value = null
  adForm.value = {
    position_id: filterPositionId.value || null,
    title: '',
    subtitle: '',
    image_url: '',
    link_url: '',
    order: 0,
    status: 1,
    start_time: null,
    end_time: null
  }
  fileList.value = []
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEditing.value = true
  currentId.value = record.id
  adForm.value = {
    position_id: record.position_id,
    title: record.title || '',
    subtitle: record.subtitle || '',
    image_url: record.image_url,
    link_url: record.link_url || '',
    order: record.order || 0,
    status: record.status,
    start_time: record.start_time ? dayjs(record.start_time) : null,
    end_time: record.end_time ? dayjs(record.end_time) : null
  }
  if (record.image_url) {
    fileList.value = [{
      uid: '-1',
      name: record.image_url.split('/').pop(),
      status: 'done',
      url: record.image_url,
      response: { url: record.image_url }
    }]
  } else {
    fileList.value = []
  }
  modalVisible.value = true
}

const handlePreview = (record) => {
  previewImage.value = record.image_url
  previewVisible.value = true
}

const formatDate = (date) => {
  if (!date) return null
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const handleSubmit = async () => {
  if (!adForm.value.image_url) {
    message.error('请上传图片')
    return
  }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    let url = `${config.API_BASE_URL}/ads`
    let method = 'POST'
    
    const data = {
      ...adForm.value,
      start_time: formatDate(adForm.value.start_time),
      end_time: formatDate(adForm.value.end_time)
    }
    
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
      body: JSON.stringify(data)
    })
    
    const result = await response.json()
    if (response.ok) {
      message.success(result.message || (isEditing.value ? '更新成功' : '创建成功'))
      modalVisible.value = false
      fetchAds()
    } else {
      message.error(result.error || '操作失败')
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
    const response = await fetch(`${config.API_BASE_URL}/ads/${currentId.value}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    if (response.ok) {
      message.success(data.message || '删除成功')
      deleteModalVisible.value = false
      fetchAds()
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
  fetchAds()
}

onMounted(() => {
  fetchPositions()
  fetchAds()
})
</script>

<style scoped>
.toolbar {
  margin-bottom: 16px;
}

.no-image-placeholder {
  width: 80px;
  height: 40px;
  background: #f5f5f5;
  border: 1px dashed #d9d9d9;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 12px;
}
</style>
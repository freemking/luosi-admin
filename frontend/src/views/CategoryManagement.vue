<template>
  <div class="category-management">
    <a-page-header title="产品分类管理" style="padding: 0 0 16px 0;" />
    <div class="toolbar">
      <a-space>
        <a-button type="primary" @click="showAddModal">
          <PlusOutlined />
          新建分类
        </a-button>
      </a-space>
    </div>
    <a-card :bordered="false">
      <a-skeleton :loading="loading" active>
        <a-table
          :data-source="categories"
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
            <template v-if="column.key === 'icon'">
              <a-image
                v-if="record.icon"
                :src="record.icon"
                :width="40"
                :height="40"
                style="object-fit: cover"
              />
              <div v-else class="no-image-placeholder small">
                <span>暂无</span>
              </div>
            </template>
            <template v-else-if="column.key === 'image'">
              <a-image
                v-if="record.image_url"
                :src="record.image_url"
                :width="80"
                :height="40"
                :preview="{ visible: false }"
                @click="handlePreview(record.image_url)"
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
      :title="isEditing ? '编辑分类' : '新建分类'"
      @ok="handleSubmit"
      :confirmLoading="submitting"
      ok-text="确认"
      cancel-text="取消"
      width="700px"
    >
      <a-form :model="categoryForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="分类名称" name="name" :rules="[{ required: true, message: '请输入分类名称' }]">
              <a-input v-model:value="categoryForm.name" placeholder="分类名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="分类标识" name="slug" :rules="[{ required: true, message: '请输入分类标识' }]">
              <a-input v-model:value="categoryForm.slug" placeholder="分类标识（英文）" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="描述" name="description">
          <a-textarea v-model:value="categoryForm.description" placeholder="分类描述（可选）" :rows="3" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="图标" name="icon">
              <a-upload
                v-model:file-list="iconFileList"
                :action="config.getUploadUrl('categories')"
                list-type="picture-card"
                :max-count="1"
                :headers="uploadHeaders"
                name="image"
                @change="handleIconUploadChange"
              >
                <div v-if="iconFileList.length < 1">
                  <PlusOutlined />
                  <div style="margin-top: 8px">上传图标</div>
                </div>
              </a-upload>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="图片" name="image_url">
              <a-upload
                v-model:file-list="imageFileList"
                :action="config.getUploadUrl('categories')"
                list-type="picture-card"
                :max-count="1"
                :headers="uploadHeaders"
                name="image"
                @change="handleImageUploadChange"
              >
                <div v-if="imageFileList.length < 1">
                  <PlusOutlined />
                  <div style="margin-top: 8px">上传图片</div>
                </div>
              </a-upload>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="排序" name="order">
              <a-input-number v-model:value="categoryForm.order" :min="0" style="width: 100%" />
              <div class="ant-form-item-extra">数字越小越靠前</div>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="状态" name="status">
              <a-radio-group v-model:value="categoryForm.status">
                <a-radio :value="1">启用</a-radio>
                <a-radio :value="0">禁用</a-radio>
              </a-radio-group>
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
        description="确定要删除此分类吗？删除后无法恢复。"
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
const categories = ref([])
const iconFileList = ref([])
const imageFileList = ref([])

const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`
}))

const categoryForm = ref({
  name: '',
  slug: '',
  description: '',
  icon: '',
  image_url: '',
  order: 0,
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
    title: '图标',
    key: 'icon',
    width: 80
  },
  {
    title: '图片',
    key: 'image',
    width: 100
  },
  {
    title: '分类名称',
    dataIndex: 'name',
    key: 'name',
    width: 150
  },
  {
    title: '分类标识',
    dataIndex: 'slug',
    key: 'slug',
    width: 120
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    width: 200,
    ellipsis: true
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

const fetchCategories = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const url = `${config.API_BASE_URL}/categories?page=${currentPage.value}&page_size=${pageSize.value}`
    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    if (response.ok) {
      categories.value = data.data
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

const handleIconUploadChange = ({ fileList: newFileList }) => {
  iconFileList.value = newFileList
  if (newFileList.length > 0) {
    categoryForm.value.icon = newFileList[0].response?.url || newFileList[0].url || ''
  } else {
    categoryForm.value.icon = ''
  }
}

const handleImageUploadChange = ({ fileList: newFileList }) => {
  imageFileList.value = newFileList
  if (newFileList.length > 0) {
    categoryForm.value.image_url = newFileList[0].response?.url || newFileList[0].url || ''
  } else {
    categoryForm.value.image_url = ''
  }
}

const showAddModal = () => {
  isEditing.value = false
  currentId.value = null
  categoryForm.value = {
    name: '',
    slug: '',
    description: '',
    icon: '',
    image_url: '',
    order: 0,
    status: 1
  }
  iconFileList.value = []
  imageFileList.value = []
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEditing.value = true
  currentId.value = record.id
  categoryForm.value = {
    name: record.name || '',
    slug: record.slug || '',
    description: record.description || '',
    icon: record.icon || '',
    image_url: record.image_url || '',
    order: record.order || 0,
    status: record.status
  }
  if (record.icon) {
    iconFileList.value = [{
      uid: '-1',
      name: record.icon.split('/').pop(),
      status: 'done',
      url: record.icon,
      response: { url: record.icon }
    }]
  } else {
    iconFileList.value = []
  }
  if (record.image_url) {
    imageFileList.value = [{
      uid: '-2',
      name: record.image_url.split('/').pop(),
      status: 'done',
      url: record.image_url,
      response: { url: record.image_url }
    }]
  } else {
    imageFileList.value = []
  }
  modalVisible.value = true
}

const handlePreview = (url) => {
  previewImage.value = url
  previewVisible.value = true
}

const handleSubmit = async () => {
  if (!categoryForm.value.name) {
    message.error('请输入分类名称')
    return
  }
  if (!categoryForm.value.slug) {
    message.error('请输入分类标识')
    return
  }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    let url = `${config.API_BASE_URL}/categories`
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
      body: JSON.stringify(categoryForm.value)
    })
    
    const result = await response.json()
    if (response.ok) {
      message.success(result.message || (isEditing.value ? '更新成功' : '创建成功'))
      modalVisible.value = false
      fetchCategories()
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
    const response = await fetch(`${config.API_BASE_URL}/categories/${currentId.value}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    if (response.ok) {
      message.success(data.message || '删除成功')
      deleteModalVisible.value = false
      fetchCategories()
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
  fetchCategories()
}

onMounted(() => {
  fetchCategories()
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

.no-image-placeholder.small {
  width: 40px;
  height: 40px;
  font-size: 10px;
}
</style>
<template>
  <div class="news-management">
    <a-page-header title="新闻管理" style="padding: 0 0 16px 0;" />
    <div class="toolbar">
      <a-button type="primary" @click="showAddModal">
        <PlusOutlined />
        新建新闻
      </a-button>
    </div>
    <a-card :bordered="false">
      <a-skeleton :loading="loading" active>
        <a-table
          :data-source="news"
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
            <template v-if="column.key === 'cover_image'">
              <img 
                v-if="record.cover_image" 
                :src="record.cover_image" 
                class="cover-thumb"
                alt="封面"
              />
              <span v-else style="color: #666666">暂无封面</span>
            </template>
            <template v-if="column.key === 'status'">
              <a-tag :color="record.status === 1 ? 'success' : 'default'">
                {{ record.status === 1 ? '已发布' : '草稿' }}
              </a-tag>
            </template>
            <template v-if="column.key === 'publish_date'">
              {{ formatDate(record.publish_date) }}
            </template>
            <template v-if="column.key === 'action'">
              <a-space size="small">
                <a-button size="small" @click="handleEdit(record)">编辑</a-button>
                <a-button size="small" danger @click="showDeleteModal(record)">删除</a-button>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-skeleton>
    </a-card>

    <!-- Add/Edit Modal -->
    <a-modal
      v-model:open="modalVisible"
      :title="isEditing ? '编辑新闻' : '新建新闻'"
      @ok="handleSubmit"
      :confirmLoading="submitting"
      width="900px"
      destroyOnClose
    >
      <a-form
        :model="newsForm"
        layout="vertical"
        :colon="false"
      >
        <a-form-item
          label="标题"
          name="title"
          :rules="[{ required: true, message: '请输入标题' }]"
        >
          <a-input v-model:value="newsForm.title" placeholder="请输入标题" />
        </a-form-item>
        
        <a-form-item label="封面图片">
          <a-upload
            v-model:file-list="fileList"
            :action="config.getUploadUrl('news')"
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

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="发布日期">
              <a-date-picker 
                v-model:value="publishDateValue" 
                style="width: 100%"
                valueFormat="YYYY-MM-DD"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="状态">
              <a-select v-model:value="newsForm.status">
                <a-select-option :value="1">已发布</a-select-option>
                <a-select-option :value="0">草稿</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="摘要">
          <a-textarea 
            v-model:value="newsForm.summary" 
            placeholder="请输入摘要"
            :rows="2"
            :maxlength="500"
            showCount
          />
        </a-form-item>

        <a-form-item label="内容">
          <Ckeditor
            v-model="newsForm.content"
            :editor="editor"
            :config="editorConfig"
            :key="editorKey"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Delete Modal -->
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
        description="确定要删除此新闻吗？删除后无法恢复。"
        type="warning"
        show-icon
      />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useNewsStore } from '../stores/auth'
import { message } from 'ant-design-vue'
import config from '../config'
import dayjs from 'dayjs'
import { Ckeditor } from '@ckeditor/ckeditor5-vue'
import ClassicEditor from '@ckeditor/ckeditor5-build-classic'

const newsStore = useNewsStore()
const loading = ref(true)
const modalVisible = ref(false)
const deleteModalVisible = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const isEditing = ref(false)
const currentId = ref(null)

// CKEditor
const editorKey = ref(0)
const editor = ClassicEditor

const editorConfig = {
  toolbar: [
    'heading', '|', 'bold', 'italic', 'link', 'bulletedList', 'numberedList', '|',
    'outdent', 'indent', '|', 'blockQuote', 'insertTable', 'undo', 'redo'
  ],
  table: {
    contentToolbar: ['tableColumn', 'tableRow', 'mergeTableCells']
  }
}

// Pagination
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// Upload
const fileList = ref([])

const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`
}))

const newsForm = ref({
  title: '',
  cover_image: '',
  publish_date: '',
  summary: '',
  content: '',
  status: 1
})

const publishDateValue = ref(null)

// Watch date picker changes
watch(publishDateValue, (val) => {
  newsForm.value.publish_date = val || ''
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
    title: '标题',
    dataIndex: 'title',
    key: 'title',
    width: 200,
    ellipsis: true
  },
  {
    title: '封面',
    key: 'cover_image',
    width: 100
  },
  {
    title: '发布日期',
    dataIndex: 'publish_date',
    key: 'publish_date',
    width: 120
  },
  {
    title: '状态',
    key: 'status',
    width: 100
  },
  {
    title: '修改时间',
    dataIndex: 'updated_at',
    key: 'updated_at',
    width: 180
  },
  {
    title: '操作',
    key: 'action',
    width: 120,
    fixed: 'right'
  }
]

const news = computed(() => newsStore.news)

const fetchNews = async () => {
  try {
    loading.value = true
    const result = await newsStore.getNews(currentPage.value, pageSize.value)
    total.value = result.total
  } catch (err) {
    message.error('获取新闻列表失败')
  } finally {
    loading.value = false
  }
}

const handleUploadChange = ({ fileList: newFileList }) => {
  fileList.value = newFileList
  if (newFileList.length > 0) {
    // Use response.url (relative path) for newly uploaded files
    // or the existing url for files from edit mode
    newsForm.value.cover_image = newFileList[0].response?.url || newFileList[0].url || ''
  } else {
    newsForm.value.cover_image = ''
  }
}

const showAddModal = () => {
  isEditing.value = false
  currentId.value = null
  newsForm.value = {
    title: '',
    cover_image: '',
    publish_date: '',
    summary: '',
    content: '',
    status: 1
  }
  publishDateValue.value = null
  fileList.value = []
  modalVisible.value = true
}

const handleEdit = async (record) => {
  isEditing.value = true
  currentId.value = record.id
  try {
    const newsItem = await newsStore.getNewsItem(record.id)
    newsForm.value = {
      title: newsItem.title || '',
      cover_image: newsItem.cover_image || '',
      publish_date: newsItem.publish_date || '',
      summary: newsItem.summary || '',
      content: newsItem.content || '',
      status: newsItem.status
    }
    publishDateValue.value = newsItem.publish_date ? dayjs(newsItem.publish_date) : null
    
    // Set file list for cover image
    // Store full URL in 'url' for display, but the form already has the full URL
    // which will be sent to backend (backend will convert to relative path)
    if (newsItem.cover_image) {
      fileList.value = [{
        uid: '-1',
        name: newsItem.cover_image.split('/').pop(),
        status: 'done',
        url: newsItem.cover_image,
        response: { url: newsItem.cover_image }
      }]
    } else {
      fileList.value = []
    }
    
    modalVisible.value = true
  } catch (err) {
    message.error('获取新闻详情失败')
  }
}

const handleSubmit = async () => {
  if (!newsForm.value.title) {
    message.error('请输入标题')
    return
  }
  
  try {
    submitting.value = true
    const submitData = { ...newsForm.value }
    
    if (isEditing.value) {
      await newsStore.updateNews(currentId.value, submitData)
      message.success('新闻更新成功')
    } else {
      await newsStore.createNews(submitData)
      message.success('新闻创建成功')
    }
    modalVisible.value = false
  } catch (err) {
    message.error(newsStore.error || '保存新闻失败')
    return
  } finally {
    submitting.value = false
  }
  
  // Refresh list outside of try-catch to avoid showing error if refresh fails
  fetchNews()
}

const showDeleteModal = (record) => {
  currentId.value = record.id
  deleteModalVisible.value = true
}

const handleDelete = async () => {
  try {
    deleting.value = true
    await newsStore.deleteNews(currentId.value)
    message.success('新闻删除成功')
    deleteModalVisible.value = false
    fetchNews()
  } catch (err) {
    message.error(newsStore.error || '删除新闻失败')
  } finally {
    deleting.value = false
  }
}

const handlePaginationChange = (current, size) => {
  currentPage.value = current
  pageSize.value = size
  fetchNews()
}

// Format date to YYYY-MM-DD
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  // Handle if dateStr is already in YYYY-MM-DD format
  if (typeof dateStr === 'string' && dateStr.match(/^\d{4}-\d{2}-\d{2}$/)) {
    return dateStr
  }
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return '-'
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

onMounted(() => {
  fetchNews()
})
</script>

<style scoped lang="less">
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700&family=Source+Sans+3:wght@400;500;600&display=swap');

@primary: #1e3a5f;
@primary-light: #2d5a8a;
@accent: #c77b30;

.news-management {
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

  :deep(.ant-form-item-label label) {
    font-family: 'Outfit', sans-serif;
    font-weight: 500;
    color: #151c24;
    font-size: 13px;
  }

  :deep(.ant-upload-wrapper) {
    .ant-upload-select {
      border-radius: 6px;
      border-color: #e8ecf0;
      transition: all 0.2s ease;

      &:hover {
        border-color: @primary;
      }
    }
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

.cover-thumb {
  width: 60px;
  height: 40px;
  object-fit: cover;
  border-radius: 4px;
  border: 1px solid #e8ecf0;
}

@media (max-width: 768px) {
  .toolbar {
    text-align: left;
  }
}
</style>

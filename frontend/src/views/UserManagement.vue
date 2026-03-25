<template>
  <div class="user-management">
    <a-page-header title="用户管理" style="padding: 0 0 16px 0;" />
    <div class="toolbar">
      <a-button type="primary" @click="showAddModal">
        <PlusOutlined />
        新建用户
      </a-button>
    </div>
    <a-card :bordered="false">
      <a-skeleton :loading="loading" active>
        <a-table
          :data-source="users"
          :loading="loading"
          :columns="columns"
          rowKey="id"
          :scroll="{ x: 600 }"
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
            <template v-if="column.key === 'actions'">
              <a-button type="primary" size="small" style="margin-right: 8px" @click="handleEdit(record)">
                编辑
              </a-button>
              <a-button danger size="small" @click="showDeleteModal(record)">
                删除
              </a-button>
            </template>
          </template>
        </a-table>
      </a-skeleton>
    </a-card>

    <a-modal
      v-model:open="addModalVisible"
      :title="isEditing ? '编辑用户' : '新建用户'"
      @ok="handleSubmit"
      :confirmLoading="submitting"
      width="520px"
      destroyOnClose
    >
      <a-form
        :model="userForm"
        layout="horizontal"
        :colon="false"
        label-col="{ span: 6 }"
        wrapper-col="{ span: 16 }"
      >
        <a-form-item
          label="用户名"
          name="username"
          :rules="[{ required: true, message: '请输入用户名' }]"
        >
          <a-input v-model="userForm.username" placeholder="请输入用户名" />
        </a-form-item>
        <a-form-item
          label="密码"
          name="password"
          :rules="[{ required: !isEditing, message: '请输入密码' }]"
        >
          <a-input-password v-model="userForm.password" placeholder="请输入密码" />
        </a-form-item>
        <a-form-item
          label="角色"
          name="role"
          :rules="[{ required: true, message: '请选择角色' }]"
        >
          <a-select v-model="userForm.role" placeholder="请选择角色">
            <a-select-option value="user">普通用户</a-select-option>
          </a-select>
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
        description="确定要删除此用户吗？删除后无法恢复。"
        type="warning"
        show-icon
      />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '../stores/auth'
import { message } from 'ant-design-vue'

const userStore = useUserStore()
const loading = ref(true)
const addModalVisible = ref(false)
const deleteModalVisible = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const isEditing = ref(false)
const currentId = ref(null)

const userForm = ref({
  username: '',
  password: '',
  role: 'user'
})

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 80,
    fixed: 'left'
  },
  {
    title: '用户名',
    dataIndex: 'username',
    key: 'username',
    width: 150
  },
  {
    title: '角色',
    dataIndex: 'role',
    key: 'role',
    width: 100
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    width: 180
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    fixed: 'right'
  }
]

const users = computed(() => {
  return userStore.users.filter(user => user.role !== 'super')
})

const fetchUsers = async () => {
  try {
    loading.value = true
    await userStore.getUsers()
  } catch (err) {
    message.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const showAddModal = () => {
  isEditing.value = false
  currentId.value = null
  userForm.value = {
    username: '',
    password: '',
    role: 'user'
  }
  addModalVisible.value = true
}

const handleEdit = (record) => {
  isEditing.value = true
  currentId.value = record.id
  userForm.value = {
    username: record.username,
    password: '',
    role: record.role
  }
  addModalVisible.value = true
}

const handleSubmit = async () => {
  try {
    submitting.value = true
    if (isEditing.value) {
      await userStore.updateUser(currentId.value, userForm.value)
      message.success('用户更新成功')
    } else {
      await userStore.createUser(userForm.value)
      message.success('用户创建成功')
    }
    addModalVisible.value = false
    fetchUsers()
  } catch (err) {
    message.error(userStore.error || '保存用户失败')
  } finally {
    submitting.value = false
  }
}

const showDeleteModal = (record) => {
  currentId.value = record.id
  deleteModalVisible.value = true
}

const handleDelete = async () => {
  try {
    deleting.value = true
    await userStore.deleteUser(currentId.value)
    message.success('用户删除成功')
    deleteModalVisible.value = false
    fetchUsers()
  } catch (err) {
    message.error(userStore.error || '删除用户失败')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped lang="less">
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700&family=Source+Sans+3:wght@400;500;600&display=swap');

@primary: #1e3a5f;
@primary-light: #2d5a8a;
@accent: #c77b30;

.user-management {
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

  :deep(.ant-form-item-label label) {
    font-family: 'Outfit', sans-serif;
    font-weight: 500;
    color: #151c24;
    font-size: 13px;
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

@media (max-width: 768px) {
  .toolbar {
    text-align: left;
  }
}
</style>

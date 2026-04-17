<template>
  <div class="product-detail">
    <a-page-header 
      :title="isEditing ? '编辑产品' : '新建产品'" 
      style="padding: 0 0 16px 0;"
    >
      <template #extra>
        <a-button @click="goBack">
          <ArrowLeftOutlined />
          返回列表
        </a-button>
      </template>
    </a-page-header>
    
    <a-card bordered="false">
      <a-form
        :model="productForm"
        layout="horizontal"
        :colon="false"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item
          label="产品名称"
          name="name"
          :rules="[{ required: true, message: '请输入产品名称' }]"
        >
          <a-input v-model:value="productForm.name" placeholder="请输入产品名称" />
        </a-form-item>
        <a-form-item
          label="产品分类"
          name="category"
          :rules="[{ required: true, message: '请选择产品分类' }]"
        >
          <a-select v-model:value="productForm.category" placeholder="请选择产品分类">
            <a-select-option value="Blind Rivet">Blind Rivet</a-select-option>
            <a-select-option value="Insert Nut">Insert Nut</a-select-option>
            <a-select-option value="Self Clinching Fasteners">Self Clinching Fasteners</a-select-option>
            <a-select-option value="Tools">Tools</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
          label="规格"
          name="standard"
        >
          <a-input v-model:value="productForm.standard" placeholder="请输入规格" />
        </a-form-item>
        <a-form-item
          label="材质"
          name="material"
        >
          <a-input v-model:value="productForm.material" placeholder="请输入材质" />
        </a-form-item>
        <a-form-item
          label="SEO 标题"
          name="seoTitle"
        >
          <a-input v-model:value="productForm.seoTitle" placeholder="请输入SEO标题" />
        </a-form-item>
        <a-form-item
          label="SEO 关键词"
          name="seoKeywords"
        >
          <a-textarea v-model:value="productForm.seoKeywords" placeholder="请输入SEO关键词" :rows="3" />
        </a-form-item>
        <a-form-item
          label="SEO 描述"
          name="seoDescription"
        >
          <a-textarea v-model:value="productForm.seoDescription" placeholder="请输入SEO描述" :rows="3" />
        </a-form-item>
        <a-form-item label="产品图片">
          <a-upload
            v-model:file-list="fileList"
            :action="config.getUploadUrl('products')"
            list-type="picture-card"
            :multiple="true"
            :max-count="9"
            :headers="uploadHeaders"
            name="image"
            @preview="handlePreview"
            @remove="handleRemove"
          >
            <div v-if="fileList.length < 9">
              <PlusOutlined />
              <div style="margin-top: 8px">上传</div>
            </div>
          </a-upload>
          <a-modal :open="previewVisible" :footer="null" @cancel="previewVisible = false">
            <img alt="example" style="width: 100%" :src="previewImage" />
          </a-modal>
        </a-form-item>
        <a-form-item
          label="产品描述"
          name="description"
        >
          <div v-if="loading">Loading...</div>
          <Ckeditor
            v-else
            v-model="productForm.description"
            :editor="editor"
            :config="editorConfig"
            :key="editorKey"
          />
        </a-form-item>
      </a-form>
      
      <div class="action-buttons">
        <a-button @click="goBack">取消</a-button>
        <a-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEditing ? '保存修改' : '创建产品' }}
        </a-button>
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { PlusOutlined, ArrowLeftOutlined } from '@ant-design/icons-vue'
import { useProductStore } from '../stores/auth'
import { message } from 'ant-design-vue'
import config from '../config'
import { Ckeditor } from '@ckeditor/ckeditor5-vue'
import ClassicEditor from '@ckeditor/ckeditor5-build-classic'

// CKEditor 图片上传适配器
class UploadAdapter {
  constructor(loader) {
    this.loader = loader
  }

  upload() {
    return this.loader.file.then(file => {
      return new Promise((resolve, reject) => {
        const formData = new FormData()
        formData.append('image', file)
        
        fetch(config.getUploadUrl('products'), {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
          body: formData
        })
        .then(response => response.json())
         .then(data => {
           if (data.full_url) {
             resolve({ default: data.full_url })
           } else if (data.url) {
             resolve({ default: data.url })
           } else {
             reject(data.error || '上传失败')
           }
         })
        .catch(error => {
          reject(error.message || '上传失败')
        })
      })
    })
  }

  abort() {
    // 可以实现中止上传逻辑
  }
}

// 设置编辑器插件
function MyCKEditorPlugins(editor) {
  editor.plugins.get('FileRepository').createUploadAdapter = (loader) => {
    return new UploadAdapter(loader)
  }
}

const router = useRouter()
const route = useRoute()
const productStore = useProductStore()

const submitting = ref(false)
const loading = ref(true)
const productForm = ref({
  name: '',
  description: '',
  category: '',
  standard: '',
  material: '',
  seoTitle: '',
  seoKeywords: '',
  seoDescription: '',
  images: []
})

// CKEditor
const editorKey = ref(0)
const editor = ClassicEditor

// 添加自定义插件
ClassicEditor.builtinPlugins.push(MyCKEditorPlugins)

const editorConfig = {
  toolbar: [
    'heading', '|', 'bold', 'italic', 'link', 'bulletedList', 'numberedList', '|',
    'outdent', 'indent', '|', 'blockQuote', 'insertTable', 'imageUpload', 'undo', 'redo'
  ],
  table: {
    contentToolbar: ['tableColumn', 'tableRow', 'mergeTableCells']
  },
  image: {
    toolbar: [
      'imageTextAlternative',
      'toggleImageCaption',
      '|',
      'imageStyle:inline',
      'imageStyle:block',
      'imageStyle:side'
    ]
  }
}

const fileList = ref([])
const previewVisible = ref(false)
const previewImage = ref('')

const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`
}))

const isEditing = computed(() => !!route.params.id)

const fetchProduct = async () => {
  if (isEditing.value) {
    try {
      const product = await productStore.getProduct(route.params.id)
      productForm.value.name = product.name || ''
      productForm.value.category = product.category || ''
      productForm.value.standard = product.standard || ''
      productForm.value.material = product.material || ''
      productForm.value.description = product.description || ''
      productForm.value.seoTitle = product.seo_title || ''
      productForm.value.seoKeywords = product.seo_keywords || ''
      productForm.value.seoDescription = product.seo_description || ''
      productForm.value.images = product.images ? [...product.images] : []
      
      if (product.images && product.images.length > 0) {
        fileList.value = product.images.map((img, index) => ({
          uid: String(-index - 1),
          name: img.image_url.split('/').pop(),
          status: 'done',
          url: img.image_url,
          response: { url: img.image_url }
        }))
      }
    } catch (err) {
      message.error('获取产品信息失败')
      goBack()
    }
  }
  loading.value = false
}

const handlePreview = async (file) => {
  if (!file.url && !file.preview) {
    file.preview = await getBase64(file.originFileObj)
  }
  previewImage.value = file.url || file.preview
  previewVisible.value = true
}

const getBase64 = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => resolve(reader.result)
    reader.onerror = (error) => reject(error)
  })
}

const handleRemove = (file) => {
  const index = fileList.value.indexOf(file)
  const newFileList = fileList.value.slice()
  newFileList.splice(index, 1)
  fileList.value = newFileList
}

const handleSubmit = async () => {
  try {
    submitting.value = true
    
    // Process images - handle both newly uploaded and existing images
    const images = fileList.value
      .filter(file => file.status === 'done')
      .map((file, index) => {
        // For newly uploaded files, use response.url (relative path)
        // For existing files from edit, use response.url or extract from url property
        let imageUrl = file.response?.url || file.url || ''
        return {
          image_url: imageUrl,
          order: index
        }
      })
      .filter(img => img.image_url) // Filter out empty URLs
    
    const submitData = {
      ...productForm.value,
      images
    }
    
    if (isEditing.value) {
      await productStore.updateProduct(route.params.id, submitData)
      message.success('产品更新成功')
    } else {
      await productStore.createProduct(submitData)
      message.success('产品创建成功')
    }
  } catch (err) {
    message.error(productStore.error || '保存产品失败')
    return
  } finally {
    submitting.value = false
  }
  
  // Navigate back outside of try-catch
  goBack()
}

const goBack = () => {
  router.push('/products')
}

onMounted(() => {
  fetchProduct()
})
</script>

<style scoped lang="less">
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700&family=Source+Sans+3:wght@400;500;600&display=swap');

@primary: #1e3a5f;
@primary-light: #2d5a8a;
@accent: #c77b30;

.product-detail {
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

  :deep(.ant-form-item-label label) {
    font-family: 'Outfit', sans-serif;
    font-weight: 500;
    color: #151c24;
    font-size: 13px;
  }

  :deep(.ant-input),
  :deep(.ant-input textarea) {
    font-family: 'Source Sans 3', sans-serif;
    border-radius: 4px;
    border-color: #d1d9e0;
    transition: all 0.2s ease;

    &:hover {
      border-color: @primary;
    }

    &:focus, &.ant-input-focused {
      border-color: @primary;
      box-shadow: 0 0 0 3px rgba(30, 58, 95, 0.1);
    }
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

.action-buttons {
  margin-top: 28px;
  text-align: right;
  padding-top: 20px;
  border-top: 1px solid #e8ecf0;
  
  :deep(.ant-btn) {
    margin-left: 10px;
    border-radius: 4px;
    font-family: 'Outfit', sans-serif;
    font-weight: 500;
    min-width: 100px;
  }
}

.image-item {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  align-items: center;
}

@media (max-width: 768px) {
  .image-item {
    flex-direction: column;
    align-items: stretch;

    .ant-input {
      width: 100% !important;
    }
  }

  .action-buttons {
    text-align: center;

    :deep(.ant-btn) {
      margin: 4px;
    }
  }
}
</style>
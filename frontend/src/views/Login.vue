<template>
  <div class="login-page">
    <div class="login-background"></div>
    <a-row 
      justify="center" 
      align="middle" 
      style="min-height: 100vh; position: relative; z-index: 1; padding: 16px;"
    >
      <a-col 
        :xs="24" 
        :sm="18" 
        :md="14" 
        :lg="10" 
        :xl="8"
      >
        <div class="login-header">
          <h1>管理后台</h1>
          <p>Admin Management System</p>
        </div>
        <a-card 
          :bordered="false" 
          class="login-card"
          size="default"
        >
          <a-form
            ref="loginFormRef"
            :model="loginForm"
            @submit.prevent="handleLogin"
            :colon="false"
            :label-align="'left'"
          >
            <a-form-item
              label="用户名"
              name="username"
              :rules="[{ required: true, message: '请输入用户名' }]"
              class="form-item-spacing"
            >
              <a-input 
                v-model:value="loginForm.username"
                placeholder="请输入用户名"
                size="large"
                style="width: 100%"
              />
            </a-form-item>

            <a-form-item
              label="密&nbsp&nbsp&nbsp&nbsp码"
              name="password"
              :rules="[{ required: true, message: '请输入密码' }]"
              class="form-item-spacing"
            >
              <a-input-password 
                v-model:value="loginForm.password"
                placeholder="请输入密码"
                size="large"
                @keyup.enter="handleLogin"
                style="width: 100%"
              />
            </a-form-item>

            <a-form-item
              label="验证码"
              name="captcha"
              :rules="[{ required: true, message: '请输入验证码', whitespace: true }]"
              class="form-item-spacing"
            >
              <a-row :gutter="8" align="middle" style="width: 100%">
                <a-col :span="14">
                  <a-input 
                    v-model:value="loginForm.captcha"
                    placeholder="请输入验证码"
                    size="large"
                    style="width: 100%"
                  />
                </a-col>
                <a-col :span="10">
                  <img 
                    :src="captchaImage" 
                    @click="refreshCaptcha" 
                    alt="验证码" 
                    class="captcha-img"
                  />
                </a-col>
              </a-row>
            </a-form-item>

            <div v-if="error" class="error-message">
              {{ error }}
            </div>

            <a-form-item class="form-item-spacing submit-btn">
              <a-button 
                type="primary" 
                html-type="submit" 
                :loading="loading" 
                :disabled="loading"
                block
                size="large"
              >
                {{ loading ? '登录中...' : '登录' }}
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { message } from 'ant-design-vue'

const router = useRouter()
const authStore = useAuthStore()

const loginFormRef = ref(null)
const loginForm = ref({
  username: 'admin',
  password: '1qaz@WSX',
  captcha: ''
})

const captchaId = ref('')
const captchaImage = ref('')
const loading = ref(false)
const error = ref('')

const refreshCaptcha = async () => {
  try {
    const data = await authStore.getCaptcha()
    captchaId.value = data.captcha_id
    captchaImage.value = data.captcha
    loginForm.value.captcha = ''
  } catch (err) {
    message.error('获取验证码失败')
    console.error('Failed to refresh captcha:', err)
  }
}

const handleLogin = async () => {
  if (!loginFormRef.value) {
    return
  }
  try {
    await loginFormRef.value.validate()
  } catch (err) {
    // 检查是否是因为outOfDate导致的验证失败
    if (err.outOfDate) {
      try {
        // 重新验证
        await loginFormRef.value.validate()
      } catch (secondErr) {
        return
      }
    } else {
      return
    }
  }
  error.value = ''
  try {
    loading.value = true
    await authStore.login(
      loginForm.value.username,
      loginForm.value.password,
      loginForm.value.captcha.trim(),
      captchaId.value
    )
    message.success('登录成功')
    router.push('/')
  } catch (err) {
    error.value = authStore.error || '登录失败'
    message.error(error.value)
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refreshCaptcha()
})
</script>

<style scoped lang="less">
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700&family=Source+Sans+3:wght@400;500;600&display=swap');

@primary: #1e3a5f;
@primary-light: #2d5a8a;
@primary-dark: #0f2538;
@accent: #c77b30;
@accent-light: #e5a55d;
@steel-gradient: linear-gradient(135deg, #2d3e50 0%, #1e3a5f 50%, #0f2538 100%);

.login-page {
  width: 100vw;
  height: 100vh;
  position: fixed;
  top: 0;
  left: 0;
  overflow: hidden;
  background-color: #f4f6f8;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: @steel-gradient;
  z-index: 0;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: 
      radial-gradient(ellipse at 20% 80%, rgba(199, 123, 48, 0.12) 0%, transparent 50%),
      radial-gradient(ellipse at 80% 20%, rgba(61, 124, 184, 0.1) 0%, transparent 50%);
    pointer-events: none;
  }

  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.02'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
    opacity: 0.6;
    pointer-events: none;
  }
}

.login-header {
  text-align: center;
  color: #fff;
  margin-bottom: 32px;
  padding: 0 8px;

  h1 {
    font-family: 'Outfit', sans-serif;
    font-size: 36px;
    margin: 0 0 10px 0;
    font-weight: 700;
    text-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    line-height: 1.2;
    letter-spacing: -0.5px;
  }

  p {
    font-size: 14px;
    margin: 0;
    opacity: 0.85;
    letter-spacing: 2px;
    color: rgba(255, 255, 255, 0.85);
    font-weight: 400;
  }
}

.login-card {
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  padding: 8px;
  max-width: 440px;
  margin: 0 auto;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  animation: cardSlideUp 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);

  @keyframes cardSlideUp {
    from {
      opacity: 0;
      transform: translateY(30px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  :deep(.ant-card-body) {
    padding: 32px;
  }
}

.form-item-spacing {
  margin-bottom: 20px;

  :deep(.ant-form-item-label) {
    padding-bottom: 8px;

    label {
      font-family: 'Outfit', sans-serif;
      font-size: 13px;
      font-weight: 500;
      color: #151c24;
      letter-spacing: 0.3px;
    }
  }

  :deep(.ant-input) {
    font-family: 'Source Sans 3', sans-serif;
    border-radius: 4px;
    border-color: #d1d9e0;
    padding: 12px 14px;
    height: auto;
    transition: all 0.2s ease;

    &:hover {
      border-color: @primary;
    }

    &:focus, &.ant-input-focused {
      border-color: @primary;
      box-shadow: 0 0 0 3px rgba(30, 58, 95, 0.1);
    }
  }
}

.captcha-img {
  width: 100%;
  height: 42px;
  object-fit: cover;
  cursor: pointer;
  border-radius: 4px;
  border: 1px solid #d1d9e0;
  transition: all 0.2s ease;
  background: #f4f6f8;

  &:hover {
    border-color: @primary;
    box-shadow: 0 0 0 3px rgba(30, 58, 95, 0.1);
  }
}

.error-message {
  color: #c0392b;
  margin-bottom: 18px;
  text-align: center;
  font-size: 13px;
  line-height: 1.5;
  padding: 12px;
  background: rgba(192, 57, 43, 0.08);
  border-radius: 4px;
  border: 1px solid rgba(192, 57, 43, 0.15);
}

.submit-btn {
  margin-bottom: 0;
  margin-top: 28px;

  :deep(.ant-btn) {
    font-family: 'Outfit', sans-serif;
    font-size: 14px;
    font-weight: 600;
    letter-spacing: 0.5px;
    height: auto;
    padding: 14px 24px;
    border-radius: 4px;
    background: @primary;
    border-color: @primary;
    transition: all 0.25s ease;

    &:hover {
      background: @primary-light;
      border-color: @primary-light;
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(30, 58, 95, 0.25);
    }

    &:active {
      transform: translateY(0);
    }
  }
}

@media (max-width: 768px) {
  .login-header {
    margin-bottom: 24px;
    
    h1 {
      font-size: 28px;
    }
  }

  .login-card {
    :deep(.ant-card-body) {
      padding: 24px;
    }
  }

  :deep(.ant-form-item-label) {
    padding-bottom: 6px;
  }
}
</style>

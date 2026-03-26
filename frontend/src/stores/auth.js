import { defineStore } from 'pinia'
import axios from 'axios'
import config from '../config'

const apiClient = axios.create({
  baseURL: config.API_BASE_URL
})

apiClient.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user')) || null,
    token: localStorage.getItem('token'),
    loading: false,
    error: null
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    isSuperAdmin: (state) => state.user && state.user.role === 'super'
  },
  actions: {
    async login(username, password, captcha, captchaId) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.post('/login', {
          username,
          password,
          captcha,
          captcha_id: captchaId
        })
        this.token = response.data.token
        this.user = response.data.user
        localStorage.setItem('token', response.data.token)
        localStorage.setItem('user', JSON.stringify(response.data.user))
        return response.data
      } catch (error) {
        this.error = error.response?.data?.error || 'Login failed'
        throw error
      } finally {
        this.loading = false
      }
    },
    async getCaptcha() {
      try {
        const response = await apiClient.get('/captcha')
        return response.data
      } catch (error) {
        console.error('Failed to get captcha:', error)
        throw error
      }
    },
    async getUserInfo() {
      try {
        const response = await apiClient.get('/user/info')
        this.user = response.data.user
        localStorage.setItem('user', JSON.stringify(response.data.user))
        return response.data.user
      } catch (error) {
        console.error('Failed to get user info:', error)
        throw error
      }
    },
    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  }
})

export const useProductStore = defineStore('product', {
  state: () => ({
    products: [],
    loading: false,
    error: null
  }),
  actions: {
    async getProducts(page = 1, pageSize = 10) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get('/products', {
          params: {
            page,
            pageSize
          }
        })
        this.products = response.data.products
        return {
          products: response.data.products,
          total: response.data.total
        }
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get products'
        throw error
      } finally {
        this.loading = false
      }
    },
    async getProduct(id) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get(`/products/${id}`)
        return response.data.product
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get product'
        throw error
      } finally {
        this.loading = false
      }
    },
    async createProduct(product) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.post('/products', product)
        this.products.push(response.data.product)
        return response.data.product
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to create product'
        throw error
      } finally {
        this.loading = false
      }
    },
    async updateProduct(id, product) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.put(`/products/${id}`, product)
        const index = this.products.findIndex(p => p.id === id)
        if (index !== -1) {
          this.products[index] = response.data.product
        }
        return response.data.product
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to update product'
        throw error
      } finally {
        this.loading = false
      }
    },
    async deleteProduct(id) {
      this.loading = true
      this.error = null
      try {
        await apiClient.delete(`/products/${id}`)
        this.products = this.products.filter(p => p.id !== id)
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to delete product'
        throw error
      } finally {
        this.loading = false
      }
    },
    async getProductCount() {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get('/products/count')
        return response.data.count
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get product count'
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})

export const useUserStore = defineStore('user', {
  state: () => ({
    users: [],
    loading: false,
    error: null
  }),
  actions: {
    async getUsers() {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get('/users')
        this.users = response.data.users
        return response.data.users
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get users'
        throw error
      } finally {
        this.loading = false
      }
    },
    async getUser(id) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get(`/users/${id}`)
        return response.data.user
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get user'
        throw error
      } finally {
        this.loading = false
      }
    },
    async createUser(user) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.post('/users', user)
        this.users.push(response.data.user)
        return response.data.user
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to create user'
        throw error
      } finally {
        this.loading = false
      }
    },
    async updateUser(id, user) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.put(`/users/${id}`, user)
        const index = this.users.findIndex(u => u.id === id)
        if (index !== -1) {
          this.users[index] = response.data.user
        }
        return response.data.user
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to update user'
        throw error
      } finally {
        this.loading = false
      }
    },
    async deleteUser(id) {
      this.loading = true
      this.error = null
      try {
        await apiClient.delete(`/users/${id}`)
        this.users = this.users.filter(u => u.id !== id)
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to delete user'
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})

export const useFeedbackStore = defineStore('feedback', {
  state: () => ({
    feedbacks: [],
    loading: false,
    error: null
  }),
  actions: {
    async getFeedbacks() {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get('/feedbacks')
        this.feedbacks = response.data.feedbacks
        return response.data.feedbacks
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get feedbacks'
        throw error
      } finally {
        this.loading = false
      }
    },
    async getFeedback(id) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get(`/feedbacks/${id}`)
        return response.data.feedback
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get feedback'
        throw error
      } finally {
        this.loading = false
      }
    },
    async getFeedbackCount() {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get('/feedbacks/count')
        return response.data.count
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get feedback count'
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})

export const useNewsStore = defineStore('news', {
  state: () => ({
    news: [],
    loading: false,
    error: null
  }),
  actions: {
    async getNews(page = 1, pageSize = 10) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get('/news', {
          params: { page, pageSize }
        })
        this.news = response.data.news
        return {
          news: response.data.news,
          total: response.data.total
        }
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get news'
        throw error
      } finally {
        this.loading = false
      }
    },
    async getNewsItem(id) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get(`/news/${id}`)
        return response.data.news
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get news'
        throw error
      } finally {
        this.loading = false
      }
    },
    async createNews(newsData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.post('/news', newsData)
        this.news.push(response.data.news)
        return response.data.news
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to create news'
        throw error
      } finally {
        this.loading = false
      }
    },
    async updateNews(id, newsData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.put(`/news/${id}`, newsData)
        const index = this.news.findIndex(n => n.id === id)
        if (index !== -1) {
          this.news[index] = response.data.news
        }
        return response.data.news
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to update news'
        throw error
      } finally {
        this.loading = false
      }
    },
    async deleteNews(id) {
      this.loading = true
      this.error = null
      try {
        await apiClient.delete(`/news/${id}`)
        this.news = this.news.filter(n => n.id !== id)
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to delete news'
        throw error
      } finally {
        this.loading = false
      }
    },
    async getNewsCount() {
      this.loading = true
      this.error = null
      try {
        const response = await apiClient.get('/news/count')
        return response.data.count
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get news count'
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})

import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Layout from '../components/Layout.vue'
import Dashboard from '../views/Dashboard.vue'
import UserManagement from '../views/UserManagement.vue'
import ProductManagement from '../views/ProductManagement.vue'
import ProductDetail from '../views/ProductDetail.vue'
import FeedbackManagement from '../views/FeedbackManagement.vue'
import NewsManagement from '../views/NewsManagement.vue'
import AdPositionManagement from '../views/AdPositionManagement.vue'
import AdManagement from '../views/AdManagement.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/',
    name: 'Layout',
    component: Layout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: Dashboard
      },
      {
        path: 'users',
        name: 'UserManagement',
        component: UserManagement
      },
      {
        path: 'products',
        name: 'ProductManagement',
        component: ProductManagement
      },
      {
        path: 'products/create',
        name: 'ProductCreate',
        component: ProductDetail,
        meta: { title: '新建产品' }
      },
      {
        path: 'products/:id',
        name: 'ProductEdit',
        component: ProductDetail,
        meta: { title: '编辑产品' }
      },
      {
        path: 'feedbacks',
        name: 'FeedbackManagement',
        component: FeedbackManagement
      },
      {
        path: 'news',
        name: 'NewsManagement',
        component: NewsManagement
      },
      {
        path: 'ad-positions',
        name: 'AdPositionManagement',
        component: AdPositionManagement
      },
      {
        path: 'ads',
        name: 'AdManagement',
        component: AdManagement
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!token) {
      next({ name: 'Login' })
    } else {
      if (to.matched.some(record => record.meta.requiresSuperAdmin)) {
        const user = JSON.parse(localStorage.getItem('user'))
        if (user && user.role === 'super') {
          next()
        } else {
          next({ name: 'Dashboard' })
        }
      } else {
        next()
      }
    }
  } else {
    next()
  }
})

export default router

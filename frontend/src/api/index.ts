import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 120000, // 爬虫执行需要较长超时时间
})

// 响应拦截器
api.interceptors.response.use(
  response => {
    if (response.data.code !== 0) {
      return Promise.reject(new Error(response.data.message))
    }
    return response.data
  },
  error => {
    return Promise.reject(error)
  }
)

export default api
import api from './index'

// 国标政策API
export const nationalApi = {
  list: (params: { page?: number; pageSize?: number; keyword?: string; status?: string; checkStatus?: number }) =>
    api.get('/national', { params }),

  get: (id: number) => api.get(`/national/${id}`),

  create: (data: any) => api.post('/national', data),

  update: (id: number, data: any) => api.put(`/national/${id}`, data),

  delete: (id: number) => api.delete(`/national/${id}`),

  import: (data: any[]) => api.post('/national/import', data),

  check: (id: number, checkStatus: number) => api.put(`/national/${id}/check`, { check_status: checkStatus }),
}

// 行标政策API
export const industryApi = {
  list: (params: { page?: number; pageSize?: number; keyword?: string; status?: string; checkStatus?: number }) =>
    api.get('/industry', { params }),

  get: (id: number) => api.get(`/industry/${id}`),

  create: (data: any) => api.post('/industry', data),

  update: (id: number, data: any) => api.put(`/industry/${id}`, data),

  delete: (id: number) => api.delete(`/industry/${id}`),

  import: (data: any[]) => api.post('/industry/import', data),

  check: (id: number, checkStatus: number) => api.put(`/industry/${id}/check`, { check_status: checkStatus }),
}

// 地标政策API
export const localApi = {
  list: (params: { page?: number; pageSize?: number; keyword?: string; status?: string; checkStatus?: number }) =>
    api.get('/local', { params }),

  get: (id: number) => api.get(`/local/${id}`),

  create: (data: any) => api.post('/local', data),

  update: (id: number, data: any) => api.put(`/local/${id}`, data),

  delete: (id: number) => api.delete(`/local/${id}`),

  import: (data: any[]) => api.post('/local/import', data),

  check: (id: number, checkStatus: number) => api.put(`/local/${id}/check`, { check_status: checkStatus }),
}

// 分类API
export const categoryApi = {
  list: (params: { page?: number; pageSize?: number; keyword?: string }) =>
    api.get('/categories', { params }),

  get: (id: number) => api.get(`/categories/${id}`),

  create: (data: any) => api.post('/categories', data),

  update: (id: number, data: any) => api.put(`/categories/${id}`, data),

  delete: (id: number) => api.delete(`/categories/${id}`),

  count: (id: number) => api.get(`/categories/${id}/count`),

  all: () => api.get('/categories', { params: { page: 1, pageSize: 100 } }),
}

// 爬虫API
export const crawlerApi = {
  execute: (script: string, query: string) => api.post('/crawler/execute', { script, query }),

  logs: () => api.get('/crawler/logs'),

  import: (data: any[], targetType: string) => api.post('/crawler/import', { data, target_type: targetType }),
}

// 推荐政策API
export const recommendApi = {
  list: (params?: { keyword?: string; policy_type?: string; page?: number; pageSize?: number }) =>
    api.get('/recommends', { params }),

  create: (data: any) => api.post('/recommends', data),

  update: (id: number, data: any) => api.put(`/recommends/${id}`, data),

  delete: (id: number) => api.delete(`/recommends/${id}`),
}

// 前端展示API
export const policyApi = {
  listAll: (params: { page?: number; pageSize?: number; keyword?: string }) =>
    api.get('/policies', { params }),

  search: (keyword: string) => api.get('/search', { params: { keyword } }),
}

// H5移动端API
export const mobileApi = {
  // 鎷取首页数据
  getHomeData: () => api.get('/home'),

  // 按类型获取政策列表
  getListByType: (type: string, params: { page?: number; pageSize?: number; categoryId?: number; keyword?: string }) =>
    api.get(`/policies/${type}`, { params }),

  // 按类型获取政策详情
  getDetail: (type: string, id: number) =>
    api.get(`/policies/${type}/${id}`),

  // 获取分类列表
  getCategories: () =>
    api.get('/categories'),

  // 搜索
  search: (keyword: string) =>
    api.get('/search', { params: { keyword } }),
}

export default {
  nationalApi,
  industryApi,
  localApi,
  categoryApi,
  crawlerApi,
  recommendApi,
  policyApi,
  mobileApi,
}
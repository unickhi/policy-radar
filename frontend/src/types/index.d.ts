/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

interface PageData<T = any> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

interface PageResponse<T = any> {
  code: number
  message: string
  data: PageData<T>
}

// 政策基础类型
interface PolicyBase {
  id: number
  standard_no: string
  standard_name: string
  publish_date: string
  implement_date: string
  status: string
  download_url: string
  check_status: number // 0待核验 1已核验 2不通过
  category_id: number
  created_at: string
  updated_at: string
}

// 国标政策
interface NationalStandard extends PolicyBase {
  link1: string
  link2: string
  english_name: string
  nature: string
  category: string
  is_adopted: string
  ccs_code: string
  ics_code: string
  department: string
  technical_dept: string
  publisher: string
  description: string
  standard_type: string
}

// 行标政策
interface IndustryStandard extends PolicyBase {
  detail_link: string
  revision_type: string
  ccs_code: string
  ics_code: string
  technical_owner: string
  approve_dept: string
  industry_class: string
  standard_class: string
  replace_standard: string
  standard_type: string
}

// 地标政策
interface LocalStandard extends PolicyBase {
  detail_link: string
  nature: string
  ccs_code: string
  ics_code: string
  department: string
  publisher: string
  description: string
  standard_type: string
}

// 政策分类
interface PolicyCategory {
  id: number
  name: string
  code: string
  description: string
  created_at: string
  updated_at: string
}

// 推荐政策
interface PolicyRecommend {
  id: number
  policy_id: number
  policy_type: string
  title: string
  content: string
  sort: number
  created_at: string
  updated_at: string
}
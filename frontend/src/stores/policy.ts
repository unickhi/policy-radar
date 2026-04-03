import { defineStore } from 'pinia'
import { nationalApi, industryApi, localApi, categoryApi } from '@/api/policy'

export const usePolicyStore = defineStore('policy', {
  state: () => ({
    nationalList: [] as any[],
    nationalTotal: 0,
    industryList: [] as any[],
    industryTotal: 0,
    localList: [] as any[],
    localTotal: 0,
    categoryList: [] as any[],
    loading: false,
  }),

  actions: {
    async fetchNational(params: any) {
      this.loading = true
      try {
        const res = await nationalApi.list(params) as any
        this.nationalList = res.data.list
        this.nationalTotal = res.data.total
      } finally {
        this.loading = false
      }
    },

    async fetchIndustry(params: any) {
      this.loading = true
      try {
        const res = await industryApi.list(params) as any
        this.industryList = res.data.list
        this.industryTotal = res.data.total
      } finally {
        this.loading = false
      }
    },

    async fetchLocal(params: any) {
      this.loading = true
      try {
        const res = await localApi.list(params) as any
        this.localList = res.data.list
        this.localTotal = res.data.total
      } finally {
        this.loading = false
      }
    },

    async fetchCategories(params: any) {
      this.loading = true
      try {
        const res = await categoryApi.list(params) as any
        this.categoryList = res.data.list
      } finally {
        this.loading = false
      }
    },
  }
})
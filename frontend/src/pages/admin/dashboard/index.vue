<template>
  <div class="flex flex-col gap-4">
    <!-- 统计卡片 -->
    <div class="grid grid-cols-5 gap-4">
      <div class="bg-white p-5 shadow-sm flex items-center gap-4">
        <img src="/standard.svg" class="w-12 h-12" alt="国标" />
        <div>
          <div class="text-gray-500 text-sm mb-1">国标政策</div>
          <div class="text-2xl font-bold text-gray-800">{{ stats.national }}</div>
          <div class="text-xs text-gray-400">国家标准数据</div>
        </div>
      </div>
      <div class="bg-white p-5 shadow-sm flex items-center gap-4">
        <img src="/dianli.svg" class="w-12 h-12" alt="行标" />
        <div>
          <div class="text-gray-500 text-sm mb-1">行标政策</div>
          <div class="text-2xl font-bold text-gray-800">{{ stats.industry }}</div>
          <div class="text-xs text-gray-400">行业标准数据</div>
        </div>
      </div>
      <div class="bg-white p-5 shadow-sm flex items-center gap-4">
        <img src="/area.svg" class="w-12 h-12" alt="地标" />
        <div>
          <div class="text-gray-500 text-sm mb-1">地标政策</div>
          <div class="text-2xl font-bold text-gray-800">{{ stats.local }}</div>
          <div class="text-xs text-gray-400">地方标准数据</div>
        </div>
      </div>
      <div class="bg-white p-5 shadow-sm flex items-center gap-4">
        <img src="/fenlei.svg" class="w-12 h-12" alt="分类" />
        <div>
          <div class="text-gray-500 text-sm mb-1">政策分类</div>
          <div class="text-2xl font-bold text-gray-800">{{ stats.category }}</div>
          <div class="text-xs text-gray-400">分类总数</div>
        </div>
      </div>
      <div class="bg-white p-5 shadow-sm flex items-center gap-4">
        <img src="/recomend.svg" class="w-12 h-12" alt="推荐" />
        <div>
          <div class="text-gray-500 text-sm mb-1">推荐政策</div>
          <div class="text-2xl font-bold text-gray-800">{{ stats.recommend }}</div>
          <div class="text-xs text-gray-400">前台推荐数</div>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="grid grid-cols-3 gap-4">
      <!-- 政策类型占比 -->
      <div class="bg-white p-4 shadow-sm">
        <div class="text-gray-700 font-medium mb-3">政策类型分布</div>
        <div ref="pieChartRef" class="h-64"></div>
      </div>

      <!-- 核验状态 -->
      <div class="bg-white p-4 shadow-sm">
        <div class="text-gray-700 font-medium mb-3">核验状态分布</div>
        <div ref="checkChartRef" class="h-64"></div>
      </div>

      <!-- 标准状态 -->
      <div class="bg-white p-4 shadow-sm">
        <div class="text-gray-700 font-medium mb-3">标准状态分布</div>
        <div ref="statusChartRef" class="h-64"></div>
      </div>
    </div>

    <!-- 分类统计 -->
    <div class="bg-white p-4 shadow-sm">
      <div class="text-gray-700 font-medium mb-3">各分类政策数量统计</div>
      <div ref="categoryChartRef" class="h-64"></div>
    </div>

    <!-- 最近政策 -->
    <div class="bg-white p-4 shadow-sm">
      <div class="text-gray-700 font-medium mb-3">最近添加的政策</div>
      <el-table :data="recentPolicies" stripe border size="small">
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag size="small" :type="row.type === '国标' ? 'danger' : row.type === '行标' ? 'warning' : 'success'">{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="standard_no" label="标准号" width="180" />
        <el-table-column prop="standard_name" label="标准名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="publish_date" label="发布日期" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="row.status === '现行' ? 'success' : 'warning'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="source" label="来源" width="80">
          <template #default="{ row }">
            <el-tag size="small" :type="row.source === 'crawl' ? 'info' : ''">{{ row.source === 'crawl' ? '爬取' : '人工' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { nationalApi, industryApi, localApi, categoryApi, recommendApi } from '@/api/policy'
import * as echarts from 'echarts'

const stats = ref({ national: 0, industry: 0, local: 0, category: 0, recommend: 0 })
const checkStats = ref({ pending: 0, checked: 0, rejected: 0 })
const statusStats = ref({ active: 0, upcoming: 0, obsolete: 0 })
const categoryStats = ref<{ name: string; count: number }[]>([])
const recentPolicies = ref<any[]>([])

const pieChartRef = ref<HTMLElement | null>(null)
const checkChartRef = ref<HTMLElement | null>(null)
const statusChartRef = ref<HTMLElement | null>(null)
const categoryChartRef = ref<HTMLElement | null>(null)

let pieChart: echarts.ECharts | null = null
let checkChart: echarts.ECharts | null = null
let statusChart: echarts.ECharts | null = null
let categoryChart: echarts.ECharts | null = null

const fetchStats = async () => {
  try {
    // 获取真实数据
    const [nationalRes, industryRes, localRes, categoryRes, recommendRes] = await Promise.all([
      nationalApi.list({ page: 1, pageSize: 1 }) as any,
      industryApi.list({ page: 1, pageSize: 1 }) as any,
      localApi.list({ page: 1, pageSize: 1 }) as any,
      categoryApi.list({ page: 1, pageSize: 100 }) as any,
      recommendApi.list() as any,
    ])

    stats.value = {
      national: nationalRes.data?.total || 0,
      industry: industryRes.data?.total || 0,
      local: localRes.data?.total || 0,
      category: categoryRes.data?.total || 0,
      recommend: recommendRes.data?.length || 0,
    }

    // 获取各分类的政策数量
    const categories = categoryRes.data?.list || []
    const catStats: { name: string; count: number }[] = []
    for (const cat of categories) {
      try {
        const countRes = await categoryApi.count(cat.id) as any
        const total = (countRes.data?.national || 0) + (countRes.data?.industry || 0) + (countRes.data?.local || 0)
        if (total > 0) {
          catStats.push({ name: cat.name, count: total })
        }
      } catch (e) {}
    }
    categoryStats.value = catStats

    // 获取最近政策（包含来源）
    const recentRes = await nationalApi.list({ page: 1, pageSize: 5 }) as any
    recentPolicies.value = (recentRes.data?.list || []).map((item: any) => ({
      type: '国标',
      standard_no: item.standard_no,
      standard_name: item.standard_name,
      publish_date: item.publish_date,
      status: item.status,
      source: item.source || 'manual',
    }))

    // 计算核验状态和标准状态
    const allNational = await nationalApi.list({ page: 1, pageSize: 200 }) as any
    const nationalList = allNational.data?.list || []

    let pending = 0, checked = 0, rejected = 0
    let active = 0, upcoming = 0, obsolete = 0

    nationalList.forEach((item: any) => {
      if (item.check_status === 0) pending++
      else if (item.check_status === 1) checked++
      else if (item.check_status === 2) rejected++

      if (item.status === '现行') active++
      else if (item.status === '即将实施') upcoming++
      else obsolete++
    })

    checkStats.value = { pending, checked, rejected }
    statusStats.value = { active, upcoming, obsolete }

    initCharts()
  } catch (e) {
    console.error(e)
    initCharts()
  }
}

const initCharts = () => {
  // 政策类型饼图
  if (pieChartRef.value) {
    pieChart = echarts.init(pieChartRef.value)
    pieChart.setOption({
      tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
      legend: { bottom: 0 },
      color: ['#ef4444', '#f97316', '#22c55e'],
      series: [{
        type: 'pie',
        radius: ['40%', '65%'],
        center: ['50%', '45%'],
        itemStyle: { borderColor: '#fff', borderWidth: 2 },
        data: [
          { value: stats.value.national, name: '国标' },
          { value: stats.value.industry, name: '行标' },
          { value: stats.value.local, name: '地标' },
        ]
      }]
    })
  }

  // 核验状态饼图
  if (checkChartRef.value) {
    checkChart = echarts.init(checkChartRef.value)
    checkChart.setOption({
      tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
      legend: { bottom: 0 },
      color: ['#9ca3af', '#22c55e', '#ef4444'],
      series: [{
        type: 'pie',
        radius: ['40%', '65%'],
        center: ['50%', '45%'],
        itemStyle: { borderColor: '#fff', borderWidth: 2 },
        data: [
          { value: checkStats.value.pending, name: '待核验' },
          { value: checkStats.value.checked, name: '已核验' },
          { value: checkStats.value.rejected, name: '不通过' },
        ]
      }]
    })
  }

  // 标准状态饼图
  if (statusChartRef.value) {
    statusChart = echarts.init(statusChartRef.value)
    statusChart.setOption({
      tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
      legend: { bottom: 0 },
      color: ['#22c55e', '#f97316', '#9ca3af'],
      series: [{
        type: 'pie',
        radius: ['40%', '65%'],
        center: ['50%', '45%'],
        itemStyle: { borderColor: '#fff', borderWidth: 2 },
        data: [
          { value: statusStats.value.active, name: '现行' },
          { value: statusStats.value.upcoming, name: '即将实施' },
          { value: statusStats.value.obsolete, name: '废止' },
        ]
      }]
    })
  }

  // 分类统计柱状图
  if (categoryChartRef.value) {
    categoryChart = echarts.init(categoryChartRef.value)
    categoryChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: categoryStats.value.map(c => c.name),
        axisLabel: { interval: 0, rotate: 30 }
      },
      yAxis: { type: 'value' },
      series: [{
        type: 'bar',
        data: categoryStats.value.map(c => c.count),
        itemStyle: { color: '#3b82f6' },
        barWidth: '40%'
      }]
    })
  }
}

onMounted(() => {
  fetchStats()
  window.addEventListener('resize', () => {
    pieChart?.resize()
    checkChart?.resize()
    statusChart?.resize()
    categoryChart?.resize()
  })
})

onUnmounted(() => {
  pieChart?.dispose()
  checkChart?.dispose()
  statusChart?.dispose()
  categoryChart?.dispose()
})
</script>
<template>
  <div class="flex flex-col h-full">
    <!-- 头部 -->
    <div class="bg-white p-4 mb-4 border-b">
      <div class="flex justify-between items-center">
        <div class="flex gap-3 items-center">
          <el-input v-model="searchForm.keyword" placeholder="解析标题" clearable class="w-48" @keyup.enter="handleSearch" />
          <el-select v-model="searchForm.policy_type" placeholder="政策类型" clearable class="w-28">
            <el-option label="国标" value="国标" />
            <el-option label="行标" value="行标" />
            <el-option label="地标" value="地标" />
          </el-select>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        <el-button type="primary" @click="handleAdd">新增推荐</el-button>
      </div>
    </div>

    <!-- 表格 -->
    <div class="flex-1 bg-white overflow-hidden flex flex-col">
      <el-table :data="tableData" stripe v-loading="loading" border class="flex-1">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="policy_type" label="政策类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.policy_type === '国标' ? 'danger' : row.policy_type === '行标' ? 'warning' : 'success'" size="small">{{ row.policy_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="policy_name" label="关联政策" min-width="200" show-overflow-tooltip />
        <el-table-column prop="title" label="解析标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="p-4 flex justify-end border-t">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize" :total="total" :page-sizes="[10, 20, 50]" layout="total, sizes, prev, pager, next" @size-change="fetchData" @current-change="fetchData" />
      </div>
    </div>

    <!-- 编辑抽屉 -->
    <el-drawer v-model="drawerVisible" :title="drawerTitle" size="900px">
      <el-form :model="form" label-width="120px">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="关联政策类型" required>
              <el-select v-model="form.policy_type" placeholder="请选择政策类型" class="w-full" @change="handlePolicyTypeChange">
                <el-option label="国标" value="国标" />
                <el-option label="行标" value="行标" />
                <el-option label="地标" value="地标" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序权重">
              <el-input-number v-model="form.sort" :min="0" :max="999" placeholder="数字越大越靠前" class="w-full" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="关联政策" required>
          <el-select v-model="form.policy_id" filterable placeholder="请选择关联政策" class="w-full">
            <el-option v-for="policy in policyList" :key="policy.id" :label="`${policy.standard_no} - ${policy.standard_name}`" :value="policy.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="解析标题" required>
          <el-input v-model="form.title" placeholder="请输入解析标题" />
        </el-form-item>
        <el-form-item label="解析内容" required>
          <MdEditor v-model="form.content" language="zh-CN" :previewOnly="false" style="height: 400px" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { recommendApi, nationalApi, industryApi, localApi } from '@/api/policy'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const policyList = ref<any[]>([])

const searchForm = reactive({ keyword: '', policy_type: '' })
const pagination = reactive({ page: 1, pageSize: 10 })

const drawerVisible = ref(false)
const drawerTitle = ref('新增推荐')
const form = reactive({ id: 0, policy_id: null as number | null, policy_type: '', policy_name: '', title: '', content: '', sort: 0 })

watch(() => form.policy_type, async (val) => {
  if (val) { form.policy_id = null; await fetchPolicyList(val) }
})

const fetchPolicyList = async (type: string) => {
  try {
    let res: any
    if (type === '国标') res = await nationalApi.list({ page: 1, pageSize: 100 })
    else if (type === '行标') res = await industryApi.list({ page: 1, pageSize: 100 })
    else res = await localApi.list({ page: 1, pageSize: 100 })
    policyList.value = res.data.list || []
  } catch (e) { policyList.value = [] }
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await recommendApi.list({
      keyword: searchForm.keyword,
      policy_type: searchForm.policy_type,
      page: pagination.page,
      pageSize: pagination.pageSize,
    }) as any
    tableData.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e: any) { ElMessage.error(e.message) } finally { loading.value = false }
}

onMounted(() => fetchData())

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.policy_type = ''; pagination.page = 1; fetchData() }
const handlePolicyTypeChange = () => { form.policy_id = null }

const handleAdd = () => {
  drawerTitle.value = '新增推荐'
  Object.assign(form, { id: 0, policy_id: null, policy_type: '', policy_name: '', title: '', content: '', sort: 0 })
  policyList.value = []
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  drawerTitle.value = '编辑推荐'
  Object.assign(form, row)
  if (row.policy_type) fetchPolicyList(row.policy_type)
  drawerVisible.value = true
}

const handleSave = async () => {
  if (!form.policy_id || !form.policy_type || !form.title) { ElMessage.warning('请填写完整信息'); return }
  const policy = policyList.value.find(p => p.id === form.policy_id)
  form.policy_name = policy ? `${policy.standard_no} - ${policy.standard_name}` : ''
  try {
    if (form.id) {
      await recommendApi.update(form.id, form)
    } else {
      await recommendApi.create(form)
    }
    ElMessage.success('保存成功')
    drawerVisible.value = false
    fetchData()
  } catch (e: any) { ElMessage.error(e.message) }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定删除该推荐?', '提示', { type: 'warning' })
    await recommendApi.delete(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e.message)
  }
}
</script>
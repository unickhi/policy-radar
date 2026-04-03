<template>
  <div class="flex flex-col h-full">
    <div class="bg-white p-4 mb-4 border-b">
      <div class="flex justify-between items-center">
        <div class="flex gap-3 items-center">
          <el-input v-model="searchForm.keyword" placeholder="分类名称/编码" clearable class="w-48" @keyup.enter="handleSearch" />
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        <el-button type="primary" @click="handleAdd">新增分类</el-button>
      </div>
    </div>

    <div class="flex-1 bg-white overflow-hidden flex flex-col">
      <el-table :data="tableData" stripe v-loading="loading" border class="flex-1">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="分类名称" width="200" />
        <el-table-column prop="code" label="分类编码" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="关联政策数量" width="300">
          <template #default="{ row }">
            <div class="flex gap-2">
              <el-tag size="small" type="danger">国标: {{ row.count?.national || 0 }}</el-tag>
              <el-tag size="small" type="warning">行标: {{ row.count?.industry || 0 }}</el-tag>
              <el-tag size="small" type="success">地标: {{ row.count?.local || 0 }}</el-tag>
            </div>
          </template>
        </el-table-column>
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

    <el-drawer v-model="drawerVisible" :title="drawerTitle" size="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="分类名称" required><el-input v-model="form.name" placeholder="请输入分类名称" /></el-form-item>
        <el-form-item label="分类编码" required><el-input v-model="form.code" placeholder="请输入分类编码" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入分类描述" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { categoryApi } from '@/api/policy'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const searchForm = reactive({ keyword: '' })
const pagination = reactive({ page: 1, pageSize: 10 })

const drawerVisible = ref(false)
const drawerTitle = ref('新增分类')
const form = reactive({ id: 0, name: '', code: '', description: '' })

const fetchData = async () => {
  loading.value = true
  try {
    const res = await categoryApi.list({ ...pagination, ...searchForm }) as any
    tableData.value = res.data.list
    total.value = res.data.total
    for (const item of tableData.value) {
      try { const countRes = await categoryApi.count(item.id) as any; item.count = countRes.data } catch { item.count = { national: 0, industry: 0, local: 0 } }
    }
  } catch (e: any) { ElMessage.error(e.message) } finally { loading.value = false }
}

onMounted(() => fetchData())

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; pagination.page = 1; fetchData() }

const handleAdd = () => { drawerTitle.value = '新增分类'; Object.assign(form, { id: 0, name: '', code: '', description: '' }); drawerVisible.value = true }
const handleEdit = (row: any) => { drawerTitle.value = '编辑分类'; Object.assign(form, row); drawerVisible.value = true }

const handleSave = async () => {
  try { form.id ? await categoryApi.update(form.id, form) : await categoryApi.create(form); ElMessage.success('保存成功'); drawerVisible.value = false; fetchData() }
  catch (e: any) { ElMessage.error(e.message) }
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm('确定删除该分类?', '提示', { type: 'warning' }).then(async () => { try { await categoryApi.delete(row.id); ElMessage.success('删除成功'); fetchData() } catch (e: any) { ElMessage.error(e.message) } })
}
</script>
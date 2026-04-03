<template>
  <div class="flex flex-col h-full">
    <!-- 头部 -->
    <div class="bg-white p-4 mb-4 border-b">
      <div class="flex justify-between items-center">
        <div class="flex gap-3 items-center">
          <el-input v-model="searchForm.keyword" placeholder="标准号/名称" clearable class="w-48" @keyup.enter="handleSearch" />
          <el-select v-model="searchForm.categoryId" placeholder="政策分类" clearable class="w-28">
            <el-option v-for="cat in categoryList" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
          <el-select v-model="searchForm.checkStatus" placeholder="核验状态" clearable class="w-28">
            <el-option label="待核验" :value="0" />
            <el-option label="已核验" :value="1" />
            <el-option label="不通过" :value="2" />
          </el-select>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        <div class="flex gap-2">
          <el-button type="primary" @click="handleAdd">新增</el-button>
          <el-button type="success" @click="handleImport">导入</el-button>
          <el-button @click="handleExport">导出</el-button>
        </div>
      </div>
    </div>

    <!-- 表格 -->
    <div class="flex-1 bg-white overflow-hidden flex flex-col">
      <el-table :data="tableData" stripe v-loading="loading" border class="flex-1">
        <el-table-column prop="standard_no" label="标准号" width="180" />
        <el-table-column prop="standard_name" label="标准名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="category_name" label="政策分类" width="120">
          <template #default="{ row }">
            <span v-if="row.category_id">{{ getCategoryName(row.category_id) }}</span>
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="publish_date" label="发布日期" width="120" />
        <el-table-column prop="implement_date" label="实施日期" width="120" />
        <el-table-column prop="status" label="标准状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === '现行'" type="success" size="small">{{ row.status }}</el-tag>
            <el-tag v-else type="warning" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="check_status" label="核验状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.check_status === 0" type="info" size="small">待核验</el-tag>
            <el-tag v-else-if="row.check_status === 1" type="success" size="small">已核验</el-tag>
            <el-tag v-else type="danger" size="small">不通过</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" link type="warning" @click="handleCheck(row)">核验</el-button>
            <el-button size="small" link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="p-4 flex justify-end border-t">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize" :total="total" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" @size-change="fetchData" @current-change="fetchData" />
      </div>
    </div>

    <!-- 编辑抽屉 -->
    <el-drawer v-model="drawerVisible" :title="drawerTitle" size="700px">
      <el-form :model="form" label-width="120px">
        <el-row :gutter="24">
          <el-col :span="12"><el-form-item label="标准号" required><el-input v-model="form.standard_no" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="标准名称" required><el-input v-model="form.standard_name" /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="24">
          <el-col :span="12"><el-form-item label="发布日期"><el-date-picker v-model="form.publish_date" type="date" value-format="YYYY-MM-DD" class="w-full" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="实施日期"><el-date-picker v-model="form.implement_date" type="date" value-format="YYYY-MM-DD" class="w-full" /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="24">
          <el-col :span="12"><el-form-item label="标准状态"><el-select v-model="form.status" class="w-full"><el-option label="现行" value="现行" /><el-option label="即将实施" value="即将实施" /></el-select></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="政策分类">
            <el-select v-model="form.category_id" clearable placeholder="请选择分类" class="w-full">
              <el-option v-for="cat in categoryList" :key="cat.id" :label="cat.name" :value="cat.id" />
            </el-select>
          </el-form-item></el-col>
        </el-row>
        <el-row :gutter="24">
          <el-col :span="12"><el-form-item label="标准性质"><el-select v-model="form.nature" class="w-full"><el-option label="推荐性" value="推荐性" /><el-option label="强制性" value="强制性" /></el-select></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="主管部门"><el-input v-model="form.department" /></el-form-item></el-col>
        </el-row>
        <el-form-item label="发布单位"><el-input v-model="form.publisher" /></el-form-item>
        <el-form-item label="下载链接"><el-input v-model="form.download_url" /></el-form-item>
        <el-form-item label="详情简介">
          <MdEditor v-model="form.description" language="zh-CN" :previewOnly="false" style="height: 300px" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-drawer>

    <!-- 核验抽屉 -->
    <el-drawer v-model="checkDrawerVisible" title="政策核验" size="400px">
      <el-form :model="checkForm" label-width="100px">
        <el-form-item label="核验状态">
          <el-radio-group v-model="checkForm.check_status"><el-radio :value="1">已核验</el-radio><el-radio :value="2">不通过</el-radio></el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="checkDrawerVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveCheck">确认</el-button>
      </template>
    </el-drawer>

    <!-- 导入抽屉 -->
    <el-drawer v-model="importDrawerVisible" title="批量导入" size="600px">
      <el-upload drag accept=".xlsx,.xls" :auto-upload="false" :on-change="handleFileChange" :limit="1">
        <el-icon class="el-icon--upload text-4xl text-gray-300"><UploadFilled /></el-icon>
        <div class="el-upload__text">拖拽文件或 <em>点击上传</em></div>
      </el-upload>
      <div v-if="importPreview.length > 0" class="mt-4">
        <div class="text-gray-600 mb-2">数据预览 (共 {{ importPreview.length }} 条)</div>
        <el-table :data="importPreview.slice(0, 5)" stripe border max-height="300">
          <el-table-column prop="standard_no" label="标准号" width="150" />
          <el-table-column prop="standard_name" label="标准名称" min-width="200" />
        </el-table>
      </div>
      <template #footer>
        <el-button @click="importDrawerVisible = false">取消</el-button>
        <el-button type="primary" @click="handleDoImport" :loading="importing" :disabled="importPreview.length === 0">确认导入 {{ importPreview.length }} 条</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { localApi, categoryApi } from '@/api/policy'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import * as XLSX from 'xlsx'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const categoryList = ref<any[]>([])
const searchForm = reactive({ keyword: '', checkStatus: undefined as number | undefined, categoryId: undefined as number | undefined })
const pagination = reactive({ page: 1, pageSize: 10 })

const drawerVisible = ref(false)
const drawerTitle = ref('新增政策')
const form = reactive<any>({ id: 0, standard_no: '', standard_name: '', publish_date: '', implement_date: '', status: '现行', nature: '推荐性', department: '', publisher: '', description: '', download_url: '', category_id: undefined as number | undefined })

const checkDrawerVisible = ref(false)
const checkForm = reactive({ id: 0, check_status: 1 })

const importDrawerVisible = ref(false)
const importPreview = ref<any[]>([])
const importing = ref(false)

const fetchCategoryList = async () => {
  try {
    const res = await categoryApi.all() as any
    categoryList.value = res.data.list || []
  } catch (e) { categoryList.value = [] }
}

const getCategoryName = (categoryId: number) => {
  const cat = categoryList.value.find(c => c.id === categoryId)
  return cat ? cat.name : '-'
}

const fetchData = async () => {
  loading.value = true
  try { const res = await localApi.list({ ...pagination, ...searchForm }) as any; tableData.value = res.data.list; total.value = res.data.total }
  catch (e: any) { ElMessage.error(e.message) } finally { loading.value = false }
}

onMounted(() => {
  fetchCategoryList()
  fetchData()
})

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.checkStatus = undefined; searchForm.categoryId = undefined; pagination.page = 1; fetchData() }

const handleAdd = () => { drawerTitle.value = '新增政策'; Object.keys(form).forEach(key => form[key] = key === 'id' ? 0 : key === 'status' ? '现行' : key === 'nature' ? '推荐性' : ''); drawerVisible.value = true }
const handleEdit = (row: any) => { drawerTitle.value = '编辑政策'; Object.assign(form, row); drawerVisible.value = true }

const handleSave = async () => {
  try { form.id ? await localApi.update(form.id, form) : await localApi.create(form); ElMessage.success('保存成功'); drawerVisible.value = false; fetchData() }
  catch (e: any) { ElMessage.error(e.message) }
}

const handleCheck = (row: any) => { checkForm.id = row.id; checkForm.check_status = 1; checkDrawerVisible.value = true }
const handleSaveCheck = async () => { try { await localApi.check(checkForm.id, checkForm.check_status); ElMessage.success('核验完成'); checkDrawerVisible.value = false; fetchData() } catch (e: any) { ElMessage.error(e.message) } }

const handleDelete = (row: any) => {
  ElMessageBox.confirm('确定删除?', '提示', { type: 'warning' }).then(async () => { try { await localApi.delete(row.id); ElMessage.success('删除成功'); fetchData() } catch (e: any) { ElMessage.error(e.message) } })
}

const handleImport = () => { importPreview.value = []; importDrawerVisible.value = true }

const handleFileChange = async (file: any) => {
  if (!file.raw) return
  try {
    const data = await new Promise<any[]>((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = (e) => { try { const wb = XLSX.read(e.target?.result, { type: 'binary' }); resolve(XLSX.utils.sheet_to_json(wb.Sheets[wb.SheetNames[0]])) } catch (err) { reject(err) } }
      reader.onerror = reject; reader.readAsBinaryString(file.raw)
    })
    importPreview.value = data.map((row: any) => ({ standard_no: row['标准号'] || '', standard_name: row['标准中文名称'] || '', publish_date: row['发布日期'] || '', implement_date: row['实施日期'] || '', status: row['标准状态'] || '现行' }))
    ElMessage.success(`解析成功，共 ${importPreview.value.length} 条`)
  } catch (e: any) { ElMessage.error('解析失败: ' + e.message) }
}

const handleDoImport = async () => {
  importing.value = true
  try { await localApi.import(importPreview.value); ElMessage.success('导入成功'); importDrawerVisible.value = false; fetchData() }
  catch (e: any) { ElMessage.error('导入失败') } finally { importing.value = false }
}

const handleExport = () => { ElMessage.info('导出功能开发中') }
</script>
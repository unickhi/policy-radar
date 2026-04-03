<template>
  <div class="flex flex-col h-full">
    <!-- 头部搜索+工具栏 -->
    <div class="bg-white p-4 mb-4 border-b">
      <div class="flex justify-between items-center">
        <div class="flex gap-3 items-center">
          <el-input v-model="searchForm.keyword" placeholder="标准号/名称" clearable class="w-48" @keyup.enter="handleSearch" />
          <el-select v-model="searchForm.status" placeholder="标准状态" clearable class="w-28">
            <el-option label="现行" value="现行" />
            <el-option label="即将实施" value="即将实施" />
            <el-option label="废止" value="废止" />
          </el-select>
          <el-select v-model="searchForm.categoryId" placeholder="政策分类" clearable class="w-28">
            <el-option v-for="cat in categoryList" :key="cat.id" :label="cat.name" :value="cat.id" />
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
            <el-tag v-else-if="row.status === '即将实施'" type="warning" size="small">{{ row.status }}</el-tag>
            <el-tag v-else type="info" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="nature" label="性质" width="80" />
        <el-table-column prop="check_status" label="核验状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.check_status === 0" type="info" size="small">待核验</el-tag>
            <el-tag v-else-if="row.check_status === 1" type="success" size="small">已核验</el-tag>
            <el-tag v-else type="danger" size="small">不通过</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="download_url" label="下载链接" width="100">
          <template #default="{ row }">
            <el-link v-if="row.download_url" :href="row.download_url" target="_blank" type="primary">下载</el-link>
            <span v-else class="text-gray-300">-</span>
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
    <el-drawer v-model="drawerVisible" :title="drawerTitle" size="800px">
      <el-form :model="form" label-width="120px">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="标准号" required><el-input v-model="form.standard_no" placeholder="如：GB/T 1094.7-2024" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="标准名称" required><el-input v-model="form.standard_name" placeholder="请输入标准中文名称" /></el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="英文名称"><el-input v-model="form.english_name" placeholder="请输入英文标准名称" /></el-form-item>
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="发布日期"><el-date-picker v-model="form.publish_date" type="date" value-format="YYYY-MM-DD" class="w-full" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="实施日期"><el-date-picker v-model="form.implement_date" type="date" value-format="YYYY-MM-DD" class="w-full" /></el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="标准状态">
              <el-select v-model="form.status" class="w-full">
                <el-option label="现行" value="现行" />
                <el-option label="即将实施" value="即将实施" />
                <el-option label="废止" value="废止" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="政策分类">
              <el-select v-model="form.category_id" clearable placeholder="请选择分类" class="w-full">
                <el-option v-for="cat in categoryList" :key="cat.id" :label="cat.name" :value="cat.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="标准性质">
              <el-select v-model="form.nature" class="w-full">
                <el-option label="推荐性" value="推荐性" />
                <el-option label="强制性" value="强制性" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="类别">
              <el-select v-model="form.category" class="w-full">
                <el-option label="推标" value="推标" />
                <el-option label="强标" value="强标" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="是否采标">
              <el-select v-model="form.is_adopted" class="w-full">
                <el-option label="采" value="采" />
                <el-option label="非采" value="非采" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="中国分类号"><el-input v-model="form.ccs_code" placeholder="如：K41" /></el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="国际分类号"><el-input v-model="form.ics_code" placeholder="如：29.180" /></el-form-item>
        <el-form-item label="主管部门"><el-input v-model="form.department" placeholder="如：中国电器工业协会、中国电力企业联合会" /></el-form-item>
        <el-form-item label="归口部门"><el-input v-model="form.technical_dept" placeholder="如：全国变压器标准化技术委员会、中国电器工业协会" /></el-form-item>
        <el-form-item label="发布单位"><el-input v-model="form.publisher" placeholder="如：国家市场监督管理总局、国家标准化管理委员会" /></el-form-item>
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="链接1"><el-input v-model="form.link1" placeholder="std.samr.gov.cn详情链接" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="链接2"><el-input v-model="form.link2" placeholder="openstd.samr.gov.cn详情链接" /></el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="下载链接"><el-input v-model="form.download_url" placeholder="政策文件下载地址" /></el-form-item>
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
          <el-radio-group v-model="checkForm.check_status">
            <el-radio :value="1">已核验</el-radio>
            <el-radio :value="2">不通过</el-radio>
          </el-radio-group>
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
        <template #tip><div class="el-upload__tip">支持 Excel 文件</div></template>
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
import { nationalApi, categoryApi } from '@/api/policy'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import * as XLSX from 'xlsx'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const categoryList = ref<any[]>([])

const searchForm = reactive({ keyword: '', status: '', categoryId: undefined as number | undefined })
const pagination = reactive({ page: 1, pageSize: 10 })

const drawerVisible = ref(false)
const drawerTitle = ref('新增政策')
const form = reactive<any>({ id: 0, standard_no: '', standard_name: '', english_name: '', link1: '', link2: '', publish_date: '', implement_date: '', status: '现行', nature: '推荐性', category: '推标', is_adopted: '', ccs_code: '', ics_code: '', department: '', technical_dept: '', publisher: '', description: '', download_url: '', category_id: undefined as number | undefined })

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
  try {
    const res = await nationalApi.list({ ...pagination, ...searchForm }) as any
    tableData.value = res.data.list
    total.value = res.data.total
  } catch (e: any) { ElMessage.error(e.message) } finally { loading.value = false }
}

onMounted(() => {
  fetchCategoryList()
  fetchData()
})

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.status = ''; searchForm.categoryId = undefined; pagination.page = 1; fetchData() }

const handleAdd = () => {
  drawerTitle.value = '新增政策'
  Object.keys(form).forEach(key => form[key] = key === 'id' ? 0 : key === 'status' ? '现行' : key === 'nature' ? '推荐性' : key === 'category' ? '推标' : '')
  drawerVisible.value = true
}

const handleEdit = (row: any) => { drawerTitle.value = '编辑政策'; Object.assign(form, row); drawerVisible.value = true }

const handleSave = async () => {
  try { form.id ? await nationalApi.update(form.id, form) : await nationalApi.create(form); ElMessage.success('保存成功'); drawerVisible.value = false; fetchData() }
  catch (e: any) { ElMessage.error(e.message) }
}

const handleCheck = (row: any) => { checkForm.id = row.id; checkForm.check_status = 1; checkDrawerVisible.value = true }
const handleSaveCheck = async () => { try { await nationalApi.check(checkForm.id, checkForm.check_status); ElMessage.success('核验完成'); checkDrawerVisible.value = false; fetchData() } catch (e: any) { ElMessage.error(e.message) } }

const handleDelete = (row: any) => {
  ElMessageBox.confirm('确定删除?', '提示', { type: 'warning' }).then(async () => { try { await nationalApi.delete(row.id); ElMessage.success('删除成功'); fetchData() } catch (e: any) { ElMessage.error(e.message) } })
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
    importPreview.value = data.map((row: any) => ({
      standard_no: row['标准号'] || '', standard_name: row['标准中文名称'] || '', english_name: row['英文标准名称'] || '',
      link1: row['链接1'] || '', link2: row['链接2'] || '', publish_date: row['发布日期'] || '', implement_date: row['实施日期'] || '',
      status: row['标准状态'] || '现行', nature: row['标准性质'] || '推荐性', category: row['类别'] || '推标',
      is_adopted: row['是否采标'] || '', ccs_code: row['中国标准分类号'] || '', ics_code: row['国际标准分类号'] || '',
      department: row['主管部门'] || '', technical_dept: row['归口部门'] || '', publisher: row['发布单位'] || '', description: row['详情简介'] || '',
    }))
    ElMessage.success(`解析成功，共 ${importPreview.value.length} 条`)
  } catch (e: any) { ElMessage.error('解析失败: ' + e.message) }
}

const handleDoImport = async () => {
  importing.value = true
  try { await nationalApi.import(importPreview.value); ElMessage.success('导入成功'); importDrawerVisible.value = false; fetchData() }
  catch (e: any) { ElMessage.error('导入失败') } finally { importing.value = false }
}

const handleExport = () => { ElMessage.info('导出功能开发中') }
</script>
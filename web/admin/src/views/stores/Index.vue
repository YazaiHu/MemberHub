<template>
  <div class="stores-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>门店列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增门店
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="门店名称">
          <el-input v-model="searchForm.name" placeholder="请输入门店名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="全部" value="" />
            <el-option label="营业中" :value="1" />
            <el-option label="已关闭" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table
        :data="tableData"
        border
        stripe
        v-loading="loading"
        :empty-text="loading ? '加载中...' : '暂无数据'"
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="门店名称" width="200" />
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="address" label="门店地址" min-width="250" />
        <el-table-column label="位置信息" width="200">
          <template #default="{ row }">
            <div>经度：{{ row.longitude }}</div>
            <div>纬度：{{ row.latitude }}</div>
          </template>
        </el-table-column>
        <el-table-column label="营业时间" width="200">
          <template #default="{ row }">
            {{ row.business_hours || '09:00-21:00' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '营业中' : '已关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-if="pagination.total > 0"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :page-sizes="[10, 20, 50]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        class="pagination"
      />
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑门店' : '新增门店'" width="700px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="门店名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入门店名称" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="门店地址" prop="address">
          <el-input v-model="form.address" placeholder="请输入门店地址" />
        </el-form-item>
        <el-form-item label="详细地址" prop="detailed_address">
          <el-input v-model="form.detailed_address" type="textarea" :rows="2" placeholder="请输入详细地址" />
        </el-form-item>
        <el-form-item label="经度" prop="longitude">
          <el-input-number v-model="form.longitude" :precision="6" :step="0.000001" style="width: 100%" />
          <div style="margin-top: 5px; color: #999; font-size: 12px;">
            可使用地图工具获取坐标：
            <a href="https://lbs.qq.com/getPoint/" target="_blank">腾讯位置服务</a>
          </div>
        </el-form-item>
        <el-form-item label="纬度" prop="latitude">
          <el-input-number v-model="form.latitude" :precision="6" :step="0.000001" style="width: 100%" />
        </el-form-item>
        <el-form-item label="营业时间" prop="business_hours">
          <el-input v-model="form.business_hours" placeholder="例如：09:00-21:00" />
        </el-form-item>
        <el-form-item label="门店图片" prop="images">
          <el-upload
            class="upload-demo"
            action="#"
            :file-list="form.images"
            list-type="picture-card"
            :auto-upload="false"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <div style="margin-top: 5px; color: #999; font-size: 12px;">支持jpg、png格式，最多5张</div>
        </el-form-item>
        <el-form-item label="门店介绍" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="4" placeholder="请输入门店介绍" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">营业中</el-radio>
            <el-radio :label="0">已关闭</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
import {
  getStoreList,
  getStoreDetail,
  createStore,
  updateStore,
  deleteStore
} from '@/api/store'

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const submitting = ref(false)
const formRef = ref(null)

const searchForm = reactive({
  name: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const form = reactive({
  id: null,
  name: '',
  phone: '',
  address: '',
  detailed_address: '',
  longitude: 0,
  latitude: 0,
  business_hours: '09:00-21:00',
  images: [],
  description: '',
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入门店名称', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  address: [{ required: true, message: '请输入门店地址', trigger: 'blur' }],
  longitude: [{ required: true, message: '请输入经度', trigger: 'blur' }],
  latitude: [{ required: true, message: '请输入纬度', trigger: 'blur' }]
}

// 加载门店列表
const loadStores = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      ...searchForm
    }
    const data = await getStoreList(params)
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('加载门店列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadStores()
}

// 重置
const handleReset = () => {
  searchForm.name = ''
  searchForm.status = ''
  pagination.page = 1
  loadStores()
}

// 分页大小改变
const handleSizeChange = (size) => {
  pagination.page_size = size
  loadStores()
}

// 页码改变
const handlePageChange = (page) => {
  pagination.page = page
  loadStores()
}

// 新增门店
const handleAdd = () => {
  Object.assign(form, {
    id: null,
    name: '',
    phone: '',
    address: '',
    detailed_address: '',
    longitude: 0,
    latitude: 0,
    business_hours: '09:00-21:00',
    images: [],
    description: '',
    status: 1
  })
  dialogVisible.value = true
}

// 编辑门店
const handleEdit = (row) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

// 删除门店
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该门店吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteStore(row.id)
      ElMessage.success('删除成功')
      loadStores()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (form.id) {
          await updateStore(form.id, form)
        } else {
          await createStore(form)
        }
        ElMessage.success(form.id ? '更新成功' : '新增成功')
        dialogVisible.value = false
        loadStores()
      } catch (error) {
        ElMessage.error(error.message || '操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

onMounted(() => {
  loadStores()
})
</script>

<style scoped>
.stores-container {
  padding: 20px;
  min-height: 100%;
  box-sizing: border-box;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.upload-demo {
  width: 100%;
}
</style>
